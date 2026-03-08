// sqrt_api_test.go v1
package cf

import "testing"

func TestSqrtApprox_PerfectSquareFastPath(t *testing.T) {
	got, err := SqrtApprox(mustRat(9, 16))
	if err != nil {
		t.Fatalf("SqrtApprox failed: %v", err)
	}

	want := mustRat(3, 4)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtApprox_Sqrt2_DefaultPolicy(t *testing.T) {
	got, err := SqrtApprox(mustRat(2, 1))
	if err != nil {
		t.Fatalf("SqrtApprox failed: %v", err)
	}

	// With the current default policy:
	// seed = 2, steps = 5
	// Newton iterates end at 886731088897/627013566048.
	want := mustRat(886731088897, 627013566048)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtApproxCF_Sqrt2_DefaultPolicy(t *testing.T) {
	cf, err := SqrtApproxCF(mustRat(2, 1))
	if err != nil {
		t.Fatalf("SqrtApproxCF failed: %v", err)
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

func TestSqrtApproxTermsAuto_Sqrt2_DefaultPolicy(t *testing.T) {
	got, err := SqrtApproxTermsAuto(mustRat(2, 1), 16)
	if err != nil {
		t.Fatalf("SqrtApproxTermsAuto failed: %v", err)
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

func TestSqrtApproxTermsAuto_RejectsNegativeDigits(t *testing.T) {
	_, err := SqrtApproxTermsAuto(mustRat(2, 1), -1)
	if err == nil {
		t.Fatalf("expected error for negative digits")
	}
}

func TestSqrtApprox_RejectsNegativeInput(t *testing.T) {
	_, err := SqrtApprox(mustRat(-2, 1))
	if err == nil {
		t.Fatalf("expected error for negative input")
	}
}

// sqrt_api_test.go v1
