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
	return specializedGCFApproxFromPrefix(
		prefixTerms,
		func() GCFSource { return NewLambertPiOver4GCFSource() },
		LambertPiOver4TailRangeAfterPrefix,
		LambertPiOver4TailLowerBoundAfterPrefix,
	)
}

// lambert_pi_tail.go v1
