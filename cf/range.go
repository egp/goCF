// range.go v1
package cf

import "fmt"

// Range represents a closed interval [Lo, Hi] over exact rationals.
//
// Invariant: Lo <= Hi (by Rational.Cmp).
type Range struct {
	Lo Rational
	Hi Rational
}

// NewRange constructs a Range and enforces Lo <= Hi.
// If lo > hi, returns an error.
func NewRange(lo, hi Rational) (Range, error) {
	if lo.Cmp(hi) > 0 {
		return Range{}, fmt.Errorf("invalid range: lo=%v > hi=%v", lo, hi)
	}
	return Range{Lo: lo, Hi: hi}, nil
}

// MustRange is a convenience for tests; it panics on invalid input.
func MustRange(lo, hi Rational) Range {
	r, err := NewRange(lo, hi)
	if err != nil {
		panic(err)
	}
	return r
}

// Width returns Hi - Lo.
func (r Range) Width() (Rational, error) {
	return r.Hi.Sub(r.Lo)
}

// Contains reports whether x is inside the closed interval [Lo, Hi].
func (r Range) Contains(x Rational) bool {
	return r.Lo.Cmp(x) <= 0 && x.Cmp(r.Hi) <= 0
}

// ApplyULFT returns the image of the interval under a ULFT:
//
//	f(x) = (A x + B) / (C x + D)
//
// Correct handling requires detecting whether the function is monotone
// over the interval (i.e., denominator does not cross 0). If the denominator
// crosses 0 within [Lo, Hi], the image is not a single interval (it is
// unbounded / split), and we return an error.
//
// When monotone, the image is simply [min(f(Lo), f(Hi)), max(f(Lo), f(Hi))].
func (r Range) ApplyULFT(t ULFT) (Range, error) {
	// Denominator at endpoints: Cx + D
	denLo := t.C*r.Lo.P + t.D*r.Lo.Q
	denHi := t.C*r.Hi.P + t.D*r.Hi.Q

	// If either endpoint maps to undefined, or denominator crosses 0, reject.
	if denLo == 0 || denHi == 0 || (denLo < 0) != (denHi < 0) {
		return Range{}, fmt.Errorf("ULFT denominator crosses 0 on range [%v,%v]", r.Lo, r.Hi)
	}

	fLo, err := t.ApplyRat(r.Lo)
	if err != nil {
		return Range{}, err
	}
	fHi, err := t.ApplyRat(r.Hi)
	if err != nil {
		return Range{}, err
	}

	if fLo.Cmp(fHi) <= 0 {
		return Range{Lo: fLo, Hi: fHi}, nil
	}
	return Range{Lo: fHi, Hi: fLo}, nil
}

// FloorBounds returns floor(Lo), floor(Hi) under the library's floor convention.
// If both are equal, callers can safely emit that digit in Gosper-style digit
// extraction logic.
//
// Note: Uses floor division for negatives (consistent with RationalCF).
func (r Range) FloorBounds() (int64, int64) {
	return floorRat(r.Lo), floorRat(r.Hi)
}

func floorRat(x Rational) int64 {
	// floor(p/q) with q>0
	a, _ := floorDivMod(x.P, x.Q)
	return a
}

// range.go v1
