// mvp_sources_core_test.go v1
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

func equalTerms(got, want []int64) bool {
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

func TestMVPSources_69DegreeSource_AsExactFiniteGCFPrefixEvaluatesExactly(t *testing.T) {
	got, err := GCFSourceConvergent(MVP69DegreeGCFSource(), 2)
	if err != nil {
		t.Fatalf("GCFSourceConvergent failed: %v", err)
	}

	want := mustRat(69, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVP69DegreeGCFSource_MatchesDirectFiniteEncoding(t *testing.T) {
	got, err := GCFSourceConvergent(MVP69DegreeGCFSource(), 2)
	if err != nil {
		t.Fatalf("GCFSourceConvergent failed: %v", err)
	}

	want, err := EvaluateFiniteGCF(NewSliceGCF(
		[2]int64{68, 1},
		[2]int64{1, 1},
	))
	if err != nil {
		t.Fatalf("EvaluateFiniteGCF failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}
