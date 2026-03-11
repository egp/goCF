// gcf_specialized_prefix.go v1
package cf

func specializedGCFApproxFromPrefix(
	prefixTerms int,
	newSrc func() GCFSource,
	tailRangeAfterPrefix func(prefixTerms int) (Range, bool, error),
	tailLowerBoundAfterPrefix func(prefixTerms int) Rational,
) (GCFApprox, error) {
	if err := requirePositivePrefixTerms("specializedGCFApproxFromPrefix:", prefixTerms); err != nil {
		return GCFApprox{}, err
	}

	src := newSrc()
	b := NewGCFBounder()

	if r, ok, err := tailRangeAfterPrefix(prefixTerms); err != nil {
		return GCFApprox{}, err
	} else if ok {
		if err := b.SetTailRange(r); err != nil {
			return GCFApprox{}, err
		}
	} else {
		if err := b.SetTailLowerBound(tailLowerBoundAfterPrefix(prefixTerms)); err != nil {
			return GCFApprox{}, err
		}
	}

	if err := ingestPrefixTermsIntoBounder(b, src, prefixTerms); err != nil {
		return GCFApprox{}, err
	}

	return gcfApproxFromBounder(b, prefixTerms, "specializedGCFApproxFromPrefix: empty source")
}

// gcf_specialized_prefix.go v1
