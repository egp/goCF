// sqrt_api_policy_test.go v2
package cf

import "testing"

func TestSqrtApproxWithPolicy_PerfectSquare(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 1,
		Tol:      mustRat(1, 10),
	}

	got, err := SqrtApproxWithPolicy(mustRat(9, 16), p)
	if err != nil {
		t.Fatalf("SqrtApproxWithPolicy failed: %v", err)
	}

	want := mustRat(3, 4)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtApproxWithPolicy_Sqrt2_ThreeSteps(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	got, err := SqrtApproxWithPolicy(mustRat(2, 1), p)
	if err != nil {
		t.Fatalf("SqrtApproxWithPolicy failed: %v", err)
	}

	want := mustRat(577, 408)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtApproxWithPolicy_ExplicitSeedOverride(t *testing.T) {
	seed := mustRat(1, 1)
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
		Seed:     &seed,
	}

	got, err := SqrtApproxWithPolicy(mustRat(2, 1), p)
	if err != nil {
		t.Fatalf("SqrtApproxWithPolicy failed: %v", err)
	}

	want := mustRat(577, 408)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtApproxCFWithPolicy_Sqrt2(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	cf, err := SqrtApproxCFWithPolicy(mustRat(2, 1), p)
	if err != nil {
		t.Fatalf("SqrtApproxCFWithPolicy failed: %v", err)
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

func TestSqrtApproxTermsWithPolicy_Sqrt2(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	got, err := SqrtApproxTermsWithPolicy(mustRat(2, 1), p, 16)
	if err != nil {
		t.Fatalf("SqrtApproxTermsWithPolicy failed: %v", err)
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

func TestSqrtApproxTermsWithPolicy_RejectsNegativeDigits(t *testing.T) {
	p := DefaultSqrtPolicy()
	_, err := SqrtApproxTermsWithPolicy(mustRat(2, 1), p, -1)
	if err == nil {
		t.Fatalf("expected error for negative digits")
	}
}

// sqrt_api_policy_test.go v2
