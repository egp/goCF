// sqrt_gcf_test.go v1
package cf

import (
	"strings"
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

func TestSqrtGCF_RejectsNegativeExactFiniteInput(t *testing.T) {
	_, err := SqrtGCF(NewSliceGCF([2]int64{-1, 1}))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "negative input") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSqrtGCF_ExactFiniteNonSquare_NotImplemented(t *testing.T) {
	_, err := SqrtGCF(NewSliceGCF([2]int64{2, 1}))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "not implemented") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// sqrt_gcf_test.go v1
