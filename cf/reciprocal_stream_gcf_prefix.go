// reciprocal_stream_gcf_prefix.go v1
package cf

import "fmt"

// ReciprocalGCFPrefixStream2 is an inspectable unary reciprocal stream over a
// bounded generalized continued-fraction prefix.
//
// Current milestone:
//   - ingest up to prefixTerms from a GCFSource
//   - use the exact rational convergent carried by GCFApprox
//   - reciprocal computed by exact rational collapse
//
// Future work:
//   - use unfinished-tail range information more aggressively
//   - progressive certification
type ReciprocalGCFPrefixStream2 struct {
	err     error
	done    bool
	started bool
	exactCF ContinuedFraction
	approx  *Rational

	inputApprox *GCFApprox

	src         GCFSource
	prefixTerms int
}

func NewReciprocalGCFPrefixStream2(src GCFSource, prefixTerms int) *ReciprocalGCFPrefixStream2 {
	return &ReciprocalGCFPrefixStream2{
		src:         src,
		prefixTerms: prefixTerms,
	}
}

func (s *ReciprocalGCFPrefixStream2) Err() error { return s.err }

func (s *ReciprocalGCFPrefixStream2) Snapshot() ReciprocalApproxStreamSnapshot {
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
	return ReciprocalApproxStreamSnapshot{
		Started:        s.started,
		PrefixTerms:    s.prefixTerms,
		Approx:         approxCopy,
		GCFInputApprox: gcfInputApproxCopy,
	}
}

func (s *ReciprocalGCFPrefixStream2) initExactCF() bool {
	if s.started {
		return s.err == nil
	}
	s.started = true

	a, err := GCFApproxFromPrefix(s.src, s.prefixTerms)
	if err != nil {
		s.err = fmt.Errorf("ReciprocalGCFPrefixStream2: %w", err)
		s.done = true
		return false
	}
	s.inputApprox = &a

	if a.Convergent.Cmp(intRat(0)) == 0 {
		s.err = fmt.Errorf("ReciprocalGCFPrefixStream2: reciprocal of zero")
		s.done = true
		return false
	}

	recip, err := intRat(1).Div(a.Convergent)
	if err != nil {
		s.err = fmt.Errorf("ReciprocalGCFPrefixStream2: %w", err)
		s.done = true
		return false
	}

	s.approx = &recip
	s.exactCF = NewRationalCF(recip)
	return true
}

func (s *ReciprocalGCFPrefixStream2) Next() (int64, bool) {
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

// reciprocal_stream_gcf_prefix.go v1
