// mvp_radicand_test.go v1
package cf

import "testing"

func TestMVPThreeOverPiSquaredPlusEAsGCFSource_RoundTripsApproxValue(t *testing.T) {
	src, err := MVPThreeOverPiSquaredPlusEAsGCFSource(4, 6)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEAsGCFSource failed: %v", err)
	}

	got, err := GCFSourceConvergent(src, 64)
	if err != nil {
		t.Fatalf("GCFSourceConvergent failed: %v", err)
	}

	want, err := MVPThreeOverPiSquaredPlusEApprox(4, 6)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEApprox failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPThreeOverPiSquaredPlusEApprox_RejectsBadBounds(t *testing.T) {
	if _, err := MVPThreeOverPiSquaredPlusEApprox(0, 4); err == nil {
		t.Fatalf("expected error for fourOverPiPrefixTerms=0")
	}
	if _, err := MVPThreeOverPiSquaredPlusEApprox(4, 0); err == nil {
		t.Fatalf("expected error for ePrefixTerms=0")
	}
}

func TestMVPThreeOverPiSquaredPlusEApprox_UsesCanonicalSources(t *testing.T) {
	got, err := MVPThreeOverPiSquaredPlusEApprox(4, 6)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEApprox failed: %v", err)
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

func TestMVPThreeOverPiSquaredPlusEApprox_AlternateFourOverPiSourceHookWorks(t *testing.T) {
	got, err := MVPThreeOverPiSquaredPlusEApproxWithFourOverPiApprox(
		MVPFourOverPiApproxBrouncker,
		4,
		6,
	)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEApproxWithFourOverPiApprox failed: %v", err)
	}

	want, err := MVPThreeOverPiSquaredPlusEApprox(4, 6)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEApprox failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPThreeOverPiSquaredPlusEApprox_LambertParityPath_IsPositiveAndExceedsE(t *testing.T) {
	got, err := MVPThreeOverPiSquaredPlusEApproxWithFourOverPiApprox(
		MVPFourOverPiApproxLambert,
		8,
		6,
	)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEApproxWithFourOverPiApprox failed: %v", err)
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

func TestMVPThreeOverPiSquaredPlusEApprox_IsPositiveAndExceedsEApprox(t *testing.T) {
	got, err := MVPThreeOverPiSquaredPlusEApprox(4, 6)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEApprox failed: %v", err)
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

func TestMVPThreeOverPiSquaredPlusEApproxSnapshot_RoundTripsApproxValue(t *testing.T) {
	got, err := MVPThreeOverPiSquaredPlusEApproxSnapshot(4, 6, 64)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEApproxSnapshot failed: %v", err)
	}

	want, err := MVPThreeOverPiSquaredPlusEApprox(4, 6)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEApprox failed: %v", err)
	}

	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

func TestMVPThreeOverPiSquaredPlusEApproxSnapshot_RejectsBadBridgeTerms(t *testing.T) {
	_, err := MVPThreeOverPiSquaredPlusEApproxSnapshot(4, 6, 0)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMVPThreeOverPiSquaredPlusEApproxSnapshot_CurrentBridgeBudgetIsStable(t *testing.T) {
	got, err := MVPThreeOverPiSquaredPlusEApproxSnapshot(4, 6, 64)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEApproxSnapshot(64) failed: %v", err)
	}

	want, err := MVPThreeOverPiSquaredPlusEApproxSnapshot(4, 6, 96)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEApproxSnapshot(96) failed: %v", err)
	}

	if got.Convergent.Cmp(want.Convergent) != 0 {
		t.Fatalf("snapshot convergent not stable: got=%v want=%v", got.Convergent, want.Convergent)
	}
}

func TestMVPThreeOverPiSquaredPlusEFiniteBridgeSource_MatchesLegacyBridgeSource(t *testing.T) {
	got, err := MVPThreeOverPiSquaredPlusEFiniteBridgeSource(4, 6)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEFiniteBridgeSource failed: %v", err)
	}

	want, err := MVPThreeOverPiSquaredPlusEAsGCFSource(4, 6)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEAsGCFSource failed: %v", err)
	}

	gotConv, err := GCFSourceConvergent(got, 64)
	if err != nil {
		t.Fatalf("GCFSourceConvergent got failed: %v", err)
	}

	wantConv, err := GCFSourceConvergent(want, 64)
	if err != nil {
		t.Fatalf("GCFSourceConvergent want failed: %v", err)
	}

	if gotConv.Cmp(wantConv) != 0 {
		t.Fatalf("got %v want %v", gotConv, wantConv)
	}
}

func TestMVPThreeOverPiSquaredPlusEFiniteBridgeSnapshot_MatchesLegacySnapshot(t *testing.T) {
	got, err := MVPThreeOverPiSquaredPlusEFiniteBridgeSnapshot(4, 6, 64)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEFiniteBridgeSnapshot failed: %v", err)
	}

	want, err := MVPThreeOverPiSquaredPlusEApproxSnapshot(4, 6, 64)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEApproxSnapshot failed: %v", err)
	}

	if got.Convergent.Cmp(want.Convergent) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want.Convergent)
	}
}

