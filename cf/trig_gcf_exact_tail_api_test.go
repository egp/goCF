// trig_gcf_exact_tail_api_test.go v3
package cf

import (
	"strings"
	"testing"
)

func TestSinBoundsDegreesFromGCFWithTail2_ArbitraryExactAngleStillWorks(t *testing.T) {
	got, err := SinBoundsDegreesFromGCFWithTail2(
		NewSliceGCF([2]int64{29, 1}),
		mustRat(1, 1),
		1,
	)
	if err != nil {
		t.Fatalf("SinBoundsDegreesFromGCFWithTail2 failed: %v", err)
	}

	want := NewRange(mustRat(1, 2), mustRat(1, 2), true, true)
	if got.Lo.Cmp(want.Lo) != 0 || got.Hi.Cmp(want.Hi) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSinApproxDegreesFromGCFWithTail2_BoundedNonPointStillReportsAsSuch(t *testing.T) {
	_, err := SinApproxDegreesFromGCFWithTail2(
		NewSliceGCF([2]int64{68, 1}),
		mustRat(1, 1),
		1,
	)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "bounded non-point result") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// trig_gcf_exact_tail_api_test.go v3
