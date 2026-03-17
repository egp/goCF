// mvp_radicand_test.go v4
package cf

import "testing"

func TestMVPRadicandAssembleSnapshot_RoundTripsApproxValue(t *testing.T) {
	got, err := MVPRadicandAssembleSnapshot(4, 6, 64)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleSnapshot failed: %v", err)
	}

	want, err := MVPRadicandAssembleConvergent(4, 6)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleConvergent failed: %v", err)
	}

	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

func TestMVPRadicandAssembleConvergent_RejectsBadBounds(t *testing.T) {
	if _, err := MVPRadicandAssembleConvergent(0, 4); err == nil {
		t.Fatalf("expected error for fourOverPiPrefixTerms=0")
	}
	if _, err := MVPRadicandAssembleConvergent(4, 0); err == nil {
		t.Fatalf("expected error for ePrefixTerms=0")
	}
}

func TestMVPRadicandAssembleConvergent_UsesCanonicalSources(t *testing.T) {
	got, err := MVPRadicandAssembleConvergent(4, 6)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleConvergent failed: %v", err)
	}

	fourOverPi, err := GCFSourceConvergent(NewBrouncker4OverPiGCFSource(), 4)
	if err != nil {
		t.Fatalf("GCFSourceConvergent Brouncker failed: %v", err)
	}
	eApprox, err := GCFSourceConvergent(NewECFGSource(), 6)
	if err != nil {
		t.Fatalf("GCFSourceConvergent e failed: %v", err)
	}

	fourOverPiSq, err := fourOverPi.Mul(fourOverPi)
	if err != nil {
		t.Fatalf("Mul failed: %v", err)
	}
	threeOverPiSq, err := mustRat(3, 16).Mul(fourOverPiSq)
	if err != nil {
		t.Fatalf("Mul scale failed: %v", err)
	}
	want, err := threeOverPiSq.Add(eApprox)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPRadicandAssembleConvergent_AlternateFourOverPiSourceHookWorks(t *testing.T) {
	got, err := MVPRadicandAssembleConvergentWithFourOverPiApprox(
		MVPFourOverPiApproxBrouncker,
		4,
		6,
	)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleConvergentWithFourOverPiApprox failed: %v", err)
	}

	want, err := MVPRadicandAssembleConvergent(4, 6)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleConvergent failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPRadicandAssembleConvergent_LambertParityPath_IsPositiveAndExceedsE(t *testing.T) {
	got, err := MVPRadicandAssembleConvergentWithFourOverPiApprox(
		MVPFourOverPiApproxLambert,
		8,
		6,
	)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleConvergentWithFourOverPiApprox failed: %v", err)
	}

	eApprox, err := GCFSourceConvergent(NewECFGSource(), 6)
	if err != nil {
		t.Fatalf("GCFSourceConvergent e failed: %v", err)
	}

	if got.Cmp(eApprox) <= 0 {
		t.Fatalf("got %v want > eApprox %v", got, eApprox)
	}
	if got.Cmp(intRat(0)) <= 0 {
		t.Fatalf("got %v want positive", got)
	}
}

func TestMVPRadicandAssembleConvergent_IsPositiveAndExceedsEApprox(t *testing.T) {
	got, err := MVPRadicandAssembleConvergent(4, 6)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleConvergent failed: %v", err)
	}

	eApprox, err := GCFSourceConvergent(NewECFGSource(), 6)
	if err != nil {
		t.Fatalf("GCFSourceConvergent e failed: %v", err)
	}

	if got.Cmp(eApprox) <= 0 {
		t.Fatalf("got %v want > eApprox %v", got, eApprox)
	}
	if got.Cmp(intRat(0)) <= 0 {
		t.Fatalf("got %v want positive", got)
	}
}

func TestMVPRadicandAssembleSnapshot_RejectsBadSnapshotTerms(t *testing.T) {
	_, err := MVPRadicandAssembleSnapshot(4, 6, 0)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMVPRadicandAssembleSnapshot_CurrentSnapshotBudgetIsStable(t *testing.T) {
	got, err := MVPRadicandAssembleSnapshot(4, 6, 64)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleSnapshot(64) failed: %v", err)
	}

	want, err := MVPRadicandAssembleSnapshot(4, 6, 96)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleSnapshot(96) failed: %v", err)
	}

	if got.Convergent.Cmp(want.Convergent) != 0 {
		t.Fatalf("snapshot convergent not stable: got=%v want=%v", got.Convergent, want.Convergent)
	}
}

