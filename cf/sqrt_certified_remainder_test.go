// sqrt_certified_remainder_test.go v1
package cf

import "testing"

func TestCertifiedRemainderRange_Sqrt2FirstDigit(t *testing.T) {
	// A conservative enclosure for sqrt(2) with certified first digit 1.
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

// sqrt_certified_remainder_test.go v1
