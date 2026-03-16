// mvp_sources_test.go v9
package cf

import "testing"

func collectPQ(src GCFSource, n int) [][2]int64 {
	out := make([][2]int64, 0, n)
	for i := 0; i < n; i++ {
		p, q, ok := src.NextPQ()
		if !ok {
			break
		}
		out = append(out, [2]int64{p, q})
	}
	return out
}

func equalPQ(got, want [][2]int64) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range want {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

func TestMVPSources_DefaultFourOverPiFamilyIsBrouncker(t *testing.T) {
	if MVPDefaultFourOverPiFamily != MVPFourOverPiFamilyBrouncker {
		t.Fatalf("got %q want %q", MVPDefaultFourOverPiFamily, MVPFourOverPiFamilyBrouncker)
	}
}

func TestMVPSources_ReciprocalPiUsesBrouncker4OverPi(t *testing.T) {
	got := collectPQ(MVPReciprocalPiGCFSource(), 4)
	want := [][2]int64{
		{1, 1},
		{2, 9},
		{2, 25},
		{2, 49},
	}
	if !equalPQ(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPSources_EUsesECFGSource(t *testing.T) {
	got := collectPQ(MVPEGCFSource(), 7)
	want := [][2]int64{
		{2, 1},
		{1, 1},
		{2, 1},
		{1, 1},
		{1, 1},
		{4, 1},
		{1, 1},
	}
	if !equalPQ(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPSources_69DegreeSourceWithTailEvaluatesExactly(t *testing.T) {
	got, _, err := EvalGCFWithTailExact(
		MVP69DegreeGCFSource(),
		MVP69DegreeTail(),
		1,
	)
	if err != nil {
		t.Fatalf("EvalGCFWithTailExact failed: %v", err)
	}

	want := mustRat(69, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPFourOverPiApproxFuncForFamily_Brouncker(t *testing.T) {
	fn, err := MVPFourOverPiApproxFuncForFamily(MVPFourOverPiFamilyBrouncker)
	if err != nil {
		t.Fatalf("MVPFourOverPiApproxFuncForFamily failed: %v", err)
	}

	got, err := fn(4)
	if err != nil {
		t.Fatalf("brouncker fn failed: %v", err)
	}

	want, err := MVPFourOverPiApproxBrouncker(4)
	if err != nil {
		t.Fatalf("MVPFourOverPiApproxBrouncker failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPFourOverPiApproxFuncForFamily_Lambert(t *testing.T) {
	fn, err := MVPFourOverPiApproxFuncForFamily(MVPFourOverPiFamilyLambert)
	if err != nil {
		t.Fatalf("MVPFourOverPiApproxFuncForFamily failed: %v", err)
	}

	got, err := fn(8)
	if err != nil {
		t.Fatalf("lambert fn failed: %v", err)
	}

	want, err := MVPFourOverPiApproxLambert(8)
	if err != nil {
		t.Fatalf("MVPFourOverPiApproxLambert failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPFourOverPiApproxFuncForFamily_RejectsUnknown(t *testing.T) {
	_, err := MVPFourOverPiApproxFuncForFamily(MVPFourOverPiFamily("bogus"))
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMVPFourOverPiApproxBrouncker_MatchesDirectConvergent(t *testing.T) {
	got, err := MVPFourOverPiApproxBrouncker(4)
	if err != nil {
		t.Fatalf("MVPFourOverPiApproxBrouncker failed: %v", err)
	}

	want, err := GCFSourceConvergent(NewBrouncker4OverPiGCFSource(), 4)
	if err != nil {
		t.Fatalf("GCFSourceConvergent failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPFourOverPiApproxLambert_IsPositiveAndGreaterThanOne(t *testing.T) {
	got, err := MVPFourOverPiApproxLambert(8)
	if err != nil {
		t.Fatalf("MVPFourOverPiApproxLambert failed: %v", err)
	}

	if got.Cmp(intRat(1)) <= 0 {
		t.Fatalf("got %v want > 1", got)
	}
	if got.Cmp(intRat(2)) >= 0 {
		t.Fatalf("got %v want < 2", got)
	}
}

func TestMVPFourOverPiApproxWithSource_BrounckerPath(t *testing.T) {
	brounckerSrc := func() GCFSource { return NewBrouncker4OverPiGCFSource() }

	got, err := MVPFourOverPiApproxWithSource(brounckerSrc, 4)
	if err != nil {
		t.Fatalf("MVPFourOverPiApproxWithSource failed: %v", err)
	}

	want, err := GCFSourceConvergent(NewBrouncker4OverPiGCFSource(), 4)
	if err != nil {
		t.Fatalf("GCFSourceConvergent failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

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

// mvp_sources_test.go v9
