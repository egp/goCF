// sqrt_source_api2.go v1
package cf

// SqrtApproxFromApproxRangeSeed2 takes a bundled CFApprox, derives a seed from
// its range when needed, and returns a bounded rational sqrt approximation
// under the supplied policy.
//
// If p.Seed is already set, it is honored and the range-derived seed is not used.
func SqrtApproxFromApproxRangeSeed2(a CFApprox, p SqrtPolicy2) (Rational, error) {
	return sqrtApproxFromApproxRangeSeedCanonical(a, p)
}

// SqrtApproxCFFromApproxRangeSeed2 returns a ContinuedFraction source for the
// bounded sqrt approximation produced by SqrtApproxFromApproxRangeSeed2.
func SqrtApproxCFFromApproxRangeSeed2(a CFApprox, p SqrtPolicy2) (ContinuedFraction, error) {
	return sqrtApproxCFFromApproxRangeSeedCanonical(a, p)
}

// SqrtApproxTermsFromApproxRangeSeed2 returns up to digits CF terms for the
// bounded sqrt approximation produced by SqrtApproxCFFromApproxRangeSeed2.
func SqrtApproxTermsFromApproxRangeSeed2(a CFApprox, p SqrtPolicy2, digits int) ([]int64, error) {
	return sqrtApproxTermsFromApproxRangeSeedCanonical(a, p, digits)
}

// sqrt_source_api2.go v1
