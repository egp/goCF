// lambert_pi_tail.go v4
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
// Current v4 support:
//   - prefixTerms == 0: tail is the whole Lambert object, conservatively in [3/4, 1]
//   - prefixTerms == 1: remaining tail starts at 1 + 1/(3 + 4/(5 + ...)),
//     conservatively in [1, 4/3]
//   - prefixTerms == 2: remaining tail starts at 3 + 4/(5 + 9/(7 + ...)),
//     conservatively in [3, 5]
//   - prefixTerms == 3: remaining tail starts at 5 + 9/(7 + 16/(9 + ...)),
//     conservatively in [5, 34/5]
//   - prefixTerms >= 4: no tighter interval currently provided
//
// Lambert uses source-specific prefix evidence mainly to improve early-digit
// cadence on the infinite source. Prefix-2 and prefix-3 specializations give
// visibly tighter unfinished-value ranges than generic lower-bound-only fallback.
func LambertPiOver4TailRangeAfterPrefix(prefixTerms int) (Range, bool, error) {
	if prefixTerms < 0 {
		return Range{}, false, fmt.Errorf("LambertPiOver4TailRangeAfterPrefix: negative prefixTerms %d", prefixTerms)
	}

	switch prefixTerms {
	case 0:
		return NewRange(mustRat(3, 4), mustRat(1, 1), true, true), true, nil
	case 1:
		return NewRange(mustRat(1, 1), mustRat(4, 3), true, true), true, nil
	case 2:
		// Remaining tail is:
		//   3 + 4/u
		// where u > 0, and in particular u >= 1.
		//
		// We keep the simple conservative interval [3,5].
		return NewRange(mustRat(3, 1), mustRat(5, 1), true, true), true, nil
	case 3:
		// Remaining tail is:
		//   5 + 9/u
		// with u positive. A simple conservative enclosure is [5, 34/5].
		return NewRange(mustRat(5, 1), mustRat(34, 5), true, true), true, nil
	default:
		return Range{}, false, nil
	}
}

// LambertPiOver4ApproxFromPrefix ingests up to prefixTerms terms from Lambert's
// pi/4 GCF source and returns a GCFApprox using source-provided tail evidence.
func LambertPiOver4ApproxFromPrefix(prefixTerms int) (GCFApprox, error) {
	return specializedGCFApproxFromPrefixUsingSourceEvidence(
		prefixTerms,
		func() GCFSource { return NewLambertPiOver4GCFSource() },
	)
}

// lambert_pi_tail.go v4
