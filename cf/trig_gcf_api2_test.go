// trig_gcf_api2_test.go v5
package cf

import (
	"strings"
	"testing"
)

func TestSinBoundsDegreesFromGCFPrefix2_Infinite69AngleIsNotExactAtPrefix1(t *testing.T) {
	_, err := SinBoundsDegreesFromGCFPrefix2(
		AdaptCFToGCF(NewPeriodicCF([]int64{69}, []int64{1})),
		1,
	)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "angle not exact") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSinApproxDegreesFromGCFPrefix2_69IsStillBoundedNonPoint(t *testing.T) {
	_, err := SinApproxDegreesFromGCFPrefix2(
		AdaptCFToGCF(NewPeriodicCF([]int64{69}, []int64{1})),
		1,
	)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "angle not exact") {
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

// trig_gcf_api2_test.go v5
