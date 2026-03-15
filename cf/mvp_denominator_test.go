// mvp_denominator_test.go v6
package cf

import (
	"strings"
	"testing"
)

func TestMVPDenominatorBoundsDefault_UsesDegreesByDefault(t *testing.T) {
	got, err := MVPDenominatorBoundsDefault()
	if err != nil {
		t.Fatalf("MVPDenominatorBoundsDefault failed: %v", err)
	}

	want := NewRange(mustRat(29, 1540), mustRat(1, 15), true, true)
	if got.Lo.Cmp(want.Lo) != 0 || got.Hi.Cmp(want.Hi) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
	if got.Contains(intRat(0)) {
		t.Fatalf("got %v want zero excluded", got)
	}
}

func TestMVPDenominatorBounds_RejectsRadiansForMVP(t *testing.T) {
	_, err := MVPDenominatorBounds(
		DefaultSqrtPolicy2(),
		Radians(mustRat(69, 1)),
	)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "degrees") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMVPDenominatorBounds_Accepts69DegreeBoundAndExcludesZero(t *testing.T) {
	got, err := MVPDenominatorBounds(
		DefaultSqrtPolicy2(),
		Degrees(mustRat(69, 1)),
	)
	if err != nil {
		t.Fatalf("MVPDenominatorBounds failed: %v", err)
	}

	want := NewRange(mustRat(29, 1540), mustRat(1, 15), true, true)
	if got.Lo.Cmp(want.Lo) != 0 || got.Hi.Cmp(want.Hi) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
	if got.Contains(intRat(0)) {
		t.Fatalf("got %v want zero excluded", got)
	}
}

func TestMVPDenominatorApprox_CurrentlyReportsBoundedNonPoint(t *testing.T) {
	_, err := MVPDenominatorApproxDefault()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "bounded non-point result") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Full target formula intentionally remains in test code for now.
// This test fixes only the denominator shape:
//
//	tanh(sqrt(5)) - sin(69°)
//
// mvp_denominator_test.go v6
