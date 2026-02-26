// range.go v3
package cf

import "fmt"

// Range represents a closed interval [Lo, Hi] over exact rationals.
//
// Semantics:
//   - Inside (non-empty) range: Lo <= Hi (includes Lo==Hi for exact).
//   - Outside (complement) range: Lo > Hi, representing (-∞, Hi] ∪ [Lo, ∞).
//
// Note: Some operations (Width, FloorBounds, ApplyULFT) are only defined for
// inside ranges because their results are not representable as a single Range
// for outside ranges.
type Range struct {
	Lo Rational
	Hi Rational
}

// NewRange constructs a Range without enforcing ordering.
// Use r.IsInside() to test non-emptiness.
func NewRange(lo, hi Rational) Range {
	return Range{Lo: lo, Hi: hi}
}

// MustRange is a convenience for tests and callers that expect a non-empty range.
// It panics if lo > hi.
func MustRange(lo, hi Rational) Range {
	r := Range{Lo: lo, Hi: hi}
	if !r.IsInside() {
		panic(fmt.Errorf("invalid (outside) range: lo=%v > hi=%v", lo, hi))
	}
	return r
}

// IsInside reports whether the range is non-empty: Lo <= Hi.
func (r Range) IsInside() bool {
	return r.Lo.Cmp(r.Hi) <= 0
}

// IsOutside reports whether the range is a complement range: Lo > Hi.
func (r Range) IsOutside() bool {
	return r.Lo.Cmp(r.Hi) > 0
}

// Contains reports membership under the Range semantics.
//
// Inside:  Lo <= x <= Hi
// Outside: x <= Hi  OR  x >= Lo    (i.e., outside the open interval (Hi,Lo))
func (r Range) Contains(x Rational) bool {
	if r.IsInside() {
		return r.Lo.Cmp(x) <= 0 && x.Cmp(r.Hi) <= 0
	}
	// Outside (Lo > Hi): (-∞,Hi] ∪ [Lo,∞)
	return x.Cmp(r.Hi) <= 0 || x.Cmp(r.Lo) >= 0
}

// ContainsZero reports whether the Range contains 0 under the same semantics.
func (r Range) ContainsZero() bool {
	return r.Contains(mustRat(0, 1))
}

// Width returns Hi - Lo. Outside ranges return an error (unbounded / not representable).
func (r Range) Width() (Rational, error) {
	if !r.IsInside() {
		return Rational{}, fmt.Errorf("width undefined for outside range [%v,%v]", r.Lo, r.Hi)
	}
	return r.Hi.Sub(r.Lo)
}

// ApplyULFT returns the image of the interval under a ULFT:
//
//	f(x) = (A x + B) / (C x + D)
//
// Defined only for inside ranges in this version.
// If denominator crosses 0 within [Lo, Hi], the image is not a single interval.
func (r Range) ApplyULFT(t ULFT) (Range, error) {
	if !r.IsInside() {
		return Range{}, fmt.Errorf("ApplyULFT undefined for outside range [%v,%v]", r.Lo, r.Hi)
	}

	// Denominator at endpoints: Cx + D
	denLo := t.C*r.Lo.P + t.D*r.Lo.Q
	denHi := t.C*r.Hi.P + t.D*r.Hi.Q

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
// Defined only for inside ranges in this version.
func (r Range) FloorBounds() (int64, int64, error) {
	if !r.IsInside() {
		return 0, 0, fmt.Errorf("FloorBounds undefined for outside range [%v,%v]", r.Lo, r.Hi)
	}
	return floorRat(r.Lo), floorRat(r.Hi), nil
}

func floorRat(x Rational) int64 {
	// floor(p/q) with q>0
	a, _ := floorDivMod(x.P, x.Q)
	return a
}

// range.go v3
