// gcf_stream.go v5
package cf

import (
	"fmt"
	"math/big"
)

type GCFStreamOptions struct{}

type LowerBoundRayMinPrefixGCFSource interface {
	LowerBoundRayMinPrefix() int
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

	if bounded, ok := src.(PositiveTailLowerBoundedGCFSource); ok {
		lb := bounded.TailLowerBound()
		s.lower = &lb
	}

	return s
}

func (s *GCFStream) Err() error { return s.err }

func (s *GCFStream) canEmitFromCurrentPrefixEvidence() bool {
	return s.prefixTerms > s.lastEmitPrefixTerms
}

func (s *GCFStream) currentTailImageRange() (Range, bool, error) {
	if ranged, ok := s.src.(interface{ TailRange() Range }); ok {
		r := ranged.TailRange()
		img, err := r.ApplyULFT(s.t)
		if err != nil {
			return Range{}, false, err
		}
		return img, true, nil
	}

	if s.canUseGenericLowerBoundEmission() {
		img, err := ApplyULFTToTailRay(s.t, *s.lower)
		if err != nil {
			return Range{}, false, err
		}
		return img, true, nil
	}

	return Range{}, false, nil
}

func (s *GCFStream) canUseGenericLowerBoundEmission() bool {
	if s.lower == nil {
		return false
	}

	minPrefix := 0
	if delayed, ok := s.src.(LowerBoundRayMinPrefixGCFSource); ok {
		minPrefix = delayed.LowerBoundRayMinPrefix()
	}
	return s.prefixTerms >= minPrefix
}
func (s *GCFStream) currentCertifiedTailDigit() (int64, bool, error) {
	if !s.ingestedAny {
		return 0, false, nil
	}
	if !s.canEmitFromCurrentPrefixEvidence() {
		return 0, false, nil
	}

	r, ok, err := s.currentTailImageRange()
	if err != nil || !ok {
		return 0, false, err
	}

	lo, hi, err := r.FloorBounds()
	if err != nil {
		return 0, false, err
	}
	if lo != hi {
		return 0, false, nil
	}

	return lo, true, nil
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

// gcf_stream.go v5