func TestMVPRadicandAssembleSnapshotWithFourOverPiApprox_LambertRoundTripsAlternateApprox(t *testing.T) {
	got, err := MVPRadicandAssembleSnapshotWithFourOverPiApprox(
		MVPFourOverPiApproxLambert,
		8,
		6,
	)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleSnapshotWithFourOverPiApprox failed: %v", err)
	}

	want, err := MVPRadicandAssembleConvergentWithFourOverPiApprox(
		MVPFourOverPiApproxLambert,
		8,
		6,
	)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleConvergentWithFourOverPiApprox failed: %v", err)
	}

	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

func TestMVPRadicandRootValueFromSnapshot_LambertSnapshotPathIsUsableBySqrt(t *testing.T) {
	a, err := MVPRadicandAssembleSnapshotWithFourOverPiApprox(
		MVPFourOverPiApproxLambert,
		8,
		6,
	)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleSnapshotWithFourOverPiApprox failed: %v", err)
	}

	got, err := MVPRadicandRootValueFromSnapshot(a, DefaultSqrtPolicy2())
	if err != nil {
		t.Fatalf("MVPRadicandRootValueFromSnapshot failed: %v", err)
	}

	if got.Cmp(intRat(1)) <= 0 {
		t.Fatalf("got %v want > 1", got)
	}
}

func TestMVPRadicandAssembleSnapshot_MatchesDirectApproxAsPointSnapshot(t *testing.T) {
	got, err := MVPRadicandAssembleSnapshot(4, 6, 64)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleSnapshot failed: %v", err)
	}

	want, err := MVPRadicandAssembleConvergent(4, 6)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleConvergent failed: %v", err)
	}

	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
	if got.Range == nil {
		t.Fatalf("expected point range")
	}
	if got.Range.Lo.Cmp(want) != 0 || got.Range.Hi.Cmp(want) != 0 {
		t.Fatalf("got range %v want point %v", *got.Range, want)
	}
}

func TestMVPRadicandAssembleSnapshotWithFourOverPiApprox_LambertMatchesAlternateApprox(t *testing.T) {
	got, err := MVPRadicandAssembleSnapshotWithFourOverPiApprox(
		MVPFourOverPiApproxLambert,
		8,
		6,
	)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleSnapshotWithFourOverPiApprox failed: %v", err)
	}

	want, err := MVPRadicandAssembleConvergentWithFourOverPiApprox(
		MVPFourOverPiApproxLambert,
		8,
		6,
	)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleConvergentWithFourOverPiApprox failed: %v", err)
	}

	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
	if got.Range == nil {
		t.Fatalf("expected point range")
	}
	if got.Range.Lo.Cmp(want) != 0 || got.Range.Hi.Cmp(want) != 0 {
		t.Fatalf("got range %v want point %v", *got.Range, want)
	}
}

func TestMVPRadicandRootValueWithSnapshotTerms_IgnoresSnapshotBudgetOnDirectSnapshotPath(t *testing.T) {
	got, err := MVPRadicandRootValueWithSnapshotTerms(4, 6, DefaultSqrtPolicy2(), 64)
	if err != nil {
		t.Fatalf("MVPRadicandRootValueWithSnapshotTerms(64) failed: %v", err)
	}

	want, err := MVPRadicandRootValueWithSnapshotTerms(4, 6, DefaultSqrtPolicy2(), 96)
	if err != nil {
		t.Fatalf("MVPRadicandRootValueWithSnapshotTerms(96) failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("direct snapshot path should ignore snapshot budget: got=%v want=%v", got, want)
	}
}

func TestMVPRadicandSnapshotParts_DefaultFourOverPiSnapshotMatchesCanonicalConvergent(t *testing.T) {
	got, err := MVPRadicandDefaultFourOverPiSnapshot(4)
	if err != nil {
		t.Fatalf("MVPRadicandDefaultFourOverPiSnapshot failed: %v", err)
	}

	want, err := MVPDefaultFourOverPiApproxSnapshot(4)
	if err != nil {
		t.Fatalf("MVPDefaultFourOverPiApproxSnapshot failed: %v", err)
	}

	if got.Convergent.Cmp(want.Convergent) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want.Convergent)
	}
}

func TestMVPRadicandSnapshotParts_DefaultESnapshotMatchesCanonicalConvergent(t *testing.T) {
	got, err := MVPRadicandDefaultEApproxSnapshot(6)
	if err != nil {
		t.Fatalf("MVPRadicandDefaultEApproxSnapshot failed: %v", err)
	}

	want, err := MVPDefaultEApproxSnapshot(6)
	if err != nil {
		t.Fatalf("MVPDefaultEApproxSnapshot failed: %v", err)
	}

	if got.Convergent.Cmp(want.Convergent) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want.Convergent)
	}
}

