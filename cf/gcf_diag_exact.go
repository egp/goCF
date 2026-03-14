// gcf_diag_exact.go v2
package cf

import "fmt"

// ApplyComposedGCFDiagBLFTToTailExact boundedly ingests src into base,
// requires source exhaustion within the bound, then applies the resulting
// diagonal transform exactly to xTail.
func ApplyComposedGCFDiagBLFTToTailExact(
	base DiagBLFT,
	src GCFSource,
	xTail Rational,
	maxIngestTerms int,
) (Rational, int, error) {
	cur := base
	ingested := 0

	for {
		if maxIngestTerms >= 0 && ingested >= maxIngestTerms {
			return Rational{}, ingested,
				fmt.Errorf("ApplyComposedGCFDiagBLFTToTailExact: exceeded MaxIngestTerms=%d before source exhaustion", maxIngestTerms)
		}

		p, q, ok := src.NextPQ()
		if !ok {
			break
		}

		var err error
		cur, err = cur.IngestGCF(p, q)
		if err != nil {
			return Rational{}, ingested, err
		}
		ingested++
	}

	out, err := cur.ApplyRat(xTail)
	if err != nil {
		return Rational{}, ingested, err
	}
	return out, ingested, nil
}

// gcf_diag_exact.go v2
