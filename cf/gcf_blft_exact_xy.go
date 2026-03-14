// gcf_blft_exact_xy.go v2
package cf

import "fmt"

// ApplyComposedGCFXYBLFTToTailsExact boundedly ingests xSrc into the x-side
// and ySrc into the y-side of base, requires both sources to exhaust within
// their bounds, then applies the resulting BLFT exactly to (xTail, yTail).
func ApplyComposedGCFXYBLFTToTailsExact(
	base BLFT,
	xSrc GCFSource,
	xTail Rational,
	maxXIngestTerms int,
	ySrc GCFSource,
	yTail Rational,
	maxYIngestTerms int,
) (Rational, int, int, error) {
	cur := base
	xIngested := 0
	yIngested := 0

	for {
		if maxXIngestTerms >= 0 && xIngested >= maxXIngestTerms {
			return Rational{}, xIngested, yIngested,
				fmt.Errorf("ApplyComposedGCFXYBLFTToTailsExact: exceeded MaxIngestTerms=%d before x-source exhaustion", maxXIngestTerms)
		}

		p, q, ok := xSrc.NextPQ()
		if !ok {
			break
		}

		var err error
		cur, err = cur.IngestGCFX(p, q)
		if err != nil {
			return Rational{}, xIngested, yIngested, err
		}
		xIngested++
	}

	for {
		if maxYIngestTerms >= 0 && yIngested >= maxYIngestTerms {
			return Rational{}, xIngested, yIngested,
				fmt.Errorf("ApplyComposedGCFXYBLFTToTailsExact: exceeded MaxIngestTerms=%d before y-source exhaustion", maxYIngestTerms)
		}

		p, q, ok := ySrc.NextPQ()
		if !ok {
			break
		}

		var err error
		cur, err = cur.IngestGCFY(p, q)
		if err != nil {
			return Rational{}, xIngested, yIngested, err
		}
		yIngested++
	}

	out, err := cur.ApplyRat(xTail, yTail)
	if err != nil {
		return Rational{}, xIngested, yIngested, err
	}
	return out, xIngested, yIngested, nil
}

// gcf_blft_exact_xy.go v2
