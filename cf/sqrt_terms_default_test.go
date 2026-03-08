// sqrt_terms_default_test.go v1
package cf

import "testing"

func TestSqrtApproxTermsDefault_Sqrt2_ThreeSteps(t *testing.T) {
	got, err := SqrtApproxTermsDefault(mustRat(2, 1), 3, 16)
	if err != nil {
		t.Fatalf("SqrtApproxTermsDefault failed: %v", err)
	}

	want := []int64{1, 2, 2, 2, 2, 2, 2, 2} // 577/408
	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d full=%v", i, got[i], want[i], got)
		}
	}
}

func TestSqrtApproxTermsDefault_PerfectSquareFastPath(t *testing.T) {
	got, err := SqrtApproxTermsDefault(mustRat(9, 16), 5, 8)
	if err != nil {
		t.Fatalf("SqrtApproxTermsDefault failed: %v", err)
	}

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

func TestSqrtApproxTermsDefault_RejectsNegativeInput(t *testing.T) {
	_, err := SqrtApproxTermsDefault(mustRat(-2, 1), 3, 8)
	if err == nil {
		t.Fatalf("expected error for negative input")
	}
}

// sqrt_terms_default_test.go v1
