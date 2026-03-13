// gcf_stream.go v16
package cf

import (
	"fmt"
	"math/big"
)

type GCFStreamOptions struct{}

type GCFStream struct {
	src GCFSource
	t   ULFT

	lower                *Rational
	tailEvidenceOverride *GCFTailEvidence
	tailEvidenceFresh    bool
	tail                 ContinuedFraction
	ingestedAny          bool
	prefixTerms          int
	lastEmitPrefixTerms  int
	srcDone              bool
	done                 bool
	err                  error
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
		lastEmitPrefixTerms: -1,
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

func (s *GCFStream) applyPostEmitTailEvidence(d int64) error {
	postSrc, ok := s.src.(PostEmitTailEvidenceGCFSource)
	if !ok {
		s.tailEvidenceOverride = nil
		s.tailEvidenceFresh = false
		return nil
	}

	ev, ok := postSrc.PostEmitTailEvidence(d)
	if !ok {
		s.tailEvidenceOverride = nil
		s.tailEvidenceFresh = false
		return nil
	}

	if err := validateTailEvidence(fmt.Sprintf("GCFStream post-emit evidence from source %T", s.src), ev); err != nil {
		return err
	}

	s.tailEvidenceOverride = &ev
	s.tailEvidenceFresh = true
	if ev.LowerBound != nil {
		lb := *ev.LowerBound
		s.lower = &lb
	} else {
		s.lower = nil
	}
	return nil
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

	if err := s.applyPostEmitTailEvidence(d); err != nil {
		s.err = err
		s.done = true
		return 0, false, true
	}

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
	s.tailEvidenceOverride = nil
	s.tailEvidenceFresh = false
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

// gcf_stream.go v16
