// sqrt_until_exact_default.go v1
package cf

// SqrtApproxRationalUntilExactDefault performs bounded Newton iteration for
// sqrt(x), using DefaultSqrtSeed(x), and stops early if exact convergence occurs.
func SqrtApproxRationalUntilExactDefault(x Rational, maxSteps int) (Rational, bool, error) {
	seed, err := DefaultSqrtSeed(x)
	if err != nil {
		return Rational{}, false, err
	}
	return SqrtApproxRationalUntilExact(x, seed, maxSteps)
}

// sqrt_until_exact_default.go v1
