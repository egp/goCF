// sqrt_residual.go v1
package cf

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

// sqrt_residual.go v1
