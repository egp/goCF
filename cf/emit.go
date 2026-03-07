// emit.go v4
package cf

import (
	"fmt"
	"math/big"
)

// SafeDigit attempts to determine a Gosper-safe next continued-fraction digit
// from an interval and a ULFT.
//
// It computes the image interval r' = t(r) and compares floor(r'.Lo) and
// floor(r'.Hi). If they are equal, that integer is safe to emit.
//
// Returns:
//
//	(digit, true, nil)  if safe
//	(_,     false, nil) if not safe (need to refine the source interval)
//	(_,     false, err) if the transform is not interval-safe on r (e.g. denom crosses 0)
func SafeDigit(t ULFT, r Range) (int64, bool, error) {
	img, err := r.ApplyULFT(t)
	if err != nil {
		return 0, false, err
	}
	lo, hi, err := img.FloorBounds()
	if err != nil {
		return 0, false, err
	}
	if lo == hi {
		return lo, true, nil
	}
	return 0, false, nil
}

// EmitDigit updates a ULFT after emitting a continued-fraction digit a.
//
// New transform:
//
//	T' = [[C, D], [A - a*C, B - a*D]]
func EmitDigit(t ULFT, a int64) (ULFT, error) {
	if err := t.Validate(); err != nil {
		return ULFT{}, err
	}

	ai := big.NewInt(a)

	ac := new(big.Int).Mul(ai, t.C)
	ad := new(big.Int).Mul(ai, t.D)

	cNew := new(big.Int).Sub(t.A, ac)
	dNew := new(big.Int).Sub(t.B, ad)

	out := ULFT{
		A: new(big.Int).Set(t.C),
		B: new(big.Int).Set(t.D),
		C: cNew,
		D: dNew,
	}
	if err := out.Validate(); err != nil {
		return ULFT{}, err
	}
	return out, nil
}

// Determinant returns det = A*D - B*C.
//
// Since the public API still returns int64, this returns ErrOverflow if the
// exact determinant does not fit in int64.
func (t ULFT) Determinant() (int64, error) {
	ad := new(big.Int).Mul(t.A, t.D)
	bc := new(big.Int).Mul(t.B, t.C)
	det := new(big.Int).Sub(ad, bc)

	if !det.IsInt64() {
		return 0, ErrOverflow
	}
	return det.Int64(), nil
}

func (t ULFT) Validate() error {
	// We allow many ULFTs, but disallow the all-zero denominator row.
	if t.C.Sign() == 0 && t.D.Sign() == 0 {
		return fmt.Errorf("invalid ULFT: denominator is identically 0: %v", t)
	}
	return nil
}

// emit.go v4