func TestMVPThreeOverPiSquaredPlusEAsGCFSource_LegacyNameMatchesFiniteBridgeSource(t *testing.T) {
	got, err := MVPThreeOverPiSquaredPlusEAsGCFSource(4, 6)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEAsGCFSource failed: %v", err)
	}

	want, err := MVPThreeOverPiSquaredPlusEFiniteBridgeSource(4, 6)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEFiniteBridgeSource failed: %v", err)
	}

	gotConv, err := GCFSourceConvergent(got, 64)
	if err != nil {
		t.Fatalf("GCFSourceConvergent got failed: %v", err)
	}

	wantConv, err := GCFSourceConvergent(want, 64)
	if err != nil {
		t.Fatalf("GCFSourceConvergent want failed: %v", err)
	}

	if gotConv.Cmp(wantConv) != 0 {
		t.Fatalf("got %v want %v", gotConv, wantConv)
	}
}

func TestMVPThreeOverPiSquaredPlusEApproxSnapshot_LegacyNameMatchesFiniteBridgeSnapshot(t *testing.T) {
	got, err := MVPThreeOverPiSquaredPlusEApproxSnapshot(4, 6, 64)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEApproxSnapshot failed: %v", err)
	}

	want, err := MVPThreeOverPiSquaredPlusEFiniteBridgeSnapshot(4, 6, 64)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEFiniteBridgeSnapshot failed: %v", err)
	}

	if got.Convergent.Cmp(want.Convergent) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want.Convergent)
	}
}

func TestMVPThreeOverPiSquaredPlusERadicandSource_CurrentlyMatchesFiniteBridge(t *testing.T) {
	got, err := MVPThreeOverPiSquaredPlusERadicandSource(4, 6)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusERadicandSource failed: %v", err)
	}

	want, err := MVPThreeOverPiSquaredPlusEFiniteBridgeSource(4, 6)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEFiniteBridgeSource failed: %v", err)
	}

	gotConv, err := GCFSourceConvergent(got, 64)
	if err != nil {
		t.Fatalf("GCFSourceConvergent got failed: %v", err)
	}

	wantConv, err := GCFSourceConvergent(want, 64)
	if err != nil {
		t.Fatalf("GCFSourceConvergent want failed: %v", err)
	}

	if gotConv.Cmp(wantConv) != 0 {
		t.Fatalf("got %v want %v", gotConv, wantConv)
	}
}

func TestMVPThreeOverPiSquaredPlusEFiniteBridgeSourceWithFourOverPiApprox_LambertRoundTripsAlternateApprox(t *testing.T) {
	src, err := MVPThreeOverPiSquaredPlusEFiniteBridgeSourceWithFourOverPiApprox(
		MVPFourOverPiApproxLambert,
		8,
		6,
	)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEFiniteBridgeSourceWithFourOverPiApprox failed: %v", err)
	}

	got, err := GCFSourceConvergent(src, 64)
	if err != nil {
		t.Fatalf("GCFSourceConvergent failed: %v", err)
	}

	want, err := MVPThreeOverPiSquaredPlusEApproxWithFourOverPiApprox(
		MVPFourOverPiApproxLambert,
		8,
		6,
	)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEApproxWithFourOverPiApprox failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPThreeOverPiSquaredPlusEFiniteBridgeSnapshotWithFourOverPiApprox_LambertRoundTripsAlternateApprox(t *testing.T) {
	got, err := MVPThreeOverPiSquaredPlusEFiniteBridgeSnapshotWithFourOverPiApprox(
		MVPFourOverPiApproxLambert,
		8,
		6,
		64,
	)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEFiniteBridgeSnapshotWithFourOverPiApprox failed: %v", err)
	}

	want, err := MVPThreeOverPiSquaredPlusEApproxWithFourOverPiApprox(
		MVPFourOverPiApproxLambert,
		8,
		6,
	)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEApproxWithFourOverPiApprox failed: %v", err)
	}

	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

func TestMVPNumeratorApproxFromRadicandApprox_LambertBridgePathIsUsableBySqrt(t *testing.T) {
	a, err := MVPThreeOverPiSquaredPlusEFiniteBridgeSnapshotWithFourOverPiApprox(
		MVPFourOverPiApproxLambert,
		8,
		6,
		MVPNumeratorBridgePrefixTerms,
	)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEFiniteBridgeSnapshotWithFourOverPiApprox failed: %v", err)
	}

	got, err := MVPNumeratorApproxFromRadicandApprox(a, DefaultSqrtPolicy2())
	if err != nil {
		t.Fatalf("MVPNumeratorApproxFromRadicandApprox failed: %v", err)
	}

	if got.Cmp(intRat(1)) <= 0 {
		t.Fatalf("got %v want > 1", got)
	}
}

// EOF
