// sqrt_source_prefix_api2.go v1
package cf

import "fmt"

// SqrtApproxFromSourceRangeSeed2 consumes a finite prefix of src, converts that
// prefix to a CFApprox, then returns a bounded rational sqrt approximation
// under the supplied policy.
func SqrtApproxFromSourceRangeSeed2(src ContinuedFraction, prefixTerms int, p SqrtPolicy2) (Rational, error) {
	return sqrtApproxFromSourceRangeSeedCanonical(src, prefixTerms, p)
}

// SqrtApproxCFFromSourceRangeSeed2 consumes a finite prefix of src, converts
// that prefix to a CFApprox, then returns a ContinuedFraction source for the
// bounded sqrt approximation under the supplied policy.
func SqrtApproxCFFromSourceRangeSeed2(src ContinuedFraction, prefixTerms int, p SqrtPolicy2) (ContinuedFraction, error) {
	return sqrtApproxCFFromSourceRangeSeedCanonical(src, prefixTerms, p)
}

// SqrtApproxTermsFromSourceRangeSeed2 consumes a finite prefix of src, converts
// that prefix to a CFApprox, then returns up to digits CF terms for the bounded
// sqrt approximation under the supplied policy.
func SqrtApproxTermsFromSourceRangeSeed2(src ContinuedFraction, prefixTerms int, p SqrtPolicy2, digits int) ([]int64, error) {
	return sqrtApproxTermsFromSourceRangeSeedCanonical(src, prefixTerms, p, digits)
}

// SqrtApproxFromSourceRangeSeedDefault2 is the default-policy wrapper around
// SqrtApproxFromSourceRangeSeed2.
func SqrtApproxFromSourceRangeSeedDefault2(src ContinuedFraction, prefixTerms int) (Rational, error) {
	return SqrtApproxFromSourceRangeSeed2(src, prefixTerms, DefaultSqrtPolicy2())
}

// SqrtApproxCFFromSourceRangeSeedDefault2 is the default-policy wrapper around
// SqrtApproxCFFromSourceRangeSeed2.
func SqrtApproxCFFromSourceRangeSeedDefault2(src ContinuedFraction, prefixTerms int) (ContinuedFraction, error) {
	return SqrtApproxCFFromSourceRangeSeed2(src, prefixTerms, DefaultSqrtPolicy2())
}

// SqrtApproxTermsFromSourceRangeSeedDefault2 is the default-policy wrapper around
// SqrtApproxTermsFromSourceRangeSeed2.
func SqrtApproxTermsFromSourceRangeSeedDefault2(src ContinuedFraction, prefixTerms int, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromSourceRangeSeedDefault2: negative digits %d", digits)
	}
	return SqrtApproxTermsFromSourceRangeSeed2(src, prefixTerms, DefaultSqrtPolicy2(), digits)
}

// sqrt_source_prefix_api2.go v1
