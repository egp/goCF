// gcf_blft_exact.go v2
package cf

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

	ingested, err := ingestGCFBounded(
		"ApplyComposedGCFXBLFTToTailsExact",
		xSrc,
		maxIngestTerms,
		func(p, q int64) error {
			var ierr error
			cur, ierr = cur.IngestGCFX(p, q)
			return ierr
		},
	)
	if err != nil {
		return Rational{}, ingested, err
	}

	out, err := cur.ApplyRat(xTail, y)
	if err != nil {
		return Rational{}, ingested, err
	}
	return out, ingested, nil
}

// gcf_blft_exact.go v2
