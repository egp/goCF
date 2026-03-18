// cf/sqrt_unary_newton.go v4
package cf

import "fmt"

// SqrtUnaryNewtonStepExact applies exactly one Newton update for solving y^2 = x:
//
//	y' = (y + x/y) / 2
//
// Domain for the current real positive sqrt operator kernel:
//   - x > 0
//   - y > 0
func SqrtUnaryNewtonStepExact(x Rational, y Rational) (Rational, error) {
	if x.Cmp(intRat(0)) <= 0 {
		return Rational{}, fmt.Errorf("SqrtUnaryNewtonStepExact: nonpositive input %v", x)
	}
	if y.Cmp(intRat(0)) <= 0 {
		return Rational{}, fmt.Errorf("SqrtUnaryNewtonStepExact: nonpositive iterate %v", y)
	}

	xy, err := x.Div(y)
	if err != nil {
		return Rational{}, err
	}

	sum, err := y.Add(xy)
	if err != nil {
		return Rational{}, err
	}

	return sum.Div(mustRat(2, 1))
}

// cf/sqrt_unary_newton.go v4
