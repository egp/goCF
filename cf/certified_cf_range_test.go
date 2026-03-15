// certified_cf_range_test.go v2
package cf

import "testing"

func TestCertifyCFDigitsFromRange_Sqrt2LikeRange_CertifiesPrefix(t *testing.T) {
	// Conservative enclosure for sqrt(2) tight enough to certify a useful prefix.
	r := NewRange(mustRat(181, 128), mustRat(362, 255), true, true)

	got, rest, err := CertifyCFDigitsFromRange(r, 8)
	if err != nil {
		t.Fatalf("CertifyCFDigitsFromRange failed: %v", err)
	}

	wantPrefix := []int64{1, 2}
	if len(got) < len(wantPrefix) {
		t.Fatalf("len(got)=%d want at least %d got=%v", len(got), len(wantPrefix), got)
	}
	for i := range wantPrefix {
		if got[i] != wantPrefix[i] {
			t.Fatalf("got[%d]=%d want=%d full=%v", i, got[i], wantPrefix[i], got)
		}
	}

	if !rest.IsInside() {
		t.Fatalf("expected inside remainder range, got %v", rest)
	}
}

func TestCertifyCFDigitsFromRange_ExactInteger_CertifiesOneDigit(t *testing.T) {
	r := NewRange(mustRat(2, 1), mustRat(2, 1), true, true)

	got, _, err := CertifyCFDigitsFromRange(r, 8)
	if err != nil {
		t.Fatalf("CertifyCFDigitsFromRange failed: %v", err)
	}

	want := []int64{2}
	if len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestCertifyCFDigitsFromRange_UncertifiedRangeYieldsNoDigits(t *testing.T) {
	r := NewRange(mustRat(1, 1), mustRat(2, 1), true, true)

	got, rest, err := CertifyCFDigitsFromRange(r, 8)
	if err != nil {
		t.Fatalf("CertifyCFDigitsFromRange failed: %v", err)
	}

	if len(got) != 0 {
		t.Fatalf("expected no digits, got %v", got)
	}
	if rest.Lo.Cmp(r.Lo) != 0 || rest.Hi.Cmp(r.Hi) != 0 {
		t.Fatalf("expected unchanged rest range, got %v want %v", rest, r)
	}
}

func TestCertifyCFDigitsFromRange_RejectsBadArgs(t *testing.T) {
	r := NewRange(mustRat(2, 1), mustRat(1, 1), true, true)

	_, _, err := CertifyCFDigitsFromRange(r, 8)
	if err == nil {
		t.Fatalf("expected error for outside range")
	}

	r2 := NewRange(mustRat(1, 1), mustRat(2, 1), true, true)
	_, _, err = CertifyCFDigitsFromRange(r2, -1)
	if err == nil {
		t.Fatalf("expected error for negative maxDigits")
	}
}

// certified_cf_range_test.go v2
