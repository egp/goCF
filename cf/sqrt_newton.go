// sqrt_newton.go v3
package cf

import (
	"fmt"
	"math/big"
)

// RationalSqrtExact reports whether x has an exact rational square root.
// If so, it returns that root (the nonnegative one).
func RationalSqrtExact(x Rational) (Rational, bool, error) {
	if x.r.Sign() < 0 {
		return Rational{}, false, fmt.Errorf("sqrt of negative rational: %v", x)
	}

	num, den := x.ratNumDen()
	sn, ok := exactSqrtBig(num)
	if !ok {
		return Rational{}, false, nil
	}
	sd, ok := exactSqrtBig(den)
	if !ok {
		return Rational{}, false, nil
	}

	r, err := newRationalBig(sn, sd)
	if err != nil {
		return Rational{}, false, err
	}
	return r, true, nil
}

// NewtonSqrtStep computes one exact Newton iteration for sqrt(x):
//
//	y_(n+1) = (y_n + x / y_n) / 2
//
// Preconditions:
//   - x >= 0
//   - y != 0
//   - y > 0 (we stay on the nonnegative branch)
func NewtonSqrtStep(x, y Rational) (Rational, error) {
	if x.r.Sign() < 0 {
		return Rational{}, fmt.Errorf("sqrt of negative rational: %v", x)
	}
	if y.r.Sign() == 0 {
		return Rational{}, fmt.Errorf("NewtonSqrtStep: zero seed")
	}
	if y.r.Sign() < 0 {
		return Rational{}, fmt.Errorf("NewtonSqrtStep: negative seed %v", y)
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

// NewtonSqrtIterates returns successive exact Newton iterates for sqrt(x),
// starting from seed. The returned slice has length steps.
//
// Behavior:
//   - steps < 0 => error
//   - if x has an exact rational square root, that exact root is returned
//     repeated steps times
//   - otherwise, the sequence is produced by repeated NewtonSqrtStep calls
func NewtonSqrtIterates(x, seed Rational, steps int) ([]Rational, error) {
	if steps < 0 {
		return nil, fmt.Errorf("NewtonSqrtIterates: negative steps %d", steps)
	}
	if x.r.Sign() < 0 {
		return nil, fmt.Errorf("sqrt of negative rational: %v", x)
	}
	if steps == 0 {
		return []Rational{}, nil
	}

	if root, ok, err := RationalSqrtExact(x); err != nil {
		return nil, err
	} else if ok {
		out := make([]Rational, steps)
		for i := range out {
			out[i] = root
		}
		return out, nil
	}

	if seed.r.Sign() == 0 {
		return nil, fmt.Errorf("NewtonSqrtIterates: zero seed")
	}
	if seed.r.Sign() < 0 {
		return nil, fmt.Errorf("NewtonSqrtIterates: negative seed %v", seed)
	}

	out := make([]Rational, 0, steps)
	y := seed
	for i := 0; i < steps; i++ {
		next, err := NewtonSqrtStep(x, y)
		if err != nil {
			return nil, err
		}
		out = append(out, next)
		y = next
	}
	return out, nil
}

// SqrtApproxRational returns a single exact rational approximation to sqrt(x)
// after the requested number of Newton iterations from the given seed.
//
// Behavior:
//   - steps < 0 => error
//   - steps == 0 => returns seed (after domain/seed validation, unless x is exact square)
//   - if x has an exact rational square root, returns that exact root immediately
func SqrtApproxRational(x, seed Rational, steps int) (Rational, error) {
	if steps < 0 {
		return Rational{}, fmt.Errorf("SqrtApproxRational: negative steps %d", steps)
	}
	if x.r.Sign() < 0 {
		return Rational{}, fmt.Errorf("sqrt of negative rational: %v", x)
	}

	if root, ok, err := RationalSqrtExact(x); err != nil {
		return Rational{}, err
	} else if ok {
		return root, nil
	}

	if seed.r.Sign() == 0 {
		return Rational{}, fmt.Errorf("SqrtApproxRational: zero seed")
	}
	if seed.r.Sign() < 0 {
		return Rational{}, fmt.Errorf("SqrtApproxRational: negative seed %v", seed)
	}
	if steps == 0 {
		return seed, nil
	}

	ys, err := NewtonSqrtIterates(x, seed, steps)
	if err != nil {
		return Rational{}, err
	}
	return ys[len(ys)-1], nil
}

// SqrtResidual returns the exact residual y^2 - x.
func SqrtResidual(x, y Rational) (Rational, error) {
	y2, err := y.Mul(y)
	if err != nil {
		return Rational{}, err
	}
	return y2.Sub(x)
}

// SqrtResidualAbs returns the exact absolute residual |y^2 - x|.
func SqrtResidualAbs(x, y Rational) (Rational, error) {
	r, err := SqrtResidual(x, y)
	if err != nil {
		return Rational{}, err
	}
	if r.Cmp(intRat(0)) < 0 {
		return intRat(0).Sub(r)
	}
	return r, nil
}

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

// SqrtApproxRationalUntilExactDefault performs bounded Newton iteration for
// sqrt(x), using DefaultSqrtSeed(x), and stops early if exact convergence occurs.
func SqrtApproxRationalUntilExactDefault(x Rational, maxSteps int) (Rational, bool, error) {
	seed, err := DefaultSqrtSeed(x)
	if err != nil {
		return Rational{}, false, err
	}
	return SqrtApproxRationalUntilExact(x, seed, maxSteps)
}

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

// SqrtApproxRationalUntilResidualDefault performs bounded Newton iteration for
// sqrt(x), using DefaultSqrtSeed(x), and stops early once |y^2 - x| <= tol.
func SqrtApproxRationalUntilResidualDefault(x Rational, maxSteps int, tol Rational) (Rational, bool, error) {
	seed, err := DefaultSqrtSeed(x)
	if err != nil {
		return Rational{}, false, err
	}
	return SqrtApproxRationalUntilResidual(x, seed, maxSteps, tol)
}

// exactSqrtBig returns sqrt(n),true iff n is a perfect square and n >= 0.
func exactSqrtBig(n *big.Int) (*big.Int, bool) {
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

// sqrt_newton.go v3
