// sqrt_api_policy.go v4
package cf

import "fmt"

// SqrtApproxWithPolicy computes a bounded rational approximation to sqrt(x)
// using the supplied policy.
func SqrtApproxWithPolicy(x Rational, p SqrtPolicy) (Rational, error) {
	if err := p.Validate(); err != nil {
		return Rational{}, err
	}

	if p.Seed != nil {
		approx, _, err := SqrtApproxRationalUntilResidual(x, *p.Seed, p.MaxSteps, p.Tol)
		return approx, err
	}

	approx, _, err := SqrtApproxRationalUntilResidualDefault(x, p.MaxSteps, p.Tol)
	return approx, err
}

// SqrtApproxCFWithPolicy returns a ContinuedFraction source for the bounded
// sqrt approximation produced by SqrtApproxWithPolicy.
func SqrtApproxCFWithPolicy(x Rational, p SqrtPolicy) (ContinuedFraction, error) {
	approx, err := SqrtApproxWithPolicy(x, p)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

// SqrtApproxTermsWithPolicy returns up to digits CF terms for the bounded
// sqrt approximation produced by SqrtApproxWithPolicy.
func SqrtApproxTermsWithPolicy(x Rational, p SqrtPolicy, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsWithPolicy: negative digits %d", digits)
	}
	cf, err := SqrtApproxCFWithPolicy(x, p)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// sqrt_api_policy.go v4
