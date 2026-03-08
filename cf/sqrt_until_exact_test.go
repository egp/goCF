// sqrt_until_exact_test.go v1
package cf

import "testing"

func TestSqrtApproxRationalUntilExact_PerfectSquareImmediate(t *testing.T) {
	x := mustRat(9, 16)
	seed := mustRat(1, 1)

	got, exact, err := SqrtApproxRationalUntilExact(x, seed, 10)
	if err != nil {
		t.Fatalf("SqrtApproxRationalUntilExact failed: %v", err)
	}
	if !exact {
		t.Fatalf("expected exact=true")
	}

	want := mustRat(3, 4)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtApproxRationalUntilExact_Sqrt2NotExactWithinThree(t *testing.T) {
	x := mustRat(2, 1)
	seed := mustRat(1, 1)

	got, exact, err := SqrtApproxRationalUntilExact(x, seed, 3)
	if err != nil {
		t.Fatalf("SqrtApproxRationalUntilExact failed: %v", err)
	}
	if exact {
		t.Fatalf("did not expect exact convergence for sqrt(2)")
	}

	want := mustRat(577, 408)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtApproxRationalUntilExact_ZeroStepsReturnsSeed(t *testing.T) {
	x := mustRat(2, 1)
	seed := mustRat(7, 5)

	got, exact, err := SqrtApproxRationalUntilExact(x, seed, 0)
	if err != nil {
		t.Fatalf("SqrtApproxRationalUntilExact failed: %v", err)
	}
	if exact {
		t.Fatalf("did not expect exact=true")
	}
	if got.Cmp(seed) != 0 {
		t.Fatalf("got %v, want seed %v", got, seed)
	}
}

func TestSqrtApproxRationalUntilExact_RejectsNegativeInput(t *testing.T) {
	_, _, err := SqrtApproxRationalUntilExact(mustRat(-2, 1), mustRat(1, 1), 3)
	if err == nil {
		t.Fatalf("expected error for negative input")
	}
}

func TestSqrtApproxRationalUntilExact_RejectsNegativeMaxSteps(t *testing.T) {
	_, _, err := SqrtApproxRationalUntilExact(mustRat(2, 1), mustRat(1, 1), -1)
	if err == nil {
		t.Fatalf("expected error for negative maxSteps")
	}
}

// sqrt_until_exact_test.go v1
