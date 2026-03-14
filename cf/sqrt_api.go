// sqrt_api.go v3
package cf

import "fmt"

type SqrtPolicy struct {
	MaxSteps int
	Tol      Rational
	Seed     *Rational
}

func DefaultSqrtPolicy() SqrtPolicy {
	return SqrtPolicy{
		MaxSteps: 5,
		Tol:      mustRat(1, 1_000_000_000_000),
		Seed:     nil,
	}
}

func (p SqrtPolicy) Validate() error {
	if p.MaxSteps < 0 {
		return fmt.Errorf("SqrtPolicy: negative MaxSteps %d", p.MaxSteps)
	}
	if p.Tol.Cmp(intRat(0)) < 0 {
		return fmt.Errorf("SqrtPolicy: negative Tol %v", p.Tol)
	}
	if p.Seed != nil {
		if p.Seed.r.Sign() == 0 {
			return fmt.Errorf("SqrtPolicy: zero Seed")
		}
		if p.Seed.r.Sign() < 0 {
			return fmt.Errorf("SqrtPolicy: negative Seed %v", *p.Seed)
		}
	}
	return nil
}

// SqrtApprox uses a simple default policy to compute a bounded rational
// approximation to sqrt(x).
//
// Current policy:
//   - seed: DefaultSqrtSeed(x)
//   - maxSteps: 5
//   - tolerance: 1 / 10^12
//
// This is a convenience API, not yet a true streaming sqrt operator.
func SqrtApprox(x Rational) (Rational, error) {
	return SqrtApprox2(x)
}

// SqrtApproxCF returns a ContinuedFraction source for the bounded default
// sqrt approximation produced by SqrtApprox.
func SqrtApproxCF(x Rational) (ContinuedFraction, error) {
	return SqrtApproxCF2(x)
}

// SqrtApproxTermsAuto returns up to digits CF terms for the bounded default
// sqrt approximation produced by SqrtApprox.
func SqrtApproxTermsAuto(x Rational, digits int) ([]int64, error) {
	return SqrtApproxTerms2(x, digits)
}

// SqrtApproxWithPolicy computes a bounded rational approximation to sqrt(x)
// using the supplied policy.
func SqrtApproxWithPolicy(x Rational, p SqrtPolicy) (Rational, error) {
	return SqrtApproxWithPolicy2(x, sqrtPolicy2FromOld(p))
}

func SqrtApproxCF2WithPolicyBridge(x Rational, p SqrtPolicy) (ContinuedFraction, error) {
	approx, err := SqrtApproxWithPolicy2(x, sqrtPolicy2FromOld(p))
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

// SqrtApproxCFWithPolicy returns a ContinuedFraction source for the bounded
// sqrt approximation produced by SqrtApproxWithPolicy.
func SqrtApproxCFWithPolicy(x Rational, p SqrtPolicy) (ContinuedFraction, error) {
	return SqrtApproxCF2WithPolicyBridge(x, p)
}

// SqrtApproxTermsWithPolicy returns up to digits CF terms for the bounded
// sqrt approximation produced by SqrtApproxWithPolicy.
func SqrtApproxTermsWithPolicy(x Rational, p SqrtPolicy, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsWithPolicy: negative digits %d", digits)
	}
	cf, err := SqrtApproxCF2WithPolicyBridge(x, p)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// SqrtApproxWithSeedAndPolicy computes a bounded rational approximation to
// sqrt(x) using the supplied explicit seed and policy.
//
// Preferred configuration path:
//   - use SqrtApproxWithPolicy with p.Seed set
//
// This compatibility wrapper remains available for now.
func SqrtApproxWithSeedAndPolicy(x, seed Rational, p SqrtPolicy) (Rational, error) {
	pp := p
	pp.Seed = &seed
	return SqrtApproxWithPolicy(x, pp)
}

// SqrtApproxCFWithSeedAndPolicy returns a ContinuedFraction source for the
// bounded sqrt approximation produced by SqrtApproxWithSeedAndPolicy.
//
// Preferred configuration path:
//   - use SqrtApproxCFWithPolicy with p.Seed set
//
// This compatibility wrapper remains available for now.
func SqrtApproxCFWithSeedAndPolicy(x, seed Rational, p SqrtPolicy) (ContinuedFraction, error) {
	pp := p
	pp.Seed = &seed
	return SqrtApproxCFWithPolicy(x, pp)
}

// SqrtApproxTermsWithSeedAndPolicy returns up to digits CF terms for the
// bounded sqrt approximation produced by SqrtApproxWithSeedAndPolicy.
//
// Preferred configuration path:
//   - use SqrtApproxTermsWithPolicy with p.Seed set
//
// This compatibility wrapper remains available for now.
func SqrtApproxTermsWithSeedAndPolicy(x, seed Rational, p SqrtPolicy, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsWithSeedAndPolicy: negative digits %d", digits)
	}
	pp := p
	pp.Seed = &seed
	return SqrtApproxTermsWithPolicy(x, pp, digits)
}

// sqrt_api.go v3
