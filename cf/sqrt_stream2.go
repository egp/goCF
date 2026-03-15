// sqrt_stream2.go v3
package cf

import "fmt"

// SqrtGCFExactTailStream2 is the new-path sqrt stream over a bounded GCF prefix
// plus exact tail evidence.
//
// Current milestone:
//   - bounded ingestion of a GCFSource
//   - exact tail evidence only
//   - delegates to the canonical bounded sqrt approximation path
//
// Future work:
//   - stronger tail evidence
//   - true progressive certification instead of exact collapse first
type SqrtGCFExactTailStream2 struct {
	err     error
	done    bool
	started bool
	exactCF ContinuedFraction
	approx  *Rational

	src            GCFSource
	tailSrc        GCFTailSource
	maxIngestTerms int
	policy         SqrtPolicy2
}

func NewSqrtGCFExactTailStream2(src GCFSource, tailSrc GCFTailSource, maxIngestTerms int, p SqrtPolicy2) *SqrtGCFExactTailStream2 {
	return &SqrtGCFExactTailStream2{
		src:            src,
		tailSrc:        tailSrc,
		maxIngestTerms: maxIngestTerms,
		policy:         p,
	}
}

func NewSqrtGCFExactTailStreamWithTail2(src GCFSource, tail Rational, maxIngestTerms int, p SqrtPolicy2) *SqrtGCFExactTailStream2 {
	return NewSqrtGCFExactTailStream2(src, NewExactTailSource(tail), maxIngestTerms, p)
}

func (s *SqrtGCFExactTailStream2) Err() error { return s.err }

func (s *SqrtGCFExactTailStream2) Snapshot() SqrtApproxStreamSnapshot {
	var approxCopy *Rational
	if s.approx != nil {
		v := *s.approx
		approxCopy = &v
	}
	return SqrtApproxStreamSnapshot{
		Started:     s.started,
		PrefixTerms: s.maxIngestTerms,
		Approx:      approxCopy,
	}
}

func (s *SqrtGCFExactTailStream2) initExactCF() bool {
	if s.started {
		return s.err == nil
	}
	s.started = true

	tail, ok := s.tailSrc.ExactTail()
	if !ok {
		s.err = fmt.Errorf("SqrtGCFExactTailStream2: tail evidence not implemented")
		s.done = true
		return false
	}

	approx, err := SqrtApproxFromGCFWithTail2(s.src, tail, s.maxIngestTerms, s.policy)
	if err != nil {
		s.err = fmt.Errorf("SqrtGCFExactTailStream2: %w", err)
		s.done = true
		return false
	}

	s.approx = &approx
	s.exactCF = NewRationalCF(approx)
	return true
}

func (s *SqrtGCFExactTailStream2) Next() (int64, bool) {
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

// Legacy transitional aliases retained for now.
type SqrtGCFStream2 = SqrtGCFExactTailStream2

func NewSqrtGCFStream2(src GCFSource, tailSrc GCFTailSource, maxIngestTerms int, p SqrtPolicy2) *SqrtGCFStream2 {
	return NewSqrtGCFExactTailStream2(src, tailSrc, maxIngestTerms, p)
}

func NewSqrtGCFStreamWithTail2(src GCFSource, tail Rational, maxIngestTerms int, p SqrtPolicy2) *SqrtGCFStream2 {
	return NewSqrtGCFExactTailStreamWithTail2(src, tail, maxIngestTerms, p)
}

// sqrt_stream2.go v3
