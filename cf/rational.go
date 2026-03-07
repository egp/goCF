// rational.go v3
package cf

import (
	"fmt"
	"math/big"
)

// Rational is an exact rational number with arbitrary precision.
//
// Representation invariant (enforced by constructors / ops here):
//   - denominator > 0
//   - fraction is reduced (gcd(num,den)=1) unless value is 0 (0/1)
type Rational struct {
	r big.Rat
}

// NewRational constructs p/q exactly (q != 0) and normalizes to reduced form with q>0.
func NewRational(p, q int64) (Rational, error) {
	if q == 0 {
		return Rational{}, fmt.Errorf("zero denominator")
	}
	var rr big.Rat
	rr.SetFrac(big.NewInt(p), big.NewInt(q)) // normalizes sign + reduces
	return Rational{r: rr}, nil
}

// newRationalBig constructs num/den exactly (den != 0) and normalizes.
func newRationalBig(num, den *big.Int) (Rational, error) {
	if den.Sign() == 0 {
		return Rational{}, fmt.Errorf("zero denominator")
	}
	var rr big.Rat
	rr.SetFrac(num, den) // normalizes sign + reduces
	return Rational{r: rr}, nil
}

func (r Rational) String() string {
	return r.r.RatString()
}

func (r Rational) Add(o Rational) (Rational, error) {
	var out big.Rat
	out.Add(&r.r, &o.r)
	return Rational{r: out}, nil
}

func (r Rational) Sub(o Rational) (Rational, error) {
	var out big.Rat
	out.Sub(&r.r, &o.r)
	return Rational{r: out}, nil
}

func (r Rational) Mul(o Rational) (Rational, error) {
	var out big.Rat
	out.Mul(&r.r, &o.r)
	return Rational{r: out}, nil
}

func (r Rational) Div(o Rational) (Rational, error) {
	if o.r.Sign() == 0 {
		return Rational{}, fmt.Errorf("division by zero")
	}
	var out big.Rat
	out.Quo(&r.r, &o.r)
	return Rational{r: out}, nil
}

func (r Rational) Cmp(o Rational) int {
	return r.r.Cmp(&o.r)
}

// ---- internal helpers used by other packages/files ----

// ratNumDen returns (num, den) pointers to *copies* (safe to mutate).
func (r Rational) ratNumDen() (num *big.Int, den *big.Int) {
	// big.Rat exposes Num()/Den() as pointers to internal big.Int values.
	// Copy them to avoid aliasing surprises.
	return new(big.Int).Set(r.r.Num()), new(big.Int).Set(r.r.Denom())
}

// rationalFromInt64 returns v/1.
// func rationalFromInt64(v int64) Rational {
// 	var rr big.Rat
// 	rr.SetInt64(v)
// 	return Rational{r: rr}
// }

// rationalZero returns 0/1.
// func rationalZero() Rational { return rationalFromInt64(0) }

// // rationalOne returns 1/1.
// func rationalOne() Rational { return rationalFromInt64(1) }

// // rationalNegOne returns -1/1.
// func rationalNegOne() Rational { return rationalFromInt64(-1) }

// rationalAbs returns |r|.
// func rationalAbs(r Rational) Rational {
// 	if r.r.Sign() >= 0 {
// 		return r
// 	}
// 	var out big.Rat
// 	out.Neg(&r.r)
// 	return Rational{r: out}
// }

// rationalIsInteger reports whether r is an integer (den==1).
// func rationalIsInteger(r Rational) bool {
// 	return r.r.IsInt()
// }

// rationalInt64Exact returns (value,true) if r is an exact int64 integer.
// func rationalInt64Exact(r Rational) (int64, bool) {
// 	if !r.r.IsInt() {
// 		return 0, false
// 	}
// 	num := r.r.Num() // already reduced, denom=1
// 	if !num.IsInt64() {
// 		return 0, false
// 	}
// 	return num.Int64(), true
// }

// rationalMin / Max
func minRat(a, b Rational) Rational {
	if a.Cmp(b) <= 0 {
		return a
	}
	return b
}
func maxRat(a, b Rational) Rational {
	if a.Cmp(b) >= 0 {
		return a
	}
	return b
}

// rational.go v3
