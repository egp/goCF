// fingerprint_test.go v1
package cf

import "testing"

func TestFingerprintULFT_CanonicalizesSignAndGCD(t *testing.T) {
	r := NewRange(mustRat(1, 2), mustRat(3, 2), true, false)

	// Same transform up to overall sign and gcd factor.
	t1 := NewULFT(2, 4, 6, 8)     // gcd 2 -> (1,2,3,4)
	t2 := NewULFT(-1, -2, -3, -4) // sign flip

	f1, err := FingerprintULFT(t1, r)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	f2, err := FingerprintULFT(t2, r)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if f1 != f2 {
		t.Fatalf("expected canonical fingerprints to match:\n f1=%q\n f2=%q", f1, f2)
	}
}

func TestFingerprintULFT_ChangesWithRangeFlags(t *testing.T) {
	tform := NewULFT(1, 0, 0, 1)

	r1 := NewRange(mustRat(0, 1), mustRat(1, 1), true, true)
	r2 := NewRange(mustRat(0, 1), mustRat(1, 1), false, true)

	f1, err := FingerprintULFT(tform, r1)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	f2, err := FingerprintULFT(tform, r2)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if f1 == f2 {
		t.Fatalf("expected different fingerprints for different endpoint flags:\n f1=%q\n f2=%q", f1, f2)
	}
}

func TestFingerprintBLFT_CanonicalizesSignAndGCD(t *testing.T) {
	rx := NewRange(mustRat(1, 2), mustRat(3, 2), true, true)
	ry := NewRange(mustRat(2, 3), mustRat(5, 3), true, true)

	// Same BLFT up to gcd/sign.
	b1 := NewBLFT(2, 4, 6, 8, 10, 12, 14, 16)     // gcd 2 -> (1,2,3,4,5,6,7,8)
	b2 := NewBLFT(-1, -2, -3, -4, -5, -6, -7, -8) // sign flip

	f1, err := FingerprintBLFT(b1, rx, ry)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	f2, err := FingerprintBLFT(b2, rx, ry)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if f1 != f2 {
		t.Fatalf("expected canonical fingerprints to match:\n f1=%q\n f2=%q", f1, f2)
	}
}

func TestFingerprintBLFT_ChangesWithInputRanges(t *testing.T) {
	rx := NewRange(mustRat(0, 1), mustRat(1, 1), true, true)
	ry1 := NewRange(mustRat(0, 1), mustRat(1, 1), true, true)
	ry2 := NewRange(mustRat(0, 1), mustRat(2, 1), true, true)

	tform := NewBLFT(1, 0, 0, 0, 0, 0, 0, 1)

	f1, err := FingerprintBLFT(tform, rx, ry1)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	f2, err := FingerprintBLFT(tform, rx, ry2)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if f1 == f2 {
		t.Fatalf("expected different fingerprints for different ranges:\n f1=%q\n f2=%q", f1, f2)
	}
}

// fingerprint_test.go v1
