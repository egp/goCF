// gcf_stream.go v4
package cf

import (
	"fmt"
	"math/big"
)

type GCFStreamOptions struct{}

type GCFStream struct {
	src GCFSource
	t   ULFT

	lower               *Rational // optional stable positive lower bound for unfinished tail
	tailRange           *Range    // optional stable explicit unfinished-tail range
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

	if bounded, ok := src.(PositiveTailLowerBoundedGCFSource); ok {
		lb := bounded.TailLowerBound()
		s.lower = &lb
	}

	// Generic stable explicit tail range, if the source can provide one.
	if ranged, ok := src.(interface{ TailRange() Range }); ok {
		r := ranged.TailRange()
		s.tailRange = &r
	}

	return s
}

func (s *GCFStream) Err() error { return s.err }

func (s *GCFStream) canEmitFromCurrentPrefixEvidence() bool {
	return s.prefixTerms > s.lastEmitPrefixTerms
}

func (s *GCFStream) currentTailImageRange() (Range, bool, error) {
	// First preference: generic stable explicit unfinished-tail range.
	if s.tailRange != nil {
		img, err := s.tailRange.ApplyULFT(s.t)
		if err != nil {
			return Range{}, false, err
		}
		return img, true, nil
	}

	// Named-source prefix-aware helpers were useful during exploration.
	// Keep them here for now, but commented out while the engine moves toward
	// generic ingestion semantics.
	//
	// if r, ok, err := s.specializedTailRange(); err != nil {
	// 	return Range{}, false, err
	// } else if ok {
	// 	img, err := r.ApplyULFT(s.t)
	// 	if err != nil {
	// 		return Range{}, false, err
	// 	}
	// 	return img, true, nil
	// }

	// Fallback: generic positive lower-bound ray, but only when we trust it
	// enough for the current source/prefix depth.
	if s.canUseGenericLowerBoundEmission() {
		img, err := ApplyULFTToTailRay(s.t, *s.lower)
		if err != nil {
			return Range{}, false, err
		}
		return img, true, nil
	}

	return Range{}, false, nil
}

func (s *GCFStream) specializedTailRange() (Range, bool, error) {
	switch s.src.(type) {
	case *LambertPiOver4GCFSource:
		return LambertPiOver4TailRangeAfterPrefix(s.prefixTerms)
	case *Brouncker4OverPiGCFSource:
		// Not yet enabled here.
		//
		// The currently available Brouncker prefix-aware tail interval helpers are
		// useful for bounded-prefix approximation work, but they are not yet
		// strong/sound enough for early digit emission in GCFStream. Falling back
		// to the generic lower-bound ray is conservative and avoids unsound digits.
		return Range{}, false, nil
	default:
		return Range{}, false, nil
	}
}

func (s *GCFStream) canUseGenericLowerBoundEmission() bool {
	if s.lower == nil {
		return false
	}

	switch s.src.(type) {
	case *Brouncker4OverPiGCFSource:
		// A single ingested Brouncker term is not enough for sound early emission
		// from the generic lower-bound ray. Require a little more prefix before
		// falling back to that weaker enclosure.
		return s.prefixTerms >= 2
	default:
		return true
	}
}
func (s *GCFStream) maybeEmitFromTailMetadata() (int64, bool, bool) {
	if !s.ingestedAny {
		return 0, false, false
	}

	// Do not emit multiple ordinary CF digits from the same GCF prefix evidence.
	// After each metadata-based emission, require at least one more ingested GCF
	// term before trying to emit again.
	if !s.canEmitFromCurrentPrefixEvidence() {
		return 0, false, false
	}

	r, ok, err := s.currentTailImageRange()
	if err != nil {
		s.err = err
		s.done = true
		return 0, false, true
	}
	if !ok {
		return 0, false, false
	}

	lo, hi, err := r.FloorBounds()
	if err != nil {
		s.err = err
		s.done = true
		return 0, false, true
	}
	if lo != hi {
		return 0, false, false
	}

	d := lo
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

func (s *GCFStream) Next() (int64, bool) {
	if s.done {
		return 0, false
	}
	if s.err != nil {
		s.done = true
		return 0, false
	}

	// Once finite-source exhaustion has been converted to an exact rational tail,
	// just delegate to that tail stream.
	if s.tail != nil {
		d, ok := s.tail.Next()
		if !ok {
			s.done = true
			return 0, false
		}
		return d, true
	}

	for {
		// Finite source exhausted: exact remaining value is T(∞).
		if s.srcDone {
			return s.finalizeFiniteTail()
		}

		// Prefer strong prefix-aware tail metadata when available; otherwise use
		// the generic lower-bound ray fallback.
		if d, ok, handled := s.maybeEmitFromTailMetadata(); handled {
			return d, ok
		}

		// Need more GCF input.
		p, q, ok := s.src.NextPQ()
		if !ok {
			s.srcDone = true
			continue
		}

		nextT, err := s.t.IngestGCF(p, q)
		if err != nil {
			s.err = err
			s.done = true
			return 0, false
		}
		s.t = nextT
		s.ingestedAny = true
		s.prefixTerms++
	}
}

func applyULFTAtInfinity(t ULFT) (Rational, error) {
	if err := t.Validate(); err != nil {
		return Rational{}, err
	}

	// For T(x) = (A x + B) / (C x + D), the value at infinity is A/C.
	if t.C.Sign() == 0 {
		return Rational{}, fmt.Errorf("applyULFTAtInfinity: undefined for C=0 in %v", t)
	}

	return newRationalBig(new(big.Int).Set(t.A), new(big.Int).Set(t.C))
}

// gcf_stream.go v4
