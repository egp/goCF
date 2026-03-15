// sqrt_stream_cf2.go v5
package cf

import "fmt"

// SqrtCFPrefixStream2 is a bounded sqrt stream over a continued-fraction source.
//
// Current milestone:
//   - ingest up to prefixTerms from a ContinuedFraction
//   - derive a seed from the resulting range when available
//   - delegate to the canonical bounded sqrt approximation path
//
// Future work:
//   - stronger progressive certification
//   - less exact-collapse, more genuine streaming
type SqrtCFPrefixStream2 struct {
	err     error
	done    bool
	started bool
	exactCF ContinuedFraction
	approx  *Rational

	inputApprox *CFApprox

	src         ContinuedFraction
	prefixTerms int
	policy      SqrtPolicy2
}

func NewSqrtCFPrefixStream2(src ContinuedFraction, prefixTerms int, p SqrtPolicy2) *SqrtCFPrefixStream2 {
	return &SqrtCFPrefixStream2{
		src:         src,
		prefixTerms: prefixTerms,
		policy:      p,
	}
}

func (s *SqrtCFPrefixStream2) Err() error { return s.err }

func (s *SqrtCFPrefixStream2) Snapshot() SqrtApproxStreamSnapshot {
	var approxCopy *Rational
	if s.approx != nil {
		v := *s.approx
		approxCopy = &v
	}

	var cfInputApproxCopy *CFApprox
	if s.inputApprox != nil {
		v := *s.inputApprox
		cfInputApproxCopy = &v
	}

	status := SqrtStreamStatusUnstarted
	if s.err != nil {
		status = SqrtStreamStatusFailed
	} else if s.started {
		status = SqrtStreamStatusBoundedCollapse
	}

	return SqrtApproxStreamSnapshot{
		Status:        status,
		Started:       s.started,
		PrefixTerms:   s.prefixTerms,
		Approx:        approxCopy,
		CFInputApprox: cfInputApproxCopy,
	}
}

func (s *SqrtCFPrefixStream2) initExactCF() bool {
	if s.started {
		return s.err == nil
	}
	s.started = true

	a, err := CFApproxFromPrefix(s.src, s.prefixTerms)
	if err != nil {
		s.err = fmt.Errorf("SqrtCFPrefixStream2: %w", err)
		s.done = true
		return false
	}
	s.inputApprox = &a

	approx, err := sqrtApproxFromApproxRangeSeedCanonical(a, s.policy)
	if err != nil {
		s.err = fmt.Errorf("SqrtCFPrefixStream2: %w", err)
		s.done = true
		return false
	}

	s.approx = &approx
	s.exactCF = NewRationalCF(approx)
	return true
}

func (s *SqrtCFPrefixStream2) Next() (int64, bool) {
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

// sqrt_stream_cf2.go v5
