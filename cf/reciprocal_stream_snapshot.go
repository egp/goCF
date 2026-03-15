// reciprocal_stream_snapshot.go v1
package cf

type ReciprocalApproxStreamSnapshot struct {
	Started     bool
	PrefixTerms int
	Approx      *Rational

	GCFInputApprox *GCFApprox
}

// reciprocal_stream_snapshot.go v1
