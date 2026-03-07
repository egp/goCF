// lft.go v3
package cf

import (
	"fmt"
	"math/big"
)

// ULFT: (A x + B) / (C x + D)
type ULFT struct {
	A, B, C, D *big.Int
}

func NewULFT(a, b, c, d *big.Int) ULFT {
	return ULFT{A: a, B: b, C: c, D: d}
}

func ratFromBigInt(z *big.Int) *big.Rat {
	return new(big.Rat).SetFrac(new(big.Int).Set(z), big.NewInt(1))
}

// ApplyRat evaluates the ULFT exactly on a rational x.
//
// (A*x + B) / (C*x + D)
func (t ULFT) ApplyRat(x Rational) (Rational, error) {
	if err := t.Validate(); err != nil {
		return Rational{}, err
	}

	// Work in big.Rat throughout.
	var ax, num, cx, den big.Rat

	ax.Mul(ratFromBigInt(t.A), &x.r)
	num.Add(&ax, ratFromBigInt(t.B))

	cx.Mul(ratFromBigInt(t.C), &x.r)
	den.Add(&cx, ratFromBigInt(t.D))

	if den.Sign() == 0 {
		return Rational{}, fmt.Errorf("ULFT ApplyRat: denominator is zero at x=%v", x)
	}

	var out big.Rat
	out.Quo(&num, &den)
	return Rational{r: out}, nil
}

func (t ULFT) String() string {
	return fmt.Sprintf("[[%s %s],[%s %s]]", t.A.String(), t.B.String(), t.C.String(), t.D.String())
}

// BLFT: (A x y + B x + C y + D) / (E x y + F x + G y + H)
type BLFT struct {
	A, B, C, D int64
	E, F, G, H int64
}

func NewBLFT(a, b, c, d, e, f, g, h int64) BLFT {
	return BLFT{A: a, B: b, C: c, D: d, E: e, F: f, G: g, H: h}
}

// ApplyRat evaluates the BLFT exactly on rationals x and y.
//
// BLFT: (A*x*y + B*x + C*y + D) / (E*x*y + F*x + G*y + H)
func (t BLFT) ApplyRat(x, y Rational) (Rational, error) {
	// Compute xy once.
	var xy big.Rat
	xy.Mul(&x.r, &y.r)

	// ---- Numerator: A*xy + B*x + C*y + D ----
	var num big.Rat
	num.SetInt64(0)

	var term big.Rat

	term.Mul(big.NewRat(t.A, 1), &xy)
	num.Add(&num, &term)

	term.Mul(big.NewRat(t.B, 1), &x.r)
	num.Add(&num, &term)

	term.Mul(big.NewRat(t.C, 1), &y.r)
	num.Add(&num, &term)

	num.Add(&num, big.NewRat(t.D, 1))

	// ---- Denominator: E*xy + F*x + G*y + H ----
	var den big.Rat
	den.SetInt64(0)

	term.Mul(big.NewRat(t.E, 1), &xy)
	den.Add(&den, &term)

	term.Mul(big.NewRat(t.F, 1), &x.r)
	den.Add(&den, &term)

	term.Mul(big.NewRat(t.G, 1), &y.r)
	den.Add(&den, &term)

	den.Add(&den, big.NewRat(t.H, 1))

	if den.Sign() == 0 {
		return Rational{}, fmt.Errorf("BLFT ApplyRat: denominator is zero at x=%v y=%v", x, y)
	}

	var out big.Rat
	out.Quo(&num, &den)
	return Rational{r: out}, nil
}

// func bigIntFitsInt64(z *big.Int) bool {
// 	if z.Sign() == 0 {
// 		return true
// 	}
// 	if z.BitLen() < 63 {
// 		return true
// 	}
// 	if z.BitLen() > 63 {
// 		return false
// 	}
// 	max := big.NewInt(int64(^uint64(0) >> 1))
// 	min := new(big.Int).Neg(new(big.Int).Add(max, big.NewInt(1)))
// 	return z.Cmp(min) >= 0 && z.Cmp(max) <= 0
// }

// lft.go v3
