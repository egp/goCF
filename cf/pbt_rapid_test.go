// pbt_rapid_test.go v1
package cf

import (
	"testing"

	"pgregory.net/rapid"
)

func genSmallInt64() *rapid.Generator[int64] {
	// Keep tiny to avoid overflow while we’re still int64.
	return rapid.Int64Range(-20, 20)
}

func genPosInt64() *rapid.Generator[int64] {
	return rapid.Int64Range(1, 20)
}

func genSmallRat() *rapid.Generator[Rational] {
	return rapid.Custom(func(t *rapid.T) Rational {
		p := genSmallInt64().Draw(t, "p")
		q := genPosInt64().Draw(t, "q")
		return mustRat(p, q)
	})
}

func genInsideRange() *rapid.Generator[Range] {
	return rapid.Custom(func(t *rapid.T) Range {
		a := genSmallRat().Draw(t, "a")
		b := genSmallRat().Draw(t, "b")
		if a.Cmp(b) <= 0 {
			return Range{Lo: a, Hi: b}
		}
		return Range{Lo: b, Hi: a}
	})
}

func genOutsideRange() *rapid.Generator[Range] {
	return rapid.Custom(func(t *rapid.T) Range {
		// Ensure Lo > Hi.
		lo := genSmallRat().Draw(t, "lo")
		hi := genSmallRat().Draw(t, "hi")
		if lo.Cmp(hi) <= 0 {
			// swap then nudge to enforce strict >
			lo, hi = hi, lo
		}
		if lo.Cmp(hi) == 0 {
			// force outside by moving lo upward
			lo = mustRat(lo.P+1, lo.Q)
		}
		if lo.Cmp(hi) <= 0 {
			// paranoia
			lo = mustRat(1, 1)
			hi = mustRat(0, 1)
		}
		return Range{Lo: lo, Hi: hi}
	})
}

func TestPBT_RangeOutsideContainsSemantics(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		r := genOutsideRange().Draw(t, "r")
		x := genSmallRat().Draw(t, "x")

		// Outside semantics: x <= Hi OR x >= Lo
		want := x.Cmp(r.Hi) <= 0 || x.Cmp(r.Lo) >= 0
		got := r.Contains(x)
		if got != want {
			t.Fatalf("outside Contains mismatch: r=[%v,%v] x=%v got=%v want=%v", r.Lo, r.Hi, x, got, want)
		}

		// ContainsZero: Hi>=0 OR Lo<=0 (equivalently: NOT (Hi<0<Lo))
		z := mustRat(0, 1)
		wantZ := z.Cmp(r.Hi) <= 0 || z.Cmp(r.Lo) >= 0
		gotZ := r.ContainsZero()
		if gotZ != wantZ {
			t.Fatalf("outside ContainsZero mismatch: r=[%v,%v] got=%v want=%v", r.Lo, r.Hi, gotZ, wantZ)
		}
	})
}

func TestPBT_BLFTRangeRejectsPoleForXMinus1(t *testing.T) {
	// z = 1/(x-1), independent of y
	// denom = x - 1  => E=0 F=1 G=0 H=-1
	tform := NewBLFT(0, 0, 0, 1, 0, 1, 0, -1)

	rapid.Check(t, func(t *rapid.T) {
		// Construct an inside x-range that straddles 1: [1-a, 1+b], with a,b>=1
		a := rapid.Int64Range(1, 10).Draw(t, "a")
		b := rapid.Int64Range(1, 10).Draw(t, "b")
		xlo := mustRat(1-a, 1)
		xhi := mustRat(1+b, 1)
		rx := MustRange(xlo, xhi)

		// Any y-range; keep simple inside.
		ry := genInsideRange().Draw(t, "ry")

		_, err := tform.ApplyBLFTRange(rx, ry)
		if err == nil {
			t.Fatalf("expected pole rejection for rx spanning 1: rx=[%v,%v] ry=[%v,%v]", rx.Lo, rx.Hi, ry.Lo, ry.Hi)
		}
	})
}

func TestPBT_BLFTRangeEnclosesSamplesWhenAccepted(t *testing.T) {
	// Use a BLFT with a relatively stable denom in our small rational domain.
	tform := NewBLFT(
		2, 1, -1, 3,
		1, 2, 1, 5,
	)

	rapid.Check(t, func(t *rapid.T) {
		rx := genInsideRange().Draw(t, "rx")
		ry := genInsideRange().Draw(t, "ry")

		out, err := tform.ApplyBLFTRange(rx, ry)
		if err != nil {
			// Rejections are OK (pole hazard); property is only enforced when accepted.
			return
		}
		if !out.IsInside() {
			t.Fatalf("expected inside output range, got outside: [%v,%v]", out.Lo, out.Hi)
		}

		// Sample a few interior points by picking endpoints and mid-ish points.
		xs := []Rational{rx.Lo, rx.Hi}
		ys := []Rational{ry.Lo, ry.Hi}

		// Add one random interior sample each (if possible).
		xs = append(xs, genSmallRat().Draw(t, "xSample"))
		ys = append(ys, genSmallRat().Draw(t, "ySample"))

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
					// If ApplyRat is undefined at a sampled point, enclosure is meaningless.
					// But ApplyBLFTRange accepted based on denom enclosure, so this should be rare.
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

// pbt_rapid_test.go v1
