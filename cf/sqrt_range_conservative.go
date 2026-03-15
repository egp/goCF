// sqrt_range_conservative.go v2
package cf

import "fmt"

// SqrtLowerBoundRational returns a rational lower bound L such that
//
//	L <= sqrt(x)
//
// Current implementation:
//   - rejects negative input
//   - exact-square fast path
//   - zero fast path
//   - otherwise not yet implemented
func SqrtLowerBoundRational(x Rational) (Rational, error) {
	if x.Cmp(intRat(0)) < 0 {
		return Rational{}, fmt.Errorf("SqrtLowerBoundRational: negative input %v", x)
	}

	if root, ok, err := RationalSqrtExact(x); err != nil {
		return Rational{}, err
	} else if ok {
		return root, nil
	}

	if x.Cmp(intRat(0)) == 0 {
		return intRat(0), nil
	}

	return Rational{}, fmt.Errorf("SqrtLowerBoundRational: non-square input not implemented")
}

// SqrtUpperBoundRational returns a rational upper bound U such that
//
//	sqrt(x) <= U
//
// Current implementation:
//   - rejects negative input
//   - exact-square fast path
//   - zero fast path
//   - otherwise not yet implemented
func SqrtUpperBoundRational(x Rational) (Rational, error) {
	if x.Cmp(intRat(0)) < 0 {
		return Rational{}, fmt.Errorf("SqrtUpperBoundRational: negative input %v", x)
	}

	if root, ok, err := RationalSqrtExact(x); err != nil {
		return Rational{}, err
	} else if ok {
		return root, nil
	}

	if x.Cmp(intRat(0)) == 0 {
		return intRat(0), nil
	}

	return Rational{}, fmt.Errorf("SqrtUpperBoundRational: non-square input not implemented")
}

// SqrtRangeConservative returns a proof-safe enclosure for sqrt(x) over a
// nonnegative inside range r.
//
// Current implementation:
//   - rejects outside / negative ranges
//   - exact-endpoint fast path only
//   - otherwise not yet implemented
func SqrtRangeConservative(r Range) (Range, error) {
	if !r.IsInside() {
		return Range{}, fmt.Errorf("SqrtRangeConservative: requires inside range; got %v", r)
	}
	if r.Lo.Cmp(intRat(0)) < 0 {
		return Range{}, fmt.Errorf("SqrtRangeConservative: negative range %v", r)
	}

	lo, err := SqrtLowerBoundRational(r.Lo)
	if err != nil {
		return Range{}, err
	}
	hi, err := SqrtUpperBoundRational(r.Hi)
	if err != nil {
		return Range{}, err
	}

	return NewRange(lo, hi, r.IncLo, r.IncHi), nil
}

// sqrt_range_conservative.go v2
