// sqrt_range.go v1
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

// sqrt_range.go v1
