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

	lower *Rational // optional stable positive lower bound for unfinished tail
	tail  ContinuedFraction

	ingestedAny bool
	srcDone     bool
	done        bool
	err         error
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
	}

	if bounded, ok := src.(PositiveTailLowerBoundedGCFSource); ok {
		lb := bounded.TailLowerBound()
		s.lower = &lb
	}

	return s
}

func (s *GCFStream) Err() error { return s.err }

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

		// If we have a stable lower bound for the unfinished tail, try to emit.
		if s.lower != nil && s.ingestedAny {
			r, err := ApplyULFTToTailRay(s.t, *s.lower)
			if err == nil {
				lo, hi, err := r.FloorBounds()
				if err == nil && lo == hi {
					d := lo
					nextT, err := EmitDigit(s.t, d)
					if err != nil {
						s.err = err
						s.done = true
						return 0, false
					}
					s.t = nextT
					return d, true
				}
			}
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
