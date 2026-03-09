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

// gcf_approx_test.go v1
