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

// ApplyComposedGCFULFTToTailExact boundedly ingests src into base, requires
// source exhaustion within the bound, then applies the resulting ULFT exactly
// to the supplied tail rational.
func ApplyComposedGCFULFTToTailExact(
	base ULFT,
	src GCFSource,
	tail Rational,
	maxIngestTerms int,
) (Rational, int, error) {
	composed, ingested, exhausted, err := ComposeGCFIntoULFTBounded(base, src, maxIngestTerms)
	if err != nil {
		return Rational{}, ingested, err
	}
	if !exhausted {
		return Rational{}, ingested, fmt.Errorf(
			"ApplyComposedGCFULFTToTailExact: internal: bounded compose returned !exhausted without error",
		)
	}

	y, err := composed.ApplyRat(tail)
	if err != nil {
		return Rational{}, ingested, err
	}
	return y, ingested, nil
}

// gcf_compose.go v1
