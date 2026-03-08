// sqrt_cfsource_test.go v1
package cf

import "testing"

func TestRationalApproxFromCFPrefix_FiniteSourceExact(t *testing.T) {
	// 355/113 = [3; 7, 16]
	got, err := RationalApproxFromCFPrefix(NewSliceCF(3, 7, 16), 10)
	if err != nil {
		t.Fatalf("RationalApproxFromCFPrefix failed: %v", err)
	}

	want := mustRat(355, 113)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestRationalApproxFromCFPrefix_InfiniteSourceConvergent(t *testing.T) {
	// sqrt(2) prefix [1;2,2] => convergent 7/5
	got, err := RationalApproxFromCFPrefix(Sqrt2CF(), 3)
	if err != nil {
		t.Fatalf("RationalApproxFromCFPrefix failed: %v", err)
	}

	want := mustRat(7, 5)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestNewSqrtApproxCFFromSource_Sqrt2Prefix3(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	cf, err := NewSqrtApproxCFFromSource(Sqrt2CF(), 3, p)
	if err != nil {
		t.Fatalf("NewSqrtApproxCFFromSource failed: %v", err)
	}

	got := collectTerms(cf, 16)

	// Prefix [1;2,2] = 7/5, so we are approximating sqrt(7/5).
	// Just assert the CF begins with 1 and remains finite/nonempty.
	if len(got) == 0 {
		t.Fatalf("expected non-empty CF")
	}
	if got[0] != 1 {
		t.Fatalf("got first digit %d, want 1; full=%v", got[0], got)
	}
}

func TestNewSqrtApproxCFFromSourceDefault_PerfectSquareFiniteSource(t *testing.T) {
	// 4 = [4]
	cf, err := NewSqrtApproxCFFromSourceDefault(NewSliceCF(4), 1)
	if err != nil {
		t.Fatalf("NewSqrtApproxCFFromSourceDefault failed: %v", err)
	}

	got := collectTerms(cf, 8)
	want := []int64{2}

	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d full=%v", i, got[i], want[i], got)
		}
	}
}

func TestSqrtApproxTermsFromSourceDefault_RejectsNegativeDigits(t *testing.T) {
	_, err := SqrtApproxTermsFromSourceDefault(Sqrt2CF(), 3, -1)
	if err == nil {
		t.Fatalf("expected error for negative digits")
	}
}

func TestRationalApproxFromCFPrefix_RejectsZeroPrefixTerms(t *testing.T) {
	_, err := RationalApproxFromCFPrefix(Sqrt2CF(), 0)
	if err == nil {
		t.Fatalf("expected error for zero prefixTerms")
	}
}

// sqrt_cfsource_test.go v1
