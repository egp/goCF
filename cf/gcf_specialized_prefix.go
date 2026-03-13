// gcf_specialized_prefix.go v3
package cf

import "fmt"

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

	if err := ingestPrefixTermsIntoBounder(b, src, prefixTerms); err != nil {
		return GCFApprox{}, err
	}

	if err := applyTailEvidenceAfterPrefix(
		b,
		src,
		prefixTerms,
		tailRangeAfterPrefix,
		tailLowerBoundAfterPrefix,
	); err != nil {
		return GCFApprox{}, err
	}

	return gcfApproxFromBounder(b, prefixTerms, "specializedGCFApproxFromPrefix: empty source")
}

func specializedGCFApproxFromPrefixUsingSourceEvidence(
	prefixTerms int,
	newSrc func() GCFSource,
) (GCFApprox, error) {
	if err := requirePositivePrefixTerms("specializedGCFApproxFromPrefixUsingSourceEvidence:", prefixTerms); err != nil {
		return GCFApprox{}, err
	}

	src := newSrc()
	b := NewGCFBounder()

	if err := ingestPrefixTermsIntoBounder(b, src, prefixTerms); err != nil {
		return GCFApprox{}, err
	}

	if err := applySourceTailEvidenceAfterPrefix(b, src, prefixTerms); err != nil {
		return GCFApprox{}, err
	}

	return gcfApproxFromBounder(
		b,
		prefixTerms,
		"specializedGCFApproxFromPrefixUsingSourceEvidence: empty source",
	)
}

func applySourceTailEvidenceAfterPrefix(
	b *GCFBounder,
	src GCFSource,
	prefixTerms int,
) error {
	evSrc, ok := src.(TailEvidenceGCFSource)
	if !ok {
		return fmt.Errorf(
			"specializedGCFApproxFromPrefixUsingSourceEvidence: source %T does not implement TailEvidence",
			src,
		)
	}

	ev := evSrc.TailEvidence()

	if ev.RangeReusable && ev.Range == nil {
		return fmt.Errorf(
			"specializedGCFApproxFromPrefixUsingSourceEvidence: source %T provides reusable tail-range policy without a tail range",
			src,
		)
	}
	if ev.LowerBoundMinPrefix < 0 {
		return fmt.Errorf(
			"specializedGCFApproxFromPrefixUsingSourceEvidence: source %T provides negative LowerBoundMinPrefix=%d",
			src,
			ev.LowerBoundMinPrefix,
		)
	}

	if ev.Range != nil {
		return b.SetTailRange(*ev.Range)
	}

	if ev.LowerBound != nil && prefixTerms >= ev.LowerBoundMinPrefix {
		return b.SetTailLowerBound(*ev.LowerBound)
	}

	return nil
}

func applyTailEvidenceAfterPrefix(
	b *GCFBounder,
	src GCFSource,
	prefixTerms int,
	tailRangeAfterPrefix func(prefixTerms int) (Range, bool, error),
	tailLowerBoundAfterPrefix func(prefixTerms int) Rational,
) error {
	if evSrc, ok := src.(TailEvidenceGCFSource); ok {
		ev := evSrc.TailEvidence()

		if ev.RangeReusable && ev.Range == nil {
			return fmt.Errorf(
				"specializedGCFApproxFromPrefix: source %T provides reusable tail-range policy without a tail range",
				src,
			)
		}
		if ev.LowerBoundMinPrefix < 0 {
			return fmt.Errorf(
				"specializedGCFApproxFromPrefix: source %T provides negative LowerBoundMinPrefix=%d",
				src,
				ev.LowerBoundMinPrefix,
			)
		}

		if ev.Range != nil {
			if err := b.SetTailRange(*ev.Range); err != nil {
				return err
			}
			return nil
		}

		if ev.LowerBound != nil && prefixTerms >= ev.LowerBoundMinPrefix {
			if err := b.SetTailLowerBound(*ev.LowerBound); err != nil {
				return err
			}
			return nil
		}

		return nil
	}

	if r, ok, err := tailRangeAfterPrefix(prefixTerms); err != nil {
		return err
	} else if ok {
		return b.SetTailRange(r)
	}

	return b.SetTailLowerBound(tailLowerBoundAfterPrefix(prefixTerms))
}

// gcf_specialized_prefix.go v3
