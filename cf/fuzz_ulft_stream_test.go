// fuzz_ulft_stream_test.go v1
package cf

import (
	"math/big"
	"testing"
)

// FuzzULFTStreamMatchesExactRationalImage checks that streaming a finite rational
// source through ULFTStream matches exact ULFT.ApplyRat(x), whenever the ULFT is
// valid and defined at x.
//
// This fuzz target is intentionally conservative:
//   - it skips zero denominators in rational construction
//   - it skips invalid ULFTs
//   - it skips ULFTs undefined at the exact rational point
//   - it bounds output collection to avoid runaway cases
func FuzzULFTStreamMatchesExactRationalImage(f *testing.F) {
	// Seeds: a few simple, high-value cases.
	f.Add(int64(1), int64(1), int64(1), int64(0), int64(0), int64(1))    // identity at 1
	f.Add(int64(3), int64(2), int64(0), int64(1), int64(1), int64(0))    // reciprocal at 3/2
	f.Add(int64(-3), int64(2), int64(1), int64(1), int64(0), int64(1))   // x+1 at -3/2
	f.Add(int64(1), int64(1), int64(-1), int64(-1), int64(0), int64(-2)) // (x+1)/2 at 1
	f.Add(int64(0), int64(1), int64(0), int64(0), int64(0), int64(1))    // constant 0 at 0

	f.Fuzz(func(t *testing.T, p, q, a, b, c, d int64) {
		if q == 0 {
			t.Skip()
		}

		x, err := NewRational(p, q)
		if err != nil {
			t.Skip()
		}

		u := NewULFT(
			big.NewInt(a),
			big.NewInt(b),
			big.NewInt(c),
			big.NewInt(d),
		)

		if err := u.Validate(); err != nil {
			t.Skip()
		}

		wantRat, err := u.ApplyRat(x)
		if err != nil {
			// Undefined at the exact point is outside this fuzz target.
			t.Skip()
		}

		s := NewULFTStream(u, NewRationalCF(x), ULFTStreamOptions{})

		var got []int64
		for i := 0; i < 64; i++ {
			digit, ok := s.Next()
			if !ok {
				break
			}
			got = append(got, digit)
		}

		if err := s.Err(); err != nil {
			t.Fatalf("ULFTStream error: %v (x=%v u=%v)", err, x, u)
		}

		want := collectTerms(NewRationalCF(wantRat), 64)

		if len(got) != len(want) {
			t.Fatalf("len mismatch: got=%v want=%v x=%v u=%v T(x)=%v", got, want, x, u, wantRat)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("digit mismatch at %d: got=%v want=%v x=%v u=%v T(x)=%v", i, got, want, x, u, wantRat)
			}
		}

		// Stable exhaustion contract.
		if _, ok := s.Next(); ok {
			t.Fatalf("expected exhausted stream to stay exhausted")
		}
		if _, ok := s.Next(); ok {
			t.Fatalf("expected exhausted stream to stay exhausted on repeated calls")
		}
	})
}

// fuzz_ulft_stream_test.go v1
