// sqrt_range_conservative_test.go v1
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

// sqrt_range_conservative_test.go v1
