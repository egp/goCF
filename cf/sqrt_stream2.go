// sqrt_stream2.go v2
package cf

import "fmt"

// SqrtGCFStream2 is the new-path sqrt stream over a bounded GCF prefix plus
// tail evidence.
//
// Current milestone:
//   - bounded ingestion of a GCFSource
//   - exact tail evidence only
//   - delegates to the canonical bounded sqrt approximation path
//
// Future work:
//   - stronger tail evidence
//   - true progressive certification instead of exact collapse first
type SqrtGCFStream2 struct {
	err     error
	done    bool
	started bool
	exactCF ContinuedFraction

	src            GCFSource
	tailSrc        GCFTailSource
	maxIngestTerms int
	policy         SqrtPolicy2
}

func NewSqrtGCFStream2(src GCFSource, tailSrc GCFTailSource, maxIngestTerms int, p SqrtPolicy2) *SqrtGCFStream2 {
	return &SqrtGCFStream2{
		src:            src,
		tailSrc:        tailSrc,
		maxIngestTerms: maxIngestTerms,
		policy:         p,
	}
}

func NewSqrtGCFStreamWithTail2(src GCFSource, tail Rational, maxIngestTerms int, p SqrtPolicy2) *SqrtGCFStream2 {
	return NewSqrtGCFStream2(src, NewExactTailSource(tail), maxIngestTerms, p)
}

func (s *SqrtGCFStream2) Err() error { return s.err }

func (s *SqrtGCFStream2) initExactCF() bool {
	if s.started {
		return s.err == nil
	}
	s.started = true

	tail, ok := s.tailSrc.ExactTail()
	if !ok {
		s.err = fmt.Errorf("SqrtGCFStream2: tail evidence not implemented")
		s.done = true
		return false
	}

	cf, err := SqrtApproxCFFromGCFWithTail2(s.src, tail, s.maxIngestTerms, s.policy)
	if err != nil {
		s.err = fmt.Errorf("SqrtGCFStream2: %w", err)
		s.done = true
		return false
	}

	s.exactCF = cf
	return true
}

func (s *SqrtGCFStream2) Next() (int64, bool) {
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

// sqrt_stream2.go v2
