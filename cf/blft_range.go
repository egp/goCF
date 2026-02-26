// blft_range.go v3
package cf

import "fmt"

// ApplyBLFTRange maps the rectangle rx × ry through the BLFT and returns a
// conservative bounding interval.
//
// v3:
//   - Computes denominator values at 4 corners and bounds denom by [min,max].
//   - If denom range contains 0 => reject (pole may exist inside rectangle).
//   - Otherwise evaluates z at corners and returns [min,max].
func (t BLFT) ApplyBLFTRange(rx, ry Range) (Range, error) {
	if !rx.IsInside() || !ry.IsInside() {
		return Range{}, fmt.Errorf("ApplyBLFTRange requires inside ranges: rx=[%v,%v] ry=[%v,%v]", rx.Lo, rx.Hi, ry.Lo, ry.Hi)
	}

	xs := []Rational{rx.Lo, rx.Hi}
	ys := []Rational{ry.Lo, ry.Hi}

	// 1) Denominator enclosure via corners
	dmin, dmax, err := t.denomCornerBounds(xs, ys)
	if err != nil {
		return Range{}, err
	}
	denRange := Range{Lo: dmin, Hi: dmax}
	if denRange.ContainsZero() {
		return Range{}, fmt.Errorf("BLFT denominator may cross 0 for rx=[%v,%v] ry=[%v,%v] (den in [%v,%v])",
			rx.Lo, rx.Hi, ry.Lo, ry.Hi, denRange.Lo, denRange.Hi)
	}

	// 2) Evaluate corners for output enclosure
	var zmin, zmax Rational
	first := true
	for _, x := range xs {
		for _, y := range ys {
			z, err := t.ApplyRat(x, y)
			if err != nil {
				// With denomRange excluding 0, corner ApplyRat errors should be rare;
				// still propagate.
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
	}
	return Range{Lo: zmin, Hi: zmax}, nil
}

// denomCornerBounds evaluates denom D(x,y)=Exy+Fx+Gy+H at the 4 corners,
// returning min and max as rationals.
func (t BLFT) denomCornerBounds(xs, ys []Rational) (Rational, Rational, error) {
	var dmin, dmax Rational
	first := true

	for _, x := range xs {
		for _, y := range ys {
			d, err := t.denomAt(x, y)
			if err != nil {
				return Rational{}, Rational{}, err
			}
			if first {
				dmin, dmax = d, d
				first = false
				continue
			}
			if d.Cmp(dmin) < 0 {
				dmin = d
			}
			if d.Cmp(dmax) > 0 {
				dmax = d
			}
		}
	}
	return dmin, dmax, nil
}

// denomAt returns the exact rational value of D(x,y)=Exy+Fx+Gy+H.
func (t BLFT) denomAt(x, y Rational) (Rational, error) {
	xy, err := x.Mul(y)
	if err != nil {
		return Rational{}, err
	}

	term1, err := (Rational{P: t.E, Q: 1}).Mul(xy)
	if err != nil {
		return Rational{}, err
	}
	term2, err := (Rational{P: t.F, Q: 1}).Mul(x)
	if err != nil {
		return Rational{}, err
	}
	term3, err := (Rational{P: t.G, Q: 1}).Mul(y)
	if err != nil {
		return Rational{}, err
	}
	term4 := Rational{P: t.H, Q: 1}

	s, err := term1.Add(term2)
	if err != nil {
		return Rational{}, err
	}
	s, err = s.Add(term3)
	if err != nil {
		return Rational{}, err
	}
	s, err = s.Add(term4)
	if err != nil {
		return Rational{}, err
	}
	return s, nil
}

// blft_range.go v3