func TestMVPRadicandSnapshotParts_ScaledSquareOfFourOverPiMatchesExistingSubexpression(t *testing.T) {
	fourOverPi, err := MVPDefaultFourOverPiApproxSnapshot(4)
	if err != nil {
		t.Fatalf("MVPDefaultFourOverPiApproxSnapshot failed: %v", err)
	}

	got, err := MVPRadicandScaledSquareOfFourOverPiApprox(fourOverPi)
	if err != nil {
		t.Fatalf("MVPRadicandScaledSquareOfFourOverPiApprox failed: %v", err)
	}

	wantFourOverPi, err := GCFSourceConvergent(NewBrouncker4OverPiGCFSource(), 4)
	if err != nil {
		t.Fatalf("GCFSourceConvergent failed: %v", err)
	}
	wantSq, err := wantFourOverPi.Mul(wantFourOverPi)
	if err != nil {
		t.Fatalf("Mul failed: %v", err)
	}
	want, err := mustRat(3, 16).Mul(wantSq)
	if err != nil {
		t.Fatalf("Mul scale failed: %v", err)
	}

	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

func TestMVPRadicandSnapshotParts_AssembleMatchesExistingApprox(t *testing.T) {
	fourOverPi, err := MVPDefaultFourOverPiApproxSnapshot(4)
	if err != nil {
		t.Fatalf("MVPDefaultFourOverPiApproxSnapshot failed: %v", err)
	}
	eSnap, err := MVPDefaultEApproxSnapshot(6)
	if err != nil {
		t.Fatalf("MVPDefaultEApproxSnapshot failed: %v", err)
	}

	got, err := MVPRadicandAssembleFromSnapshots(fourOverPi, eSnap)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleFromSnapshots failed: %v", err)
	}

	want, err := MVPRadicandAssembleConvergent(4, 6)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleConvergent failed: %v", err)
	}

	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
	if got.Range == nil {
		t.Fatalf("expected point range")
	}
	if got.Range.Lo.Cmp(want) != 0 || got.Range.Hi.Cmp(want) != 0 {
		t.Fatalf("got range %v want point %v", *got.Range, want)
	}
}

func TestMVPRadicandSnapshotParts_LambertAssemblyMatchesAlternateApprox(t *testing.T) {
	fourOverPiFn, err := MVPFourOverPiApproxFuncForFamily(MVPFourOverPiFamilyLambert)
	if err != nil {
		t.Fatalf("MVPFourOverPiApproxFuncForFamily failed: %v", err)
	}

	fourOverPi, err := MVPApproxSnapshotFromApproxFunc(fourOverPiFn, 8)
	if err != nil {
		t.Fatalf("MVPApproxSnapshotFromApproxFunc failed: %v", err)
	}
	eSnap, err := MVPDefaultEApproxSnapshot(6)
	if err != nil {
		t.Fatalf("MVPDefaultEApproxSnapshot failed: %v", err)
	}

	got, err := MVPRadicandAssembleFromSnapshots(fourOverPi, eSnap)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleFromSnapshots failed: %v", err)
	}

	want, err := MVPRadicandAssembleConvergentWithFourOverPiApprox(
		MVPFourOverPiApproxLambert,
		8,
		6,
	)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleConvergentWithFourOverPiApprox failed: %v", err)
	}

	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

func TestMVPExactScalarSnapshot_Three_IsExactPointSnapshot(t *testing.T) {
	got, err := MVPExactScalarSnapshot(3)
	if err != nil {
		t.Fatalf("MVPExactScalarSnapshot failed: %v", err)
	}

	want := mustRat(3, 1)
	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
	if got.Range == nil {
		t.Fatalf("expected point range")
	}
	if got.Range.Lo.Cmp(want) != 0 || got.Range.Hi.Cmp(want) != 0 {
		t.Fatalf("got range %v want point %v", *got.Range, want)
	}
}

func TestMVPExactScalarSnapshot_Sixteen_IsExactPointSnapshot(t *testing.T) {
	got, err := MVPExactScalarSnapshot(16)
	if err != nil {
		t.Fatalf("MVPExactScalarSnapshot failed: %v", err)
	}

	want := mustRat(16, 1)
	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
	if got.Range == nil {
		t.Fatalf("expected point range")
	}
	if got.Range.Lo.Cmp(want) != 0 || got.Range.Hi.Cmp(want) != 0 {
		t.Fatalf("got range %v want point %v", *got.Range, want)
	}
}

