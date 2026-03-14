// sqrt_core_exact.go v1
package cf

import (
	"fmt"
	"math/big"
)

// SqrtCoreRationalExact reports whether x has an exact rational square root.
// If so, it returns that root (the nonnegative one).
func SqrtCoreRationalExact(x Rational) (Rational, bool, error) {
	if x.r.Sign() < 0 {
		return Rational{}, false, fmt.Errorf("sqrt of negative rational: %v", x)
	}

	num, den := x.ratNumDen()
	sn, ok := sqrtCoreExactBig(num)
	if !ok {
		return Rational{}, false, nil
	}
	sd, ok := sqrtCoreExactBig(den)
	if !ok {
		return Rational{}, false, nil
	}

	r, err := newRationalBig(sn, sd)
	if err != nil {
		return Rational{}, false, err
	}
	return r, true, nil
}

// SqrtCoreNewtonStep computes one exact Newton iteration for sqrt(x):
//
//	y_(n+1) = (y_n + x / y_n) / 2
//
// Preconditions:
//   - x >= 0
//   - y != 0
//   - y > 0
func SqrtCoreNewtonStep(x, y Rational) (Rational, error) {
	if x.r.Sign() < 0 {
		return Rational{}, fmt.Errorf("sqrt of negative rational: %v", x)
	}
	if y.r.Sign() == 0 {
		return Rational{}, fmt.Errorf("SqrtCoreNewtonStep: zero seed")
	}
	if y.r.Sign() < 0 {
		return Rational{}, fmt.Errorf("SqrtCoreNewtonStep: negative seed %v", y)
	}

	quot, err := x.Div(y)
	if err != nil {
		return Rational{}, err
	}
	sum, err := y.Add(quot)
	if err != nil {
		return Rational{}, err
	}
	return sum.Div(intRat(2))
}

// SqrtCoreNewtonIterates returns successive exact Newton iterates for sqrt(x),
// starting from seed. The returned slice has length steps.
func SqrtCoreNewtonIterates(x, seed Rational, steps int) ([]Rational, error) {
	if steps < 0 {
		return nil, fmt.Errorf("SqrtCoreNewtonIterates: negative steps %d", steps)
	}
	if x.r.Sign() < 0 {
		return nil, fmt.Errorf("sqrt of negative rational: %v", x)
	}
	if steps == 0 {
		return []Rational{}, nil
	}

	if root, ok, err := SqrtCoreRationalExact(x); err != nil {
		return nil, err
	} else if ok {
		out := make([]Rational, steps)
		for i := range out {
			out[i] = root
		}
		return out, nil
	}

	if seed.r.Sign() == 0 {
		return nil, fmt.Errorf("SqrtCoreNewtonIterates: zero seed")
	}
	if seed.r.Sign() < 0 {
		return nil, fmt.Errorf("SqrtCoreNewtonIterates: negative seed %v", seed)
	}

	out := make([]Rational, 0, steps)
	y := seed
	for i := 0; i < steps; i++ {
		next, err := SqrtCoreNewtonStep(x, y)
		if err != nil {
			return nil, err
		}
		out = append(out, next)
		y = next
	}
	return out, nil
}

// SqrtCoreApproxRational returns a single exact rational approximation to sqrt(x)
// after the requested number of Newton iterations from the given seed.
func SqrtCoreApproxRational(x, seed Rational, steps int) (Rational, error) {
	if steps < 0 {
		return Rational{}, fmt.Errorf("SqrtCoreApproxRational: negative steps %d", steps)
	}
	if x.r.Sign() < 0 {
		return Rational{}, fmt.Errorf("sqrt of negative rational: %v", x)
	}

	if root, ok, err := SqrtCoreRationalExact(x); err != nil {
		return Rational{}, err
	} else if ok {
		return root, nil
	}

	if seed.r.Sign() == 0 {
		return Rational{}, fmt.Errorf("SqrtCoreApproxRational: zero seed")
	}
	if seed.r.Sign() < 0 {
		return Rational{}, fmt.Errorf("SqrtCoreApproxRational: negative seed %v", seed)
	}
	if steps == 0 {
		return seed, nil
	}

	ys, err := SqrtCoreNewtonIterates(x, seed, steps)
	if err != nil {
		return Rational{}, err
	}
	return ys[len(ys)-1], nil
}

