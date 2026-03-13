// gcf_stream_tail_evidence.go v1
package cf

import "fmt"

type LowerBoundRayMinPrefixGCFSource interface {
	LowerBoundRayMinPrefix() int
}

type TailRangedGCFSource interface {
	TailRange() Range
}

type ReusableTailRangeGCFSource interface {
	TailRangeReusable() bool
}

type GCFTailEvidence struct {
	LowerBound          *Rational
	Range               *Range
	RangeReusable       bool
	LowerBoundMinPrefix int
}

type TailEvidenceGCFSource interface {
	TailEvidence() GCFTailEvidence
}

type PostEmitTailEvidenceGCFSource interface {
	PostEmitTailEvidence(emittedDigit int64) (GCFTailEvidence, bool)
}

type RefinedTailEvidenceGCFSource interface {
	RefinedTailEvidence() (GCFTailEvidence, bool)
}

func validateTailEvidence(owner string, ev GCFTailEvidence) error {
	if ev.RangeReusable && ev.Range == nil {
		return fmt.Errorf("%s provides reusable tail-range policy without a tail range", owner)
	}
	if ev.LowerBoundMinPrefix < 0 {
		return fmt.Errorf("%s provides negative LowerBoundMinPrefix=%d", owner, ev.LowerBoundMinPrefix)
	}
	return nil
}

func (s *GCFStream) tailEvidence() (GCFTailEvidence, bool, error) {
	if s.tailEvidenceOverride != nil {
		ev := *s.tailEvidenceOverride
		if err := validateTailEvidence(fmt.Sprintf("GCFStream override from source %T", s.src), ev); err != nil {
			return GCFTailEvidence{}, false, err
		}
		return ev, true, nil
	}

	if evSrc, ok := s.src.(TailEvidenceGCFSource); ok {
		ev := evSrc.TailEvidence()
		if err := validateTailEvidence(fmt.Sprintf("GCFStream: source %T", s.src), ev); err != nil {
			return GCFTailEvidence{}, false, err
		}
		return ev, true, nil
	}

	var ev GCFTailEvidence

	if ranged, ok := s.src.(TailRangedGCFSource); ok {
		r := ranged.TailRange()
		ev.Range = &r
	}
	if reusable, ok := s.src.(ReusableTailRangeGCFSource); ok {
		ev.RangeReusable = reusable.TailRangeReusable()
		if ev.Range == nil && ev.RangeReusable {
			return GCFTailEvidence{}, false, fmt.Errorf(
				"GCFStream: source %T provides TailRangeReusable without TailRange",
				s.src,
			)
		}
	}
	if bounded, ok := s.src.(PositiveTailLowerBoundedGCFSource); ok {
		lb := bounded.TailLowerBound()
		ev.LowerBound = &lb
	}
	if delayed, ok := s.src.(LowerBoundRayMinPrefixGCFSource); ok {
		ev.LowerBoundMinPrefix = delayed.LowerBoundRayMinPrefix()
		if ev.LowerBoundMinPrefix < 0 {
			return GCFTailEvidence{}, false, fmt.Errorf(
				"GCFStream: source %T provides negative LowerBoundRayMinPrefix=%d",
				s.src,
				ev.LowerBoundMinPrefix,
			)
		}
	}

	if ev.Range == nil && ev.LowerBound == nil {
		return GCFTailEvidence{}, false, nil
	}
	return ev, true, nil
}

func isExactPointRange(r Range) bool {
	return r.Lo.Cmp(r.Hi) == 0
}

func classifyExplicitTailRangeFailure(r Range, err error) (usable bool, fatal error) {
	if err == nil {
		return true, nil
	}
	if !isExactPointRange(r) {
		return false, nil
	}
	return false, err
}

