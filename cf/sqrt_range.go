// sqrt_range.go v2
package cf

import "fmt"

// SqrtRangeExact returns the exact monotone image of a nonnegative inside range
// under sqrt(x), when both endpoints have exact rational square roots.
//
// Returns:
//   - (out, true, nil)  if exact endpoint square roots exist
//   - (_,   false, nil) if not yet supported exactly
//   - (_,   false, err) on invalid input
func SqrtRangeExact(r Range) (Range, bool, error) {
	if !r.IsInside() {
		return Range{}, false, fmt.Errorf("SqrtRangeExact: requires inside range; got %v", r)
	}
	if r.Lo.Cmp(intRat(0)) < 0 {
		return Range{}, false, fmt.Errorf("SqrtRangeExact: negative range %v", r)
	}

	lo, okLo, err := RationalSqrtExact(r.Lo)
	if err != nil {
		return Range{}, false, err
	}
	if !okLo {
		return Range{}, false, nil
	}

	hi, okHi, err := RationalSqrtExact(r.Hi)
	if err != nil {
		return Range{}, false, err
	}
	if !okHi {
		return Range{}, false, nil
	}

	// sqrt is monotone increasing on [0,+inf)
	return NewRange(lo, hi, r.IncLo, r.IncHi), true, nil
}

// SqrtRangeExactFromCFApprox applies SqrtRangeExact to the enclosure carried by CFApprox.
func SqrtRangeExactFromCFApprox(a CFApprox) (Range, bool, error) {
	return SqrtRangeExact(a.Range)
}

// SqrtRangeHeuristic returns a bounded rational approximation to the monotone
// image of a nonnegative inside range under sqrt(x).
//
// IMPORTANT:
//   - this is a heuristic helper
//   - it is not yet a proof-safe conservative enclosure for irrational endpoints
//   - it should not be used anywhere that requires a formally guaranteed range
//
// Current policy:
//   - first try SqrtRangeExact
//   - otherwise approximate sqrt(lo) and sqrt(hi) independently using SqrtApprox
//   - preserve endpoint inclusions
func SqrtRangeHeuristic(r Range) (Range, error) {
	if !r.IsInside() {
		return Range{}, fmt.Errorf("SqrtRangeHeuristic: requires inside range; got %v", r)
	}
	if r.Lo.Cmp(intRat(0)) < 0 {
		return Range{}, fmt.Errorf("SqrtRangeHeuristic: negative range %v", r)
	}

	if exact, ok, err := SqrtRangeExact(r); err != nil {
		return Range{}, err
	} else if ok {
		return exact, nil
	}

	lo, err := SqrtApprox(r.Lo)
	if err != nil {
		return Range{}, err
	}
	hi, err := SqrtApprox(r.Hi)
	if err != nil {
		return Range{}, err
	}

	if lo.Cmp(hi) > 0 {
		lo, hi = hi, lo
	}
	return NewRange(lo, hi, r.IncLo, r.IncHi), nil
}

// SqrtRangeHeuristicFromCFApprox applies SqrtRangeHeuristic to the enclosure
// carried by CFApprox.
func SqrtRangeHeuristicFromCFApprox(a CFApprox) (Range, error) {
	return SqrtRangeHeuristic(a.Range)
}

// sqrt_range.go v2
