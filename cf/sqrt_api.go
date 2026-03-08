// sqrt_api.go v2
package cf

import "fmt"

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
	p := DefaultSqrtPolicy()
	approx, _, err := SqrtApproxRationalUntilResidualDefault(x, p.MaxSteps, p.Tol)
	return approx, err
}

// SqrtApproxCF returns a ContinuedFraction source for the bounded default
// sqrt approximation produced by SqrtApprox.
func SqrtApproxCF(x Rational) (ContinuedFraction, error) {
	approx, err := SqrtApprox(x)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

// SqrtApproxTermsAuto returns up to digits CF terms for the bounded default
// sqrt approximation produced by SqrtApprox.
func SqrtApproxTermsAuto(x Rational, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsAuto: negative digits %d", digits)
	}
	cf, err := SqrtApproxCF(x)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// sqrt_api.go v2
