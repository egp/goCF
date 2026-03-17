// mvp_numerator_test.go v6
package cf

import "testing"

func TestMVPRadicandRootValue_RejectsBadBounds(t *testing.T) {
	if _, err := MVPRadicandRootValueDefault(0, 6); err == nil {
		t.Fatalf("expected error for fourOverPiPrefixTerms=0")
	}
	if _, err := MVPRadicandRootValueDefault(4, 0); err == nil {
		t.Fatalf("expected error for ePrefixTerms=0")
	}
}

func TestMVPRadicandRootValueFromSnapshot_UsesSnapshotUnaryPath(t *testing.T) {
	a, err := MVPRadicandSnapshot(4, 6, MVPRadicandSnapshotTerms)
	if err != nil {
		t.Fatalf("MVPRadicandSnapshot failed: %v", err)
	}

	got, err := MVPRadicandRootValueFromSnapshot(a, DefaultSqrtPolicy2())
	if err != nil {
		t.Fatalf("MVPRadicandRootValueFromSnapshot failed: %v", err)
	}

	want, err := MVPRadicandRootValueDefault(4, 6)
	if err != nil {
		t.Fatalf("MVPRadicandRootValueDefault failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPRadicandRootValue_UsesSnapshotAssembledRadicandPath(t *testing.T) {
	got, err := MVPRadicandRootValueDefault(4, 6)
	if err != nil {
		t.Fatalf("MVPRadicandRootValueDefault failed: %v", err)
	}

	a, err := MVPRadicandSnapshot(4, 6, MVPRadicandSnapshotTerms)
	if err != nil {
		t.Fatalf("MVPRadicandSnapshot failed: %v", err)
	}

	want, err := MVPRadicandRootValueFromSnapshot(a, DefaultSqrtPolicy2())
	if err != nil {
		t.Fatalf("MVPRadicandRootValueFromSnapshot failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPRadicandRootValue_CurrentDefaultUsesSharperBudgets(t *testing.T) {
	got, err := MVPRadicandRootValueCurrentDefault()
	if err != nil {
		t.Fatalf("MVPRadicandRootValueCurrentDefault failed: %v", err)
	}

	want, err := MVPRadicandRootValueDefault(
		MVPRadicandDefaultFourOverPiPrefixTerms,
		MVPRadicandDefaultEPrefixTerms,
	)
	if err != nil {
		t.Fatalf("MVPRadicandRootValueDefault failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPRadicandRootValue_UsesExplicitSnapshotBudget(t *testing.T) {
	got, err := MVPRadicandRootValueCurrentDefault()
	if err != nil {
		t.Fatalf("MVPRadicandRootValueCurrentDefault failed: %v", err)
	}

	a, err := MVPRadicandSnapshot(
		MVPRadicandDefaultFourOverPiPrefixTerms,
		MVPRadicandDefaultEPrefixTerms,
		MVPRadicandSnapshotTerms,
	)
	if err != nil {
		t.Fatalf("MVPRadicandSnapshot failed: %v", err)
	}

	want, err := MVPRadicandRootValueFromSnapshot(a, DefaultSqrtPolicy2())
	if err != nil {
		t.Fatalf("MVPRadicandRootValueFromSnapshot failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPRadicandRootCF_MatchesValue(t *testing.T) {
	gotTerms, err := MVPRadicandRootTermsDefault(4, 6, 12)
	if err != nil {
		t.Fatalf("MVPRadicandRootTermsDefault failed: %v", err)
	}

	gotValue, err := MVPRadicandRootValueDefault(4, 6)
	if err != nil {
		t.Fatalf("MVPRadicandRootValueDefault failed: %v", err)
	}
	wantTerms := collectTerms(NewRationalCF(gotValue), 12)

	if len(gotTerms) != len(wantTerms) {
		t.Fatalf("len mismatch: got=%v want=%v", gotTerms, wantTerms)
	}
	for i := range wantTerms {
		if gotTerms[i] != wantTerms[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, gotTerms, wantTerms)
		}
	}
}

func TestMVPRadicandRootValue_IsPositive(t *testing.T) {
	got, err := MVPRadicandRootValueCurrentDefault()
	if err != nil {
		t.Fatalf("MVPRadicandRootValueCurrentDefault failed: %v", err)
	}
	if got.Cmp(intRat(0)) <= 0 {
		t.Fatalf("got %v want positive", got)
	}
}

func TestMVPRadicandRootValue_ExceedsOne(t *testing.T) {
	got, err := MVPRadicandRootValueCurrentDefault()
	if err != nil {
		t.Fatalf("MVPRadicandRootValueCurrentDefault failed: %v", err)
	}
	if got.Cmp(intRat(1)) <= 0 {
		t.Fatalf("got %v want > 1", got)
	}
}

func TestMVPRadicandConvergent_MatchesGenericAssemblyConvergent(t *testing.T) {
	got, err := MVPRadicandConvergent(4, 6)
	if err != nil {
		t.Fatalf("MVPRadicandConvergent failed: %v", err)
	}

	want, err := MVPRadicandAssembleConvergent(4, 6)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleConvergent failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPRadicandSnapshot_RoundTripsRadicandConvergent(t *testing.T) {
	got, err := MVPRadicandSnapshot(4, 6, MVPRadicandSnapshotTerms)
	if err != nil {
		t.Fatalf("MVPRadicandSnapshot failed: %v", err)
	}

	want, err := MVPRadicandConvergent(4, 6)
	if err != nil {
		t.Fatalf("MVPRadicandConvergent failed: %v", err)
	}

	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

func TestMVPRadicandSnapshot_ConvergentIsStableAcrossSnapshotBudgets(t *testing.T) {
	got64, err := MVPRadicandSnapshot(4, 6, 64)
	if err != nil {
		t.Fatalf("MVPRadicandSnapshot(64) failed: %v", err)
	}

	got96, err := MVPRadicandSnapshot(4, 6, 96)
	if err != nil {
		t.Fatalf("MVPRadicandSnapshot(96) failed: %v", err)
	}

	if got64.Convergent.Cmp(got96.Convergent) != 0 {
		t.Fatalf("snapshot convergent not stable: got64=%v got96=%v", got64.Convergent, got96.Convergent)
	}
}

func TestMVPRadicandRootValue_CurrentSnapshotBudgetIsStable(t *testing.T) {
	got, err := MVPRadicandRootValueWithSnapshotTerms(4, 6, DefaultSqrtPolicy2(), 64)
	if err != nil {
		t.Fatalf("MVPRadicandRootValueWithSnapshotTerms(64) failed: %v", err)
	}

	want, err := MVPRadicandRootValueWithSnapshotTerms(4, 6, DefaultSqrtPolicy2(), 96)
	if err != nil {
		t.Fatalf("MVPRadicandRootValueWithSnapshotTerms(96) failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("root value not stable across snapshot budgets: got=%v want=%v", got, want)
	}
}

func TestMVPRadicandRootValue_SharperBudgetsRemainPositive(t *testing.T) {
	got, err := MVPRadicandRootValueDefault(8, 10)
	if err != nil {
		t.Fatalf("MVPRadicandRootValueDefault failed: %v", err)
	}
	if got.Cmp(intRat(0)) <= 0 {
		t.Fatalf("got %v want positive", got)
	}
}

func TestMVPRadicandRootValue_SharperBudgetsExceedOne(t *testing.T) {
	got, err := MVPRadicandRootValueDefault(8, 10)
	if err != nil {
		t.Fatalf("MVPRadicandRootValueDefault failed: %v", err)
	}
	if got.Cmp(intRat(1)) <= 0 {
		t.Fatalf("got %v want > 1", got)
	}
}

func TestMVPRadicandRootValue_CurrentAndSharperBudgetsAreDistinctButClose(t *testing.T) {
	current, err := MVPRadicandRootValueDefault(
		MVPRadicandDefaultFourOverPiPrefixTerms,
		MVPRadicandDefaultEPrefixTerms,
	)
	if err != nil {
		t.Fatalf("MVPRadicandRootValueDefault current failed: %v", err)
	}

	sharper, err := MVPRadicandRootValueDefault(8, 10)
	if err != nil {
		t.Fatalf("MVPRadicandRootValueDefault sharper failed: %v", err)
	}

	if current.Cmp(sharper) == 0 {
		t.Fatalf("expected sharper budgets to change the bounded rooted-radicand value")
	}

	if current.Cmp(intRat(1)) <= 0 || sharper.Cmp(intRat(1)) <= 0 {
		t.Fatalf("current=%v sharper=%v want both > 1", current, sharper)
	}
}

func TestMVPRadicandSnapshot_MatchesCanonicalRadicandAssembly(t *testing.T) {
	got, err := MVPRadicandSnapshot(4, 6, 64)
	if err != nil {
		t.Fatalf("MVPRadicandSnapshot failed: %v", err)
	}

	want, err := MVPRadicandAssembleSnapshotWithFourOverPiApprox(
		MVPDefaultFourOverPiApproxFunc(),
		4,
		6,
	)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleSnapshotWithFourOverPiApprox failed: %v", err)
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

func TestMVPRadicandSnapshot_RoundTripsCurrentRadicand(t *testing.T) {
	got, err := MVPRadicandSnapshot(4, 6, MVPRadicandSnapshotTerms)
	if err != nil {
		t.Fatalf("MVPRadicandSnapshot failed: %v", err)
	}

	want, err := MVPRadicandConvergent(4, 6)
	if err != nil {
		t.Fatalf("MVPRadicandConvergent failed: %v", err)
	}

	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

func TestMVPRadicandSnapshot_RejectsBadSnapshotTerms(t *testing.T) {
	_, err := MVPRadicandSnapshot(4, 6, 0)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMVPRadicandRootValueFromSnapshot_MatchesCurrentPath(t *testing.T) {
	a, err := MVPRadicandSnapshot(4, 6, MVPRadicandSnapshotTerms)
	if err != nil {
		t.Fatalf("MVPRadicandSnapshot failed: %v", err)
	}

	got, err := MVPRadicandRootValueFromSnapshot(a, DefaultSqrtPolicy2())
	if err != nil {
		t.Fatalf("MVPRadicandRootValueFromSnapshot failed: %v", err)
	}

	want, err := MVPRadicandRootValueDefault(4, 6)
	if err != nil {
		t.Fatalf("MVPRadicandRootValueDefault failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPRadicandSnapshot_CurrentSnapshotBudgetIsStable(t *testing.T) {
	got, err := MVPRadicandSnapshot(4, 6, 64)
	if err != nil {
		t.Fatalf("MVPRadicandSnapshot(64) failed: %v", err)
	}

	want, err := MVPRadicandSnapshot(4, 6, 96)
	if err != nil {
		t.Fatalf("MVPRadicandSnapshot(96) failed: %v", err)
	}

	if got.Convergent.Cmp(want.Convergent) != 0 {
		t.Fatalf("snapshot convergent not stable: got=%v want=%v", got.Convergent, want.Convergent)
	}
}

// mvp_numerator_test.go v6
