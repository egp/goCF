// sqrt_api_seed_policy_test.go v1
package cf

import "testing"

func TestSqrtApproxWithSeedAndPolicy_PerfectSquare(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 1,
		Tol:      mustRat(1, 10),
	}

	got, err := SqrtApproxWithSeedAndPolicy(mustRat(9, 16), mustRat(1, 1), p)
	if err != nil {
		t.Fatalf("SqrtApproxWithSeedAndPolicy failed: %v", err)
	}

	want := mustRat(3, 4)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtApproxWithSeedAndPolicy_Sqrt2_ThreeStepsFromOne(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	got, err := SqrtApproxWithSeedAndPolicy(mustRat(2, 1), mustRat(1, 1), p)
	if err != nil {
		t.Fatalf("SqrtApproxWithSeedAndPolicy failed: %v", err)
	}

	want := mustRat(577, 408)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtApproxCFWithSeedAndPolicy_Sqrt2(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	cf, err := SqrtApproxCFWithSeedAndPolicy(mustRat(2, 1), mustRat(1, 1), p)
	if err != nil {
		t.Fatalf("SqrtApproxCFWithSeedAndPolicy failed: %v", err)
	}

	got := collectTerms(cf, 16)
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

func TestSqrtApproxTermsWithSeedAndPolicy_Sqrt2(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	got, err := SqrtApproxTermsWithSeedAndPolicy(mustRat(2, 1), mustRat(1, 1), p, 16)
	if err != nil {
		t.Fatalf("SqrtApproxTermsWithSeedAndPolicy failed: %v", err)
	}

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

func TestSqrtApproxTermsWithSeedAndPolicy_RejectsNegativeDigits(t *testing.T) {
	p := DefaultSqrtPolicy()
	_, err := SqrtApproxTermsWithSeedAndPolicy(mustRat(2, 1), mustRat(1, 1), p, -1)
	if err == nil {
		t.Fatalf("expected error for negative digits")
	}
}

// sqrt_api_seed_policy_test.go v1
