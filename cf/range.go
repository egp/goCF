// range.go v12
package cf

import "fmt"

type Range struct {
	Lo    Rational
	Hi    Rational
	IncLo bool
	IncHi bool
}

func NewRange(lo, hi Rational, incLo, incHi bool) Range {
	return Range{Lo: lo, Hi: hi, IncLo: incLo, IncHi: incHi}
}

func (r Range) IsInside() bool  { return r.Lo.Cmp(r.Hi) <= 0 }
func (r Range) IsOutside() bool { return r.Lo.Cmp(r.Hi) > 0 }

func (r Range) String() string {
	kind := "inside"
	if r.IsOutside() {
		kind = "outside"
	}
	return fmt.Sprintf("Range[%v,%v]{incLo=%t,incHi=%t,%s}", r.Lo, r.Hi, r.IncLo, r.IncHi, kind)
}

func (r Range) Contains(x Rational) bool {
	cLo := x.Cmp(r.Lo) // x ? Lo
	cHi := x.Cmp(r.Hi) // x ? Hi

	if r.IsInside() {
		// Inside: Lo <= x <= Hi with endpoint inclusions.
		if cLo < 0 {
			return false
		}
		if cLo == 0 && !r.IncLo {
			return false
		}
		if cHi > 0 {
			return false
		}
		if cHi == 0 && !r.IncHi {
			return false
		}
		return true
	}

	// Outside: (-∞,Hi] ∪ [Lo,+∞)
	// i.e. x is contained if x <= Hi OR x >= Lo (with endpoint inclusions).
	if cHi < 0 {
		// x < Hi => in (-∞,Hi)
		return true
	}
	if cHi == 0 && r.IncHi {
		// x == Hi and included
		return true
	}
	if cLo > 0 {
		// x > Lo => in (Lo,∞)
		return true
	}
	if cLo == 0 && r.IncLo {
		// x == Lo and included
		return true
	}
	return false
}

// ContainsZero() ≡ Contains(0)
func (r Range) ContainsZero() bool {
	return r.Contains(mustRat(0, 1))
}

// Width returns Hi - Lo for inside ranges.
// Outside ranges do not have a finite width.
func (r Range) Width() (Rational, error) {
	if r.IsOutside() {
		return Rational{}, fmt.Errorf("Width undefined for outside range: %v", r)
	}
	return r.Hi.Sub(r.Lo)
}

// FloorBounds returns floor(Lo), floor(Hi) for inside ranges.
// Outside ranges are not a single contiguous interval and are rejected here.
func (r Range) FloorBounds() (int64, int64, error) {
	if r.IsOutside() {
		return 0, 0, fmt.Errorf("FloorBounds undefined for outside range: %v", r)
	}
	lo, err := floorRat(r.Lo)
	if err != nil {
		return 0, 0, err
	}
	hi, err := floorRat(r.Hi)
	if err != nil {
		return 0, 0, err
	}
	return lo, hi, nil
}

func floorRat(x Rational) (int64, error) {
	if x.Q == 0 {
		return 0, fmt.Errorf("floorRat: zero denominator")
	}
	p := x.P
	q := x.Q
	if q < 0 {
		p = -p
		q = -q
	}

	quo := p / q
	rem := p % q
	if rem != 0 && p < 0 {
		quo -= 1
	}
	return quo, nil
}

// ApplyULFT maps an inside range through a ULFT and returns a conservative CLOSED enclosure.
func (r Range) ApplyULFT(t ULFT) (Range, error) {
	if r.IsOutside() {
		return Range{}, fmt.Errorf("ApplyULFT requires inside range: %v", r)
	}

	// denom(x) = C*x + D; extrema at endpoints.
	dLo, err := ulftDenomAt(t, r.Lo)
	if err != nil {
		return Range{}, err
	}
	dHi, err := ulftDenomAt(t, r.Hi)
	if err != nil {
		return Range{}, err
	}

	den := NewRange(dLo, dHi, true, true)
	if den.Lo.Cmp(den.Hi) > 0 {
		den = NewRange(dHi, dLo, true, true)
	}
	if den.ContainsZero() {
		return Range{}, fmt.Errorf("ULFT denominator may cross 0 on range %v (den in [%v,%v])", r, den.Lo, den.Hi)
	}

	zLo, err := t.ApplyRat(r.Lo)
	if err != nil {
		return Range{}, err
	}
	zHi, err := t.ApplyRat(r.Hi)
	if err != nil {
		return Range{}, err
	}

	// Conservative enclosure must be CLOSED.
	if zLo.Cmp(zHi) <= 0 {
		return NewRange(zLo, zHi, true, true), nil
	}
	return NewRange(zHi, zLo, true, true), nil
}

func ulftDenomAt(t ULFT, x Rational) (Rational, error) {
	c := mustRat(t.C, 1)
	d := mustRat(t.D, 1)

	cx, err := c.Mul(x)
	if err != nil {
		return Rational{}, err
	}
	return cx.Add(d)
}

// range.go v12
