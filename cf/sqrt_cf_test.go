// sqrt_cf_test.go v1
package cf

import "testing"

func TestNewSqrtApproxCF_Sqrt2_ThreeSteps(t *testing.T) {
	x := mustRat(2, 1)
	seed := mustRat(1, 1)

	cf, err := NewSqrtApproxCF(x, seed, 3)
	if err != nil {
		t.Fatalf("NewSqrtApproxCF failed: %v", err)
	}

	got := collectTerms(cf, 16)

	// 3 Newton steps from seed 1 give 577/408 = [1; 2,2,2,2,2,2,2]
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

func TestNewSqrtApproxCF_PerfectSquareFastPath(t *testing.T) {
	x := mustRat(9, 16)
	seed := mustRat(1, 1)

	cf, err := NewSqrtApproxCF(x, seed, 5)
	if err != nil {
		t.Fatalf("NewSqrtApproxCF failed: %v", err)
	}

	got := collectTerms(cf, 8)

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

func TestNewSqrtApproxCF_RejectsNegativeInput(t *testing.T) {
	_, err := NewSqrtApproxCF(mustRat(-2, 1), mustRat(1, 1), 3)
	if err == nil {
		t.Fatalf("expected error for negative input")
	}
}

func TestNewSqrtApproxCF_RejectsZeroSeed(t *testing.T) {
	_, err := NewSqrtApproxCF(mustRat(2, 1), mustRat(0, 1), 3)
	if err == nil {
		t.Fatalf("expected error for zero seed")
	}
}

// sqrt_cf_test.go v1
