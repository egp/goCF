// diag_blft.go v2
package cf

import (
	"fmt"
	"math/big"
)

// DiagBLFT is the diagonal specialization of a BLFT on T(x,x):
//
//	(A*x^2 + B*x + C) / (D*x^2 + E*x + F)
type DiagBLFT struct {
	A, B, C *big.Int
	D, E, F *big.Int
}

func NewDiagBLFT(a, b, c, d, e, f *big.Int) DiagBLFT {
	return DiagBLFT{
		A: new(big.Int).Set(a),
		B: new(big.Int).Set(b),
		C: new(big.Int).Set(c),
		D: new(big.Int).Set(d),
		E: new(big.Int).Set(e),
		F: new(big.Int).Set(f),
	}
}

// func bigIntAdd(x, y *big.Int) *big.Int {
// 	return new(big.Int).Add(x, y)
// }

// DiagFromBLFT specializes a BLFT to the diagonal x=y.
//
// BLFT:
//
//	(Axy + Bx + Cy + D) / (Exy + Fx + Gy + H)
//
// on x=y becomes:
//
//	(A*x^2 + (B+C)*x + D) / (E*x^2 + (F+G)*x + H)
func DiagFromBLFT(t BLFT) DiagBLFT {
	return NewDiagBLFT(
		big.NewInt(t.A),
		big.NewInt(t.B+t.C),
		big.NewInt(t.D),
		big.NewInt(t.E),
		big.NewInt(t.F+t.G),
		big.NewInt(t.H),
	)
}

func (t DiagBLFT) String() string {
	return fmt.Sprintf(
		"(%s*x^2 + %s*x + %s) / (%s*x^2 + %s*x + %s)",
		t.A.String(), t.B.String(), t.C.String(),
		t.D.String(), t.E.String(), t.F.String(),
	)
}

// ApplyRat evaluates the diagonal transform exactly on a rational x.
func (t DiagBLFT) ApplyRat(x Rational) (Rational, error) {
	var x2 big.Rat
	x2.Mul(&x.r, &x.r)

	var num, den, term big.Rat
	num.SetInt64(0)
	den.SetInt64(0)

	term.Mul(ratFromBigInt(t.A), &x2)
	num.Add(&num, &term)

	term.Mul(ratFromBigInt(t.B), &x.r)
	num.Add(&num, &term)

	num.Add(&num, ratFromBigInt(t.C))

	term.Mul(ratFromBigInt(t.D), &x2)
	den.Add(&den, &term)

	term.Mul(ratFromBigInt(t.E), &x.r)
	den.Add(&den, &term)

	den.Add(&den, ratFromBigInt(t.F))

	if den.Sign() == 0 {
		return Rational{}, fmt.Errorf("DiagBLFT ApplyRat: denominator is zero at x=%v", x)
	}

	var out big.Rat
	out.Quo(&num, &den)
	return Rational{r: out}, nil
}

// ApplyRange maps an inside range through the diagonal transform.
//
// Current minimal support:
//   - exact point ranges are always supported
//   - non-point ranges are supported only when the denominator is constant nonzero
//
// For the constant-denominator case, extrema of the quadratic numerator occur at:
//   - endpoints
//   - the vertex x = -B/(2A), if A != 0 and the vertex lies inside the range
func (t DiagBLFT) ApplyRange(r Range) (Range, error) {
	if !r.IsInside() {
		return Range{}, fmt.Errorf("DiagBLFT ApplyRange requires inside range; got %v", r)
	}

	// Exact-point fast path always works.
	if r.Lo.Cmp(r.Hi) == 0 {
		z, err := t.ApplyRat(r.Lo)
		if err != nil {
			return Range{}, err
		}
		return NewRange(z, z, true, true), nil
	}

	// Minimal interval support for now: denominator must be constant nonzero.
	if t.D.Sign() != 0 || t.E.Sign() != 0 {
		return Range{}, fmt.Errorf("DiagBLFT ApplyRange: non-constant denominator not yet supported")
	}
	if t.F.Sign() == 0 {
		return Range{}, fmt.Errorf("DiagBLFT ApplyRange: denominator is identically zero")
	}

	samples := []Rational{r.Lo, r.Hi}

	// If numerator is quadratic, include the vertex when it lies inside the interval.
	if t.A.Sign() != 0 {
		twoA := new(big.Int).Mul(big.NewInt(2), t.A)
		var negB big.Int
		negB.Neg(t.B)
		vx, err := newRationalBig(&negB, twoA)
		if err != nil {
			return Range{}, err
		}
		if r.Contains(vx) {
			samples = append(samples, vx)
		}
	}

	var zmin, zmax Rational
	first := true
	for _, x := range samples {
		z, err := t.ApplyRat(x)
		if err != nil {
			return Range{}, err
		}
		if first {
			zmin, zmax = z, z
			first = false
			continue
		}
		if z.Cmp(zmin) < 0 {
			zmin = z
		}
		if z.Cmp(zmax) > 0 {
			zmax = z
		}
	}

	return NewRange(zmin, zmax, true, true), nil
}

// emitDigitDiag updates the diagonal transform after emitting digit d.
//
// If z = N/D, then z' = 1/(z-d) = D/(N-dD).
func (t DiagBLFT) emitDigitDiag(d int64) (DiagBLFT, error) {
	di := big.NewInt(d)

	A2 := new(big.Int).Set(t.D)
	B2 := new(big.Int).Set(t.E)
	C2 := new(big.Int).Set(t.F)

	dD := new(big.Int).Mul(di, t.D)
	dE := new(big.Int).Mul(di, t.E)
	dF := new(big.Int).Mul(di, t.F)

	D2 := new(big.Int).Sub(t.A, dD)
	E2 := new(big.Int).Sub(t.B, dE)
	F2 := new(big.Int).Sub(t.C, dF)

	return NewDiagBLFT(A2, B2, C2, D2, E2, F2), nil
}

// diag_blft.go v2