// SqrtCoreResidual returns the exact residual y^2 - x.
func SqrtCoreResidual(x, y Rational) (Rational, error) {
	y2, err := y.Mul(y)
	if err != nil {
		return Rational{}, err
	}
	return y2.Sub(x)
}

// SqrtCoreResidualAbs returns the exact absolute residual |y^2 - x|.
func SqrtCoreResidualAbs(x, y Rational) (Rational, error) {
	r, err := SqrtCoreResidual(x, y)
	if err != nil {
		return Rational{}, err
	}
	if r.Cmp(intRat(0)) < 0 {
		return intRat(0).Sub(r)
	}
	return r, nil
}

// SqrtCoreApproxRationalUntilExact performs Newton iterations for sqrt(x),
// stopping early if an iterate becomes exact, or after maxSteps iterations.
func SqrtCoreApproxRationalUntilExact(x, seed Rational, maxSteps int) (Rational, bool, error) {
	if maxSteps < 0 {
		return Rational{}, false, fmt.Errorf("SqrtCoreApproxRationalUntilExact: negative maxSteps %d", maxSteps)
	}
	if x.r.Sign() < 0 {
		return Rational{}, false, fmt.Errorf("sqrt of negative rational: %v", x)
	}

	if root, ok, err := SqrtCoreRationalExact(x); err != nil {
		return Rational{}, false, err
	} else if ok {
		return root, true, nil
	}

	if seed.r.Sign() == 0 {
		return Rational{}, false, fmt.Errorf("SqrtCoreApproxRationalUntilExact: zero seed")
	}
	if seed.r.Sign() < 0 {
		return Rational{}, false, fmt.Errorf("SqrtCoreApproxRationalUntilExact: negative seed %v", seed)
	}
	if maxSteps == 0 {
		return seed, false, nil
	}

	y := seed
	for i := 0; i < maxSteps; i++ {
		next, err := SqrtCoreNewtonStep(x, y)
		if err != nil {
			return Rational{}, false, err
		}
		res, err := SqrtCoreResidual(x, next)
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

// SqrtCoreApproxRationalUntilResidual performs Newton iterations for sqrt(x),
// stopping early once |y^2 - x| <= tol.
func SqrtCoreApproxRationalUntilResidual(x, seed Rational, maxSteps int, tol Rational) (Rational, bool, error) {
	if maxSteps < 0 {
		return Rational{}, false, fmt.Errorf("SqrtCoreApproxRationalUntilResidual: negative maxSteps %d", maxSteps)
	}
	if x.r.Sign() < 0 {
		return Rational{}, false, fmt.Errorf("sqrt of negative rational: %v", x)
	}
	if tol.Cmp(intRat(0)) < 0 {
		return Rational{}, false, fmt.Errorf("SqrtCoreApproxRationalUntilResidual: negative tolerance %v", tol)
	}

	if root, ok, err := SqrtCoreRationalExact(x); err != nil {
		return Rational{}, false, err
	} else if ok {
		return root, true, nil
	}

	if seed.r.Sign() == 0 {
		return Rational{}, false, fmt.Errorf("SqrtCoreApproxRationalUntilResidual: zero seed")
	}
	if seed.r.Sign() < 0 {
		return Rational{}, false, fmt.Errorf("SqrtCoreApproxRationalUntilResidual: negative seed %v", seed)
	}

	check := func(y Rational) (bool, error) {
		r, err := SqrtCoreResidualAbs(x, y)
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
		next, err := SqrtCoreNewtonStep(x, y)
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

func sqrtCoreExactBig(n *big.Int) (*big.Int, bool) {
	if n.Sign() < 0 {
		return nil, false
	}
	s := new(big.Int).Sqrt(n)
	ss := new(big.Int).Mul(s, s)
	if ss.Cmp(n) != 0 {
		return nil, false
	}
	return s, true
}

// sqrt_core_exact.go v1
