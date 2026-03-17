// mvp_target_formula_test.go v9
package cf

import (
	"fmt"
	"testing"
)

func mvpTestTargetBounds(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	sqrtPolicy SqrtPolicy2,
	angle Angle,
) (Range, error) {
	num, err := MVPRadicandRootValue(
		fourOverPiPrefixTerms,
		ePrefixTerms,
		sqrtPolicy,
	)
	if err != nil {
		return Range{}, err
	}

	den, err := mvpTestDenominatorBounds(sqrtPolicy, angle)
	if err != nil {
		return Range{}, err
	}
	if den.Contains(intRat(0)) {
		return Range{}, fmt.Errorf("mvpTestTargetBounds: denominator range contains 0: %v", den)
	}
	if den.Lo.Cmp(intRat(0)) <= 0 || den.Hi.Cmp(intRat(0)) <= 0 {
		return Range{}, fmt.Errorf("mvpTestTargetBounds: denominator range must be strictly positive: %v", den)
	}

	recipDen, err := ReciprocalRangeConservative(den)
	if err != nil {
		return Range{}, err
	}

	lo, err := num.Mul(recipDen.Lo)
	if err != nil {
		return Range{}, err
	}
	hi, err := num.Mul(recipDen.Hi)
	if err != nil {
		return Range{}, err
	}

	return NewRange(lo, hi, true, true), nil
}

func mvpTestTargetBoundsDefault() (Range, error) {
	return mvpTestTargetBounds(
		MVPRadicandDefaultFourOverPiPrefixTerms,
		MVPRadicandDefaultEPrefixTerms,
		DefaultSqrtPolicy2(),
		Degrees(mustRat(69, 1)),
	)
}

