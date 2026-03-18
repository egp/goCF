// sqrt_gcf.go v3
package cf

import (
	"fmt"
	"math/big"
)

const (
	sqrtGCFExactBootstrapMaxTerms = 256
	sqrtGCFNewtonSteps            = 4
)

// SqrtGCF is the new canonical sqrt unary entry point.
//
// Current behavior:
//   - accepts any GCFSource
//   - for exact finite nonnegative perfect-square rational input, returns the
//     exact square root as a regular CF
//   - for exact finite nonnegative non-square input, returns a Newton rational
//     approximation as a regular CF
//   - for non-terminating input, reports not implemented for now
func SqrtGCF(src GCFSource) (ContinuedFraction, error) {
	if src == nil {
		return nil, fmt.Errorf("SqrtGCF: nil src")
	}

	x, exact, err := sqrtGCFExactFiniteValue(src, sqrtGCFExactBootstrapMaxTerms)
	if err != nil {
		return nil, err
	}
	if !exact {
		return nil, fmt.Errorf("SqrtGCF: not implemented for non-terminating input")
	}

	root, ok, err := sqrtExactNonnegativeRational(x)
	if err != nil {
		return nil, err
	}
	if ok {
		return NewRationalCF(root), nil
	}

	approx, err := sqrtNewtonApprox(x, sqrtGCFNewtonSteps)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

func sqrtGCFExactFiniteValue(src GCFSource, maxTerms int) (Rational, bool, error) {
	terms := make([][2]int64, 0, 8)

	for i := 0; i < maxTerms; i++ {
		p, q, ok := src.NextPQ()
		if !ok {
			if len(terms) == 0 {
				return Rational{}, false, fmt.Errorf("SqrtGCF: empty source")
			}
			x, err := GCFSourceConvergent(NewSliceGCF(terms...), len(terms))
			if err != nil {
				return Rational{}, false, err
			}
			return x, true, nil
		}
		terms = append(terms, [2]int64{p, q})
	}

	return Rational{}, false, nil
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

// sqrt_gcf.go v3
