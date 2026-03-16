// mvp_denominator_test.go v7
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

	want := NewRange(mustRat(11, 280), mustRat(7, 150), true, true)
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

	want := NewRange(mustRat(11, 280), mustRat(7, 150), true, true)
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

func TestMVPDenominatorBounds_UsesExactFiniteGCFPrefixFor69Degrees(t *testing.T) {
	got, err := SinBoundsDegreesFromGCFPrefix2(
		MVP69DegreeGCFSource(),
		2,
	)
	if err != nil {
		t.Fatalf("SinBoundsDegreesFromGCFPrefix2 failed: %v", err)
	}

	want, err := SinBoundsDegrees(Degrees(mustRat(69, 1)))
	if err != nil {
		t.Fatalf("SinBoundsDegrees failed: %v", err)
	}

	if got.Lo.Cmp(want.Lo) != 0 || got.Hi.Cmp(want.Hi) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPDenominatorBounds_NoLongerNeedsExactTailTrigEntryPointFor69Degrees(t *testing.T) {
	got, err := MVPDenominatorBoundsDefault()
	if err != nil {
		t.Fatalf("MVPDenominatorBoundsDefault failed: %v", err)
	}

	want := NewRange(mustRat(11, 280), mustRat(7, 150), true, true)
	if got.Lo.Cmp(want.Lo) != 0 || got.Hi.Cmp(want.Hi) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVP69DegreeGCFSource_Prefix2EvaluatesExactlyTo69(t *testing.T) {
	got, err := GCFSourceConvergent(MVP69DegreeGCFSource(), 2)
	if err != nil {
		t.Fatalf("GCFSourceConvergent failed: %v", err)
	}

	want := mustRat(69, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSinBoundsDegreesFromGCFPrefix2_MVP69DegreeGCFSource_IsExactAtPrefix2(t *testing.T) {
	got, err := SinBoundsDegreesFromGCFPrefix2(MVP69DegreeGCFSource(), 2)
	if err != nil {
		t.Fatalf("SinBoundsDegreesFromGCFPrefix2 failed: %v", err)
	}

	want, err := SinBoundsDegrees(Degrees(mustRat(69, 1)))
	if err != nil {
		t.Fatalf("SinBoundsDegrees failed: %v", err)
	}

	if got.Lo.Cmp(want.Lo) != 0 || got.Hi.Cmp(want.Hi) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

// mvp_denominator_test.go v7
