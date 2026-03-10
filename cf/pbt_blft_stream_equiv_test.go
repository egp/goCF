// pbt_blft_stream_equiv_test.go v1
package cf

import (
	"testing"

	"pgregory.net/rapid"
)

func genSmallBLFTNoPoleAt(x, y Rational) *rapid.Generator[BLFT] {
	return rapid.Custom(func(t *rapid.T) BLFT {
		for {
			a := rapid.Int64Range(-5, 5).Draw(t, "a")
			b := rapid.Int64Range(-5, 5).Draw(t, "b")
			c := rapid.Int64Range(-5, 5).Draw(t, "c")
			d := rapid.Int64Range(-5, 5).Draw(t, "d")
			e := rapid.Int64Range(-5, 5).Draw(t, "e")
			f := rapid.Int64Range(-5, 5).Draw(t, "f")
			g := rapid.Int64Range(-5, 5).Draw(t, "g")
			h := rapid.Int64Range(-5, 5).Draw(t, "h")

			// avoid identically zero denominator
			if e == 0 && f == 0 && g == 0 && h == 0 {
				continue
			}

			bl := NewBLFT(a, b, c, d, e, f, g, h)
			if _, err := bl.ApplyRat(x, y); err != nil {
				continue
			}
			return bl
		}
	})
}

func TestPBT_BLFTStreamMatchesExactRationalImage(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := genSmallRat().Draw(t, "x")
		y := genSmallRat().Draw(t, "y")
		bl := genSmallBLFTNoPoleAt(x, y).Draw(t, "bl")

		wantRat, err := bl.ApplyRat(x, y)
		if err != nil {
			t.Fatalf("ApplyRat failed: %v", err)
		}

		stream := NewBLFTStream(
			bl,
			NewRationalCF(x),
			NewRationalCF(y),
			BLFTStreamOptions{},
		)

		var got []int64
		for {
			d, ok := stream.Next()
			if !ok {
				break
			}
			got = append(got, d)
			if len(got) > 32 {
				t.Fatalf("unexpectedly long output: %v", got)
			}
		}
		if err := stream.Err(); err != nil {
			t.Fatalf("BLFTStream error: %v", err)
		}

		want := collectTerms(NewRationalCF(wantRat), 32)

		if len(got) != len(want) {
			t.Fatalf("len mismatch: got=%v want=%v x=%v y=%v T=%v T(x,y)=%v", got, want, x, y, bl, wantRat)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("digit mismatch at %d: got=%v want=%v x=%v y=%v T=%v T(x,y)=%v", i, got, want, x, y, bl, wantRat)
			}
		}
	})
}

// pbt_blft_stream_equiv_test.go v1
