// cf/sqrt_unary_state.go v2
package cf

type sqrtUnaryState struct {
	x Rational
	y Rational
}

func newSqrtUnaryState(x Rational, y Rational) (*sqrtUnaryState, error) {
	if _, err := sqrtUnaryNewtonStepExact(x, y); err != nil {
		return nil, err
	}

	return &sqrtUnaryState{
		x: x,
		y: y,
	}, nil
}

func (s *sqrtUnaryState) xValue() Rational {
	return s.x
}

func (s *sqrtUnaryState) yValue() Rational {
	return s.y
}

func (s *sqrtUnaryState) updateInput(x Rational) error {
	if x.Cmp(intRat(0)) <= 0 {
		// Reuse the same domain rule/message shape as the Newton kernel.
		if _, err := sqrtUnaryNewtonStepExact(x, s.y); err != nil {
			return err
		}
	}
	s.x = x
	return nil
}

func (s *sqrtUnaryState) step() error {
	next, err := sqrtUnaryNewtonStepExact(s.x, s.y)
	if err != nil {
		return err
	}
	s.y = next
	return nil
}

// cf/sqrt_unary_state.go v2
