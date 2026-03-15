// sqrt_range_conservative.go v1
package cf

import "fmt"

// SqrtLowerBoundRational returns a rational lower bound L such that
//
//	L <= sqrt(x)
//
// Current status:
//   - stub for test-first development
func SqrtLowerBoundRational(x Rational) (Rational, error) {
	return Rational{}, fmt.Errorf("SqrtLowerBoundRational: not implemented")
}

// SqrtUpperBoundRational returns a rational upper bound U such that
//
//	sqrt(x) <= U
//
// Current status:
//   - stub for test-first development
func SqrtUpperBoundRational(x Rational) (Rational, error) {
	return Rational{}, fmt.Errorf("SqrtUpperBoundRational: not implemented")
}

// SqrtRangeConservative returns a proof-safe enclosure for sqrt(x) over a
// nonnegative inside range r.
//
// Intended contract:
//   - if x ∈ r, then sqrt(x) ∈ out
//   - exact endpoint square roots should be preserved when available
//
// Current status:
//   - stub for test-first development
func SqrtRangeConservative(r Range) (Range, error) {
	return Range{}, fmt.Errorf("SqrtRangeConservative: not implemented")
}

// sqrt_range_conservative.go v1