func (s *GCFStream) explicitTailImageRange() (Range, bool, bool, error) {
	ev, ok, err := s.tailEvidence()
	if err != nil {
		return Range{}, false, false, err
	}
	if !ok || ev.Range == nil {
		return Range{}, false, false, nil
	}

	img, applyErr := ev.Range.ApplyULFT(s.t)
	usable, fatalErr := classifyExplicitTailRangeFailure(*ev.Range, applyErr)
	if fatalErr != nil {
		return Range{}, false, false, fatalErr
	}
	if !usable {
		return Range{}, false, ev.RangeReusable, nil
	}

	return img, true, ev.RangeReusable, nil
}

func (s *GCFStream) lowerBoundRayImageRange() (Range, bool, error) {
	if !s.canUseGenericLowerBoundEmission() {
		return Range{}, false, nil
	}

	img, err := ApplyULFTToTailRay(s.t, *s.lower)
	if err != nil {
		return Range{}, false, err
	}
	return img, true, nil
}

func (s *GCFStream) canUseGenericLowerBoundEmission() bool {
	if s.lower == nil {
		return false
	}

	minPrefix := 0
	if ev, ok, err := s.tailEvidence(); err == nil && ok {
		minPrefix = ev.LowerBoundMinPrefix
	}
	return s.prefixTerms >= minPrefix
}

func certifiedFloorDigit(r Range) (int64, bool, error) {
	lo, hi, err := r.FloorBounds()
	if err != nil {
		return 0, false, err
	}
	if lo != hi {
		return 0, false, nil
	}
	return lo, true, nil
}

func (s *GCFStream) tryRefinedTailEvidence() (used bool, err error) {
	refiner, ok := s.src.(RefinedTailEvidenceGCFSource)
	if !ok {
		return false, nil
	}

	ev, ok := refiner.RefinedTailEvidence()
	if !ok {
		return false, nil
	}

	if err := validateTailEvidence(fmt.Sprintf("GCFStream refined evidence from source %T", s.src), ev); err != nil {
		return false, err
	}

	s.tailEvidenceOverride = &ev
	s.tailEvidenceFresh = true
	if ev.LowerBound != nil {
		lb := *ev.LowerBound
		s.lower = &lb
	} else {
		s.lower = nil
	}
	return true, nil
}

func (s *GCFStream) tryRefinementsUntilCertified(maxSteps int) (int64, bool, error) {
	for i := 0; i < maxSteps; i++ {
		used, err := s.tryRefinedTailEvidence()
		if err != nil {
			return 0, false, err
		}
		if !used {
			return 0, false, nil
		}

		r, ok, _, err := s.explicitTailImageRange()
		if err != nil {
			return 0, false, err
		}
		if !ok {
			continue
		}

		s.tailEvidenceFresh = false
		d, certified, err := certifiedFloorDigit(r)
		if err != nil {
			return 0, false, err
		}
		if certified {
			return d, true, nil
		}
	}
	return 0, false, nil
}

func (s *GCFStream) currentCertifiedTailDigit() (int64, bool, error) {
	if !s.ingestedAny {
		return 0, false, nil
	}

	if r, ok, reusable, err := s.explicitTailImageRange(); err != nil {
		return 0, false, err
	} else if ok {
		if s.tailEvidenceOverride != nil && s.tailEvidenceFresh {
			s.tailEvidenceFresh = false
			return certifiedFloorDigit(r)
		}
		if !reusable && !s.canEmitFromCurrentPrefixEvidence() {
			if d, ok, err := s.tryRefinementsUntilCertified(s.maxRefinementSteps); err != nil {
				return 0, false, err
			} else if ok {
				return d, true, nil
			}
			return 0, false, nil
		}
		return certifiedFloorDigit(r)
	}

	if d, ok, err := s.tryRefinementsUntilCertified(s.maxRefinementSteps); err != nil {
		return 0, false, err
	} else if ok {
		return d, true, nil
	}

	if !s.canEmitFromCurrentPrefixEvidence() {
		return 0, false, nil
	}

	if r, ok, err := s.lowerBoundRayImageRange(); err != nil {
		return 0, false, err
	} else if ok {
		return certifiedFloorDigit(r)
	}

	return 0, false, nil
}

// gcf_stream_tail_evidence.go v1
