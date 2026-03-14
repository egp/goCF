// gcf_compose.go v1
package cf

import "fmt"

// ComposeGCFIntoULFTBounded ingests up to MaxIngestTerms from src into base.
//
// Contract:
//   - If src exhausts before the bound is hit, returns (tFinal, ingested, true, nil).
//   - If the bound is hit before src exhausts, returns (zero, ingested, false, err).
//   - If any ingested term is invalid, returns (zero, ingestedSoFar, false, err).
//
// Bound semantics:
//   - maxIngestTerms < 0 : unlimited
//   - maxIngestTerms = 0 : no ingestion allowed
//   - maxIngestTerms > 0 : maximum permitted ingested source terms
func ComposeGCFIntoULFTBounded(
	base ULFT,
	src GCFSource,
	maxIngestTerms int,
) (ULFT, int, bool, error) {
	cur := base
	ingested := 0

	for {
		if maxIngestTerms >= 0 && ingested >= maxIngestTerms {
			return ULFT{}, ingested, false,
				fmt.Errorf("ComposeGCFIntoULFTBounded: exceeded MaxIngestTerms=%d before source exhaustion", maxIngestTerms)
		}

		p, q, ok := src.NextPQ()
		if !ok {
			return cur, ingested, true, nil
		}

		var err error
		cur, err = cur.IngestGCF(p, q)
		if err != nil {
			return ULFT{}, ingested, false, err
		}
		ingested++
	}
}

// gcf_compose.go v1
