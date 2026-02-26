// blft_stream_test.go v1
package cf

import "testing"

func TestBLFTStream_Golden_RationalAddition(t *testing.T) {
	// z = x + y
	add := NewBLFT(0, 1, 1, 0, 0, 0, 0, 1)

	type tc struct {
		x, y Rational
	}
	cases := []tc{
		{mustRat(1, 2), mustRat(1, 3)},  // 5/6
		{mustRat(3, 2), mustRat(7, 5)},  // 29/10
		{mustRat(-4, 3), mustRat(5, 7)}, // (-13/21)
	}

	for _, c := range cases {
		wantRat, err := add.ApplyRat(c.x, c.y)
		if err != nil {
			t.Fatal(err)
		}
		want := collectAll(NewRationalCF(wantRat))

		s := NewBLFTStream(add, NewRationalCF(c.x), NewRationalCF(c.y), BLFTStreamOptions{})
		got := collectAll(s)
		if s.Err() != nil {
			t.Fatalf("stream error for x=%v y=%v: %v", c.x, c.y, s.Err())
		}

		if !equalSlice(got, want) {
			t.Fatalf("x=%v y=%v got %v, want %v (rat=%v)", c.x, c.y, got, want, wantRat)
		}
	}
}

func TestBLFTStream_Golden_RationalGeneral(t *testing.T) {
	// z = (2xy + x + y + 1) / (xy + 2x + 3y + 4)
	tform := NewBLFT(2, 1, 1, 1, 1, 2, 3, 4)

	x := mustRat(3, 2)
	y := mustRat(7, 5)

	wantRat, err := tform.ApplyRat(x, y)
	if err != nil {
		t.Fatal(err)
	}
	want := collectAll(NewRationalCF(wantRat))

	s := NewBLFTStream(tform, NewRationalCF(x), NewRationalCF(y), BLFTStreamOptions{})
	got := collectAll(s)
	if s.Err() != nil {
		t.Fatalf("stream error: %v", s.Err())
	}

	if !equalSlice(got, want) {
		t.Fatalf("got %v, want %v (rat=%v)", got, want, wantRat)
	}
}

// blft_stream_test.go v1
