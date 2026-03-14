// sqrt_api2_test.go v1
package cf

import "testing"

func TestSqrtApprox2_PerfectSquareFastPath(t *testing.T) {
	got, err := SqrtApprox2(mustRat(9, 16))
	if err != nil {
		t.Fatalf("SqrtApprox2 failed: %v", err)
	}

	want := mustRat(3, 4)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtApprox2_Sqrt2_DefaultPolicy(t *testing.T) {
	got, err := SqrtApprox2(mustRat(2, 1))
	if err != nil {
		t.Fatalf("SqrtApprox2 failed: %v", err)
	}

	want := mustRat(886731088897, 627013566048)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtApproxCF2_Sqrt2_DefaultPolicy(t *testing.T) {
	cf, err := SqrtApproxCF2(mustRat(2, 1))
	if err != nil {
		t.Fatalf("SqrtApproxCF2 failed: %v", err)
	}

	got := collectTerms(cf, 16)
	want := []int64{1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}

	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d full=%v", i, got[i], want[i], got)
		}
	}
}

func TestSqrtApproxTerms2_Sqrt2_DefaultPolicy(t *testing.T) {
	got, err := SqrtApproxTerms2(mustRat(2, 1), 16)
	if err != nil {
		t.Fatalf("SqrtApproxTerms2 failed: %v", err)
	}

	want := []int64{1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d full=%v", i, got[i], want[i], got)
		}
	}
}

func TestSqrtApproxTerms2_RejectsNegativeDigits(t *testing.T) {
	_, err := SqrtApproxTerms2(mustRat(2, 1), -1)
	if err == nil {
		t.Fatalf("expected error for negative digits")
	}
}

func TestSqrtApprox2_RejectsNegativeInput(t *testing.T) {
	_, err := SqrtApprox2(mustRat(-2, 1))
	if err == nil {
		t.Fatalf("expected error for negative input")
	}
}

func TestSqrtApproxWithPolicy2_RejectsZeroSeedInPolicy(t *testing.T) {
	zero := mustRat(0, 1)
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
		Seed:     &zero,
	}

	_, err := SqrtApproxWithPolicy2(mustRat(2, 1), p)
	if err == nil {
		t.Fatalf("expected error for zero policy seed")
	}
}

func TestSqrtApproxWithPolicy2_RejectsNegativeSeedInPolicy(t *testing.T) {
	neg := mustRat(-1, 1)
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
		Seed:     &neg,
	}

	_, err := SqrtApproxWithPolicy2(mustRat(2, 1), p)
	if err == nil {
		t.Fatalf("expected error for negative policy seed")
	}
}

// sqrt_api2_test.go v1
