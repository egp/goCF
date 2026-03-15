// mvp_target_formula_test.go v1
package cf

import "testing"

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
}

func TestMVPTargetFormula_CurrentDenominatorBoundContainsZero(t *testing.T) {
	den, err := MVPDenominatorBoundsDefault()
	if err != nil {
		t.Fatalf("MVPDenominatorBoundsDefault failed: %v", err)
	}

	if !den.Contains(intRat(0)) {
		t.Fatalf("denominator got %v want range containing 0 at current MVP stage", den)
	}
}

func TestMVPTargetFormula_CurrentCertifiedQuotientIsNotYetAvailable(t *testing.T) {
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
	if !den.Contains(intRat(0)) {
		t.Fatalf("precondition changed: denominator no longer contains 0: %v", den)
	}

	// The full target formula is:
	//
	//   sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))
	//
	// At the current MVP stage, numerator is a bounded rational point approximation,
	// but the denominator certified range still crosses 0, so a certified quotient
	// is not yet available.
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

	wantDen := NewRange(mustRat(-1, 1), mustRat(1, 2), true, true)
	if den.Lo.Cmp(wantDen.Lo) != 0 || den.Hi.Cmp(wantDen.Hi) != 0 {
		t.Fatalf("denominator got %v want %v", den, wantDen)
	}
}

// Full target formula intentionally lives in test code only for now:
//
//	sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))
//
// mvp_target_formula_test.go v1
