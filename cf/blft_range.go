// blft_range.go v3
package cf

import "fmt"

// ApplyBLFTRange maps the rectangle rx × ry through the BLFT and returns a
// conservative bounding interval.
//
// v3: restores a proper denom-crossing guard:
//
//  1. Evaluate denominator D(x,y)=Exy+Fx+Gy+H at the 4 corners.
//  2. Form denomRange=[min,max]. If denomRange.ContainsZero() => reject (pole hazard).
//  3. Otherwise evaluate z at the 4 corners and return [min,max].
//
// This is conservative: it may reject some safe rectangles, but it must not
// accept rectangles where a pole is possible.
func (t BLFT) ApplyBLFTRange(rx, ry Range) (Range, error) {
	if !rx.IsInside() || !ry.IsInside() {
		return Range{}, fmt.Errorf("ApplyBLFTRange requires inside ranges: rx=[%v,%v] ry=[%v,%v]", rx.Lo, rx.Hi, ry.Lo, ry.Hi)
	}

	xs := []Rational{rx.Lo, rx.Hi}
	ys := []Rational{ry.Lo, ry.Hi}

	// 1) Denominator enclosure via corners.
	dmin, dmax, err := t.denomCornerBounds(xs, ys)
	if err != nil {
		return Range{}, err
	}
	denRange := Range{Lo: dmin, Hi: dmax} // guaranteed inside by construction
	if denRange.ContainsZero() {
		return Range{}, fmt.Errorf(
			"BLFT denominator may cross 0 for rx=[%v,%v] ry=[%v,%v] (den in [%v,%v])",
			rx.Lo, rx.Hi, ry.Lo, ry.Hi, denRange.Lo, denRange.Hi,
		)
	}

	// 2) Output enclosure via corners.
	var zmin, zmax Rational
	first := true
	for _, x := range xs {
		for _, y := range ys {
			z, err := t.ApplyRat(x, y)
			if err != nil {
				// Should be rare once denomRange excludes 0, but propagate if it happens.
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
//
// Uses Rational arithmetic so it stays consistent with ApplyRat semantics.
// Any overflow handling is delegated to Rational ops (checked int64 in your codebase).
func (t BLFT) denomAt(x, y Rational) (Rational, error) {
	xy, err := x.Mul(y)
	if err != nil {
		return Rational{}, err
	}

	e := Rational{P: t.E, Q: 1}
	f := Rational{P: t.F, Q: 1}
	g := Rational{P: t.G, Q: 1}
	h := Rational{P: t.H, Q: 1}

	term1, err := e.Mul(xy)
	if err != nil {
		return Rational{}, err
	}
	term2, err := f.Mul(x)
	if err != nil {
		return Rational{}, err
	}
	term3, err := g.Mul(y)
	if err != nil {
		return Rational{}, err
	}

	s, err := term1.Add(term2)
	if err != nil {
		return Rational{}, err
	}
	s, err = s.Add(term3)
	if err != nil {
		return Rational{}, err
	}
	s, err = s.Add(h)
	if err != nil {
		return Rational{}, err
	}
	return s, nil
}

// blft_range.go v3
