// sqrt_api_seed_policy.go v1
package cf

import "fmt"

// SqrtApproxWithSeedAndPolicy computes a bounded rational approximation to
// sqrt(x) using the supplied explicit seed and policy.
func SqrtApproxWithSeedAndPolicy(x, seed Rational, p SqrtPolicy) (Rational, error) {
	approx, _, err := SqrtApproxRationalUntilResidual(x, seed, p.MaxSteps, p.Tol)
	return approx, err
}

// SqrtApproxCFWithSeedAndPolicy returns a ContinuedFraction source for the
// bounded sqrt approximation produced by SqrtApproxWithSeedAndPolicy.
func SqrtApproxCFWithSeedAndPolicy(x, seed Rational, p SqrtPolicy) (ContinuedFraction, error) {
	approx, err := SqrtApproxWithSeedAndPolicy(x, seed, p)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

// SqrtApproxTermsWithSeedAndPolicy returns up to digits CF terms for the
// bounded sqrt approximation produced by SqrtApproxWithSeedAndPolicy.
func SqrtApproxTermsWithSeedAndPolicy(x, seed Rational, p SqrtPolicy, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsWithSeedAndPolicy: negative digits %d", digits)
	}
	cf, err := SqrtApproxCFWithSeedAndPolicy(x, seed, p)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// sqrt_api_seed_policy.go v1
