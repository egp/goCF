// gcf_ingest_common.go v2
package cf

import "fmt"

// ingestGCFBounded pulls terms from src, enforces maxIngestTerms, and calls
// ingest(p,q) for each emitted term.
//
// Contract:
//   - returns (ingested, nil) if src exhausts within the bound
//   - returns (ingested, err) if the bound is hit before exhaustion
//   - returns (ingestedSoFar, err) if ingest(p,q) fails
//
// Bound semantics:
//   - maxIngestTerms < 0 : unlimited
//   - maxIngestTerms = 0 : no ingestion allowed
//   - maxIngestTerms > 0 : maximum permitted ingested source terms
func ingestGCFBounded(
	context string,
	src GCFSource,
	maxIngestTerms int,
	ingest func(p, q int64) error,
) (int, error) {
	ingested := 0

	for {
		p, q, ok := src.NextPQ()
		if !ok {
			return ingested, nil
		}

		if maxIngestTerms >= 0 && ingested >= maxIngestTerms {
			return ingested, fmt.Errorf(
				"%s: exceeded MaxIngestTerms=%d before source exhaustion",
				context,
				maxIngestTerms,
			)
		}

		if err := ingest(p, q); err != nil {
			return ingested, err
		}
		ingested++
	}
}

// gcf_ingest_common.go v2
