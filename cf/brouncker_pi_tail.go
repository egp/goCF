// brouncker_pi_tail.go v3
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

func (s *Brouncker4OverPiGCFSource) LowerBoundRayMinPrefix() int {
	// Brouncker's generic lower-bound-ray fallback is too weak to trust for
	// later infinite-stream digits; rely on explicit prefix evidence instead.
	return 1 << 30
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
// Current v3 support:
//   - prefixTerms == 0: tail is the whole Brouncker object, conservatively in [1, 3/2]
//   - prefixTerms == 1: remaining tail starts at 2 + 1/(2 + 9/(2 + ...)),
//     conservatively in [2, 5/2]
//   - prefixTerms == 2: remaining tail starts at 2 + 25/(2 + 49/(2 + ...)),
//     conservatively in [2, 29/2]
//   - prefixTerms >= 3: no tighter interval currently provided
func Brouncker4OverPiTailRangeAfterPrefix(prefixTerms int) (Range, bool, error) {
	if prefixTerms < 0 {
		return Range{}, false, fmt.Errorf("Brouncker4OverPiTailRangeAfterPrefix: negative prefixTerms %d", prefixTerms)
	}

	switch prefixTerms {
	case 0:
		return NewRange(mustRat(1, 1), mustRat(3, 2), true, true), true, nil
	case 1:
		return NewRange(mustRat(2, 1), mustRat(5, 2), true, true), true, nil
	case 2:
		return NewRange(mustRat(2, 1), mustRat(29, 2), true, true), true, nil
	case 3:
		// Remaining tail is:
		//   2 + 49/u
		// where u >= 2, hence tail in [2, 2 + 49/2] = [2, 53/2].
		return NewRange(mustRat(2, 1), mustRat(53, 2), true, true), true, nil

	case 4:
		return NewRange(mustRat(2, 1), mustRat(83, 2), true, true), true, nil
	default:
		return Range{}, false, nil
	}
}

// Brouncker4OverPiApproxFromPrefix ingests up to prefixTerms terms from Brouncker's
// 4/pi GCF source and returns a GCFApprox using source-provided tail evidence.
func Brouncker4OverPiApproxFromPrefix(prefixTerms int) (GCFApprox, error) {
	return specializedGCFApproxFromPrefixUsingSourceEvidence(
		prefixTerms,
		func() GCFSource { return NewBrouncker4OverPiGCFSource() },
	)
}

// brouncker_pi_tail.go v3
