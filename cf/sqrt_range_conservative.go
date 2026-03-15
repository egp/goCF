// sqrt_range_conservative.go v4
package cf

import (
	"fmt"
	"math/big"
)

const sqrtBoundScaleBits = 16

func sqrtBoundScale() *big.Int {
	return new(big.Int).Lsh(big.NewInt(1), sqrtBoundScaleBits)
}

// SqrtLowerBoundRational returns a rational lower bound L such that
//
//	L <= sqrt(x)
//
// Current implementation:
//   - rejects negative input
//   - exact-square fast path
//   - proof-safe scaled integer lower bracket for non-squares
func SqrtLowerBoundRational(x Rational) (Rational, error) {
	if x.Cmp(intRat(0)) < 0 {
		return Rational{}, fmt.Errorf("SqrtLowerBoundRational: negative input %v", x)
	}

	if root, ok, err := RationalSqrtExact(x); err != nil {
		return Rational{}, err
	} else if ok {
		return root, nil
	}

	if x.Cmp(intRat(0)) == 0 {
		return intRat(0), nil
	}

	n, d := x.ratNumDen() // x = n/d, d > 0

	k := sqrtBoundScale()
	k2 := new(big.Int).Mul(k, k)

	// sqrt(x) = sqrt(n*d) / d
	// so floor(sqrt(n*d)*k) / (d*k) is a lower bound
	nd := new(big.Int).Mul(n, d)
	scaled := new(big.Int).Mul(nd, k2)
	a := new(big.Int).Sqrt(scaled)

	den := new(big.Int).Mul(d, k)
	return newRationalBig(a, den)
}

// SqrtUpperBoundRational returns a rational upper bound U such that
//
//	sqrt(x) <= U
//
// Current implementation:
//   - rejects negative input
//   - exact-square fast path
//   - proof-safe scaled integer upper bracket for non-squares
func SqrtUpperBoundRational(x Rational) (Rational, error) {
	if x.Cmp(intRat(0)) < 0 {
		return Rational{}, fmt.Errorf("SqrtUpperBoundRational: negative input %v", x)
	}

	if root, ok, err := RationalSqrtExact(x); err != nil {
		return Rational{}, err
	} else if ok {
		return root, nil
	}

	if x.Cmp(intRat(0)) == 0 {
		return intRat(0), nil
	}

	n, d := x.ratNumDen() // x = n/d, d > 0

	k := sqrtBoundScale()
	k2 := new(big.Int).Mul(k, k)

	nd := new(big.Int).Mul(n, d)
	scaled := new(big.Int).Mul(nd, k2)
	a := new(big.Int).Sqrt(scaled)
	a2 := new(big.Int).Mul(a, a)
	if a2.Cmp(scaled) < 0 {
		a = new(big.Int).Add(a, big.NewInt(1))
	}

	den := new(big.Int).Mul(d, k)
	return newRationalBig(a, den)
}

// SqrtRangeConservative returns a proof-safe enclosure for sqrt(x) over a
// nonnegative inside range r.
func SqrtRangeConservative(r Range) (Range, error) {
	if !r.IsInside() {
		return Range{}, fmt.Errorf("SqrtRangeConservative: requires inside range; got %v", r)
	}
	if r.Lo.Cmp(intRat(0)) < 0 {
		return Range{}, fmt.Errorf("SqrtRangeConservative: negative range %v", r)
	}

	lo, err := SqrtLowerBoundRational(r.Lo)
	if err != nil {
		return Range{}, err
	}
	hi, err := SqrtUpperBoundRational(r.Hi)
	if err != nil {
		return Range{}, err
	}

	return NewRange(lo, hi, r.IncLo, r.IncHi), nil
}

// sqrt_range_conservative.go v4
