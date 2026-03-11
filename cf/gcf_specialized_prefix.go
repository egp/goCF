// gcf_specialized_prefix.go v1
package cf

import "fmt"

func specializedGCFApproxFromPrefix(
	prefixTerms int,
	newSrc func() GCFSource,
	tailRangeAfterPrefix func(prefixTerms int) (Range, bool, error),
	tailLowerBoundAfterPrefix func(prefixTerms int) Rational,
) (GCFApprox, error) {
	if prefixTerms < 0 {
		return GCFApprox{}, fmt.Errorf("specializedGCFApproxFromPrefix: negative prefixTerms %d", prefixTerms)
	}
	if prefixTerms == 0 {
		return GCFApprox{}, fmt.Errorf("specializedGCFApproxFromPrefix: prefixTerms must be > 0")
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

	for i := 0; i < prefixTerms; i++ {
		p, q, ok := src.NextPQ()
		if !ok {
			b.Finish()
			break
		}
		if err := b.IngestPQ(p, q); err != nil {
			return GCFApprox{}, err
		}
	}

	return gcfApproxFromBounder(b, prefixTerms, "specializedGCFApproxFromPrefix: empty source")
}

// gcf_specialized_prefix.go v1
