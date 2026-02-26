// range_test.go v1
package cf

import "testing"

func TestNewRange_ValidAndInvalid(t *testing.T) {
	lo := mustRat(1, 3)
	hi := mustRat(1, 2)

	_, err := NewRange(hi, lo)
	if err == nil {
		t.Fatalf("expected error for lo>hi, got nil")
	}

	r, err := NewRange(lo, hi)
	if err != nil {
		t.Fatal(err)
	}
	if r.Lo.Cmp(lo) != 0 || r.Hi.Cmp(hi) != 0 {
		t.Fatalf("got [%v,%v], want [%v,%v]", r.Lo, r.Hi, lo, hi)
	}
}

func TestRange_Contains(t *testing.T) {
	r := MustRange(mustRat(1, 3), mustRat(1, 2))

	if !r.Contains(mustRat(1, 3)) {
		t.Fatalf("expected to contain Lo")
	}
	if !r.Contains(mustRat(1, 2)) {
		t.Fatalf("expected to contain Hi")
	}
	if !r.Contains(mustRat(2, 5)) { // 0.4 is between 0.333.. and 0.5
		t.Fatalf("expected to contain 2/5")
	}
	if r.Contains(mustRat(1, 4)) {
		t.Fatalf("did not expect to contain 1/4")
	}
	if r.Contains(mustRat(2, 3)) {
		t.Fatalf("did not expect to contain 2/3")
	}
}

func TestRange_Width(t *testing.T) {
	r := MustRange(mustRat(1, 3), mustRat(1, 2))
	w, err := r.Width()
	if err != nil {
		t.Fatal(err)
	}
	want := mustRat(1, 6)
	if w.Cmp(want) != 0 {
		t.Fatalf("width got %v, want %v", w, want)
	}
}

func TestRange_FloorBounds_Positive(t *testing.T) {
	r := MustRange(mustRat(1, 3), mustRat(5, 2)) // [0.333.., 2.5]
	lo, hi := r.FloorBounds()
	if lo != 0 || hi != 2 {
		t.Fatalf("floor bounds got (%d,%d), want (0,2)", lo, hi)
	}
}

func TestRange_FloorBounds_Negative_FloorConvention(t *testing.T) {
	// [-4/3, -1/2] = [-1.333.., -0.5]
	// floors: floor(-1.333..)=-2, floor(-0.5)=-1
	r := MustRange(mustRat(-4, 3), mustRat(-1, 2))
	lo, hi := r.FloorBounds()
	if lo != -2 || hi != -1 {
		t.Fatalf("floor bounds got (%d,%d), want (-2,-1)", lo, hi)
	}
}

func TestRange_ApplyULFT_Monotone(t *testing.T) {
	// f(x) = (2x+1)/(3x+4)
	// Over x in [1/3, 1/2], denominator is positive and monotone; image is interval.
	tform := NewULFT(2, 1, 3, 4)
	src := MustRange(mustRat(1, 3), mustRat(1, 2))

	got, err := src.ApplyULFT(tform)
	if err != nil {
		t.Fatal(err)
	}

	// Compute expected endpoints directly (then order them).
	fLo, _ := tform.ApplyRat(src.Lo)
	fHi, _ := tform.ApplyRat(src.Hi)

	wantLo, wantHi := fLo, fHi
	if wantLo.Cmp(wantHi) > 0 {
		wantLo, wantHi = wantHi, wantLo
	}

	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v], want [%v,%v]", got.Lo, got.Hi, wantLo, wantHi)
	}
}

func TestRange_ApplyULFT_DenominatorCrossesZero(t *testing.T) {
	// f(x) = 1/(x-1)  => (0*x+1)/(1*x-1)
	// Denominator crosses 0 at x=1. Range [0,2] spans it, so reject.
	tform := NewULFT(0, 1, 1, -1)
	src := MustRange(mustRat(0, 1), mustRat(2, 1))

	_, err := src.ApplyULFT(tform)
	if err == nil {
		t.Fatalf("expected error when denominator crosses zero")
	}
}

// range_test.go v1
