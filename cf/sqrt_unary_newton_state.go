// cf/sqrt_unary_newton_state.go v2
package cf

type sqrtUnaryNewtonState struct {
	x Rational
	y Rational
}

func newSqrtUnaryNewtonState(x Rational, y Rational) (*sqrtUnaryNewtonState, error) {
	if _, err := sqrtUnaryNewtonStepExact(x, y); err != nil {
		return nil, err
	}

	return &sqrtUnaryNewtonState{
		x: x,
		y: y,
	}, nil
}

func (s *sqrtUnaryNewtonState) xValue() Rational {
	return s.x
}

func (s *sqrtUnaryNewtonState) yValue() Rational {
	return s.y
}

func (s *sqrtUnaryNewtonState) step() error {
	next, err := sqrtUnaryNewtonStepExact(s.x, s.y)
	if err != nil {
		return err
	}
	s.y = next
	return nil
}

// cf/sqrt_unary_newton_state.go v2
