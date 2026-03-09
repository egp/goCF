// lft.go v4
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

// IngestGCF rewrites the ULFT after ingesting one generalized continued-fraction
// term into x, using the convention:
//
//	x = p + q/x' = (p*x' + q)/x'
//
// Therefore, if T(x) = (A*x + B)/(C*x + D), the rewritten transform in x' is:
//
//	T((p*x' + q)/x') = ((A*p + B)*x' + A*q) / ((C*p + D)*x' + C*q)
//
// Preconditions:
//   - q > 0
func (t ULFT) IngestGCF(p, q int64) (ULFT, error) {
	if q <= 0 {
		return ULFT{}, fmt.Errorf("ULFT IngestGCF: require q>0, got q=%d", q)
	}
	if err := t.Validate(); err != nil {
		return ULFT{}, err
	}

	P := big.NewInt(p)
	Q := big.NewInt(q)

	// A' = A*p + B
	Ap := new(big.Int).Mul(t.A, P)
	A2 := new(big.Int).Add(Ap, t.B)

	// B' = A*q
	B2 := new(big.Int).Mul(t.A, Q)

	// C' = C*p + D
	Cp := new(big.Int).Mul(t.C, P)
	C2 := new(big.Int).Add(Cp, t.D)

	// D' = C*q
	D2 := new(big.Int).Mul(t.C, Q)

	out := NewULFT(A2, B2, C2, D2)
	if err := out.Validate(); err != nil {
		return ULFT{}, err
	}
	return out, nil
}

// lft.go v4
