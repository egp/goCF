// sqrt_stream_snapshot.go v3
package cf

type SqrtApproxStreamSnapshot struct {
	Status      SqrtStreamStatus
	Started     bool
	PrefixTerms int
	Approx      *Rational

	// Exactly one of these may be non-nil depending on stream kind.
	CFInputApprox  *CFApprox
	GCFInputApprox *GCFApprox
}

// sqrt_stream_snapshot.go v3
