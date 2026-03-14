// sqrt_canonical_api.go v1
package cf

import "fmt"

func sqrtApproxCanonical(x Rational) (Rational, error) {
	p := DefaultSqrtPolicy2()
	return sqrtApproxWithPolicyCanonical(x, p)
}

func sqrtApproxCFCanonical(x Rational) (ContinuedFraction, error) {
	approx, err := sqrtApproxCanonical(x)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

func sqrtApproxTermsCanonical(x Rational, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("sqrtApproxTermsCanonical: negative digits %d", digits)
	}
	cf, err := sqrtApproxCFCanonical(x)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

func sqrtApproxWithPolicyCanonical(x Rational, p SqrtPolicy2) (Rational, error) {
	if err := p.Validate(); err != nil {
		return Rational{}, err
	}

	if p.Seed != nil {
		approx, _, err := SqrtCoreApproxRationalUntilResidual(x, *p.Seed, p.MaxSteps, p.Tol)
		return approx, err
	}

	seed, err := SqrtSeedDefault(x)
	if err != nil {
		return Rational{}, err
	}
	approx, _, err := SqrtCoreApproxRationalUntilResidual(x, seed, p.MaxSteps, p.Tol)
	return approx, err
}

// sqrt_canonical_api.go v1
