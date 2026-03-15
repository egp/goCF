// reciprocal_stream_snapshot.go v2
package cf

type ReciprocalApproxStreamSnapshot struct {
	Started bool
	Approx  *Rational

	GCFInputApprox *GCFApprox

	// For bounded-prefix streams, this is the configured prefix limit and also
	// the intended number of source terms to inspect.
	PrefixTerms int

	// For exact-tail streams, this is the configured ingestion cap.
	MaxIngestTerms int

	// Actual source terms consumed to build the current approximation.
	ConsumedTerms int
}

// reciprocal_stream_snapshot.go v2
