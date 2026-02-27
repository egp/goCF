// blft_stream_test.go v1
package cf

import "testing"

func blftCollectAll(s *BLFTStream, limit int) ([]int64, error) {
	out := make([]int64, 0, 16)
	for i := 0; i < limit; i++ {
		a, ok := s.Next()
		if !ok {
			return out, s.Err()
		}
		out = append(out, a)
	}
	return out, s.Err()
}

func blftCollectAllCF(s ContinuedFraction, limit int) ([]int64, bool) {
	out := make([]int64, 0, 16)
	for i := 0; i < limit; i++ {
		a, ok := s.Next()
		if !ok {
			return out, false
		}
		out = append(out, a)
	}
	return out, true
}

func TestBLFTStream_Golden_RationalAddition(t *testing.T) {
	// z = x + y
	tform := NewBLFT(
		0, 1, 1, 0, // numerator: x + y
		0, 0, 0, 1, // denom: 1
	)

	x := mustRat(1, 2)
	y := mustRat(1, 3)

	// 1/2 + 1/3 = 5/6 = [0; 1, 5]
	want := []int64{0, 1, 5}

	s := NewBLFTStream(
		tform,
		NewRationalCF(x),
		NewRationalCF(y),
		BLFTStreamOptions{MaxFinalizeDigits: 128},
	)

	got, err := blftCollectAll(s, 32)
	if err != nil {
		t.Fatalf("stream error for x=%v y=%v: %v", x, y, err)
	}

	if len(got) != len(want) {
		t.Fatalf("digits len mismatch: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("digits mismatch at i=%d: got=%v want=%v", i, got, want)
		}
	}
}

func TestBLFTStream_Golden_RationalGeneral_MatchesApplyRat(t *testing.T) {
	// A generic BLFT that’s well-defined at the chosen point.
	tform := NewBLFT(
		1, 2, 3, 4,
		0, 1, 1, 5, // denom = x + y + 5 (positive for chosen x,y)
	)

	x := mustRat(3, 2)
	y := mustRat(7, 5)

	z, err := tform.ApplyRat(x, y)
	if err != nil {
		t.Fatalf("ApplyRat failed: %v", err)
	}

	wantStream := NewRationalCF(z)
	want, _ := blftCollectAllCF(wantStream, 64)

	gotStream := NewBLFTStream(
		tform,
		NewRationalCF(x),
		NewRationalCF(y),
		BLFTStreamOptions{MaxFinalizeDigits: 128},
	)

	got, err := blftCollectAll(gotStream, 64)
	if err != nil {
		t.Fatalf("stream error: %v", err)
	}

	// For rationals, both should terminate quickly; compare full collected prefixes.
	if len(got) != len(want) {
		t.Fatalf("digits len mismatch: x=%v y=%v z=%v got=%v want=%v", x, y, z, got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("digits mismatch at i=%d: x=%v y=%v z=%v got=%v want=%v", i, x, y, z, got, want)
		}
	}
}

func TestBLFTStream_TerminatesOnConstantZero(t *testing.T) {
	// z(x,y) = 0/1 = 0 everywhere.
	tform := NewBLFT(
		0, 0, 0, 0,
		0, 0, 0, 1,
	)

	x := mustRat(1, 1)
	y := mustRat(1, 1)

	s := NewBLFTStream(
		tform,
		NewRationalCF(x),
		NewRationalCF(y),
		BLFTStreamOptions{MaxFinalizeDigits: 128},
	)

	a0, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit")
	}
	if a0 != 0 {
		t.Fatalf("expected 0, got %d", a0)
	}

	_, ok = s.Next()
	if ok {
		t.Fatalf("expected stream to terminate after single 0 digit")
	}
	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream error: %v", err)
	}
}

func TestBLFTStream_RationalPointAgreesWithRationalCF_Table(t *testing.T) {
	cases := []struct {
		name      string
		t         BLFT
		x, y      Rational
		maxDigits int
	}{
		{
			name:      "identity_like_x",
			t:         NewBLFT(0, 1, 0, 0, 0, 0, 0, 1), // z = x
			x:         mustRat(355, 113),
			y:         mustRat(1, 1),
			maxDigits: 32,
		},
		{
			name:      "reciprocal_like_1_over_x",
			t:         NewBLFT(0, 0, 0, 1, 0, 1, 0, 0), // z = 1/x (ignores y)
			x:         mustRat(7, 5),
			y:         mustRat(1, 1),
			maxDigits: 32,
		},
		{
			name:      "x_plus_1",
			t:         NewBLFT(0, 1, 0, 1, 0, 0, 0, 1), // z = x + 1
			x:         mustRat(3, 2),
			y:         mustRat(1, 1),
			maxDigits: 32,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Skip undefined point if denom hits 0 exactly.
			may, err := tc.t.DenomMayHitZero(NewRange(tc.x, tc.x, true, true), NewRange(tc.y, tc.y, true, true))
			if err != nil || may {
				return
			}

			z, err := tc.t.ApplyRat(tc.x, tc.y)
			if err != nil {
				t.Fatalf("ApplyRat failed: %v", err)
			}

			want, _ := blftCollectAllCF(NewRationalCF(z), tc.maxDigits)

			gotStream := NewBLFTStream(
				tc.t,
				NewRationalCF(tc.x),
				NewRationalCF(tc.y),
				BLFTStreamOptions{MaxFinalizeDigits: 128},
			)
			got, err := blftCollectAll(gotStream, tc.maxDigits)
			if err != nil {
				t.Fatalf("stream error: %v", err)
			}

			if len(got) != len(want) {
				t.Fatalf("digits len mismatch: x=%v y=%v z=%v got=%v want=%v", tc.x, tc.y, z, got, want)
			}
			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("digits mismatch at i=%d: x=%v y=%v z=%v got=%v want=%v", i, tc.x, tc.y, z, got, want)
				}
			}
		})
	}
}

// blft_stream_test.go v1
