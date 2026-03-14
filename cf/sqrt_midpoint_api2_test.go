// sqrt_midpoint_api2_test.go v1
package cf

import "testing"

func TestSqrtApproxCFFromSourceRangeMidpoint2_Sqrt2Prefix2(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	cf, err := SqrtApproxCFFromSourceRangeMidpoint2(Sqrt2CF(), 2, p)
	if err != nil {
		t.Fatalf("SqrtApproxCFFromSourceRangeMidpoint2 failed: %v", err)
	}

	got := collectTerms(cf, 8)
	if len(got) == 0 {
		t.Fatalf("expected non-empty CF")
	}
	if got[0] != 1 {
		t.Fatalf("got first digit %d, want 1; full=%v", got[0], got)
	}
}

func TestSqrtApproxTermsFromSourceRangeMidpointDefault2_RejectsNegativeDigits(t *testing.T) {
	_, err := SqrtApproxTermsFromSourceRangeMidpointDefault2(Sqrt2CF(), 2, -1)
	if err == nil {
		t.Fatalf("expected error for negative digits")
	}
}
