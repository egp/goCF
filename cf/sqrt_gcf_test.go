// sqrt_gcf_test.go v2
package cf

import (
	"testing"
)

func TestSqrtGCF_ExactFiniteZero(t *testing.T) {
	cf, err := SqrtGCF(NewSliceGCF([2]int64{0, 1}))
	if err != nil {
		t.Fatalf("SqrtGCF failed: %v", err)
	}

	got := collectTerms(cf, 8)
	want := []int64{0}

	if len(got) != len(want) {
		t.Fatalf("len mismatch: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, got, want)
		}
	}
}

func TestSqrtGCF_ExactFiniteOne(t *testing.T) {
	cf, err := SqrtGCF(NewSliceGCF([2]int64{1, 1}))
	if err != nil {
		t.Fatalf("SqrtGCF failed: %v", err)
	}

	got := collectTerms(cf, 8)
	want := []int64{1}

	if len(got) != len(want) {
		t.Fatalf("len mismatch: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, got, want)
		}
	}
}

func TestSqrtGCF_ExactFiniteFour(t *testing.T) {
	cf, err := SqrtGCF(NewSliceGCF([2]int64{4, 1}))
	if err != nil {
		t.Fatalf("SqrtGCF failed: %v", err)
	}

	got := collectTerms(cf, 8)
	want := []int64{2}

	if len(got) != len(want) {
		t.Fatalf("len mismatch: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, got, want)
		}
	}
}

func TestSqrtGCF_ExactFiniteNine(t *testing.T) {
	cf, err := SqrtGCF(NewSliceGCF([2]int64{9, 1}))
	if err != nil {
		t.Fatalf("SqrtGCF failed: %v", err)
	}

	got := collectTerms(cf, 8)
	want := []int64{3}

	if len(got) != len(want) {
		t.Fatalf("len mismatch: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, got, want)
		}
	}
}

func TestSqrtGCF_ExactFiniteOneQuarter(t *testing.T) {
	src := AdaptCFToGCF(NewRationalCF(mustRat(1, 4)))

	cf, err := SqrtGCF(src)
	if err != nil {
		t.Fatalf("SqrtGCF failed: %v", err)
	}

	got := collectTerms(cf, 8)
	want := collectTerms(NewRationalCF(mustRat(1, 2)), 8)

	if len(got) != len(want) {
		t.Fatalf("len mismatch: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, got, want)
		}
	}
}

func TestSqrtGCF_ExactFiniteNineSixteenths(t *testing.T) {
	src := AdaptCFToGCF(NewRationalCF(mustRat(9, 16)))

	cf, err := SqrtGCF(src)
	if err != nil {
		t.Fatalf("SqrtGCF failed: %v", err)
	}

	got := collectTerms(cf, 8)
	want := collectTerms(NewRationalCF(mustRat(3, 4)), 8)

	if len(got) != len(want) {
		t.Fatalf("len mismatch: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, got, want)
		}
	}
}

func TestSqrtGCF_RejectsNegativeExactFiniteInput(t *testing.T) {
	cf, err := SqrtGCF(NewSliceGCF([2]int64{-1, 1}))
	if err != nil {
		t.Fatalf("SqrtGCF failed: %v", err)
	}

	_, ok := cf.Next()
	if ok {
		t.Fatalf("expected no emitted term")
	}
}

func TestSqrtGCF_ExactFiniteNonSquare_ReturnsNewtonApproxCF(t *testing.T) {
	cf, err := SqrtGCF(NewSliceGCF([2]int64{2, 1}))
	if err != nil {
		t.Fatalf("SqrtGCF failed: %v", err)
	}

	got := collectTerms(cf, 8)
	want := collectTerms(NewRationalCF(mustRat(577, 408)), 8)

	if len(got) != len(want) {
		t.Fatalf("len mismatch: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, got, want)
		}
	}
}

func TestSqrtGCF_NonTerminatingInput_CurrentlyEmitsNothing(t *testing.T) {
	cf, err := SqrtGCF(AdaptCFToGCF(Sqrt2CF()))
	if err != nil {
		t.Fatalf("SqrtGCF failed: %v", err)
	}

	_, ok := cf.Next()
	if ok {
		t.Fatalf("expected no emitted term yet for non-terminating input")
	}
}

// sqrt_gcf_test.go v2