func TestMVPTargetFormula_CurrentShape_AssemblesNumeratorAndDenominator(t *testing.T) {
	num, err := MVPRadicandRootValueCurrentDefault()
	if err != nil {
		t.Fatalf("MVPRadicandRootValueCurrentDefault failed: %v", err)
	}

	den, err := mvpTestDenominatorBoundsDefault()
	if err != nil {
		t.Fatalf("mvpTestDenominatorBoundsDefault failed: %v", err)
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
	den, err := mvpTestDenominatorBoundsDefault()
	if err != nil {
		t.Fatalf("mvpTestDenominatorBoundsDefault failed: %v", err)
	}

	if den.Contains(intRat(0)) {
		t.Fatalf("denominator got %v want zero excluded", den)
	}
}

func TestMVPTargetFormula_CurrentNumeratorAndDenominatorSanity(t *testing.T) {
	num, err := MVPRadicandRootValueCurrentDefault()
	if err != nil {
		t.Fatalf("MVPRadicandRootValueCurrentDefault failed: %v", err)
	}

	den, err := mvpTestDenominatorBoundsDefault()
	if err != nil {
		t.Fatalf("mvpTestDenominatorBoundsDefault failed: %v", err)
	}

	if num.Cmp(intRat(1)) <= 0 {
		t.Fatalf("numerator got %v want > 1", num)
	}

	wantDen := NewRange(mustRat(11, 280), mustRat(7, 150), true, true)
	if den.Lo.Cmp(wantDen.Lo) != 0 || den.Hi.Cmp(wantDen.Hi) != 0 {
		t.Fatalf("denominator got %v want %v", den, wantDen)
	}
}

func TestMVPTargetBoundsDefault_IsInsideAndPositive(t *testing.T) {
	got, err := mvpTestTargetBoundsDefault()
	if err != nil {
		t.Fatalf("mvpTestTargetBoundsDefault failed: %v", err)
	}

	if !got.IsInside() {
		t.Fatalf("target range got %v want inside", got)
	}
	if got.Lo.Cmp(intRat(0)) <= 0 {
		t.Fatalf("target range got %v want strictly positive lower bound", got)
	}
}

func TestMVPTargetBoundsDefault_UsesCurrentSharperNumeratorBudgets(t *testing.T) {
	got, err := mvpTestTargetBoundsDefault()
	if err != nil {
		t.Fatalf("mvpTestTargetBoundsDefault failed: %v", err)
	}

	want, err := mvpTestTargetBounds(
		MVPRadicandDefaultFourOverPiPrefixTerms,
		MVPRadicandDefaultEPrefixTerms,
		DefaultSqrtPolicy2(),
		Degrees(mustRat(69, 1)),
	)
	if err != nil {
		t.Fatalf("mvpTestTargetBounds failed: %v", err)
	}

	if got.Lo.Cmp(want.Lo) != 0 || got.Hi.Cmp(want.Hi) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPTargetBoundsDefault_MatchesNumeratorOverDenominatorConstruction(t *testing.T) {
	got, err := mvpTestTargetBoundsDefault()
	if err != nil {
		t.Fatalf("mvpTestTargetBoundsDefault failed: %v", err)
	}

	num, err := MVPRadicandRootValueCurrentDefault()
	if err != nil {
		t.Fatalf("MVPRadicandRootValueCurrentDefault failed: %v", err)
	}
	den, err := mvpTestDenominatorBoundsDefault()
	if err != nil {
		t.Fatalf("mvpTestDenominatorBoundsDefault failed: %v", err)
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

func TestMVPTargetBounds_SharperNumeratorBudgetsRemainInsideAndPositive(t *testing.T) {
	got, err := mvpTestTargetBounds(
		8,
		10,
		DefaultSqrtPolicy2(),
		Degrees(mustRat(69, 1)),
	)
	if err != nil {
		t.Fatalf("mvpTestTargetBounds failed: %v", err)
	}

	if !got.IsInside() {
		t.Fatalf("got %v want inside", got)
	}
	if got.Lo.Cmp(intRat(0)) <= 0 {
		t.Fatalf("got %v want strictly positive lower bound", got)
	}
}

func TestMVPTargetBounds_CurrentAndSharperBudgetsOverlap(t *testing.T) {
	current, err := mvpTestTargetBoundsDefault()
	if err != nil {
		t.Fatalf("mvpTestTargetBoundsDefault failed: %v", err)
	}

	sharper, err := mvpTestTargetBounds(
		8,
		10,
		DefaultSqrtPolicy2(),
		Degrees(mustRat(69, 1)),
	)
	if err != nil {
		t.Fatalf("mvpTestTargetBounds failed: %v", err)
	}

	if current.Hi.Cmp(sharper.Lo) < 0 || sharper.Hi.Cmp(current.Lo) < 0 {
		t.Fatalf("ranges do not overlap: current=%v sharper=%v", current, sharper)
	}
}

func TestMVPTargetBounds_SharperNumeratorBudgetsDoNotWidenRange(t *testing.T) {
	current, err := mvpTestTargetBoundsDefault()
	if err != nil {
		t.Fatalf("mvpTestTargetBoundsDefault failed: %v", err)
	}

	sharper, err := mvpTestTargetBounds(
		8,
		10,
		DefaultSqrtPolicy2(),
		Degrees(mustRat(69, 1)),
	)
	if err != nil {
		t.Fatalf("mvpTestTargetBounds failed: %v", err)
	}

	currentWidth, err := current.Hi.Sub(current.Lo)
	if err != nil {
		t.Fatalf("current width failed: %v", err)
	}
	sharperWidth, err := sharper.Hi.Sub(sharper.Lo)
	if err != nil {
		t.Fatalf("sharper width failed: %v", err)
	}

	if sharperWidth.Cmp(currentWidth) > 0 {
		t.Fatalf(
			"sharper numerator budgets widened target range: current=%v sharper=%v currentWidth=%v sharperWidth=%v",
			current, sharper, currentWidth, sharperWidth,
		)
	}
}

func TestMVPTargetBounds_CurrentBridgeBudgetIsStable(t *testing.T) {
	got, err := mvpTestTargetBounds(
		MVPRadicandDefaultFourOverPiPrefixTerms,
		MVPRadicandDefaultEPrefixTerms,
		DefaultSqrtPolicy2(),
		Degrees(mustRat(69, 1)),
	)
	if err != nil {
		t.Fatalf("mvpTestTargetBounds current failed: %v", err)
	}

	want, err := mvpTestTargetBounds(
		MVPRadicandDefaultFourOverPiPrefixTerms,
		MVPRadicandDefaultEPrefixTerms,
		DefaultSqrtPolicy2(),
		Degrees(mustRat(69, 1)),
	)
	if err != nil {
		t.Fatalf("mvpTestTargetBounds want failed: %v", err)
	}

	if got.Lo.Cmp(want.Lo) != 0 || got.Hi.Cmp(want.Hi) != 0 {
		t.Fatalf("target range not stable across bridge budgets: got=%v want=%v", got, want)
	}
}

// Full target formula intentionally lives in test code only for now:
//
//	sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))
//
// mvp_target_formula_test.go v9
