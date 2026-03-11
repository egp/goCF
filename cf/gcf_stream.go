// gcf_stream.go v3
package cf

import (
	"fmt"
	"math/big"
)

type GCFStreamOptions struct{}

type GCFStream struct {
	src GCFSource
	t   ULFT

	tail    ContinuedFraction
	srcDone bool
	done    bool
	err     error
}

func NewGCFStream(src GCFSource, opts GCFStreamOptions) *GCFStream {
	return &GCFStream{
		src: src,
		t: NewULFT(
			big.NewInt(1),
			big.NewInt(0),
			big.NewInt(0),
			big.NewInt(1),
		),
	}
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

	if s.tail == nil {
		if err := s.finishFiniteSourceToTail(); err != nil {
			s.err = err
			s.done = true
			return 0, false
		}
	}

	if s.tail == nil {
		s.err = fmt.Errorf("GCFStream: no tail")
		s.done = true
		return 0, false
	}

	d, ok := s.tail.Next()
	if !ok {
		s.done = true
		return 0, false
	}
	return d, true
}

func (s *GCFStream) finishFiniteSourceToTail() error {
	if s.srcDone {
		return nil
	}

	ingestedAny := false

	for {
		p, q, ok := s.src.NextPQ()
		if !ok {
			break
		}
		ingestedAny = true

		nextT, err := s.t.IngestGCF(p, q)
		if err != nil {
			return err
		}
		s.t = nextT
	}

	if !ingestedAny {
		return fmt.Errorf("GCFStream: empty source")
	}

	// Finite-tail GCF semantics:
	// after ingesting all terms into the composed ULFT T, the last term contributes
	// just p_last, which corresponds to evaluating the remaining tail at infinity.
	y, err := applyULFTAtInfinity(s.t)
	if err != nil {
		return err
	}

	s.tail = NewRationalCF(y)
	s.srcDone = true
	return nil
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

// gcf_stream.go v3
