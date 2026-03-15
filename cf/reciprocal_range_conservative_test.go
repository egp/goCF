// reciprocal_range_conservative_test.go v1
package cf

import "testing"

func TestReciprocalRangeConservative_ExactPoint(t *testing.T) {
	r := NewRange(mustRat(2, 1), mustRat(2, 1), true, true)

	got, err := ReciprocalRangeConservative(r)
	if err != nil {
		t.Fatalf("ReciprocalRangeConservative failed: %v", err)
	}

	want := mustRat(1, 2)
	if got.Lo.Cmp(want) != 0 || got.Hi.Cmp(want) != 0 {
		t.Fatalf("got [%v,%v] want exact [%v,%v]", got.Lo, got.Hi, want, want)
	}
}

func TestReciprocalRangeConservative_PositiveRange(t *testing.T) {
	r := NewRange(mustRat(2, 1), mustRat(3, 1), true, true)

	got, err := ReciprocalRangeConservative(r)
	if err != nil {
		t.Fatalf("ReciprocalRangeConservative failed: %v", err)
	}

	wantLo := mustRat(1, 3)
	wantHi := mustRat(1, 2)
	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v] want [%v,%v]", got.Lo, got.Hi, wantLo, wantHi)
	}
}

func TestReciprocalRangeConservative_RejectsOutsideRange(t *testing.T) {
	r := NewRange(mustRat(2, 1), mustRat(1, 1), true, true)

	_, err := ReciprocalRangeConservative(r)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestReciprocalRangeConservative_RejectsZeroOrNegative(t *testing.T) {
	tests := []Range{
		NewRange(mustRat(0, 1), mustRat(1, 1), true, true),
		NewRange(mustRat(-1, 1), mustRat(1, 1), true, true),
		NewRange(mustRat(1, 1), mustRat(0, 1), true, true),
	}

	for i, r := range tests {
		_, err := ReciprocalRangeConservative(r)
		if err == nil {
			t.Fatalf("case %d: expected error for range %v", i, r)
		}
	}
}

// reciprocal_range_conservative_test.go v1
