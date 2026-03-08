// sqrt_api_seed_policy.go v3
package cf

import "fmt"

// SqrtApproxWithSeedAndPolicy computes a bounded rational approximation to
// sqrt(x) using the supplied explicit seed and policy.
//
// v3: implemented as a thin wrapper over SqrtApproxWithPolicy by injecting
// the explicit seed into the policy.
func SqrtApproxWithSeedAndPolicy(x, seed Rational, p SqrtPolicy) (Rational, error) {
	pp := p
	pp.Seed = &seed
	return SqrtApproxWithPolicy(x, pp)
}

// SqrtApproxCFWithSeedAndPolicy returns a ContinuedFraction source for the
// bounded sqrt approximation produced by SqrtApproxWithSeedAndPolicy.
//
// v3: thin wrapper over SqrtApproxCFWithPolicy.
func SqrtApproxCFWithSeedAndPolicy(x, seed Rational, p SqrtPolicy) (ContinuedFraction, error) {
	pp := p
	pp.Seed = &seed
	return SqrtApproxCFWithPolicy(x, pp)
}

// SqrtApproxTermsWithSeedAndPolicy returns up to digits CF terms for the
// bounded sqrt approximation produced by SqrtApproxWithSeedAndPolicy.
//
// v3: thin wrapper over SqrtApproxTermsWithPolicy.
func SqrtApproxTermsWithSeedAndPolicy(x, seed Rational, p SqrtPolicy, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsWithSeedAndPolicy: negative digits %d", digits)
	}
	pp := p
	pp.Seed = &seed
	return SqrtApproxTermsWithPolicy(x, pp, digits)
}

// sqrt_api_seed_policy.go v3
