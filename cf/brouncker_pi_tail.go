// brouncker_pi_tail.go v6
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
// Current v6 support:
//   - prefixTerms == 0: full Brouncker value, conservatively in [15/13, 105/76]
//   - prefixTerms == 1: remaining tail starts at 2 + 9/(2 + 25/(2 + ...)),
//     conservatively in [76/29, 13/2]
//   - prefixTerms == 2: remaining tail starts at 2 + 25/(2 + 49/(2 + ...)),
//     conservatively in [156/53, 29/2]
//   - prefixTerms == 3: remaining tail starts at 2 + 49/(2 + 81/(2 + ...)),
//     conservatively in [215/83, 53/2]
//   - prefixTerms == 4: remaining tail starts at 2 + 81/(2 + 121/(2 + ...)),
//     conservatively in [2, 83/2]
//   - prefixTerms >= 5: no tighter interval currently provided
//
// These are base Brouncker tail evidences. They are safe, but for infinite
// streaming they may still be too weak to certify later digits quickly.
// Stronger same-state candidate evidences can be derived by deeper bounded
// lookahead into the Brouncker recurrence; see
// Brouncker4OverPiTailLookaheadRangeAfterPrefix.
func Brouncker4OverPiTailRangeAfterPrefix(prefixTerms int) (Range, bool, error) {
	if prefixTerms < 0 {
		return Range{}, false, fmt.Errorf("Brouncker4OverPiTailRangeAfterPrefix: negative prefixTerms %d", prefixTerms)
	}

	switch prefixTerms {
	case 0:
		// Full value is:
		//   1 + 1/t1
		// where t1 is the prefix-1 tail in [76/29, 13/2].
		//
		// Since x -> 1 + 1/x is decreasing on x > 0:
		//   lo = 1 + 1/(13/2) = 15/13
		//   hi = 1 + 1/(76/29) = 105/76
		return NewRange(mustRat(15, 13), mustRat(105, 76), true, true), true, nil

	case 1:
		// Remaining tail is:
		//   2 + 9/u
		// where u is the prefix-2 tail.
		//
		// Using prefix-2 upper bound u <= 29/2 and positivity u >= 2:
		//   lo = 2 + 9/(29/2) = 76/29
		//   hi = 2 + 9/2       = 13/2
		return NewRange(mustRat(76, 29), mustRat(13, 2), true, true), true, nil

	case 2:
		// Remaining tail is:
		//   2 + 25/u
		// where u is the prefix-3 tail.
		//
		// Using prefix-3 upper bound u <= 53/2 and positivity u >= 2:
		//   lo = 2 + 25/(53/2) = 156/53
		//   hi = 2 + 25/2      = 29/2
		return NewRange(mustRat(156, 53), mustRat(29, 2), true, true), true, nil

	case 3:
		// Remaining tail is:
		//   2 + 49/u
		// where u is the prefix-4 tail.
		//
		// Using prefix-4 upper bound u <= 83/2 and positivity u >= 2:
		//   lo = 2 + 49/(83/2) = 215/83
		//   hi = 2 + 49/2      = 53/2
		return NewRange(mustRat(215, 83), mustRat(53, 2), true, true), true, nil

	case 4:
		// Remaining tail is:
		//   2 + 81/u
		// where u >= 2, hence tail in [2, 2 + 81/2] = [2, 83/2].
		return NewRange(mustRat(2, 1), mustRat(83, 2), true, true), true, nil

	default:
		return Range{}, false, nil
	}
}

type brounckerLookaheadBound struct {
	value  Rational
	finite bool
}

func finiteBrounckerLookaheadBound(num, den int64) brounckerLookaheadBound {
	return brounckerLookaheadBound{
		value:  mustRat(num, den),
		finite: true,
	}
}

func addTwoPlusSquareOverBound(square int64, b brounckerLookaheadBound) (brounckerLookaheadBound, error) {
	// 2 + square / b
	sq := mustRat(square, 1)

	term, err := sq.Div(b.value)
	if err != nil {
		return brounckerLookaheadBound{}, err
	}

	sum, err := mustRat(2, 1).Add(term)
	if err != nil {
		return brounckerLookaheadBound{}, err
	}

	return brounckerLookaheadBound{
		value:  sum,
		finite: true,
	}, nil
}

// Brouncker4OverPiTailLookaheadRangeAfterPrefix returns a stronger same-state
// conservative tail interval by looking ahead a fixed finite number of Brouncker
// denominator terms and then closing the still-unfinished tail with only the
// positive lower bound >= 2.
//
// Intended use:
//   - base TailEvidence() remains simple and cheap
//   - CandidateTailEvidence() can use this helper to offer stronger same-state
//     alternate enclosures before further ingestion
//
// Domain:
//   - prefixTerms >= 1
//   - lookaheadTerms >= 1
//
// Construction:
//
//	If the current tail starts at odd^2, then each backward step applies
//
//	    x = 2 + odd^2 / tail
//
//	over a current positive interval [L,U]. Since this map is decreasing on
//	positive tails, the image interval is:
//
//	    [2 + odd^2/U, 2 + odd^2/L]
//
//	with the initial unfinished tail closed conservatively by [2, +∞).
func Brouncker4OverPiTailLookaheadRangeAfterPrefix(prefixTerms int, lookaheadTerms int) (Range, bool, error) {
	if prefixTerms < 1 {
		return Range{}, false, fmt.Errorf(
			"Brouncker4OverPiTailLookaheadRangeAfterPrefix: prefixTerms must be >= 1, got %d",
			prefixTerms,
		)
	}
	if lookaheadTerms < 1 {
		return Range{}, false, fmt.Errorf(
			"Brouncker4OverPiTailLookaheadRangeAfterPrefix: lookaheadTerms must be >= 1, got %d",
			lookaheadTerms,
		)
	}

	lo := finiteBrounckerLookaheadBound(2, 1)
	hi := brounckerLookaheadBound{finite: false} // +∞

	for j := lookaheadTerms - 1; j >= 0; j-- {
		odd := int64(2*(prefixTerms+j) + 1)
		square := odd * odd

		nextLo := finiteBrounckerLookaheadBound(2, 1)
		if hi.finite {
			var err error
			nextLo, err = addTwoPlusSquareOverBound(square, hi)
			if err != nil {
				return Range{}, false, err
			}
		}

		nextHi, err := addTwoPlusSquareOverBound(square, lo)
		if err != nil {
			return Range{}, false, err
		}

		lo = nextLo
		hi = nextHi
	}

	return NewRange(lo.value, hi.value, true, true), true, nil
}

// Brouncker4OverPiApproxFromPrefix ingests up to prefixTerms terms from Brouncker's
// 4/pi GCF source and returns a GCFApprox using source-provided tail evidence.
func Brouncker4OverPiApproxFromPrefix(prefixTerms int) (GCFApprox, error) {
	return specializedGCFApproxFromPrefixUsingSourceEvidence(
		prefixTerms,
		func() GCFSource { return NewBrouncker4OverPiGCFSource() },
	)
}

// brouncker_pi_tail.go v6
