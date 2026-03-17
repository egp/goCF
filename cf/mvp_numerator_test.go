// mvp_numerator_test.go v4
package cf

import "testing"

func TestMVPNumeratorApprox_RejectsBadBounds(t *testing.T) {
	if _, err := MVPNumeratorApproxDefault(0, 6); err == nil {
		t.Fatalf("expected error for fourOverPiPrefixTerms=0")
	}
	if _, err := MVPNumeratorApproxDefault(4, 0); err == nil {
		t.Fatalf("expected error for ePrefixTerms=0")
	}
}

func TestMVPNumeratorApproxFromRadicandApprox_UsesSnapshotUnaryPath(t *testing.T) {
	a, err := MVPNumeratorRadicandApproxSnapshot(4, 6, MVPNumeratorBridgePrefixTerms)
	if err != nil {
		t.Fatalf("MVPNumeratorRadicandApproxSnapshot failed: %v", err)
	}

	got, err := MVPNumeratorApproxFromRadicandApprox(a, DefaultSqrtPolicy2())
	if err != nil {
		t.Fatalf("MVPNumeratorApproxFromRadicandApprox failed: %v", err)
	}

	want, err := MVPNumeratorApproxDefault(4, 6)
	if err != nil {
		t.Fatalf("MVPNumeratorApproxDefault failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPNumeratorApprox_UsesSnapshotAssembledRadicandPath(t *testing.T) {
	got, err := MVPNumeratorApproxDefault(4, 6)
	if err != nil {
		t.Fatalf("MVPNumeratorApproxDefault failed: %v", err)
	}

	a, err := MVPNumeratorRadicandApproxSnapshot(4, 6, MVPNumeratorBridgePrefixTerms)
	if err != nil {
		t.Fatalf("MVPNumeratorRadicandApproxSnapshot failed: %v", err)
	}

	want, err := MVPNumeratorApproxFromRadicandApprox(a, DefaultSqrtPolicy2())
	if err != nil {
		t.Fatalf("MVPNumeratorApproxFromRadicandApprox failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPNumeratorApprox_CurrentDefaultUsesSharperBudgets(t *testing.T) {
	got, err := MVPNumeratorApproxCurrentDefault()
	if err != nil {
		t.Fatalf("MVPNumeratorApproxCurrentDefault failed: %v", err)
	}

	want, err := MVPNumeratorApproxDefault(
		MVPDefaultFourOverPiPrefixTerms,
		MVPDefaultEPrefixTerms,
	)
	if err != nil {
		t.Fatalf("MVPNumeratorApproxDefault failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPNumeratorApprox_UsesExplicitSnapshotBudget(t *testing.T) {
	got, err := MVPNumeratorApproxCurrentDefault()
	if err != nil {
		t.Fatalf("MVPNumeratorApproxCurrentDefault failed: %v", err)
	}

	a, err := MVPNumeratorRadicandApproxSnapshot(
		MVPDefaultFourOverPiPrefixTerms,
		MVPDefaultEPrefixTerms,
		MVPNumeratorBridgePrefixTerms,
	)
	if err != nil {
		t.Fatalf("MVPNumeratorRadicandApproxSnapshot failed: %v", err)
	}

	want, err := MVPNumeratorApproxFromRadicandApprox(a, DefaultSqrtPolicy2())
	if err != nil {
		t.Fatalf("MVPNumeratorApproxFromRadicandApprox failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPNumeratorApproxCF_MatchesApprox(t *testing.T) {
	gotTerms, err := MVPNumeratorApproxTermsDefault(4, 6, 12)
	if err != nil {
		t.Fatalf("MVPNumeratorApproxTermsDefault failed: %v", err)
	}

	gotApprox, err := MVPNumeratorApproxDefault(4, 6)
	if err != nil {
		t.Fatalf("MVPNumeratorApproxDefault failed: %v", err)
	}
	wantTerms := collectTerms(NewRationalCF(gotApprox), 12)

	if len(gotTerms) != len(wantTerms) {
		t.Fatalf("len mismatch: got=%v want=%v", gotTerms, wantTerms)
	}
	for i := range wantTerms {
		if gotTerms[i] != wantTerms[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, gotTerms, wantTerms)
		}
	}
}

func TestMVPNumeratorApprox_IsPositive(t *testing.T) {
	got, err := MVPNumeratorApproxCurrentDefault()
	if err != nil {
		t.Fatalf("MVPNumeratorApproxCurrentDefault failed: %v", err)
	}
	if got.Cmp(intRat(0)) <= 0 {
		t.Fatalf("got %v want positive", got)
	}
}

func TestMVPNumeratorApprox_ExceedsOne(t *testing.T) {
	got, err := MVPNumeratorApproxCurrentDefault()
	if err != nil {
		t.Fatalf("MVPNumeratorApproxCurrentDefault failed: %v", err)
	}
	if got.Cmp(intRat(1)) <= 0 {
		t.Fatalf("got %v want > 1", got)
	}
}

// Full target formula intentionally stays in test code for now.
// This test fixes only the numerator shape:
//
//	sqrt(3/pi^2 + e)
//

func TestMVPNumeratorRadicandApprox_MatchesExistingSubexpression(t *testing.T) {
	got, err := MVPNumeratorRadicandApprox(4, 6)
	if err != nil {
		t.Fatalf("MVPNumeratorRadicandApprox failed: %v", err)
	}

	want, err := MVPThreeOverPiSquaredPlusEApprox(4, 6)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEApprox failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPNumeratorRadicandApproxSnapshot_RoundTripsRadicandApprox(t *testing.T) {
	got, err := MVPNumeratorRadicandApproxSnapshot(4, 6, MVPNumeratorBridgePrefixTerms)
	if err != nil {
		t.Fatalf("MVPNumeratorRadicandApproxSnapshot failed: %v", err)
	}

	want, err := MVPNumeratorRadicandApprox(4, 6)
	if err != nil {
		t.Fatalf("MVPNumeratorRadicandApprox failed: %v", err)
	}

	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

func TestMVPNumeratorRadicandApproxSnapshot_ConvergentIsStableAcrossSnapshotBudgets(t *testing.T) {
	got64, err := MVPNumeratorRadicandApproxSnapshot(4, 6, 64)
	if err != nil {
		t.Fatalf("MVPNumeratorRadicandApproxSnapshot(64) failed: %v", err)
	}

	got96, err := MVPNumeratorRadicandApproxSnapshot(4, 6, 96)
	if err != nil {
		t.Fatalf("MVPNumeratorRadicandApproxSnapshot(96) failed: %v", err)
	}

	if got64.Convergent.Cmp(got96.Convergent) != 0 {
		t.Fatalf("snapshot convergent not stable: got64=%v got96=%v", got64.Convergent, got96.Convergent)
	}
}

func TestMVPNumeratorApprox_CurrentBridgeBudgetIsStable(t *testing.T) {
	got, err := MVPNumeratorApproxWithBridgeTerms(4, 6, DefaultSqrtPolicy2(), 64)
	if err != nil {
		t.Fatalf("MVPNumeratorApproxWithBridgeTerms(64) failed: %v", err)
	}

	want, err := MVPNumeratorApproxWithBridgeTerms(4, 6, DefaultSqrtPolicy2(), 96)
	if err != nil {
		t.Fatalf("MVPNumeratorApproxWithBridgeTerms(96) failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("numerator not stable across bridge budgets: got=%v want=%v", got, want)
	}
}

func TestMVPNumeratorApprox_SharperRadicandBudgetsRemainPositive(t *testing.T) {
	got, err := MVPNumeratorApproxDefault(8, 10)
	if err != nil {
		t.Fatalf("MVPNumeratorApproxDefault failed: %v", err)
	}
	if got.Cmp(intRat(0)) <= 0 {
		t.Fatalf("got %v want positive", got)
	}
}

func TestMVPNumeratorApprox_SharperRadicandBudgetsExceedOne(t *testing.T) {
	got, err := MVPNumeratorApproxDefault(8, 10)
	if err != nil {
		t.Fatalf("MVPNumeratorApproxDefault failed: %v", err)
	}
	if got.Cmp(intRat(1)) <= 0 {
		t.Fatalf("got %v want > 1", got)
	}
}

func TestMVPNumeratorApprox_CurrentAndSharperBudgetsAreDistinctButClose(t *testing.T) {
	current, err := MVPNumeratorApproxDefault(
		MVPDefaultFourOverPiPrefixTerms,
		MVPDefaultEPrefixTerms,
	)
	if err != nil {
		t.Fatalf("MVPNumeratorApproxDefault current failed: %v", err)
	}

	sharper, err := MVPNumeratorApproxDefault(8, 10)
	if err != nil {
		t.Fatalf("MVPNumeratorApproxDefault sharper failed: %v", err)
	}

	if current.Cmp(sharper) == 0 {
		t.Fatalf("expected sharper budgets to change the bounded numerator approximation")
	}

	// Both approximations should still describe the same coarse MVP shape.
	if current.Cmp(intRat(1)) <= 0 || sharper.Cmp(intRat(1)) <= 0 {
		t.Fatalf("current=%v sharper=%v want both > 1", current, sharper)
	}
}

func TestMVPNumeratorRadicandApproxSnapshot_MatchesCanonicalRadicandAssembly(t *testing.T) {
	got, err := MVPNumeratorRadicandApproxSnapshot(4, 6, 64)
	if err != nil {
		t.Fatalf("MVPNumeratorRadicandApproxSnapshot failed: %v", err)
	}

	want, err := MVPThreeOverPiSquaredPlusERadicandApproxSnapshotWithFourOverPiApprox(
		MVPDefaultFourOverPiApproxFunc(),
		4,
		6,
	)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusERadicandApproxSnapshotWithFourOverPiApprox failed: %v", err)
	}

	if got.Convergent.Cmp(want.Convergent) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want.Convergent)
	}
	if got.HasRange() != want.HasRange() {
		t.Fatalf("range presence mismatch: got=%v want=%v", got.HasRange(), want.HasRange())
	}
	if got.HasRange() {
		if got.Range.Lo.Cmp(want.Range.Lo) != 0 || got.Range.Hi.Cmp(want.Range.Hi) != 0 {
			t.Fatalf("got range %v want %v", *got.Range, *want.Range)
		}
	}
}

func TestMVPNumeratorRadicandApproxSnapshot_RoundTripsCurrentRadicand(t *testing.T) {
	got, err := MVPNumeratorRadicandApproxSnapshot(4, 6, MVPNumeratorBridgePrefixTerms)
	if err != nil {
		t.Fatalf("MVPNumeratorRadicandApproxSnapshot failed: %v", err)
	}

	want, err := MVPNumeratorRadicandApprox(4, 6)
	if err != nil {
		t.Fatalf("MVPNumeratorRadicandApprox failed: %v", err)
	}

	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

func TestMVPNumeratorRadicandApproxSnapshot_RejectsBadBridgeTerms(t *testing.T) {
	_, err := MVPNumeratorRadicandApproxSnapshot(4, 6, 0)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMVPNumeratorApproxFromRadicandApprox_MatchesCurrentPath(t *testing.T) {
	a, err := MVPNumeratorRadicandApproxSnapshot(4, 6, MVPNumeratorBridgePrefixTerms)
	if err != nil {
		t.Fatalf("MVPNumeratorRadicandApproxSnapshot failed: %v", err)
	}

	got, err := MVPNumeratorApproxFromRadicandApprox(a, DefaultSqrtPolicy2())
	if err != nil {
		t.Fatalf("MVPNumeratorApproxFromRadicandApprox failed: %v", err)
	}

	want, err := MVPNumeratorApproxDefault(4, 6)
	if err != nil {
		t.Fatalf("MVPNumeratorApproxDefault failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPNumeratorRadicandApproxSnapshot_CurrentBridgeBudgetIsStable(t *testing.T) {
	got, err := MVPNumeratorRadicandApproxSnapshot(4, 6, 64)
	if err != nil {
		t.Fatalf("MVPNumeratorRadicandApproxSnapshot(64) failed: %v", err)
	}

	want, err := MVPNumeratorRadicandApproxSnapshot(4, 6, 96)
	if err != nil {
		t.Fatalf("MVPNumeratorRadicandApproxSnapshot(96) failed: %v", err)
	}

	if got.Convergent.Cmp(want.Convergent) != 0 {
		t.Fatalf("snapshot convergent not stable: got=%v want=%v", got.Convergent, want.Convergent)
	}
}

func TestMVPThreeOverPiSquaredPlusEApproxSnapshot_IsThinLegacyWrapper(t *testing.T) {
	got, err := MVPThreeOverPiSquaredPlusEApproxSnapshot(4, 6, 64)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEApproxSnapshot failed: %v", err)
	}

	want, err := MVPThreeOverPiSquaredPlusERadicandSnapshot(4, 6, 64)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusERadicandSnapshot failed: %v", err)
	}

	if got.Convergent.Cmp(want.Convergent) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want.Convergent)
	}
}

// mvp_numerator_test.go v4
