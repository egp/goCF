// gcf_stream.go v10
package cf

import (
	"fmt"
	"math/big"
)

type GCFStreamOptions struct{}

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

type GCFStream struct {
	src GCFSource
	t   ULFT

	lower               *Rational // optional stable positive lower bound for unfinished tail
	tail                ContinuedFraction
	ingestedAny         bool
	prefixTerms         int
	lastEmitPrefixTerms int
	srcDone             bool
	done                bool
	err                 error
}

func NewGCFStream(src GCFSource, opts GCFStreamOptions) *GCFStream {
	s := &GCFStream{
		src: src,
		t: NewULFT(
			big.NewInt(1),
			big.NewInt(0),
			big.NewInt(0),
			big.NewInt(1),
		),
		lastEmitPrefixTerms: -1, // no emissions yet
	}

	if evSrc, ok := src.(TailEvidenceGCFSource); ok {
		ev := evSrc.TailEvidence()
		if ev.LowerBound != nil {
			lb := *ev.LowerBound
			s.lower = &lb
		}
	} else if bounded, ok := src.(PositiveTailLowerBoundedGCFSource); ok {
		lb := bounded.TailLowerBound()
		s.lower = &lb
	}

	return s
}

func (s *GCFStream) Err() error { return s.err }

func (s *GCFStream) canEmitFromCurrentPrefixEvidence() bool {
	return s.prefixTerms > s.lastEmitPrefixTerms
}

func (s *GCFStream) tailEvidence() (GCFTailEvidence, bool, error) {
	if evSrc, ok := s.src.(TailEvidenceGCFSource); ok {
		ev := evSrc.TailEvidence()
		if ev.RangeReusable && ev.Range == nil {
			return GCFTailEvidence{}, false, fmt.Errorf(
				"GCFStream: source %T provides reusable tail-range policy without a tail range",
				s.src,
			)
		}
		if ev.LowerBoundMinPrefix < 0 {
			return GCFTailEvidence{}, false, fmt.Errorf(
				"GCFStream: source %T provides negative LowerBoundMinPrefix=%d",
				s.src,
				ev.LowerBoundMinPrefix,
			)
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

func (s *GCFStream) explicitTailImageRange() (Range, bool, bool, error) {
	ev, ok, err := s.tailEvidence()
	if err != nil {
		return Range{}, false, false, err
	}
	if !ok || ev.Range == nil {
		return Range{}, false, false, nil
	}

	img, err := ev.Range.ApplyULFT(s.t)
	if err != nil {
		// A non-point reusable/explicit tail range crossing a pole means this
		// evidence is not currently usable for certification. It is not an
		// immediate hard stream error; callers may fall back or ingest more.
		if !isExactPointRange(*ev.Range) {
			return Range{}, false, ev.RangeReusable, nil
		}
		return Range{}, false, false, err
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

func (s *GCFStream) currentCertifiedTailDigit() (int64, bool, error) {
	if !s.ingestedAny {
		return 0, false, nil
	}

	if r, ok, reusable, err := s.explicitTailImageRange(); err != nil {
		return 0, false, err
	} else if ok {
		if !reusable && !s.canEmitFromCurrentPrefixEvidence() {
			return 0, false, nil
		}
		return certifiedFloorDigit(r)
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

func (s *GCFStream) emitCertifiedDigit(d int64) (int64, bool, bool) {
	nextT, err := EmitDigit(s.t, d)
	if err != nil {
		s.err = err
		s.done = true
		return 0, false, true
	}
	s.t = nextT
	s.lastEmitPrefixTerms = s.prefixTerms
	return d, true, true
}

func (s *GCFStream) maybeEmitFromTailMetadata() (int64, bool, bool) {
	d, ok, err := s.currentCertifiedTailDigit()
	if err != nil {
		s.err = err
		s.done = true
		return 0, false, true
	}
	if !ok {
		return 0, false, false
	}
	return s.emitCertifiedDigit(d)
}

func (s *GCFStream) finalizeFiniteTail() (int64, bool) {
	if !s.ingestedAny {
		s.err = fmt.Errorf("GCFStream: empty source")
		s.done = true
		return 0, false
	}

	y, err := applyULFTAtInfinity(s.t)
	if err != nil {
		s.err = err
		s.done = true
		return 0, false
	}

	s.tail = NewRationalCF(y)

	d, ok := s.tail.Next()
	if !ok {
		s.done = true
		return 0, false
	}
	return d, true
}

func (s *GCFStream) nextFromFiniteTail() (int64, bool, bool) {
	if s.tail == nil {
		return 0, false, false
	}

	d, ok := s.tail.Next()
	if !ok {
		s.done = true
		return 0, false, true
	}
	return d, true, true
}

func (s *GCFStream) ingestNextTerm() bool {
	p, q, ok := s.src.NextPQ()
	if !ok {
		s.srcDone = true
		return false
	}

	nextT, err := s.t.IngestGCF(p, q)
	if err != nil {
		s.err = err
		s.done = true
		return false
	}

	s.t = nextT
	s.ingestedAny = true
	s.prefixTerms++
	return true
}

func (s *GCFStream) advanceUntilDigitOrFinish() (int64, bool) {
	for {
		if s.srcDone {
			return s.finalizeFiniteTail()
		}

		if d, ok, handled := s.maybeEmitFromTailMetadata(); handled {
			return d, ok
		}

		if !s.ingestNextTerm() {
			if s.done || s.err != nil {
				return 0, false
			}
		}
	}
}

func (s *GCFStream) Next() (int64, bool) {
	if s.done {
		return 0, false
	}
	if s.err != nil {
		s.done = true
		return 0, false
	}

	if d, ok, handled := s.nextFromFiniteTail(); handled {
		return d, ok
	}

	return s.advanceUntilDigitOrFinish()
}

func applyULFTAtInfinity(t ULFT) (Rational, error) {
	if err := t.Validate(); err != nil {
		return Rational{}, err
	}

	if t.C.Sign() == 0 {
		return Rational{}, fmt.Errorf("applyULFTAtInfinity: undefined for C=0 in %v", t)
	}

	return newRationalBig(new(big.Int).Set(t.A), new(big.Int).Set(t.C))
}

// gcf_stream.go v10
