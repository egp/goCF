// gcf_blft_exact_y.go v2
package cf

import "fmt"

// ApplyComposedGCFYBLFTToTailsExact boundedly ingests ySrc into the y-side of
// base, requires source exhaustion within the bound, then applies the resulting
// BLFT exactly to (x, yTail).
func ApplyComposedGCFYBLFTToTailsExact(
	base BLFT,
	x Rational,
	ySrc GCFSource,
	yTail Rational,
	maxIngestTerms int,
) (Rational, int, error) {
	cur := base
	ingested := 0

	for {
		if maxIngestTerms >= 0 && ingested >= maxIngestTerms {
			return Rational{}, ingested,
				fmt.Errorf("ApplyComposedGCFYBLFTToTailsExact: exceeded MaxIngestTerms=%d before source exhaustion", maxIngestTerms)
		}

		p, q, ok := ySrc.NextPQ()
		if !ok {
			break
		}

		var err error
		cur, err = cur.IngestGCFY(p, q)
		if err != nil {
			return Rational{}, ingested, err
		}
		ingested++
	}

	out, err := cur.ApplyRat(x, yTail)
	if err != nil {
		return Rational{}, ingested, err
	}
	return out, ingested, nil
}

// gcf_blft_exact_y.go v2