func TestMVPRadicandScaleFactorSnapshot_IsThreeSixteenths(t *testing.T) {
	got, err := MVPRadicandScaleFactorSnapshot()
	if err != nil {
		t.Fatalf("MVPRadicandScaleFactorSnapshot failed: %v", err)
	}

	want := mustRat(3, 16)
	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
	if got.Range == nil {
		t.Fatalf("expected point range")
	}
	if got.Range.Lo.Cmp(want) != 0 || got.Range.Hi.Cmp(want) != 0 {
		t.Fatalf("got range %v want point %v", *got.Range, want)
	}
}

func TestMVPRadicandScaledSquareOfFourOverPiApprox_UsesExplicitScaleFactorSnapshot(t *testing.T) {
	fourOverPi, err := MVPDefaultFourOverPiApproxSnapshot(4)
	if err != nil {
		t.Fatalf("MVPDefaultFourOverPiApproxSnapshot failed: %v", err)
	}

	got, err := MVPRadicandScaledSquareOfFourOverPiApprox(fourOverPi)
	if err != nil {
		t.Fatalf("MVPRadicandScaledSquareOfFourOverPiApprox failed: %v", err)
	}

	scale, err := MVPRadicandScaleFactorSnapshot()
	if err != nil {
		t.Fatalf("MVPRadicandScaleFactorSnapshot failed: %v", err)
	}

	sq, err := fourOverPi.Convergent.Mul(fourOverPi.Convergent)
	if err != nil {
		t.Fatalf("Mul failed: %v", err)
	}
	want, err := scale.Convergent.Mul(sq)
	if err != nil {
		t.Fatalf("Mul scale failed: %v", err)
	}

	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

func TestMVPExactScalarGCFSource_ThreeEvaluatesExactly(t *testing.T) {
	src, prefixTerms, err := MVPExactScalarGCFSource(3)
	if err != nil {
		t.Fatalf("MVPExactScalarGCFSource failed: %v", err)
	}
	if prefixTerms <= 0 {
		t.Fatalf("got prefixTerms=%d want > 0", prefixTerms)
	}

	got, err := GCFSourceConvergent(src, prefixTerms)
	if err != nil {
		t.Fatalf("GCFSourceConvergent failed: %v", err)
	}

	want := mustRat(3, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPExactScalarGCFSource_SixteenEvaluatesExactly(t *testing.T) {
	src, prefixTerms, err := MVPExactScalarGCFSource(16)
	if err != nil {
		t.Fatalf("MVPExactScalarGCFSource failed: %v", err)
	}
	if prefixTerms <= 0 {
		t.Fatalf("got prefixTerms=%d want > 0", prefixTerms)
	}

	got, err := GCFSourceConvergent(src, prefixTerms)
	if err != nil {
		t.Fatalf("GCFSourceConvergent failed: %v", err)
	}

	want := mustRat(16, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPExactScalarSnapshotFromSource_ThreeMatchesLegacyScalarSnapshot(t *testing.T) {
	got, err := MVPExactScalarSnapshotFromSource(3)
	if err != nil {
		t.Fatalf("MVPExactScalarSnapshotFromSource failed: %v", err)
	}

	want, err := MVPExactScalarSnapshot(3)
	if err != nil {
		t.Fatalf("MVPExactScalarSnapshot failed: %v", err)
	}

	if got.Convergent.Cmp(want.Convergent) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want.Convergent)
	}
	if got.Range == nil {
		t.Fatalf("expected point range")
	}
	if want.Range == nil {
		t.Fatalf("expected legacy point range")
	}
	if got.Range.Lo.Cmp(want.Range.Lo) != 0 || got.Range.Hi.Cmp(want.Range.Hi) != 0 {
		t.Fatalf("got range %v want %v", *got.Range, *want.Range)
	}
}

func TestMVPExactScalarSnapshotFromSource_SixteenMatchesLegacyScalarSnapshot(t *testing.T) {
	got, err := MVPExactScalarSnapshotFromSource(16)
	if err != nil {
		t.Fatalf("MVPExactScalarSnapshotFromSource failed: %v", err)
	}

	want, err := MVPExactScalarSnapshot(16)
	if err != nil {
		t.Fatalf("MVPExactScalarSnapshot failed: %v", err)
	}

	if got.Convergent.Cmp(want.Convergent) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want.Convergent)
	}
	if got.Range == nil {
		t.Fatalf("expected point range")
	}
	if want.Range == nil {
		t.Fatalf("expected legacy point range")
	}
	if got.Range.Lo.Cmp(want.Range.Lo) != 0 || got.Range.Hi.Cmp(want.Range.Hi) != 0 {
		t.Fatalf("got range %v want %v", *got.Range, *want.Range)
	}
}

func TestMVPRadicandScaleFactorSnapshot_RemainsThreeSixteenthsWhenBuiltFromScalarSources(t *testing.T) {
	three, err := MVPExactScalarSnapshotFromSource(3)
	if err != nil {
		t.Fatalf("MVPExactScalarSnapshotFromSource(3) failed: %v", err)
	}
	sixteen, err := MVPExactScalarSnapshotFromSource(16)
	if err != nil {
		t.Fatalf("MVPExactScalarSnapshotFromSource(16) failed: %v", err)
	}

	got, err := three.Convergent.Div(sixteen.Convergent)
	if err != nil {
		t.Fatalf("Div failed: %v", err)
	}

	want := mustRat(3, 16)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPRadicandAssembleConvergentWithFourOverPiApprox_UsesSnapshotAssemblyPath(t *testing.T) {
	fourOverPiFn, err := MVPFourOverPiApproxFuncForFamily(MVPFourOverPiFamilyLambert)
	if err != nil {
		t.Fatalf("MVPFourOverPiApproxFuncForFamily failed: %v", err)
	}

	got, err := MVPRadicandAssembleConvergentWithFourOverPiApprox(
		fourOverPiFn,
		8,
		6,
	)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleConvergentWithFourOverPiApprox failed: %v", err)
	}

	fourOverPi, err := MVPApproxSnapshotFromApproxFunc(fourOverPiFn, 8)
	if err != nil {
		t.Fatalf("MVPApproxSnapshotFromApproxFunc failed: %v", err)
	}
	eApprox, err := MVPRadicandDefaultEApproxSnapshot(6)
	if err != nil {
		t.Fatalf("MVPRadicandDefaultEApproxSnapshot failed: %v", err)
	}
	wantSnap, err := MVPRadicandAssembleFromSnapshots(fourOverPi, eApprox)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleFromSnapshots failed: %v", err)
	}

	if got.Cmp(wantSnap.Convergent) != 0 {
		t.Fatalf("got %v want %v", got, wantSnap.Convergent)
	}
}

