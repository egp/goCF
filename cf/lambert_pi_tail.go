// lambert_pi_tail.go v1
package cf

import "fmt"

// LambertPiOver4TailLowerBoundAfterPrefix returns a conservative positive lower
// bound for the unfinished Lambert pi/4 tail after consuming prefixTerms terms.
//
// Current conservative rule:
//   - for all prefixTerms >= 0, remaining Lambert tails are positive
//   - we use the stable lower bound 1
func LambertPiOver4TailLowerBoundAfterPrefix(prefixTerms int) Rational {
	return mustRat(1, 1)
}

// LambertPiOver4TailRangeAfterPrefix returns a tighter conservative inside range
// for the unfinished Lambert pi/4 tail after consuming prefixTerms terms, when
// such a tighter interval is currently implemented.
//
// Return values:
//   - (r, true, nil)  => tighter prefix-aware tail interval available
//   - (_, false, nil) => no tighter interval available yet; callers should fall
//     back to lower-bound-only ray semantics
//   - (_, false, err) => invalid input
//
// Current v1 support:
//   - prefixTerms == 0: tail is the whole Lambert object, conservatively in [3/4, 1]
//   - prefixTerms == 1: remaining tail starts at 1 + 1/(3 + ...), conservatively in [1, 4/3]
//   - prefixTerms >= 2: no tighter interval currently provided
func LambertPiOver4TailRangeAfterPrefix(prefixTerms int) (Range, bool, error) {
	if prefixTerms < 0 {
		return Range{}, false, fmt.Errorf("LambertPiOver4TailRangeAfterPrefix: negative prefixTerms %d", prefixTerms)
	}

	switch prefixTerms {
	case 0:
		// Lambert pi/4 lies in (0,1); use conservative closed interval [3/4, 1].
		return NewRange(mustRat(3, 4), mustRat(1, 1), true, true), true, nil
	case 1:
		// Remaining tail after consuming (0,1) is:
		// 1 + 1/(3 + 4/(5 + ...))
		// Conservatively bounded in [1, 4/3].
		return NewRange(mustRat(1, 1), mustRat(4, 3), true, true), true, nil
	default:
		return Range{}, false, nil
	}
}

// LambertPiOver4ApproxFromPrefix ingests up to prefixTerms terms from Lambert's
// pi/4 GCF source and returns a GCFApprox.
//
// It prefers a prefix-aware tighter tail interval when currently available;
// otherwise it falls back to the generic lower-bound-only enclosure path.
func LambertPiOver4ApproxFromPrefix(prefixTerms int) (GCFApprox, error) {
	if prefixTerms < 0 {
		return GCFApprox{}, fmt.Errorf("LambertPiOver4ApproxFromPrefix: negative prefixTerms %d", prefixTerms)
	}
	if prefixTerms == 0 {
		return GCFApprox{}, fmt.Errorf("LambertPiOver4ApproxFromPrefix: prefixTerms must be > 0")
	}

	src := NewLambertPiOver4GCFSource()
	b := NewGCFBounder()

	// Prefer tighter prefix-aware tail interval when we currently know one.
	if r, ok, err := LambertPiOver4TailRangeAfterPrefix(prefixTerms); err != nil {
		return GCFApprox{}, err
	} else if ok {
		if err := b.SetTailRange(r); err != nil {
			return GCFApprox{}, err
		}
	} else {
		if err := b.SetTailLowerBound(LambertPiOver4TailLowerBoundAfterPrefix(prefixTerms)); err != nil {
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
		return GCFApprox{}, fmt.Errorf("LambertPiOver4ApproxFromPrefix: empty source")
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

// lambert_pi_tail.go v1
