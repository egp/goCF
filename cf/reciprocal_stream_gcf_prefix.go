// reciprocal_stream_gcf_prefix.go v2
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
	core reciprocalExactCollapseCore

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

func (s *ReciprocalGCFPrefixStream2) Err() error { return s.core.Err() }

func (s *ReciprocalGCFPrefixStream2) Snapshot() ReciprocalApproxStreamSnapshot {
	var approxCopy *Rational
	if s.core.approx != nil {
		v := *s.core.approx
		approxCopy = &v
	}
	var gcfInputApproxCopy *GCFApprox
	if s.inputApprox != nil {
		v := *s.inputApprox
		gcfInputApproxCopy = &v
	}
	return ReciprocalApproxStreamSnapshot{
		Started:        s.core.started,
		Approx:         approxCopy,
		GCFInputApprox: gcfInputApproxCopy,
		PrefixTerms:    s.prefixTerms,
		ConsumedTerms: func() int {
			if s.inputApprox != nil {
				return s.inputApprox.PrefixTerms
			}
			return 0
		}(),
	}
}

func (s *ReciprocalGCFPrefixStream2) evalReciprocal() (Rational, error) {
	a, err := GCFApproxFromPrefix(s.src, s.prefixTerms)
	if err != nil {
		return Rational{}, fmt.Errorf("ReciprocalGCFPrefixStream2: %w", err)
	}
	s.inputApprox = &a

	if a.Convergent.Cmp(intRat(0)) == 0 {
		return Rational{}, fmt.Errorf("ReciprocalGCFPrefixStream2: reciprocal of zero")
	}

	recip, err := intRat(1).Div(a.Convergent)
	if err != nil {
		return Rational{}, fmt.Errorf("ReciprocalGCFPrefixStream2: %w", err)
	}
	return recip, nil
}

func (s *ReciprocalGCFPrefixStream2) Next() (int64, bool) {
	return s.core.Next(s.evalReciprocal)
}

// reciprocal_stream_gcf_prefix.go v2
