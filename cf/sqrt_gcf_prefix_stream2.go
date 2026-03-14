// sqrt_gcf_prefix_stream2.go v3
package cf

import "fmt"

// SqrtGCFPrefixStream2 is a bounded sqrt stream built from a finite GCF prefix.
//
// Current milestone:
//   - ingest up to prefixTerms from a GCFSource
//   - use GCFApprox range when available to derive a seed
//   - delegate to the canonical bounded sqrt approximation path
//
// Future work:
//   - stronger unfinished-tail handling
//   - true progressive certification
type SqrtGCFPrefixStream2 struct {
	err     error
	done    bool
	started bool
	exactCF ContinuedFraction
	approx  *Rational

	src         GCFSource
	prefixTerms int
	policy      SqrtPolicy2
}

type SqrtGCFPrefixStream2Snapshot struct {
	Started     bool
	PrefixTerms int
	Approx      *Rational
}

func NewSqrtGCFPrefixStream2(src GCFSource, prefixTerms int, p SqrtPolicy2) *SqrtGCFPrefixStream2 {
	return &SqrtGCFPrefixStream2{
		src:         src,
		prefixTerms: prefixTerms,
		policy:      p,
	}
}

func (s *SqrtGCFPrefixStream2) Err() error { return s.err }

func (s *SqrtGCFPrefixStream2) Snapshot() SqrtGCFPrefixStream2Snapshot {
	var approxCopy *Rational
	if s.approx != nil {
		v := *s.approx
		approxCopy = &v
	}
	return SqrtGCFPrefixStream2Snapshot{
		Started:     s.started,
		PrefixTerms: s.prefixTerms,
		Approx:      approxCopy,
	}
}

func (s *SqrtGCFPrefixStream2) initExactCF() bool {
	if s.started {
		return s.err == nil
	}
	s.started = true

	approx, err := SqrtApproxFromGCFSourceRangeSeed2(s.src, s.prefixTerms, s.policy)
	if err != nil {
		s.err = fmt.Errorf("SqrtGCFPrefixStream2: %w", err)
		s.done = true
		return false
	}

	s.approx = &approx
	s.exactCF = NewRationalCF(approx)
	return true
}

func (s *SqrtGCFPrefixStream2) Next() (int64, bool) {
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

// sqrt_gcf_prefix_stream2.go v3
