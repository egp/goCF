// sqrt_cf_tol_test.go v1
package cf

import "testing"

func TestNewSqrtApproxCFUntilResidual_Sqrt2(t *testing.T) {
	cf, err := NewSqrtApproxCFUntilResidual(
		mustRat(2, 1),
		mustRat(1, 1),
		3,
		mustRat(1, 1000),
	)
	if err != nil {
		t.Fatalf("NewSqrtApproxCFUntilResidual failed: %v", err)
	}

	got := collectTerms(cf, 16)
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

func TestNewSqrtApproxCFUntilResidualDefault_Sqrt2(t *testing.T) {
	cf, err := NewSqrtApproxCFUntilResidualDefault(
		mustRat(2, 1),
		3,
		mustRat(1, 1000),
	)
	if err != nil {
		t.Fatalf("NewSqrtApproxCFUntilResidualDefault failed: %v", err)
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

func TestSqrtApproxTermsUntilResidual_Sqrt2(t *testing.T) {
	got, err := SqrtApproxTermsUntilResidual(
		mustRat(2, 1),
		mustRat(1, 1),
		3,
		mustRat(1, 1000),
		16,
	)
	if err != nil {
		t.Fatalf("SqrtApproxTermsUntilResidual failed: %v", err)
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

func TestSqrtApproxTermsUntilResidualDefault_Sqrt2(t *testing.T) {
	got, err := SqrtApproxTermsUntilResidualDefault(
		mustRat(2, 1),
		3,
		mustRat(1, 1000),
		16,
	)
	if err != nil {
		t.Fatalf("SqrtApproxTermsUntilResidualDefault failed: %v", err)
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

// sqrt_cf_tol_test.go v1
