// sqrt_gcf.go v1
package cf

import (
	"fmt"
	"math/big"
)

const sqrtGCFExactBootstrapMaxTerms = 256

// SqrtGCF is the new canonical sqrt unary entry point.
//
// Bootstrap v1 behavior:
//   - accepts any GCFSource
//   - currently supports exact finite GCF inputs whose value is a nonnegative
//     perfect-square integer
//   - returns a regular CF for the exact square root in those cases
//   - returns "not implemented" for all other inputs for now
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

	root, ok, err := sqrtExactNonnegativeIntegerRational(x)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("SqrtGCF: not implemented for non-square input %v", x)
	}

	return NewRationalCF(root), nil
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

func sqrtExactNonnegativeIntegerRational(x Rational) (Rational, bool, error) {
	num, den := x.ratNumDen()

	if x.Cmp(intRat(0)) < 0 {
		return Rational{}, false, fmt.Errorf("SqrtGCF: negative input %v", x)
	}
	if den.Cmp(big.NewInt(1)) != 0 {
		return Rational{}, false, nil
	}

	root, ok := sqrtExactBigInt(num)
	if !ok {
		return Rational{}, false, nil
	}

	r, err := newRationalBig(root, big.NewInt(1))
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

// sqrt_gcf.go v1
