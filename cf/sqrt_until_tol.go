// sqrt_until_tol.go v1
package cf

import "fmt"

// SqrtApproxRationalUntilResidual performs Newton iterations for sqrt(x),
// stopping early once the exact absolute residual satisfies:
//
//	|y^2 - x| <= tol
//
// Returns:
//   - the last iterate examined/produced
//   - satisfied=true if the tolerance condition was met
//   - error on invalid input
//
// Behavior:
//   - maxSteps < 0 => error
//   - tol < 0      => error
//   - if x has an exact rational square root, returns it immediately with satisfied=true
//   - if maxSteps == 0, returns seed with satisfied determined from the residual
func SqrtApproxRationalUntilResidual(x, seed Rational, maxSteps int, tol Rational) (Rational, bool, error) {
	if maxSteps < 0 {
		return Rational{}, false, fmt.Errorf("SqrtApproxRationalUntilResidual: negative maxSteps %d", maxSteps)
	}
	if x.r.Sign() < 0 {
		return Rational{}, false, fmt.Errorf("sqrt of negative rational: %v", x)
	}
	if tol.Cmp(intRat(0)) < 0 {
		return Rational{}, false, fmt.Errorf("SqrtApproxRationalUntilResidual: negative tolerance %v", tol)
	}

	if root, ok, err := RationalSqrtExact(x); err != nil {
		return Rational{}, false, err
	} else if ok {
		return root, true, nil
	}

	if seed.r.Sign() == 0 {
		return Rational{}, false, fmt.Errorf("SqrtApproxRationalUntilResidual: zero seed")
	}
	if seed.r.Sign() < 0 {
		return Rational{}, false, fmt.Errorf("SqrtApproxRationalUntilResidual: negative seed %v", seed)
	}

	check := func(y Rational) (bool, error) {
		r, err := SqrtResidualAbs(x, y)
		if err != nil {
			return false, err
		}
		return r.Cmp(tol) <= 0, nil
	}

	if maxSteps == 0 {
		ok, err := check(seed)
		if err != nil {
			return Rational{}, false, err
		}
		return seed, ok, nil
	}

	y := seed
	for i := 0; i < maxSteps; i++ {
		next, err := NewtonSqrtStep(x, y)
		if err != nil {
			return Rational{}, false, err
		}
		ok, err := check(next)
		if err != nil {
			return Rational{}, false, err
		}
		if ok {
			return next, true, nil
		}
		y = next
	}
	return y, false, nil
}

// sqrt_until_tol.go v1
