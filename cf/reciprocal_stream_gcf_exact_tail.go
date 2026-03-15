// reciprocal_stream_gcf_exact_tail.go v6
package cf

import "fmt"

type ReciprocalApproxStream interface {
	ContinuedFraction
	Err() error
	Snapshot() ReciprocalApproxStreamSnapshot
}

// ReciprocalGCFExactTailStream2 is an inspectable unary reciprocal stream over a
// generalized continued-fraction source with exact tail evidence.
//
// Current milestone:
//   - bounded GCF ingestion
//   - exact tail evidence only
//   - reciprocal computed by exact rational collapse
//
// Future work:
//   - richer tail evidence
//   - progressive certification
type ReciprocalGCFExactTailStream2 struct {
	core reciprocalExactCollapseCore

	src            GCFSource
	tailSrc        GCFTailSource
	maxIngestTerms int
	consumedTerms  int
}

func NewReciprocalGCFExactTailStream2(src GCFSource, tailSrc GCFTailSource, maxIngestTerms int) *ReciprocalGCFExactTailStream2 {
	return &ReciprocalGCFExactTailStream2{
		src:            src,
		tailSrc:        tailSrc,
		maxIngestTerms: maxIngestTerms,
	}
}

func NewReciprocalGCFExactTailStreamWithTail2(src GCFSource, tail Rational, maxIngestTerms int) *ReciprocalGCFExactTailStream2 {
	return NewReciprocalGCFExactTailStream2(src, NewExactTailSource(tail), maxIngestTerms)
}

func (s *ReciprocalGCFExactTailStream2) Err() error { return s.core.Err() }

func (s *ReciprocalGCFExactTailStream2) unaryClass() unaryStreamClass {
	return unaryStreamClass{
		Operator: unaryOperatorReciprocal,
		Input:    unaryInputGCFExact,
		Progress: unaryProgressExactCollapse,
	}
}

func (s *ReciprocalGCFExactTailStream2) Snapshot() ReciprocalApproxStreamSnapshot {
	var approxCopy *Rational
	if s.core.approx != nil {
		v := *s.core.approx
		approxCopy = &v
	}
	return ReciprocalApproxStreamSnapshot{
		Started:        s.core.state.started,
		Approx:         approxCopy,
		MaxIngestTerms: s.maxIngestTerms,
		ConsumedTerms:  s.consumedTerms,
	}
}

func (s *ReciprocalGCFExactTailStream2) evalReciprocal() (Rational, error) {
	tail, ok := s.tailSrc.ExactTail()
	if !ok {
		return Rational{}, fmt.Errorf("ReciprocalGCFExactTailStream2: tail evidence not implemented")
	}

	x, consumed, err := EvalGCFWithTailExact(s.src, tail, s.maxIngestTerms)
	if err != nil {
		return Rational{}, fmt.Errorf("ReciprocalGCFExactTailStream2: %w", err)
	}
	s.consumedTerms = consumed

	if x.Cmp(intRat(0)) == 0 {
		return Rational{}, fmt.Errorf("ReciprocalGCFExactTailStream2: reciprocal of zero")
	}

	recip, err := intRat(1).Div(x)
	if err != nil {
		return Rational{}, fmt.Errorf("ReciprocalGCFExactTailStream2: %w", err)
	}
	return recip, nil
}

func (s *ReciprocalGCFExactTailStream2) Next() (int64, bool) {
	return s.core.Next(s.evalReciprocal)
}

// reciprocal_stream_gcf_exact_tail.go v6
