// gcf_specialized_prefix.go v1
package cf

import "fmt"

// specializedGCFApproxFromPrefix is a shared helper for named GCF families that
// want prefix-aware enclosure logic tighter than the generic source-level path.
//
// The caller provides:
//   - a fresh source constructor
//   - a prefix-aware optional tail-range function
//   - a prefix-aware tail-lower-bound function
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

	if !b.HasValue() {
		return GCFApprox{}, fmt.Errorf("specializedGCFApproxFromPrefix: empty source")
	}

	conv, err := b.Convergent()
	if err != nil {
		return GCFApprox{}, err
	}

	var rp *Range
	if r, ok, err := b.Range(); err != nil {
		return GCFApprox{}, err
	} else if ok {
		rr := r
		rp = &rr
	}

	return GCFApprox{
		Convergent:  conv,
		Range:       rp,
		PrefixTerms: prefixTerms,
	}, nil
}

// gcf_specialized_prefix.go v1
