// sqrt_source_prefix_api2_test.go v1
package cf

import "testing"

func TestSqrtApproxFromSourceRangeSeed2_Sqrt2Prefix2(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	got, err := SqrtApproxFromSourceRangeSeed2(Sqrt2CF(), 2, p)
	if err != nil {
		t.Fatalf("SqrtApproxFromSourceRangeSeed2 failed: %v", err)
	}

	if got.Cmp(intRat(0)) <= 0 {
		t.Fatalf("expected positive approximation, got %v", got)
	}
}

func TestSqrtApproxCFFromSourceRangeSeed2_Sqrt2Prefix2(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	cf, err := SqrtApproxCFFromSourceRangeSeed2(Sqrt2CF(), 2, p)
	if err != nil {
		t.Fatalf("SqrtApproxCFFromSourceRangeSeed2 failed: %v", err)
	}

	got := collectTerms(cf, 8)
	if len(got) == 0 {
		t.Fatalf("expected non-empty CF")
	}
	if got[0] != 1 {
		t.Fatalf("got first digit %d want 1 full=%v", got[0], got)
	}
}

func TestSqrtApproxTermsFromSourceRangeSeed2_Sqrt2Prefix2(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	got, err := SqrtApproxTermsFromSourceRangeSeed2(Sqrt2CF(), 2, p, 8)
	if err != nil {
		t.Fatalf("SqrtApproxTermsFromSourceRangeSeed2 failed: %v", err)
	}

	if len(got) == 0 {
		t.Fatalf("expected non-empty terms")
	}
	if got[0] != 1 {
		t.Fatalf("got first digit %d want 1 full=%v", got[0], got)
	}
}

func TestSqrtApproxTermsFromSourceRangeSeedDefault2_RejectsNegativeDigits(t *testing.T) {
	_, err := SqrtApproxTermsFromSourceRangeSeedDefault2(Sqrt2CF(), 2, -1)
	if err == nil {
		t.Fatalf("expected error for negative digits")
	}
}

// sqrt_source_prefix_api2_test.go v1
