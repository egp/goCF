// trig_gcf_exact_tail_api_test.go v2
package cf

import (
	"strings"
	"testing"
)

func TestSinBoundsDegreesFromGCFWithTail2_69Angle(t *testing.T) {
	got, err := SinBoundsDegreesFromGCFWithTail2(
		MVP69DegreeGCFSource(),
		MVP69DegreeTail(),
		1,
	)
	if err != nil {
		t.Fatalf("SinBoundsDegreesFromGCFWithTail2 failed: %v", err)
	}

	want := NewRange(mustRat(14, 15), mustRat(131, 140), true, true)
	if got.Lo.Cmp(want.Lo) != 0 || got.Hi.Cmp(want.Hi) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSinApproxDegreesFromGCFWithTail2_69IsStillBoundedNonPoint(t *testing.T) {
	_, err := SinApproxDegreesFromGCFWithTail2(
		MVP69DegreeGCFSource(),
		MVP69DegreeTail(),
		1,
	)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "bounded non-point result") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// trig_gcf_exact_tail_api_test.go v2
