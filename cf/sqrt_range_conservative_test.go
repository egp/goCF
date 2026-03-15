// sqrt_range_conservative_test.go v3
package cf

import "testing"

func TestSqrtLowerBoundRational_ExactSquare(t *testing.T) {
	got, err := SqrtLowerBoundRational(mustRat(9, 16))
	if err != nil {
		t.Fatalf("SqrtLowerBoundRational failed: %v", err)
	}

	want := mustRat(3, 4)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSqrtUpperBoundRational_ExactSquare(t *testing.T) {
	got, err := SqrtUpperBoundRational(mustRat(9, 16))
	if err != nil {
		t.Fatalf("SqrtUpperBoundRational failed: %v", err)
	}

	want := mustRat(3, 4)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSqrtLowerBoundRational_RejectsNegative(t *testing.T) {
	_, err := SqrtLowerBoundRational(mustRat(-1, 1))
	if err == nil {
		t.Fatalf("expected error for negative input")
	}
}

func TestSqrtUpperBoundRational_RejectsNegative(t *testing.T) {
	_, err := SqrtUpperBoundRational(mustRat(-1, 1))
	if err == nil {
		t.Fatalf("expected error for negative input")
	}
}

func TestSqrtLowerBoundRational_NonSquareIsBelowRoot(t *testing.T) {
	x := mustRat(2, 1)

	got, err := SqrtLowerBoundRational(x)
	if err != nil {
		t.Fatalf("SqrtLowerBoundRational failed: %v", err)
	}

	g2, err := got.Mul(got)
	if err != nil {
		t.Fatalf("Mul failed: %v", err)
	}
	if g2.Cmp(x) > 0 {
		t.Fatalf("lower bound squared exceeded x: got^2=%v x=%v", g2, x)
	}
}

func TestSqrtUpperBoundRational_NonSquareIsAboveRoot(t *testing.T) {
	x := mustRat(2, 1)

	got, err := SqrtUpperBoundRational(x)
	if err != nil {
		t.Fatalf("SqrtUpperBoundRational failed: %v", err)
	}

	g2, err := got.Mul(got)
	if err != nil {
		t.Fatalf("Mul failed: %v", err)
	}
	if g2.Cmp(x) < 0 {
		t.Fatalf("upper bound squared fell below x: got^2=%v x=%v", g2, x)
	}
}

func TestSqrtRangeConservative_ExactEndpoints(t *testing.T) {
	r := NewRange(mustRat(1, 4), mustRat(9, 16), true, true)

	got, err := SqrtRangeConservative(r)
	if err != nil {
		t.Fatalf("SqrtRangeConservative failed: %v", err)
	}

	wantLo := mustRat(1, 2)
	wantHi := mustRat(3, 4)
	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v] want [%v,%v]", got.Lo, got.Hi, wantLo, wantHi)
	}
	if got.IncLo != r.IncLo || got.IncHi != r.IncHi {
		t.Fatalf("got inclusions (%v,%v) want (%v,%v)", got.IncLo, got.IncHi, r.IncLo, r.IncHi)
	}
}

func TestSqrtRangeConservative_NonSquareEndpointsIsConservative(t *testing.T) {
	r := NewRange(mustRat(4, 3), mustRat(3, 2), true, true)

	got, err := SqrtRangeConservative(r)
	if err != nil {
		t.Fatalf("SqrtRangeConservative failed: %v", err)
	}

	lo2, err := got.Lo.Mul(got.Lo)
	if err != nil {
		t.Fatalf("Mul lo failed: %v", err)
	}
	hi2, err := got.Hi.Mul(got.Hi)
	if err != nil {
		t.Fatalf("Mul hi failed: %v", err)
	}

	if lo2.Cmp(r.Lo) > 0 {
		t.Fatalf("lower enclosure invalid: lo^2=%v original lo=%v", lo2, r.Lo)
	}
	if hi2.Cmp(r.Hi) < 0 {
		t.Fatalf("upper enclosure invalid: hi^2=%v original hi=%v", hi2, r.Hi)
	}
}

func TestSqrtRangeConservative_Sqrt2PrefixTwoTermsCertifiesFirstDigit(t *testing.T) {
	r := NewRange(mustRat(4, 3), mustRat(3, 2), true, true)

	got, err := SqrtRangeConservative(r)
	if err != nil {
		t.Fatalf("SqrtRangeConservative failed: %v", err)
	}

	lo, hi, err := got.FloorBounds()
	if err != nil {
		t.Fatalf("FloorBounds failed: %v", err)
	}
	if lo != 1 || hi != 1 {
		t.Fatalf("got floor bounds (%d,%d), want (1,1); range=%v", lo, hi, got)
	}
}

func TestSqrtRangeConservative_RejectsNegativeRange(t *testing.T) {
	r := NewRange(mustRat(-1, 1), mustRat(1, 1), true, true)

	_, err := SqrtRangeConservative(r)
	if err == nil {
		t.Fatalf("expected error for negative range")
	}
}

func TestSqrtRangeConservative_RejectsOutsideRange(t *testing.T) {
	r := NewRange(mustRat(2, 1), mustRat(1, 1), true, true)

	_, err := SqrtRangeConservative(r)
	if err == nil {
		t.Fatalf("expected error for outside range")
	}
}

// sqrt_range_conservative_test.go v3
