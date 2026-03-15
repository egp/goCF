// sqrt_certified_remainder_test.go v2
package cf

import "testing"

func TestShiftRangeByInt(t *testing.T) {
	r := NewRange(mustRat(3, 2), mustRat(7, 4), true, true)

	got, err := ShiftRangeByInt(r, 1)
	if err != nil {
		t.Fatalf("ShiftRangeByInt failed: %v", err)
	}

	wantLo := mustRat(1, 2)
	wantHi := mustRat(3, 4)
	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v] want [%v,%v]", got.Lo, got.Hi, wantLo, wantHi)
	}
}

func TestCertifiedRemainderRange_Sqrt2FirstDigit(t *testing.T) {
	// Conservative enclosure for sqrt(2) with certified first digit 1.
	r := NewRange(mustRat(181, 128), mustRat(362, 255), true, true)

	got, err := CertifiedRemainderRange(r, 1)
	if err != nil {
		t.Fatalf("CertifiedRemainderRange failed: %v", err)
	}

	lo, hi, err := got.FloorBounds()
	if err != nil {
		t.Fatalf("FloorBounds failed: %v", err)
	}

	// remainder for sqrt(2) is 1/(sqrt(2)-1) = sqrt(2)+1, whose floor is 2.
	if lo != 2 || hi != 2 {
		t.Fatalf("got floor bounds (%d,%d) want (2,2); range=%v", lo, hi, got)
	}
}

func TestCertifiedRemainderRange_RejectsUncertifiedDigit(t *testing.T) {
	r := NewRange(mustRat(1, 1), mustRat(2, 1), true, true)

	_, err := CertifiedRemainderRange(r, 1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestCertifiedRemainderRange_RejectsOutsideRange(t *testing.T) {
	r := NewRange(mustRat(2, 1), mustRat(1, 1), true, true)

	_, err := CertifiedRemainderRange(r, 1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

// sqrt_certified_remainder_test.go v2
