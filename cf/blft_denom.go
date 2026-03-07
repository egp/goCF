// blft_denom.go v1
package cf

import "fmt"

// DenomRange returns the exact denominator range over the rectangle rx×ry,
// using the fact that D(x,y)=Exy+Fx+Gy+H is bilinear, so extrema occur at corners.
//
// Preconditions:
//   - rx and ry must be inside (Lo<=Hi)
//
// Postconditions:
//   - returned range is inside (Lo<=Hi)
func (t BLFT) DenomRange(rx, ry Range) (Range, error) {
	if !rx.IsInside() || !ry.IsInside() {
		return Range{}, fmt.Errorf("DenomRange requires inside ranges: rx=[%v,%v] ry=[%v,%v]", rx.Lo, rx.Hi, ry.Lo, ry.Hi)
	}

	xs := []Rational{rx.Lo, rx.Hi}
	ys := []Rational{ry.Lo, ry.Hi}

	var dmin, dmax Rational
	first := true

	for _, x := range xs {
		for _, y := range ys {
			d, err := t.denomAt(x, y)
			if err != nil {
				return Range{}, err
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

	return Range{Lo: dmin, Hi: dmax}, nil
}

// DenomMayHitZero reports whether the denominator may be zero anywhere in rx×ry.
// This is a conservative pole guard used by BLFT range mapping and streaming.
func (t BLFT) DenomMayHitZero(rx, ry Range) (bool, error) {
	dr, err := t.DenomRange(rx, ry)
	if err != nil {
		return false, err
	}
	return dr.ContainsZero(), nil
}

// denomAt returns the exact rational value of D(x,y)=Exy+Fx+Gy+H.
func (t BLFT) denomAt(x, y Rational) (Rational, error) {
	xy, err := x.Mul(y)
	if err != nil {
		return Rational{}, err
	}

	e := intRat(t.E)
	f := intRat(t.F)
	g := intRat(t.G)
	h := intRat(t.H)

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

// blft_denom.go v1
