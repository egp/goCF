// sqrt_seed_compare_source_test.go v1
package cf

import "testing"

func TestRangeSeedSourceWrapper_AndConvergentSourceWrapper_BothProduceCF(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	cfRange, err := NewSqrtApproxCFFromSourceRangeSeed(Sqrt2CF(), 2, p)
	if err != nil {
		t.Fatalf("NewSqrtApproxCFFromSourceRangeSeed failed: %v", err)
	}
	gotRange := collectTerms(cfRange, 8)
	if len(gotRange) == 0 {
		t.Fatalf("expected non-empty range-seeded CF")
	}

	cfConv, err := NewSqrtApproxCFFromSource(Sqrt2CF(), 2, p)
	if err != nil {
		t.Fatalf("NewSqrtApproxCFFromSource failed: %v", err)
	}
	gotConv := collectTerms(cfConv, 8)
	if len(gotConv) == 0 {
		t.Fatalf("expected non-empty convergent-seeded CF")
	}

	if gotRange[0] != 1 {
		t.Fatalf("range-seeded first digit = %d, want 1; full=%v", gotRange[0], gotRange)
	}
	if gotConv[0] != 1 {
		t.Fatalf("convergent-seeded first digit = %d, want 1; full=%v", gotConv[0], gotConv)
	}
}

// sqrt_seed_compare_source_test.go v1
