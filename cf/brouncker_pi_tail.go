// brouncker_pi_tail.go v1
package cf

import "fmt"

// Brouncker4OverPiTailLowerBoundAfterPrefix returns a conservative positive lower
// bound for the unfinished Brouncker 4/pi tail after consuming prefixTerms terms.
//
// Current conservative rule:
//   - for all prefixTerms >= 0, remaining Brouncker tails are positive
//   - we use the stable lower bound 1
func Brouncker4OverPiTailLowerBoundAfterPrefix(prefixTerms int) Rational {
	return mustRat(1, 1)
}

// Brouncker4OverPiTailRangeAfterPrefix returns a tighter conservative inside range
// for the unfinished Brouncker 4/pi tail after consuming prefixTerms terms, when
// such a tighter interval is currently implemented.
//
// Return values:
//   - (r, true, nil)  => tighter prefix-aware tail interval available
//   - (_, false, nil) => no tighter interval available yet; callers should fall
//     back to lower-bound-only ray semantics
//   - (_, false, err) => invalid input
//
// Current v1 support:
//   - prefixTerms == 0: tail is the whole Brouncker object, conservatively in [1, 3/2]
//   - prefixTerms == 1: remaining tail starts at 2 + 1/(2 + 9/(2 + ...)),
//     conservatively in [2, 5/2]
//   - prefixTerms >= 2: no tighter interval currently provided
func Brouncker4OverPiTailRangeAfterPrefix(prefixTerms int) (Range, bool, error) {
	if prefixTerms < 0 {
		return Range{}, false, fmt.Errorf("Brouncker4OverPiTailRangeAfterPrefix: negative prefixTerms %d", prefixTerms)
	}

	switch prefixTerms {
	case 0:
		return NewRange(mustRat(1, 1), mustRat(3, 2), true, true), true, nil
	case 1:
		return NewRange(mustRat(2, 1), mustRat(5, 2), true, true), true, nil
	default:
		return Range{}, false, nil
	}
}

// Brouncker4OverPiApproxFromPrefix ingests up to prefixTerms terms from Brouncker's
// 4/pi GCF source and returns a GCFApprox.
//
// It prefers a prefix-aware tighter tail interval when currently available;
// otherwise it falls back to the generic lower-bound-only enclosure path.
func Brouncker4OverPiApproxFromPrefix(prefixTerms int) (GCFApprox, error) {
	if prefixTerms < 0 {
		return GCFApprox{}, fmt.Errorf("Brouncker4OverPiApproxFromPrefix: negative prefixTerms %d", prefixTerms)
	}
	if prefixTerms == 0 {
		return GCFApprox{}, fmt.Errorf("Brouncker4OverPiApproxFromPrefix: prefixTerms must be > 0")
	}

	src := NewBrouncker4OverPiGCFSource()
	b := NewGCFBounder()

	if r, ok, err := Brouncker4OverPiTailRangeAfterPrefix(prefixTerms); err != nil {
		return GCFApprox{}, err
	} else if ok {
		if err := b.SetTailRange(r); err != nil {
			return GCFApprox{}, err
		}
	} else {
		if err := b.SetTailLowerBound(Brouncker4OverPiTailLowerBoundAfterPrefix(prefixTerms)); err != nil {
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
		return GCFApprox{}, fmt.Errorf("Brouncker4OverPiApproxFromPrefix: empty source")
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

// brouncker_pi_tail.go v1
