// blft_range.go v1
package cf

import "fmt"

// ApplyBLFTRange maps the rectangle rx × ry through the BLFT and returns a
// conservative bounding interval.
//
// This first implementation is deliberately conservative:
//
//  1. Requires rx and ry to be inside (non-empty).
//  2. Computes the raw denominator D(x,y)=Exy+Fx+Gy+H at the 4 corners.
//     If any corner has D=0, or corner signs are not all equal, returns error.
//  3. Evaluates the BLFT at the 4 corners and returns [min, max].
func (t BLFT) ApplyBLFTRange(rx, ry Range) (Range, error) {
	if !rx.IsInside() || !ry.IsInside() {
		return Range{}, fmt.Errorf("ApplyBLFTRange on outside input range: rx=[%v,%v] ry=[%v,%v]", rx.Lo, rx.Hi, ry.Lo, ry.Hi)
	}

	xs := []Rational{rx.Lo, rx.Hi}
	ys := []Rational{ry.Lo, ry.Hi}

	// Denominator sign check at corners (raw, pre-normalization).
	var sgn int // -1 or +1 once set

	for _, x := range xs {
		for _, y := range ys {
			d, err := t.denRaw(x, y)
			if err != nil {
				return Range{}, err
			}
			if d == 0 {
				return Range{}, fmt.Errorf("BLFT denominator is 0 at corner (x=%v,y=%v)", x, y)
			}
			cur := 1
			if d < 0 {
				cur = -1
			}
			if sgn == 0 {
				sgn = cur
			} else if cur != sgn {
				return Range{}, fmt.Errorf("BLFT denominator sign differs across corners (likely crosses 0) for rx=[%v,%v] ry=[%v,%v]", rx.Lo, rx.Hi, ry.Lo, ry.Hi)
			}
		}
	}

	// Evaluate corners.
	var zmin, zmax Rational
	first := true

	for _, x := range xs {
		for _, y := range ys {
			z, err := t.ApplyRat(x, y)
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
	}

	return Range{Lo: zmin, Hi: zmax}, nil
}

// denRaw computes the BLFT denominator value (pre-normalization) for rationals x=p/q, y=u/v:
//
// D = E*p*u + F*p*v + G*u*q + H*q*v
//
// Uses checked ops and returns ErrOverflow on overflow.
func (t BLFT) denRaw(x, y Rational) (int64, error) {
	p, q := x.P, x.Q
	u, v := y.P, y.Q

	epu, ok := mul64(t.E, p)
	if !ok {
		return 0, ErrOverflow
	}
	epu, ok = mul64(epu, u)
	if !ok {
		return 0, ErrOverflow
	}

	fpv, ok := mul64(t.F, p)
	if !ok {
		return 0, ErrOverflow
	}
	fpv, ok = mul64(fpv, v)
	if !ok {
		return 0, ErrOverflow
	}

	guq, ok := mul64(t.G, u)
	if !ok {
		return 0, ErrOverflow
	}
	guq, ok = mul64(guq, q)
	if !ok {
		return 0, ErrOverflow
	}

	hqv, ok := mul64(t.H, q)
	if !ok {
		return 0, ErrOverflow
	}
	hqv, ok = mul64(hqv, v)
	if !ok {
		return 0, ErrOverflow
	}

	d, ok := add64(epu, fpv)
	if !ok {
		return 0, ErrOverflow
	}
	d, ok = add64(d, guq)
	if !ok {
		return 0, ErrOverflow
	}
	d, ok = add64(d, hqv)
	if !ok {
		return 0, ErrOverflow
	}

	return d, nil
}

// blft_range.go v1
