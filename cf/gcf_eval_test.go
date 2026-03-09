// gcf_eval_test.go v1
package cf

import "testing"

func TestEvaluateFiniteGCF_Empty(t *testing.T) {
	_, err := EvaluateFiniteGCF(NewSliceGCF())
	if err == nil {
		t.Fatalf("expected error for empty source")
	}
}

func TestEvaluateFiniteGCF_SingleTerm(t *testing.T) {
	got, err := EvaluateFiniteGCF(NewSliceGCF(
		[2]int64{3, 2},
	))
	if err != nil {
		t.Fatalf("EvaluateFiniteGCF failed: %v", err)
	}

	// Finite convention: last term contributes just p_last.
	want := mustRat(3, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestEvaluateFiniteGCF_TwoTerms(t *testing.T) {
	got, err := EvaluateFiniteGCF(NewSliceGCF(
		[2]int64{3, 2},
		[2]int64{5, 7},
	))
	if err != nil {
		t.Fatalf("EvaluateFiniteGCF failed: %v", err)
	}

	// x = 3 + 2/5 = 17/5
	want := mustRat(17, 5)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestEvaluateFiniteGCF_ThreeTerms(t *testing.T) {
	got, err := EvaluateFiniteGCF(NewSliceGCF(
		[2]int64{1, 1},
		[2]int64{2, 1},
		[2]int64{2, 1},
	))
	if err != nil {
		t.Fatalf("EvaluateFiniteGCF failed: %v", err)
	}

	// x = 1 + 1/(2 + 1/2) = 1 + 1/(5/2) = 7/5
	want := mustRat(7, 5)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestEvaluateFiniteGCF_RejectsBadQ(t *testing.T) {
	_, err := EvaluateFiniteGCF(NewSliceGCF(
		[2]int64{3, 0},
	))
	if err == nil {
		t.Fatalf("expected error for q=0")
	}
}

// gcf_eval_test.go v1
