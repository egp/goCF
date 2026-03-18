// cf/sqrt_unary_operator.go v2
package cf

import "fmt"

type sqrtUnaryOperatorSnapshot struct {
	HasInputApprox bool
	InputApprox    *GCFApprox
	CurrentY       *Rational
}

type sqrtUnaryOperator struct {
	src           GCFSource
	initialY      Rational
	currentY      Rational
	bufferedTerms [][2]int64
	inputApprox   *GCFApprox
}

func newSqrtUnaryOperator(src GCFSource, initialY Rational) (*sqrtUnaryOperator, error) {
	if src == nil {
		return nil, fmt.Errorf("newSqrtUnaryOperator: nil src")
	}
	if initialY.Cmp(intRat(0)) <= 0 {
		return nil, fmt.Errorf("newSqrtUnaryOperator: nonpositive iterate %v", initialY)
	}

	return &sqrtUnaryOperator{
		src:      src,
		initialY: initialY,
		currentY: initialY,
	}, nil
}

func (s *sqrtUnaryOperator) snapshot() sqrtUnaryOperatorSnapshot {
	y := s.currentY

	var a *GCFApprox
	if s.inputApprox != nil {
		cp := *s.inputApprox
		a = &cp
	}

	return sqrtUnaryOperatorSnapshot{
		HasInputApprox: s.inputApprox != nil,
		InputApprox:    a,
		CurrentY:       &y,
	}
}

func (s *sqrtUnaryOperator) ingestOneAndRefine() error {
	p, q, ok := s.src.NextPQ()
	if !ok {
		return fmt.Errorf("sqrtUnaryOperator.ingestOneAndRefine: source exhausted")
	}

	s.bufferedTerms = append(s.bufferedTerms, [2]int64{p, q})

	a, err := GCFApproxFromPrefix(NewSliceGCF(s.bufferedTerms...), len(s.bufferedTerms))
	if err != nil {
		return err
	}
	s.inputApprox = &a

	state, err := newSqrtUnaryState(a.Convergent, s.currentY)
	if err != nil {
		return err
	}
	if err := state.step(); err != nil {
		return err
	}

	s.currentY = state.yValue()
	return nil
}

// cf/sqrt_unary_operator.go v2
