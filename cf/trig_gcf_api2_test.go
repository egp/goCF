// trig_gcf_api2_test.go v1
package cf

import (
	"strings"
	"testing"
)

func TestCFGCFAdapter_Radicand_ForwardsMetadata(t *testing.T) {
	src := AdaptCFToGCF(Sqrt5CF())

	qr, ok := src.(interface {
		Radicand() (int64, bool)
	})
	if !ok {
		t.Fatalf("expected adapted source to expose Radicand")
	}

	n, ok := qr.Radicand()
	if !ok {
		t.Fatalf("expected Radicand ok=true")
	}
	if n != 5 {
		t.Fatalf("got radicand %d want 5", n)
	}
}

func TestSinBoundsDegreesFromGCFPrefix2_Finite69Angle(t *testing.T) {
	got, err := SinBoundsDegreesFromGCFPrefix2(
		AdaptCFToGCF(NewSliceCF(69)),
		1,
	)
	if err != nil {
		t.Fatalf("SinBoundsDegreesFromGCFPrefix2 failed: %v", err)
	}

	want := NewRange(mustRat(14, 15), mustRat(131, 140), true, true)
	if got.Lo.Cmp(want.Lo) != 0 || got.Hi.Cmp(want.Hi) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSinApproxDegreesFromGCFPrefix2_69IsStillBoundedNonPoint(t *testing.T) {
	_, err := SinApproxDegreesFromGCFPrefix2(
		AdaptCFToGCF(NewSliceCF(69)),
		1,
	)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "bounded non-point result") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTanhBoundsSpecialFromGCF2_AdaptedSqrt5CF(t *testing.T) {
	got, err := TanhBoundsSpecialFromGCF2(AdaptCFToGCF(Sqrt5CF()))
	if err != nil {
		t.Fatalf("TanhBoundsSpecialFromGCF2 failed: %v", err)
	}

	want := NewRange(mustRat(39, 40), mustRat(49, 50), true, true)
	if got.Lo.Cmp(want.Lo) != 0 || got.Hi.Cmp(want.Hi) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestTanhBoundsSpecialFromGCF2_RejectsNonQuadraticMetadata(t *testing.T) {
	_, err := TanhBoundsSpecialFromGCF2(AdaptCFToGCF(NewSliceCF(5)))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "quadratic-radical metadata") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// trig_gcf_api2_test.go v1
