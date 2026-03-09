// gcf_eval.go v1
package cf

import (
	"fmt"
	"math/big"
)

// EvaluateFiniteGCF computes the exact rational value of a finite generalized
// continued fraction using the convention:
//
//	x = p + q/x'
//
// For a finite source, the terminal convention is:
//
//	the final emitted term contributes just p_last
//
// Equivalently, if the finite prefix is composed into a ULFT T, this returns
// T(∞), which for T(x) = (A*x + B)/(C*x + D) is A/C, provided C != 0.
//
// This implementation composes terms forward and does not buffer the whole
// input stream.
func EvaluateFiniteGCF(src GCFSource) (Rational, error) {
	// Start from identity transform.
	// T(x) = x  =>  [[1 0],[0 1]]
	A := big.NewInt(1)
	B := big.NewInt(0)
	C := big.NewInt(0)
	D := big.NewInt(1)

	seenAny := false

	for {
		p, q, ok := src.NextPQ()
		if !ok {
			break
		}
		seenAny = true

		if q <= 0 {
			return Rational{}, fmt.Errorf("EvaluateFiniteGCF: require q>0, got q=%d", q)
		}

		P := big.NewInt(p)
		Q := big.NewInt(q)

		// Ingest one term x = p + q/x' = (p*x' + q)/x'
		//
		// If current T(x) = (A*x + B)/(C*x + D), then after substitution:
		//
		// T((p*x' + q)/x') = ((A*p + B)*x' + A*q) / ((C*p + D)*x' + C*q)
		Ap := new(big.Int).Mul(A, P)
		A2 := new(big.Int).Add(Ap, B)

		B2 := new(big.Int).Mul(A, Q)

		Cp := new(big.Int).Mul(C, P)
		C2 := new(big.Int).Add(Cp, D)

		D2 := new(big.Int).Mul(C, Q)

		A, B, C, D = A2, B2, C2, D2
	}

	if !seenAny {
		return Rational{}, fmt.Errorf("EvaluateFiniteGCF: empty source")
	}

	// Finite-tail convention => evaluate T(∞) = A/C.
	if C.Sign() == 0 {
		return Rational{}, fmt.Errorf("EvaluateFiniteGCF: invalid finite termination (C=0)")
	}

	return newRationalBig(new(big.Int).Set(A), new(big.Int).Set(C))
}

// gcf_eval.go v1
