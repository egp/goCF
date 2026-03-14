// gcf_blft_exact.go v2
package cf

import "fmt"

// ApplyComposedGCFXBLFTToTailsExact boundedly ingests xSrc into the x-side of
// base, requires source exhaustion within the bound, then applies the resulting
// BLFT exactly to (xTail, y).
func ApplyComposedGCFXBLFTToTailsExact(
	base BLFT,
	xSrc GCFSource,
	xTail Rational,
	y Rational,
	maxIngestTerms int,
) (Rational, int, error) {
	cur := base
	ingested := 0

	for {
		if maxIngestTerms >= 0 && ingested >= maxIngestTerms {
			return Rational{}, ingested,
				fmt.Errorf("ApplyComposedGCFXBLFTToTailsExact: exceeded MaxIngestTerms=%d before source exhaustion", maxIngestTerms)
		}

		p, q, ok := xSrc.NextPQ()
		if !ok {
			break
		}

		var err error
		cur, err = cur.IngestGCFX(p, q)
		if err != nil {
			return Rational{}, ingested, err
		}
		ingested++
	}

	out, err := cur.ApplyRat(xTail, y)
	if err != nil {
		return Rational{}, ingested, err
	}
	return out, ingested, nil
}

// gcf_blft_exact.go v2
