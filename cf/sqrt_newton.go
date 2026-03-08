// sqrt_newton.go v2
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

// sqrt_newton.go v2
