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

func TestNewSqrtApproxCFFromSourceRangeMidpoint_Sqrt2Prefix2(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	cf, err := NewSqrtApproxCFFromSourceRangeMidpoint(Sqrt2CF(), 2, p)
	if err != nil {
		t.Fatalf("NewSqrtApproxCFFromSourceRangeMidpoint failed: %v", err)
	}

	got := collectTerms(cf, 8)
	if len(got) == 0 {
		t.Fatalf("expected non-empty CF")
	}
	if got[0] != 1 {
		t.Fatalf("got first digit %d, want 1; full=%v", got[0], got)
	}
}

func TestSqrtApproxTermsFromSourceRangeMidpointDefault_RejectsNegativeDigits(t *testing.T) {
	_, err := SqrtApproxTermsFromSourceRangeMidpointDefault(Sqrt2CF(), 2, -1)
	if err == nil {
		t.Fatalf("expected error for negative digits")
	}
}

// sqrt_cf_test.go v1
