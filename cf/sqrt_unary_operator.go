// cf/sqrt_unary_operator.go v14
package cf

import (
	"fmt"
	"math/big"
)

type sqrtUnaryOperatorSnapshot struct {
	HasInputApprox bool
	InputApprox    *GCFApprox
	CurrentY       *Rational
	Residual       *sqrtUnaryResidualSnapshot
	SqrtEnclosure  *Range
	ForcedDigit    *big.Int
}

type sqrtUnaryOperator struct {
	src                GCFSource
	initialY           Rational
	currentY           Rational
	policy             sqrtUnaryRefinementPolicy
	prefixState        *gcfPrefixState
	inputApprox        *GCFApprox
	currentResidual    *sqrtUnaryResidualSnapshot
	currentEnclosure   *Range
	currentForcedDigit *big.Int
}

func newSqrtUnaryOperator(src GCFSource, initialY Rational, policy sqrtUnaryRefinementPolicy) (*sqrtUnaryOperator, error) {
	if src == nil {
		return nil, fmt.Errorf("newSqrtUnaryOperator: nil src")
	}
	if initialY.Cmp(intRat(0)) <= 0 {
		return nil, fmt.Errorf("newSqrtUnaryOperator: nonpositive iterate %v", initialY)
	}
	if err := policy.validate(); err != nil {
		return nil, err
	}

	return &sqrtUnaryOperator{
		src:         src,
		initialY:    initialY,
		currentY:    initialY,
		policy:      policy,
		prefixState: newGcfPrefixState(),
	}, nil
}

func (s *sqrtUnaryOperator) snapshot() sqrtUnaryOperatorSnapshot {
	y := s.currentY

	var a *GCFApprox
	if s.inputApprox != nil {
		cp := *s.inputApprox
		a = &cp
	}

	var r *sqrtUnaryResidualSnapshot
	if s.currentResidual != nil {
		cp := *s.currentResidual
		r = &cp
	}

	var e *Range
	if s.currentEnclosure != nil {
		cp := *s.currentEnclosure
		e = &cp
	}

	var d *big.Int
	if s.currentForcedDigit != nil {
		d = new(big.Int).Set(s.currentForcedDigit)
	}

	return sqrtUnaryOperatorSnapshot{
		HasInputApprox: s.inputApprox != nil,
		InputApprox:    a,
		CurrentY:       &y,
		Residual:       r,
		SqrtEnclosure:  e,
		ForcedDigit:    d,
	}
}

func (s *sqrtUnaryOperator) ingestOneAndRefine() error {
	p, q, ok := s.src.NextPQ()
	if !ok {
		return fmt.Errorf("sqrtUnaryOperator.ingestOneAndRefine: source exhausted")
	}

	if err := s.prefixState.ingestOne(p, q); err != nil {
		return err
	}

	a := s.prefixState.approx()
	s.inputApprox = &a

	state, err := newSqrtUnaryState(a.Convergent, s.currentY)
	if err != nil {
		return err
	}
	for i := 0; i < s.policy.StepsPerInput; i++ {
		if err := state.step(); err != nil {
			return err
		}
	}

	s.currentY = state.yValue()

	resid, err := state.residualSnapshot()
	if err != nil {
		return err
	}
	s.currentResidual = &resid

	if a.Range != nil {
		enclosure, err := sqrtUnaryRangeEnclosureFromInputRange(*a.Range, s.currentY)
		if err != nil {
			return err
		}
		s.currentEnclosure = &enclosure

		d, ok, err := sqrtUnaryNextDigitIfForced(enclosure)
		if err != nil {
			return err
		}
		if ok {
			s.currentForcedDigit = d
		} else {
			s.currentForcedDigit = nil
		}
	} else {
		s.currentEnclosure = nil
		s.currentForcedDigit = nil
	}

	return nil
}

func (s *sqrtUnaryOperator) forceFirstDigitWithin(maxIngests int) (*big.Int, bool, error) {
	if maxIngests < 0 {
		return nil, false, fmt.Errorf(
			"sqrtUnaryOperator.forceFirstDigitWithin: negative maxIngests %d",
			maxIngests,
		)
	}

	if s.currentForcedDigit != nil {
		return new(big.Int).Set(s.currentForcedDigit), true, nil
	}

	for i := 0; i < maxIngests; i++ {
		if err := s.ingestOneAndRefine(); err != nil {
			return nil, false, err
		}
		if s.currentForcedDigit != nil {
			return new(big.Int).Set(s.currentForcedDigit), true, nil
		}
	}

	return nil, false, nil
}

// cf/sqrt_unary_operator.go v14
