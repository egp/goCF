// pbt_blf_denom_guard_test.go v2
package cf

import (
	"testing"

	"pgregory.net/rapid"
)

func genTinyI64() *rapid.Generator[int64]    { return rapid.Int64Range(-9, 9) }
func genTinyPosI64() *rapid.Generator[int64] { return rapid.Int64Range(1, 9) }

func genTinyRat() *rapid.Generator[Rational] {
	return rapid.Custom(func(t *rapid.T) Rational {
		p := genTinyI64().Draw(t, "p")
		q := genTinyPosI64().Draw(t, "q")
		return mustRat(p, q)
	})
}

func genInsideRangeTiny() *rapid.Generator[Range] {
	return rapid.Custom(func(t *rapid.T) Range {
		a := genTinyRat().Draw(t, "a")
		b := genTinyRat().Draw(t, "b")
		if a.Cmp(b) <= 0 {
			return Range{Lo: a, Hi: b}
		}
		return Range{Lo: b, Hi: a}
	})
}

func denomAtNaive(tform BLFT, x, y Rational) (Rational, error) {
	xy, err := x.Mul(y)
	if err != nil {
		return Rational{}, err
	}
	e := Rational{P: tform.E, Q: 1}
	f := Rational{P: tform.F, Q: 1}
	g := Rational{P: tform.G, Q: 1}
	h := Rational{P: tform.H, Q: 1}

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

func TestPBT_BLF_DenomCornerBoundsExcludeZeroImpliesNoPoleSamples(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		E := genTinyI64().Draw(t, "E")
		F := genTinyI64().Draw(t, "F")
		G := genTinyI64().Draw(t, "G")
		H := genTinyI64().Draw(t, "H")

		// numerator irrelevant for denom guard
		tform := NewBLFT(0, 0, 0, 1, E, F, G, H)

		rx := genInsideRangeTiny().Draw(t, "rx")
		ry := genInsideRangeTiny().Draw(t, "ry")

		xs := []Rational{rx.Lo, rx.Hi}
		ys := []Rational{ry.Lo, ry.Hi}

		var dmin, dmax Rational
		first := true
		for _, x := range xs {
			for _, y := range ys {
				d, err := denomAtNaive(tform, x, y)
				if err != nil {
					return // overflow etc: ignore this trial
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

		denRange := Range{Lo: dmin, Hi: dmax}
		if denRange.ContainsZero() {
			return // pole hazard cases: allowed to reject
		}

		samplesX := []Rational{rx.Lo, rx.Hi, genTinyRat().Draw(t, "x1"), genTinyRat().Draw(t, "x2")}
		samplesY := []Rational{ry.Lo, ry.Hi, genTinyRat().Draw(t, "y1"), genTinyRat().Draw(t, "y2")}

		for _, x := range samplesX {
			if !rx.Contains(x) {
				continue
			}
			for _, y := range samplesY {
				if !ry.Contains(y) {
					continue
				}
				d, err := denomAtNaive(tform, x, y)
				if err != nil {
					t.Fatalf("denomAt failed: %v", err)
				}
				if d.P == 0 {
					t.Fatalf("pole found despite denomRange excluding 0: E=%d F=%d G=%d H=%d rx=[%v,%v] ry=[%v,%v] at x=%v y=%v",
						E, F, G, H, rx.Lo, rx.Hi, ry.Lo, ry.Hi, x, y)
				}

				bounds := Range{Lo: dmin, Hi: dmax} // <— FIX: composite literal assigned, not inline in if
				if !bounds.Contains(d) {
					t.Fatalf("denom out of corner bounds: den=%v not in [%v,%v]", d, dmin, dmax)
				}
			}
		}
	})
}

// pbt_blf_denom_guard_test.go v2
