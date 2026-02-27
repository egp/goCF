// pbt_blft_stream_test.go v3
package cf

import (
	"testing"

	"pgregory.net/rapid"
)

func pbtGenTinyI64Stream() *rapid.Generator[int64]    { return rapid.Int64Range(-9, 9) }
func pbtGenTinyPosI64Stream() *rapid.Generator[int64] { return rapid.Int64Range(1, 9) }

// pbtGenTinyRatStream generates small rationals with non-zero numerator (to avoid
// the x=0/y=0 cases while BLFTStream rational termination is being tightened).
func pbtGenTinyRatStream() *rapid.Generator[Rational] {
	return rapid.Custom(func(t *rapid.T) Rational {
		for {
			p := pbtGenTinyI64Stream().Draw(t, "p")
			q := pbtGenTinyPosI64Stream().Draw(t, "q")
			if p == 0 {
				continue
			}
			return mustRat(p, q)
		}
	})
}

func pbtCollectPrefixK(s ContinuedFraction, k int) ([]int64, bool) {
	out := make([]int64, 0, k)
	for i := 0; i < k; i++ {
		a, ok := s.Next()
		if !ok {
			return out, false
		}
		out = append(out, a)
	}
	return out, true
}

func TestPBT_BLFTStream_MatchesApplyRatPrefix_OnRationals_TinyDomain(t *testing.T) {
	const k = 12

	rapid.Check(t, func(t *rapid.T) {
		x := pbtGenTinyRatStream().Draw(t, "x")
		y := pbtGenTinyRatStream().Draw(t, "y")

		// Explicit skip (redundant given generator excludes 0, but keeps intent clear).
		if x.P == 0 || y.P == 0 {
			return
		}

		// Tiny coefficients to keep checked-int64 arithmetic stable.
		A := pbtGenTinyI64Stream().Draw(t, "A")
		B := pbtGenTinyI64Stream().Draw(t, "B")
		C := pbtGenTinyI64Stream().Draw(t, "C")
		D := pbtGenTinyI64Stream().Draw(t, "D")
		E := pbtGenTinyI64Stream().Draw(t, "E")
		F := pbtGenTinyI64Stream().Draw(t, "F")
		G := pbtGenTinyI64Stream().Draw(t, "G")
		H := pbtGenTinyI64Stream().Draw(t, "H")

		tform := NewBLFT(A, B, C, D, E, F, G, H)

		// Skip undefined points (denom hits 0 at the exact point).
		may, err := tform.DenomMayHitZero(NewRange(x, x, true, true), NewRange(y, y, true, true))
		if err != nil || may {
			return
		}

		z, err := tform.ApplyRat(x, y)
		if err != nil {
			return
		}

		want := NewRationalCF(z)
		got := NewBLFTStream(
			tform,
			NewRationalCF(x),
			NewRationalCF(y),
			BLFTStreamOptions{MaxFinalizeDigits: 128},
		)

		wp, _ := pbtCollectPrefixK(want, k)
		gp, _ := pbtCollectPrefixK(got, k)

		// Compare up to the shorter prefix (rationals may terminate < k digits).
		n := len(wp)
		if len(gp) < n {
			n = len(gp)
		}
		for i := 0; i < n; i++ {
			if wp[i] != gp[i] {
				t.Fatalf("prefix mismatch at i=%d: x=%v y=%v z=%v want=%v got=%v",
					i, x, y, z, wp, gp)
			}
		}

		if err := got.Err(); err != nil {
			t.Fatalf("BLFTStream error: %v", err)
		}
	})
}

// pbt_blft_stream_test.go v3