func TestMVPRadicandAssembleConvergent_DefaultPathUsesSnapshotAssembly(t *testing.T) {
	got, err := MVPRadicandAssembleConvergent(4, 6)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleConvergent failed: %v", err)
	}

	fourOverPi, err := MVPRadicandDefaultFourOverPiSnapshot(4)
	if err != nil {
		t.Fatalf("MVPRadicandDefaultFourOverPiSnapshot failed: %v", err)
	}
	eApprox, err := MVPRadicandDefaultEApproxSnapshot(6)
	if err != nil {
		t.Fatalf("MVPRadicandDefaultEApproxSnapshot failed: %v", err)
	}
	wantSnap, err := MVPRadicandAssembleFromSnapshots(fourOverPi, eApprox)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleFromSnapshots failed: %v", err)
	}

	if got.Cmp(wantSnap.Convergent) != 0 {
		t.Fatalf("got %v want %v", got, wantSnap.Convergent)
	}
}

func TestMVPRadicandAssembleSnapshotWithFourOverPiApprox_UsesAssembleFromSnapshotsPath(t *testing.T) {
	got, err := MVPRadicandAssembleSnapshotWithFourOverPiApprox(
		MVPFourOverPiApproxLambert,
		8,
		6,
	)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleSnapshotWithFourOverPiApprox failed: %v", err)
	}

	fourOverPi, err := MVPApproxSnapshotFromApproxFunc(MVPFourOverPiApproxLambert, 8)
	if err != nil {
		t.Fatalf("MVPApproxSnapshotFromApproxFunc failed: %v", err)
	}
	eApprox, err := MVPRadicandDefaultEApproxSnapshot(6)
	if err != nil {
		t.Fatalf("MVPRadicandDefaultEApproxSnapshot failed: %v", err)
	}
	want, err := MVPRadicandAssembleFromSnapshots(fourOverPi, eApprox)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleFromSnapshots failed: %v", err)
	}

	if got.Convergent.Cmp(want.Convergent) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want.Convergent)
	}
}

// mvp_radicand_test.go v4
