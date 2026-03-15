// mvp_sources_test.go v5
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

// mvp_sources_test.go v5
