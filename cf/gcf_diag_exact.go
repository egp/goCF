// gcf_diag_exact.go v2
package cf

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

	ingested, err := ingestGCFBounded(
		"ApplyComposedGCFDiagBLFTToTailExact",
		src,
		maxIngestTerms,
		func(p, q int64) error {
			var ierr error
			cur, ierr = cur.IngestGCF(p, q)
			return ierr
		},
	)
	if err != nil {
		return Rational{}, ingested, err
	}

	out, err := cur.ApplyRat(xTail)
	if err != nil {
		return Rational{}, ingested, err
	}
	return out, ingested, nil
}

// gcf_diag_exact.go v2
