// blft_range_test.go v1
package cf

import "testing"

func TestBLFTRange_PointRangesMatchApplyRat(t *testing.T) {
	// z = (2xy + x + y + 1) / (xy + 2x + 3y + 4)
	tform := NewBLFT(2, 1, 1, 1, 1, 2, 3, 4)

	x := mustRat(3, 2)
	y := mustRat(7, 5)

	want, err := tform.ApplyRat(x, y)
	if err != nil {
		t.Fatal(err)
	}

	rx := MustRange(x, x)
	ry := MustRange(y, y)

	gotR, err := tform.ApplyBLFTRange(rx, ry)
	if err != nil {
		t.Fatal(err)
	}
	if !gotR.IsInside() {
		t.Fatalf("expected inside range")
	}
	if gotR.Lo.Cmp(want) != 0 || gotR.Hi.Cmp(want) != 0 {
		t.Fatalf("got [%v,%v], want exact [%v,%v]", gotR.Lo, gotR.Hi, want, want)
	}
}

func TestBLFTRange_AdditionLikeTransform(t *testing.T) {
	// z = x + y  => (0*xy + 1*x + 1*y + 0) / (0*xy + 0*x + 0*y + 1)
	add := NewBLFT(0, 1, 1, 0, 0, 0, 0, 1)

	rx := MustRange(mustRat(1, 1), mustRat(2, 1))   // [1,2]
	ry := MustRange(mustRat(10, 1), mustRat(20, 1)) // [10,20]

	got, err := add.ApplyBLFTRange(rx, ry)
	if err != nil {
		t.Fatal(err)
	}

	want := MustRange(mustRat(11, 1), mustRat(22, 1))
	if got.Lo.Cmp(want.Lo) != 0 || got.Hi.Cmp(want.Hi) != 0 {
		t.Fatalf("got [%v,%v], want [%v,%v]", got.Lo, got.Hi, want.Lo, want.Hi)
	}
}

func TestBLFTRange_DenominatorCrossesZeroRejected(t *testing.T) {
	// z = 1/(x-1)  independent of y:
	// numerator: 1  => D=1
	// denom: x-1    => F=1, H=-1
	tform := NewBLFT(0, 0, 0, 1, 0, 1, 0, -1)

	rx := MustRange(mustRat(0, 1), mustRat(2, 1)) // spans x=1
	ry := MustRange(mustRat(0, 1), mustRat(1, 1)) // any

	_, err := tform.ApplyBLFTRange(rx, ry)
	if err == nil {
		t.Fatalf("expected error due to denom sign change across corners")
	}
}

func TestBLFTRange_OutsideInputRejected(t *testing.T) {
	add := NewBLFT(0, 1, 1, 0, 0, 0, 0, 1)

	rx := NewRange(mustRat(2, 1), mustRat(1, 1)) // outside
	ry := MustRange(mustRat(0, 1), mustRat(1, 1))

	_, err := add.ApplyBLFTRange(rx, ry)
	if err == nil {
		t.Fatalf("expected error for outside input range")
	}
}

// blft_range_test.go v1
