// gcf_approx_test.go v1
package cf

import "testing"

func TestGCFApproxFromPrefix_FiniteSource(t *testing.T) {
	got, err := GCFApproxFromPrefix(NewSliceGCF(
		[2]int64{3, 2},
		[2]int64{5, 7},
	), 10)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}

	want := mustRat(17, 5)
	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
	if got.PrefixTerms != 10 {
		t.Fatalf("PrefixTerms=%d want 10", got.PrefixTerms)
	}
}

func TestGCFApproxFromPrefix_BoundedInfiniteSource(t *testing.T) {
	g := NewPeriodicGCF(
		[][2]int64{{1, 1}},
		[][2]int64{{2, 1}},
	)

	got, err := GCFApproxFromPrefix(g, 3)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}

	want := mustRat(7, 5)
	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
	if got.PrefixTerms != 3 {
		t.Fatalf("PrefixTerms=%d want 3", got.PrefixTerms)
	}
}

func TestGCFApproxFromPrefix_RejectsZeroPrefixTerms(t *testing.T) {
	_, err := GCFApproxFromPrefix(NewSliceGCF([2]int64{3, 2}), 0)
	if err == nil {
		t.Fatalf("expected error for zero prefixTerms")
	}
}

func TestGCFApproxFromPrefix_RejectsEmptySource(t *testing.T) {
	_, err := GCFApproxFromPrefix(NewSliceGCF(), 3)
	if err == nil {
		t.Fatalf("expected error for empty source")
	}
}

func TestGCFApproxCF_AndTerms(t *testing.T) {
	a, err := GCFApproxFromPrefix(NewSliceGCF(
		[2]int64{3, 2},
		[2]int64{5, 7},
	), 10)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}

	// convergent = 17/5 = [3; 2, 2]
	got, err := GCFApproxTerms(a, 8)
	if err != nil {
		t.Fatalf("GCFApproxTerms failed: %v", err)
	}

	want := []int64{3, 2, 2}
	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d full=%v", i, got[i], want[i], got)
		}
	}
}

func TestGCFApproxTerms_RejectsNegativeDigits(t *testing.T) {
	a, err := GCFApproxFromPrefix(NewSliceGCF(
		[2]int64{3, 2},
	), 1)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}

	_, err = GCFApproxTerms(a, -1)
	if err == nil {
		t.Fatalf("expected error for negative digits")
	}
}

func TestGCFSourceConvergent(t *testing.T) {
	got, err := GCFSourceConvergent(NewSliceGCF(
		[2]int64{3, 2},
		[2]int64{5, 7},
	), 10)
	if err != nil {
		t.Fatalf("GCFSourceConvergent failed: %v", err)
	}

	want := mustRat(17, 5)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestGCFSourceTerms(t *testing.T) {
	got, err := GCFSourceTerms(NewSliceGCF(
		[2]int64{3, 2},
		[2]int64{5, 7},
	), 10, 8)
	if err != nil {
		t.Fatalf("GCFSourceTerms failed: %v", err)
	}

	want := []int64{3, 2, 2} // 17/5 = [3;2,2]
	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d full=%v", i, got[i], want[i], got)
		}
	}
}

func TestGCFSourceTerms_RejectsNegativeDigits(t *testing.T) {
	_, err := GCFSourceTerms(NewSliceGCF([2]int64{3, 2}), 1, -1)
	if err == nil {
		t.Fatalf("expected error for negative digits")
	}
}

func TestGCFApproxFromPrefix_FiniteSourceCarriesExactPointRange(t *testing.T) {
	got, err := GCFApproxFromPrefix(NewSliceGCF(
		[2]int64{3, 2},
		[2]int64{5, 7},
	), 10)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}

	want := mustRat(17, 5)
	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got convergent %v want %v", got.Convergent, want)
	}
	if got.Range == nil {
		t.Fatalf("expected non-nil Range")
	}
	if got.Range.Lo.Cmp(want) != 0 || got.Range.Hi.Cmp(want) != 0 {
		t.Fatalf("got range [%v,%v] want exact [%v,%v]", got.Range.Lo, got.Range.Hi, want, want)
	}
}

func TestGCFApproxFromPrefix_BrounckerCarriesConservativeRange(t *testing.T) {
	got, err := GCFApproxFromPrefix(NewBrouncker4OverPiGCFSource(), 3)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}

	if got.Range == nil {
		t.Fatalf("expected non-nil Range")
	}

	wantLo := mustRat(7, 5)
	wantHi := mustRat(34, 23)
	if got.Range.Lo.Cmp(wantLo) != 0 || got.Range.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got range [%v,%v] want [%v,%v]", got.Range.Lo, got.Range.Hi, wantLo, wantHi)
	}
	if !got.Range.Contains(got.Convergent) {
		t.Fatalf("range %v does not contain convergent %v", *got.Range, got.Convergent)
	}
}

func TestGCFApproxFromPrefix_LambertCarriesConservativeRange(t *testing.T) {
	got, err := GCFApproxFromPrefix(NewLambertPiOver4GCFSource(), 3)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}

	if got.Range == nil {
		t.Fatalf("expected non-nil Range")
	}

	wantLo := mustRat(3, 4)
	wantHi := mustRat(7, 8)
	if got.Range.Lo.Cmp(wantLo) != 0 || got.Range.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got range [%v,%v] want [%v,%v]", got.Range.Lo, got.Range.Hi, wantLo, wantHi)
	}
	if !got.Range.Contains(got.Convergent) {
		t.Fatalf("range %v does not contain convergent %v", *got.Range, got.Convergent)
	}
}

// gcf_approx_test.go v1
