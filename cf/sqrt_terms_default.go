// sqrt_terms_default.go v1
package cf

// SqrtApproxTermsDefault computes CF terms for a bounded rational Newton
// approximation to sqrt(x), using DefaultSqrtSeed(x).
func SqrtApproxTermsDefault(x Rational, steps, digits int) ([]int64, error) {
	seed, err := DefaultSqrtSeed(x)
	if err != nil {
		return nil, err
	}
	return SqrtApproxTerms(x, seed, steps, digits)
}

// sqrt_terms_default.go v1
