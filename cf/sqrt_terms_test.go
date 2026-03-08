// sqrt_terms_test.go v1
package cf

import "testing"

func TestSqrtApproxTerms_Sqrt2_ThreeSteps(t *testing.T) {
	x := mustRat(2, 1)
	seed := mustRat(1, 1)

	got, err := SqrtApproxTerms(x, seed, 3, 16)
	if err != nil {
		t.Fatalf("SqrtApproxTerms failed: %v", err)
	}

	// 3 Newton steps from seed 1 give 577/408 = [1; 2, 2, 2, 2, 2, 2, 2, 2]
	want := []int64{1, 2, 2, 2, 2, 2, 2, 2}
	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d full=%v", i, got[i], want[i], got)
		}
	}
}

func TestSqrtApproxTerms_PerfectSquareFastPath(t *testing.T) {
	x := mustRat(9, 16)
	seed := mustRat(1, 1)

	got, err := SqrtApproxTerms(x, seed, 5, 8)
	if err != nil {
		t.Fatalf("SqrtApproxTerms failed: %v", err)
	}

	// sqrt(9/16) = 3/4 = [0; 1, 3]
	want := []int64{0, 1, 3}
	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d full=%v", i, got[i], want[i], got)
		}
	}
}

func TestSqrtApproxTerms_ZeroDigits(t *testing.T) {
	x := mustRat(2, 1)
	seed := mustRat(1, 1)

	got, err := SqrtApproxTerms(x, seed, 3, 0)
	if err != nil {
		t.Fatalf("SqrtApproxTerms failed: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

func TestSqrtApproxTerms_RejectsNegativeDigits(t *testing.T) {
	_, err := SqrtApproxTerms(mustRat(2, 1), mustRat(1, 1), 3, -1)
	if err == nil {
		t.Fatalf("expected error for negative digits")
	}
}

// sqrt_terms_test.go v1
