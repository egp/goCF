// blft_range.go v5
package cf

import "fmt"

// ApplyBLFTRange maps the rectangle rx × ry through the BLFT and returns a
// conservative bounding interval.
//
// v5: enclosure is CLOSED (IncLo=IncHi=true) because it is a conservative bound.
// denom-crossing guard delegated to DenomMayHitZero(rx,ry).
func (t BLFT) ApplyBLFTRange(rx, ry Range) (Range, error) {
	if !rx.IsInside() || !ry.IsInside() {
		return Range{}, fmt.Errorf("ApplyBLFTRange requires inside ranges: rx=[%v,%v] ry=[%v,%v]", rx.Lo, rx.Hi, ry.Lo, ry.Hi)
	}

	mayHit, err := t.DenomMayHitZero(rx, ry)
	if err != nil {
		return Range{}, err
	}
	if mayHit {
		dr, derr := t.DenomRange(rx, ry)
		if derr != nil {
			return Range{}, derr
		}
		return Range{}, fmt.Errorf(
			"BLFT denominator may cross 0 for rx=[%v,%v] ry=[%v,%v] (den in [%v,%v])",
			rx.Lo, rx.Hi, ry.Lo, ry.Hi, dr.Lo, dr.Hi,
		)
	}

	xs := []Rational{rx.Lo, rx.Hi}
	ys := []Rational{ry.Lo, ry.Hi}

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

	// Conservative enclosure must be CLOSED.
	return NewRange(zmin, zmax, true, true), nil
}

// blft_range.go v5
