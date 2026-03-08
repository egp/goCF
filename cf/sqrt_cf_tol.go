// sqrt_cf_tol.go v1
package cf

import "fmt"

// NewSqrtApproxCFUntilResidual returns a ContinuedFraction source for the bounded
// Newton approximation to sqrt(x), stopping early once |y^2 - x| <= tol.
func NewSqrtApproxCFUntilResidual(x, seed Rational, maxSteps int, tol Rational) (ContinuedFraction, error) {
	approx, _, err := SqrtApproxRationalUntilResidual(x, seed, maxSteps, tol)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

// NewSqrtApproxCFUntilResidualDefault is the default-seed wrapper around
// NewSqrtApproxCFUntilResidual.
func NewSqrtApproxCFUntilResidualDefault(x Rational, maxSteps int, tol Rational) (ContinuedFraction, error) {
	approx, _, err := SqrtApproxRationalUntilResidualDefault(x, maxSteps, tol)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

// SqrtApproxTermsUntilResidual computes CF terms for the bounded Newton
// approximation to sqrt(x), stopping early once |y^2 - x| <= tol.
func SqrtApproxTermsUntilResidual(x, seed Rational, maxSteps int, tol Rational, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsUntilResidual: negative digits %d", digits)
	}
	cf, err := NewSqrtApproxCFUntilResidual(x, seed, maxSteps, tol)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// SqrtApproxTermsUntilResidualDefault is the default-seed wrapper around
// SqrtApproxTermsUntilResidual.
func SqrtApproxTermsUntilResidualDefault(x Rational, maxSteps int, tol Rational, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsUntilResidualDefault: negative digits %d", digits)
	}
	cf, err := NewSqrtApproxCFUntilResidualDefault(x, maxSteps, tol)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// sqrt_cf_tol.go v1
