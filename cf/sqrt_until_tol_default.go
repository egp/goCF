// sqrt_until_tol_default.go v1
package cf

// SqrtApproxRationalUntilResidualDefault performs bounded Newton iteration for
// sqrt(x), using DefaultSqrtSeed(x), and stops early once |y^2 - x| <= tol.
func SqrtApproxRationalUntilResidualDefault(x Rational, maxSteps int, tol Rational) (Rational, bool, error) {
	seed, err := DefaultSqrtSeed(x)
	if err != nil {
		return Rational{}, false, err
	}
	return SqrtApproxRationalUntilResidual(x, seed, maxSteps, tol)
}

// sqrt_until_tol_default.go v1
