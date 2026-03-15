// mvp_target_formula_test.go v4
package cf

import (
	"strings"
	"testing"
)

func TestMVPTargetFormula_CurrentShape_AssemblesNumeratorAndDenominator(t *testing.T) {
	num, err := MVPNumeratorApproxDefault(4, 6)
	if err != nil {
		t.Fatalf("MVPNumeratorApproxDefault failed: %v", err)
	}

	den, err := MVPDenominatorBoundsDefault()
	if err != nil {
		t.Fatalf("MVPDenominatorBoundsDefault failed: %v", err)
	}

	if num.Cmp(intRat(0)) <= 0 {
		t.Fatalf("numerator got %v want positive", num)
	}
	if !den.IsInside() {
		t.Fatalf("denominator got %v want inside range", den)
	}
	if den.Contains(intRat(0)) {
		t.Fatalf("denominator got %v want zero excluded", den)
	}
}

func TestMVPTargetFormula_DenominatorNowExcludesZero(t *testing.T) {
	den, err := MVPDenominatorBoundsDefault()
	if err != nil {
		t.Fatalf("MVPDenominatorBoundsDefault failed: %v", err)
	}

	if den.Contains(intRat(0)) {
		t.Fatalf("denominator got %v want zero excluded", den)
	}
}

func TestMVPTargetFormula_CurrentNumeratorAndDenominatorSanity(t *testing.T) {
	num, err := MVPNumeratorApproxDefault(4, 6)
	if err != nil {
		t.Fatalf("MVPNumeratorApproxDefault failed: %v", err)
	}

	den, err := MVPDenominatorBoundsDefault()
	if err != nil {
		t.Fatalf("MVPDenominatorBoundsDefault failed: %v", err)
	}

	if num.Cmp(intRat(1)) <= 0 {
		t.Fatalf("numerator got %v want > 1", num)
	}

	wantDen := NewRange(mustRat(29, 1540), mustRat(1, 15), true, true)
	if den.Lo.Cmp(wantDen.Lo) != 0 || den.Hi.Cmp(wantDen.Hi) != 0 {
		t.Fatalf("denominator got %v want %v", den, wantDen)
	}
}

func TestMVPTargetBoundsDefault_IsInsideAndPositive(t *testing.T) {
	got, err := MVPTargetBoundsDefault()
	if err != nil {
		t.Fatalf("MVPTargetBoundsDefault failed: %v", err)
	}

	if !got.IsInside() {
		t.Fatalf("target range got %v want inside", got)
	}
	if got.Lo.Cmp(intRat(0)) <= 0 {
		t.Fatalf("target range got %v want strictly positive lower bound", got)
	}
}

func TestMVPTargetBoundsDefault_MatchesNumeratorOverDenominatorConstruction(t *testing.T) {
	got, err := MVPTargetBoundsDefault()
	if err != nil {
		t.Fatalf("MVPTargetBoundsDefault failed: %v", err)
	}

	num, err := MVPNumeratorApproxDefault(4, 6)
	if err != nil {
		t.Fatalf("MVPNumeratorApproxDefault failed: %v", err)
	}
	den, err := MVPDenominatorBoundsDefault()
	if err != nil {
		t.Fatalf("MVPDenominatorBoundsDefault failed: %v", err)
	}

	recipDen, err := ReciprocalRangeConservative(den)
	if err != nil {
		t.Fatalf("ReciprocalRangeConservative failed: %v", err)
	}

	wantLo, err := num.Mul(recipDen.Lo)
	if err != nil {
		t.Fatalf("Mul lo failed: %v", err)
	}
	wantHi, err := num.Mul(recipDen.Hi)
	if err != nil {
		t.Fatalf("Mul hi failed: %v", err)
	}

	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got %v want [%v,%v]", got, wantLo, wantHi)
	}
}

func TestMVPTargetApproxDefault_CurrentlyReportsBoundedNonPoint(t *testing.T) {
	_, err := MVPTargetApproxDefault()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "bounded non-point result") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Full target formula intentionally lives in test code only for now:
//
//	sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))
//
// mvp_target_formula_test.go v4
