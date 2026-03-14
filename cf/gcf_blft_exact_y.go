// gcf_blft_exact_y.go v2
package cf

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

	ingested, err := ingestGCFBounded(
		"ApplyComposedGCFYBLFTToTailsExact",
		ySrc,
		maxIngestTerms,
		func(p, q int64) error {
			var ierr error
			cur, ierr = cur.IngestGCFY(p, q)
			return ierr
		},
	)
	if err != nil {
		return Rational{}, ingested, err
	}

	out, err := cur.ApplyRat(x, yTail)
	if err != nil {
		return Rational{}, ingested, err
	}
	return out, ingested, nil
}

// gcf_blft_exact_y.go v2
