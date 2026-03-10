// pbt_diag_blft_stream_equiv_test.go v2
package cf

import (
	"testing"

	"pgregory.net/rapid"
)

func genSmallSupportedDiagBLFTNoPoleAt(x Rational) *rapid.Generator[DiagBLFT] {
	return rapid.Custom(func(t *rapid.T) DiagBLFT {
		for {
			a := rapid.Int64Range(-5, 5).Draw(t, "a")
			b := rapid.Int64Range(-5, 5).Draw(t, "b")
			c := rapid.Int64Range(-5, 5).Draw(t, "c")
			f := rapid.Int64Range(-5, 5).Draw(t, "f")

			// Restrict to the currently supported ApplyRange subclass:
			// constant nonzero denominator.
			if f == 0 {
				continue
			}

			u := NewDiagBLFT(
				mustBig(a),
				mustBig(b),
				mustBig(c),
				mustBig(0),
				mustBig(0),
				mustBig(f),
			)

			if _, err := u.ApplyRat(x); err != nil {
				continue
			}
			return u
		}
	})
}

func TestPBT_DiagBLFTStreamMatchesExactRationalImage(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		x := genSmallRat().Draw(t, "x")
		u := genSmallSupportedDiagBLFTNoPoleAt(x).Draw(t, "u")

		wantRat, err := u.ApplyRat(x)
		if err != nil {
			t.Fatalf("ApplyRat failed: %v", err)
		}

		stream := NewDiagBLFTStream(u, NewRationalCF(x), DiagBLFTStreamOptions{})

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
			t.Fatalf("DiagBLFTStream error: %v", err)
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

// pbt_diag_blft_stream_equiv_test.go v2
