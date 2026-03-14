// gcf_blft_exact_xy.go v2
package cf

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

	xIngested, err := ingestGCFBounded(
		"ApplyComposedGCFXYBLFTToTailsExact",
		xSrc,
		maxXIngestTerms,
		func(p, q int64) error {
			var ierr error
			cur, ierr = cur.IngestGCFX(p, q)
			return ierr
		},
	)
	if err != nil {
		return Rational{}, xIngested, 0, err
	}

	yIngested, err := ingestGCFBounded(
		"ApplyComposedGCFXYBLFTToTailsExact",
		ySrc,
		maxYIngestTerms,
		func(p, q int64) error {
			var ierr error
			cur, ierr = cur.IngestGCFY(p, q)
			return ierr
		},
	)
	if err != nil {
		return Rational{}, xIngested, yIngested, err
	}

	out, err := cur.ApplyRat(xTail, yTail)
	if err != nil {
		return Rational{}, xIngested, yIngested, err
	}
	return out, xIngested, yIngested, nil
}

// gcf_blft_exact_xy.go v2
