// sqrt_gcf_prefix_stream2.go v5
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

	inputApprox *GCFApprox
	status      SqrtStreamStatus

	src         GCFSource
	prefixTerms int
	policy      SqrtPolicy2
}

func NewSqrtGCFPrefixStream2(src GCFSource, prefixTerms int, p SqrtPolicy2) *SqrtGCFPrefixStream2 {
	return &SqrtGCFPrefixStream2{
		src:         src,
		prefixTerms: prefixTerms,
		policy:      p,
		status:      SqrtStreamStatusUnstarted,
	}
}

func (s *SqrtGCFPrefixStream2) Err() error { return s.err }

func (s *SqrtGCFPrefixStream2) Snapshot() SqrtApproxStreamSnapshot {
	var approxCopy *Rational
	if s.approx != nil {
		v := *s.approx
		approxCopy = &v
	}

	var gcfInputApproxCopy *GCFApprox
	if s.inputApprox != nil {
		v := *s.inputApprox
		gcfInputApproxCopy = &v
	}

	return SqrtApproxStreamSnapshot{
		Status:         s.status,
		Started:        s.started,
		PrefixTerms:    s.prefixTerms,
		Approx:         approxCopy,
		GCFInputApprox: gcfInputApproxCopy,
	}
}

func (s *SqrtGCFPrefixStream2) initExactCF() bool {
	if s.started {
		return s.err == nil
	}
	s.started = true

	a, err := GCFApproxFromPrefix(s.src, s.prefixTerms)
	if err != nil {
		s.err = fmt.Errorf("SqrtGCFPrefixStream2: %w", err)
		s.done = true
		s.status = SqrtStreamStatusFailed
		return false
	}
	s.inputApprox = &a

	approx, err := SqrtApproxFromGCFApproxRangeSeed2(a, s.policy)
	if err != nil {
		s.err = fmt.Errorf("SqrtGCFPrefixStream2: %w", err)
		s.done = true
		s.status = SqrtStreamStatusFailed
		return false
	}

	s.approx = &approx
	s.exactCF = NewRationalCF(approx)

	if a.Range != nil && a.Range.Lo.Cmp(a.Range.Hi) == 0 {
		s.status = SqrtStreamStatusExactInput
	} else {
		s.status = SqrtStreamStatusBoundedCollapse
	}
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

// sqrt_gcf_prefix_stream2.go v5
