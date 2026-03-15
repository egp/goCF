// reciprocal_stream_gcf_exact_tail.go v1
package cf

import "fmt"

type ReciprocalApproxStreamSnapshot struct {
	Started     bool
	PrefixTerms int
	Approx      *Rational
}

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
	err     error
	done    bool
	started bool
	exactCF ContinuedFraction
	approx  *Rational

	src            GCFSource
	tailSrc        GCFTailSource
	maxIngestTerms int
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

func (s *ReciprocalGCFExactTailStream2) Err() error { return s.err }

func (s *ReciprocalGCFExactTailStream2) Snapshot() ReciprocalApproxStreamSnapshot {
	var approxCopy *Rational
	if s.approx != nil {
		v := *s.approx
		approxCopy = &v
	}
	return ReciprocalApproxStreamSnapshot{
		Started:     s.started,
		PrefixTerms: s.maxIngestTerms,
		Approx:      approxCopy,
	}
}

func (s *ReciprocalGCFExactTailStream2) initExactCF() bool {
	if s.started {
		return s.err == nil
	}
	s.started = true

	tail, ok := s.tailSrc.ExactTail()
	if !ok {
		s.err = fmt.Errorf("ReciprocalGCFExactTailStream2: tail evidence not implemented")
		s.done = true
		return false
	}

	x, _, err := EvalGCFWithTailExact(s.src, tail, s.maxIngestTerms)
	if err != nil {
		s.err = fmt.Errorf("ReciprocalGCFExactTailStream2: %w", err)
		s.done = true
		return false
	}
	if x.Cmp(intRat(0)) == 0 {
		s.err = fmt.Errorf("ReciprocalGCFExactTailStream2: reciprocal of zero")
		s.done = true
		return false
	}

	recip, err := intRat(1).Div(x)
	if err != nil {
		s.err = fmt.Errorf("ReciprocalGCFExactTailStream2: %w", err)
		s.done = true
		return false
	}

	s.approx = &recip
	s.exactCF = NewRationalCF(recip)
	return true
}

func (s *ReciprocalGCFExactTailStream2) Next() (int64, bool) {
	if s.done {
		return 0, false
	}
	if s.err != nil {
		s.done = true
		return 0, false
	}
	if !s.initExactCF() {
		return 0, false
	}

	d, ok := s.exactCF.Next()
	if !ok {
		s.done = true
		return 0, false
	}
	return d, true
}

// reciprocal_stream_gcf_exact_tail.go v1
