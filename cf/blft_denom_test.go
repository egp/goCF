// blft_denom_test.go v1
package cf

import "testing"

func TestBLFTDenomRange_CornerBoundsExact(t *testing.T) {
	// Denom: D(x,y)= 2xy + 3x + 5y + 7
	tform := NewBLFT(0, 0, 0, 1, 2, 3, 5, 7)

	rx := NewRange(mustRat(1, 2), mustRat(5, 2), true, true) // [0.5,2.5]
	ry := NewRange(mustRat(1, 3), mustRat(7, 3), true, true) // [0.333..,2.333..]

	dr, err := tform.DenomRange(rx, ry)
	if err != nil {
		t.Fatal(err)
	}
	if !dr.IsInside() {
		t.Fatalf("expected inside denom range, got [%v,%v]", dr.Lo, dr.Hi)
	}

	// Compute denom at corners and ensure min/max match DenomRange.
	corners := []struct {
		x, y Rational
	}{
		{rx.Lo, ry.Lo},
		{rx.Lo, ry.Hi},
		{rx.Hi, ry.Lo},
		{rx.Hi, ry.Hi},
	}

	var dmin, dmax Rational
	for i, c := range corners {
		d, err := tform.denomAt(c.x, c.y)
		if err != nil {
			t.Fatal(err)
		}
		if i == 0 {
			dmin, dmax = d, d
			continue
		}
		if d.Cmp(dmin) < 0 {
			dmin = d
		}
		if d.Cmp(dmax) > 0 {
			dmax = d
		}
	}

	if dr.Lo.Cmp(dmin) != 0 || dr.Hi.Cmp(dmax) != 0 {
		t.Fatalf("DenomRange got [%v,%v], want [%v,%v]", dr.Lo, dr.Hi, dmin, dmax)
	}
}

func TestBLFTDenomMayHitZero_SpansRootXMinus1(t *testing.T) {
	// z = 1/(x-1), denom = x-1
	tform := NewBLFT(0, 0, 0, 1, 0, 1, 0, -1)

	rx := NewRange(mustRat(0, 1), mustRat(2, 1), true, true)  // spans x=1
	ry := NewRange(mustRat(-3, 1), mustRat(4, 1), true, true) // any

	may, err := tform.DenomMayHitZero(rx, ry)
	if err != nil {
		t.Fatal(err)
	}
	if !may {
		t.Fatalf("expected denom may hit 0 for rx spanning 1")
	}
}

func TestBLFTDenomMayHitZero_DefinitelyNonZero(t *testing.T) {
	// denom is constant 5
	tform := NewBLFT(1, 0, 0, 0, 0, 0, 0, 5)

	rx := NewRange(mustRat(-5, 1), mustRat(5, 1), true, true)
	ry := NewRange(mustRat(-7, 1), mustRat(9, 1), true, true)

	may, err := tform.DenomMayHitZero(rx, ry)
	if err != nil {
		t.Fatal(err)
	}
	if may {
		t.Fatalf("expected denom never hits 0 when constant 5")
	}

	dr, err := tform.DenomRange(rx, ry)
	if err != nil {
		t.Fatal(err)
	}
	if dr.Lo.Cmp(mustRat(5, 1)) != 0 || dr.Hi.Cmp(mustRat(5, 1)) != 0 {
		t.Fatalf("expected denom range [5,5], got [%v,%v]", dr.Lo, dr.Hi)
	}
}

func TestBLFTDenomAt_ConstantDenominator(t *testing.T) {
	// Denominator = 5
	bl := NewBLFT(0, 0, 0, 0, 0, 0, 0, 5)

	got, err := bl.denomAt(mustRat(2, 3), mustRat(-7, 5))
	if err != nil {
		t.Fatalf("denomAt failed: %v", err)
	}

	want := mustRat(5, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestBLFTDenomAt_GeneralCase(t *testing.T) {
	// Denominator = 2xy + 3x + 5y + 7
	bl := NewBLFT(0, 0, 0, 0, 2, 3, 5, 7)

	x := mustRat(1, 2)
	y := mustRat(2, 3)

	got, err := bl.denomAt(x, y)
	if err != nil {
		t.Fatalf("denomAt failed: %v", err)
	}

	// 2*(1/2)*(2/3) + 3*(1/2) + 5*(2/3) + 7
	// = 2/3 + 3/2 + 10/3 + 7 = 25/2
	want := mustRat(25, 2)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestBLFTDenomMayHitZero_FalseForStrictlyPositiveConstant(t *testing.T) {
	bl := NewBLFT(0, 0, 0, 0, 0, 0, 0, 3)

	rx := NewRange(mustRat(-2, 1), mustRat(5, 1), true, true)
	ry := NewRange(mustRat(-7, 1), mustRat(4, 1), true, true)

	got, err := bl.DenomMayHitZero(rx, ry)
	if err != nil {
		t.Fatalf("DenomMayHitZero failed: %v", err)
	}
	if got {
		t.Fatalf("expected false")
	}
}

func TestBLFTDenomMayHitZero_TrueWhenCornerRangeIncludesZero(t *testing.T) {
	// Denominator = y - 1
	bl := NewBLFT(0, 0, 0, 0, 0, 0, 1, -1)

	rx := NewRange(mustRat(0, 1), mustRat(0, 1), true, true)
	ry := NewRange(mustRat(0, 1), mustRat(2, 1), true, true)

	got, err := bl.DenomMayHitZero(rx, ry)
	if err != nil {
		t.Fatalf("DenomMayHitZero failed: %v", err)
	}
	if !got {
		t.Fatalf("expected true")
	}
}

func TestBLFTDenomMayHitZero_TrueAtExactPointZero(t *testing.T) {
	// Denominator = x + y - 1
	bl := NewBLFT(0, 0, 0, 0, 0, 1, 1, -1)

	rx := NewRange(mustRat(1, 2), mustRat(1, 2), true, true)
	ry := NewRange(mustRat(1, 2), mustRat(1, 2), true, true)

	got, err := bl.DenomMayHitZero(rx, ry)
	if err != nil {
		t.Fatalf("DenomMayHitZero failed: %v", err)
	}
	if !got {
		t.Fatalf("expected true")
	}
}

func TestBLFTDenomMayHitZero_FalseOnSeparatedPositiveRange(t *testing.T) {
	// Denominator = x + y + 1, clearly positive on the chosen rectangle.
	bl := NewBLFT(0, 0, 0, 0, 0, 1, 1, 1)

	rx := NewRange(mustRat(0, 1), mustRat(2, 1), true, true)
	ry := NewRange(mustRat(0, 1), mustRat(3, 1), true, true)

	got, err := bl.DenomMayHitZero(rx, ry)
	if err != nil {
		t.Fatalf("DenomMayHitZero failed: %v", err)
	}
	if got {
		t.Fatalf("expected false")
	}
}

// blft_denom_test.go v1
