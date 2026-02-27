// pbt_rapid_test.go v2
package cf

import (
	"testing"

	"pgregory.net/rapid"
)

func genSmallInt64() *rapid.Generator[int64] { return rapid.Int64Range(-20, 20) }
func genPosInt64() *rapid.Generator[int64]   { return rapid.Int64Range(1, 20) }

func genSmallRat() *rapid.Generator[Rational] {
	return rapid.Custom(func(t *rapid.T) Rational {
		p := genSmallInt64().Draw(t, "p")
		q := genPosInt64().Draw(t, "q")
		return mustRat(p, q)
	})
}

func genBool() *rapid.Generator[bool] { return rapid.Bool() }

func genInsideRange() *rapid.Generator[Range] {
	return rapid.Custom(func(t *rapid.T) Range {
		a := genSmallRat().Draw(t, "a")
		b := genSmallRat().Draw(t, "b")
		incLo := genBool().Draw(t, "incLo")
		incHi := genBool().Draw(t, "incHi")

		if a.Cmp(b) <= 0 {
			return NewRange(a, b, incLo, incHi)
		}
		return NewRange(b, a, incLo, incHi)
	})
}

func genOutsideRange() *rapid.Generator[Range] {
	return rapid.Custom(func(t *rapid.T) Range {
		lo := genSmallRat().Draw(t, "lo")
		hi := genSmallRat().Draw(t, "hi")
		incLo := genBool().Draw(t, "incLo")
		incHi := genBool().Draw(t, "incHi")

		// Ensure Lo > Hi.
		if lo.Cmp(hi) <= 0 {
			lo, hi = hi, lo
		}
		if lo.Cmp(hi) == 0 {
			lo = mustRat(lo.P+1, lo.Q)
		}
		if lo.Cmp(hi) <= 0 {
			lo = mustRat(1, 1)
			hi = mustRat(0, 1)
		}
		return NewRange(lo, hi, incLo, incHi)
	})
}

func outsideWantContains(r Range, x Rational) bool {
	// Outside semantics: (-∞,Hi] ∪ [Lo,∞) with open/closed endpoints.
	cHi := x.Cmp(r.Hi)
	if cHi < 0 {
		return true
	}
	if cHi == 0 && r.IncHi {
		return true
	}

	cLo := x.Cmp(r.Lo)
	if cLo > 0 {
		return true
	}
	if cLo == 0 && r.IncLo {
		return true
	}

	return false
}

func TestPBT_RangeOutsideContainsSemantics(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		r := genOutsideRange().Draw(t, "r")
		x := genSmallRat().Draw(t, "x")

		want := outsideWantContains(r, x)
		got := r.Contains(x)
		if got != want {
			t.Fatalf("outside Contains mismatch: r=[%v,%v] incLo=%v incHi=%v x=%v got=%v want=%v",
				r.Lo, r.Hi, r.IncLo, r.IncHi, x, got, want)
		}

		// Unified rule: ContainsZero() == Contains(0)
		z := mustRat(0, 1)
		wantZ := outsideWantContains(r, z)
		gotZ := r.ContainsZero()
		if gotZ != wantZ {
			t.Fatalf("outside ContainsZero mismatch: r=[%v,%v] incLo=%v incHi=%v got=%v want=%v",
				r.Lo, r.Hi, r.IncLo, r.IncHi, gotZ, wantZ)
		}
	})
}

func TestPBT_BLFTRangeRejectsPoleForXMinus1(t *testing.T) {
	// z = 1/(x-1), independent of y
	// denom = x - 1  => E=0 F=1 G=0 H=-1
	tform := NewBLFT(0, 0, 0, 1, 0, 1, 0, -1)

	rapid.Check(t, func(t *rapid.T) {
		a := rapid.Int64Range(1, 10).Draw(t, "a")
		b := rapid.Int64Range(1, 10).Draw(t, "b")
		xlo := mustRat(1-a, 1)
		xhi := mustRat(1+b, 1)
		rx := NewRange(xlo, xhi, true, true)

		ry := genInsideRange().Draw(t, "ry")

		_, err := tform.ApplyBLFTRange(rx, ry)
		if err == nil {
			t.Fatalf("expected pole rejection for rx spanning 1: rx=[%v,%v] ry=[%v,%v]", rx.Lo, rx.Hi, ry.Lo, ry.Hi)
		}
	})
}

func TestPBT_BLFTRangeEnclosesSamplesWhenAccepted(t *testing.T) {
	tform := NewBLFT(
		2, 1, -1, 3,
		1, 2, 1, 5,
	)

	rapid.Check(t, func(t *rapid.T) {
		rx := genInsideRange().Draw(t, "rx")
		ry := genInsideRange().Draw(t, "ry")

		out, err := tform.ApplyBLFTRange(rx, ry)
		if err != nil {
			return
		}
		if !out.IsInside() {
			t.Fatalf("expected inside output range, got outside: [%v,%v]", out.Lo, out.Hi)
		}

		xs := []Rational{rx.Lo, rx.Hi, genSmallRat().Draw(t, "xSample")}
		ys := []Rational{ry.Lo, ry.Hi, genSmallRat().Draw(t, "ySample")}

		for _, x := range xs {
			if !rx.Contains(x) {
				continue
			}
			for _, y := range ys {
				if !ry.Contains(y) {
					continue
				}
				z, err := tform.ApplyRat(x, y)
				if err != nil {
					t.Fatalf("ApplyRat failed at (x=%v,y=%v): %v", x, y, err)
				}
				if !out.Contains(z) {
					t.Fatalf("enclosure failed: rx=[%v,%v] ry=[%v,%v] x=%v y=%v z=%v out=[%v,%v]",
						rx.Lo, rx.Hi, ry.Lo, ry.Hi, x, y, z, out.Lo, out.Hi)
				}
			}
		}
	})
}

// pbt_rapid_test.go v2
