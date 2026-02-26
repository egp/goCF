// emit.go v2
package cf

import "fmt"

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
	lo, hi := img.FloorBounds()
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
//
// This version is checked: returns ErrOverflow on int64 overflow.
func EmitDigit(t ULFT, a int64) (ULFT, error) {
	if err := t.Validate(); err != nil {
		return ULFT{}, err
	}

	ac, ok := mul64(a, t.C)
	if !ok {
		return ULFT{}, ErrOverflow
	}
	ad, ok := mul64(a, t.D)
	if !ok {
		return ULFT{}, ErrOverflow
	}

	cNew, ok := sub64(t.A, ac)
	if !ok {
		return ULFT{}, ErrOverflow
	}
	dNew, ok := sub64(t.B, ad)
	if !ok {
		return ULFT{}, ErrOverflow
	}

	out := ULFT{A: t.C, B: t.D, C: cNew, D: dNew}
	if err := out.Validate(); err != nil {
		return ULFT{}, err
	}
	return out, nil
}

// VerifyULFTDeterminant is a small helper for debugging invariants.
// det = A*D - B*C (checked).
func (t ULFT) Determinant() (int64, error) {
	ad, ok := mul64(t.A, t.D)
	if !ok {
		return 0, ErrOverflow
	}
	bc, ok := mul64(t.B, t.C)
	if !ok {
		return 0, ErrOverflow
	}
	det, ok := sub64(ad, bc)
	if !ok {
		return 0, ErrOverflow
	}
	return det, nil
}

func (t ULFT) Validate() error {
	// We allow many ULFTs, but disallow the all-zero denominator row.
	if t.C == 0 && t.D == 0 {
		return fmt.Errorf("invalid ULFT: denominator is identically 0: %v", t)
	}
	return nil
}

// emit.go v2
