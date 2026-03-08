// sqrt_until_exact.go v1
package cf

import "fmt"

// SqrtApproxRationalUntilExact performs Newton iterations for sqrt(x),
// stopping early if an iterate becomes exact (residual == 0), or after
// maxSteps iterations.
//
// Returns:
//   - the last iterate examined/produced
//   - exact=true if that iterate is an exact square root of x
//   - error on invalid input
//
// Behavior:
//   - maxSteps < 0 => error
//   - if x has an exact rational square root, returns it immediately with exact=true
//   - otherwise, if maxSteps == 0, returns seed with exact=false after validation
func SqrtApproxRationalUntilExact(x, seed Rational, maxSteps int) (Rational, bool, error) {
	if maxSteps < 0 {
		return Rational{}, false, fmt.Errorf("SqrtApproxRationalUntilExact: negative maxSteps %d", maxSteps)
	}
	if x.r.Sign() < 0 {
		return Rational{}, false, fmt.Errorf("sqrt of negative rational: %v", x)
	}

	if root, ok, err := RationalSqrtExact(x); err != nil {
		return Rational{}, false, err
	} else if ok {
		return root, true, nil
	}

	if seed.r.Sign() == 0 {
		return Rational{}, false, fmt.Errorf("SqrtApproxRationalUntilExact: zero seed")
	}
	if seed.r.Sign() < 0 {
		return Rational{}, false, fmt.Errorf("SqrtApproxRationalUntilExact: negative seed %v", seed)
	}
	if maxSteps == 0 {
		return seed, false, nil
	}

	y := seed
	for i := 0; i < maxSteps; i++ {
		next, err := NewtonSqrtStep(x, y)
		if err != nil {
			return Rational{}, false, err
		}
		res, err := SqrtResidual(x, next)
		if err != nil {
			return Rational{}, false, err
		}
		if res.Cmp(intRat(0)) == 0 {
			return next, true, nil
		}
		y = next
	}
	return y, false, nil
}

// sqrt_until_exact.go v1
