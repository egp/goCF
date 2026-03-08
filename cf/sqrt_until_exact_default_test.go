// sqrt_until_exact_default_test.go v1
package cf

import "testing"

func TestSqrtApproxRationalUntilExactDefault_PerfectSquareImmediate(t *testing.T) {
	got, exact, err := SqrtApproxRationalUntilExactDefault(mustRat(9, 16), 10)
	if err != nil {
		t.Fatalf("SqrtApproxRationalUntilExactDefault failed: %v", err)
	}
	if !exact {
		t.Fatalf("expected exact=true")
	}

	want := mustRat(3, 4)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtApproxRationalUntilExactDefault_Sqrt2NotExactWithinThree(t *testing.T) {
	got, exact, err := SqrtApproxRationalUntilExactDefault(mustRat(2, 1), 3)
	if err != nil {
		t.Fatalf("SqrtApproxRationalUntilExactDefault failed: %v", err)
	}
	if exact {
		t.Fatalf("did not expect exact convergence for sqrt(2)")
	}

	want := mustRat(577, 408)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtApproxRationalUntilExactDefault_RejectsNegativeInput(t *testing.T) {
	_, _, err := SqrtApproxRationalUntilExactDefault(mustRat(-2, 1), 3)
	if err == nil {
		t.Fatalf("expected error for negative input")
	}
}

// sqrt_until_exact_default_test.go v1
