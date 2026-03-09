// pbt_ulft_stream_equiv_test.go v1
package cf

import (
	"testing"

	"pgregory.net/rapid"
)

func genSmallNonzeroRat() *rapid.Generator[Rational] {
	return rapid.Custom(func(t *rapid.T) Rational {
		p := rapid.Int64Range(-9, 9).Draw(t, "p")
		if p == 0 {
			p = 1
		}
		q := rapid.Int64Range(1, 9).Draw(t, "q")
		return mustRat(p, q)
	})
}

func genSmallULFTNoPoleAt(x Rational) *rapid.Generator[ULFT] {
	return rapid.Custom(func(t *rapid.T) ULFT {
		for {
			a := rapid.Int64Range(-5, 5).Draw(t, "a")
			b := rapid.Int64Range(-5, 5).Draw(t, "b")
			c := rapid.Int64Range(-5, 5).Draw(t, "c")
			d := rapid.Int64Range(-5, 5).Draw(t, "d")

			u := NewULFT(mustBig(a), mustBig(b), mustBig(c), mustBig(d))
			if err := u.Validate(); err != nil {
				continue
			}
			if _, err := u.ApplyRat(x); err != nil {
				continue
			}
			return u
		}
	})
}

func TestPBT_ULFTStreamMatchesExactRationalImage(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := genSmallNonzeroRat().Draw(t, "x")
		u := genSmallULFTNoPoleAt(x).Draw(t, "u")

		wantRat, err := u.ApplyRat(x)
		if err != nil {
			t.Fatalf("ApplyRat failed: %v", err)
		}

		stream := NewULFTStream(u, NewRationalCF(x), ULFTStreamOptions{})
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
			t.Fatalf("ULFTStream error: %v", err)
		}

		want := collectTerms(NewRationalCF(wantRat), 32)

		if len(got) != len(want) {
			t.Fatalf("len mismatch: got=%v want=%v x=%v T=%v T(x)=%v", got, want, x, u, wantRat)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("digit mismatch at %d: got=%v want=%v x=%v T=%v T(x)=%v", i, got, want, x, u, wantRat)
			}
		}
	})
}

// pbt_ulft_stream_equiv_test.go v1
