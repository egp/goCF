// sqrt_gcf.go v12
package cf

import (
	"fmt"
	"math/big"
)

const (
	sqrtGCFExactBootstrapTermBudget = 256
	sqrtGCFNewtonSteps              = 4
	sqrtGCFPerNextIngestBudget      = 32
)

type sqrtBootstrapCFStream struct {
	src            GCFSource
	prepared       bool
	preparing      bool
	out            ContinuedFraction
	err            error
	bufferedTerms  [][2]int64
	inputExhausted bool
	ingestedTerms  int
}

func SqrtGCF(src GCFSource) (ContinuedFraction, error) {
	if src == nil {
		return nil, fmt.Errorf("SqrtGCF: nil src")
	}
	return &sqrtBootstrapCFStream{src: src}, nil
}

func (s *sqrtBootstrapCFStream) Next() (int64, bool) {
	if s.err != nil {
		return 0, false
	}
	if s.out != nil {
		return s.out.Next()
	}
	if !s.prepared {
		if err := s.prepare(); err != nil {
			s.err = err
			return 0, false
		}
	}
	if s.out == nil {
		return 0, false
	}
	return s.out.Next()
}

func (s *sqrtBootstrapCFStream) Err() error {
	return s.err
}

func (s *sqrtBootstrapCFStream) prepare() error {
	if s.prepared {
		return nil
	}
	if s.preparing {
		return fmt.Errorf("sqrtBootstrapCFStream.prepare: re-entrant prepare")
	}
	s.preparing = true
	defer func() { s.preparing = false }()

	if cf, ok, err := sqrtExactRootCFViaSourceMetadata(s.src); err != nil {
		s.prepared = true
		return err
	} else if ok {
		s.out = cf
		s.prepared = true
		return nil
	}

	if err := s.ingestMore(sqrtGCFPerNextIngestBudget); err != nil {
		s.prepared = true
		return err
	}
	if !s.inputExhausted {
		// Live unresolved state: not an error.
		return nil
	}

	x, err := GCFSourceConvergent(NewSliceGCF(s.bufferedTerms...), len(s.bufferedTerms))
	if err != nil {
		s.prepared = true
		return err
	}

	root, ok, err := sqrtExactNonnegativeRational(x)
	if err != nil {
		s.prepared = true
		return err
	}
	if ok {
		s.out = NewRationalCF(root)
		s.prepared = true
		return nil
	}

	state, err := newSqrtBootstrapState(x)
	if err != nil {
		s.prepared = true
		return err
	}
	for i := 0; i < sqrtGCFNewtonSteps; i++ {
		if err := state.Step(); err != nil {
			s.prepared = true
			return err
		}
	}
	cf, err := state.CF()
	if err != nil {
		s.prepared = true
		return err
	}

	s.out = cf
	s.prepared = true
	return nil
}

func (s *sqrtBootstrapCFStream) ingestMore(limit int) error {
	if limit < 0 {
		return fmt.Errorf("sqrtBootstrapCFStream.ingestMore: negative limit %d", limit)
	}
	if s.inputExhausted {
		return nil
	}
	if len(s.bufferedTerms) >= sqrtGCFExactBootstrapTermBudget {
		return fmt.Errorf(
			"SqrtGCF: exact finite value not available within bootstrap term budget %d",
			sqrtGCFExactBootstrapTermBudget,
		)
	}

	for i := 0; i < limit; i++ {
		if len(s.bufferedTerms) >= sqrtGCFExactBootstrapTermBudget {
			return fmt.Errorf(
				"SqrtGCF: exact finite value not available within bootstrap term budget %d",
				sqrtGCFExactBootstrapTermBudget,
			)
		}
		p, q, ok := s.src.NextPQ()
		if !ok {
			s.inputExhausted = true
			if len(s.bufferedTerms) == 0 {
				return fmt.Errorf("SqrtGCF: empty source")
			}
			return nil
		}
		s.bufferedTerms = append(s.bufferedTerms, [2]int64{p, q})
		s.ingestedTerms++
	}
	return nil
}

func sqrtExactRootCFViaSourceMetadata(src GCFSource) (ContinuedFraction, bool, error) {
	qr, ok := src.(QuadraticRadicalSource)
	if !ok {
		return nil, false, nil
	}

	n, ok := qr.Radicand()
	if !ok {
		return nil, false, nil
	}
	if n < 0 {
		return nil, false, fmt.Errorf("SqrtGCF: negative radicand metadata %d", n)
	}

	root, ok := sqrtExactBigInt(big.NewInt(n))
	if !ok {
		return nil, false, nil
	}

	r, err := newRationalBig(root, big.NewInt(1))
	if err != nil {
		return nil, false, err
	}
	return NewRationalCF(r), true, nil
}

func sqrtExactNonnegativeRational(x Rational) (Rational, bool, error) {
	if x.Cmp(intRat(0)) < 0 {
		return Rational{}, false, fmt.Errorf("SqrtGCF: negative input %v", x)
	}

	num, den := x.ratNumDen()
	numRoot, ok := sqrtExactBigInt(num)
	if !ok {
		return Rational{}, false, nil
	}
	denRoot, ok := sqrtExactBigInt(den)
	if !ok {
		return Rational{}, false, nil
	}

	r, err := newRationalBig(numRoot, denRoot)
	if err != nil {
		return Rational{}, false, err
	}
	return r, true, nil
}

func sqrtExactBigInt(n *big.Int) (*big.Int, bool) {
	if n.Sign() < 0 {
		return nil, false
	}
	if n.Sign() == 0 {
		return big.NewInt(0), true
	}

	root := new(big.Int).Sqrt(n)
	sq := new(big.Int).Mul(root, root)
	if sq.Cmp(n) != 0 {
		return nil, false
	}
	return root, true
}

func sqrtNewtonStep(x Rational, y Rational) (Rational, error) {
	if y.Cmp(intRat(0)) == 0 {
		return Rational{}, fmt.Errorf("sqrtNewtonStep: zero iterate")
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

func sqrtNewtonApprox(x Rational, steps int) (Rational, error) {
	if x.Cmp(intRat(0)) < 0 {
		return Rational{}, fmt.Errorf("sqrtNewtonApprox: negative input %v", x)
	}
	if steps < 0 {
		return Rational{}, fmt.Errorf("sqrtNewtonApprox: negative steps %d", steps)
	}
	if x.Cmp(intRat(0)) == 0 {
		return intRat(0), nil
	}

	y := intRat(1)
	for i := 0; i < steps; i++ {
		next, err := sqrtNewtonStep(x, y)
		if err != nil {
			return Rational{}, err
		}
		y = next
	}
	return y, nil
}

type sqrtBootstrapState struct {
	x Rational
	y Rational
}

func newSqrtBootstrapState(x Rational) (*sqrtBootstrapState, error) {
	if x.Cmp(intRat(0)) < 0 {
		return nil, fmt.Errorf("newSqrtBootstrapState: negative input %v", x)
	}
	return &sqrtBootstrapState{
		x: x,
		y: intRat(1),
	}, nil
}

func (s *sqrtBootstrapState) Step() error {
	next, err := sqrtNewtonStep(s.x, s.y)
	if err != nil {
		return err
	}
	s.y = next
	return nil
}

func (s *sqrtBootstrapState) CF() (ContinuedFraction, error) {
	return NewRationalCF(s.y), nil
}

// sqrt_gcf.go v12
