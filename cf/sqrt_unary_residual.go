// cf/sqrt_unary_residual.go v2
package cf

type sqrtUnaryResidualSnapshot struct {
	X        Rational
	Y        Rational
	YSquared Rational
	Residual Rational
}

func sqrtUnaryResidualExact(x Rational, y Rational) (Rational, error) {
	yy, err := y.Mul(y)
	if err != nil {
		return Rational{}, err
	}
	return yy.Sub(x)
}

func (s *sqrtUnaryState) residualSnapshot() (sqrtUnaryResidualSnapshot, error) {
	yy, err := s.y.Mul(s.y)
	if err != nil {
		return sqrtUnaryResidualSnapshot{}, err
	}
	r, err := yy.Sub(s.x)
	if err != nil {
		return sqrtUnaryResidualSnapshot{}, err
	}

	return sqrtUnaryResidualSnapshot{
		X:        s.x,
		Y:        s.y,
		YSquared: yy,
		Residual: r,
	}, nil
}

// cf/sqrt_unary_residual.go v2
