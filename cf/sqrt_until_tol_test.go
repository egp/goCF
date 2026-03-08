// sqrt_until_tol_test.go v1
package cf

import "testing"

func TestSqrtApproxRationalUntilResidual_PerfectSquareImmediate(t *testing.T) {
	got, ok, err := SqrtApproxRationalUntilResidual(
		mustRat(9, 16),
		mustRat(1, 1),
		10,
		mustRat(0, 1),
	)
	if err != nil {
		t.Fatalf("SqrtApproxRationalUntilResidual failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected satisfied=true")
	}
	want := mustRat(3, 4)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtApproxRationalUntilResidual_Sqrt2_SatisfiesLooseTolerance(t *testing.T) {
	got, ok, err := SqrtApproxRationalUntilResidual(
		mustRat(2, 1),
		mustRat(1, 1),
		2,
		mustRat(1, 100),
	)
	if err != nil {
		t.Fatalf("SqrtApproxRationalUntilResidual failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected satisfied=true")
	}

	// After 2 Newton steps from seed 1, sqrt(2) approx is 17/12.
	want := mustRat(17, 12)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtApproxRationalUntilResidual_Sqrt2_DoesNotSatisfyTightToleranceSoon(t *testing.T) {
	got, ok, err := SqrtApproxRationalUntilResidual(
		mustRat(2, 1),
		mustRat(1, 1),
		2,
		mustRat(1, 1000),
	)
	if err != nil {
		t.Fatalf("SqrtApproxRationalUntilResidual failed: %v", err)
	}
	if ok {
		t.Fatalf("did not expect satisfied=true")
	}

	want := mustRat(17, 12)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtApproxRationalUntilResidual_ZeroStepsChecksSeed(t *testing.T) {
	seed := mustRat(3, 2)
	got, ok, err := SqrtApproxRationalUntilResidual(
		mustRat(2, 1),
		seed,
		0,
		mustRat(1, 3), // residual at 3/2 is 1/4 <= 1/3
	)
	if err != nil {
		t.Fatalf("SqrtApproxRationalUntilResidual failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected satisfied=true at seed")
	}
	if got.Cmp(seed) != 0 {
		t.Fatalf("got %v, want seed %v", got, seed)
	}
}

func TestSqrtApproxRationalUntilResidual_RejectsNegativeTolerance(t *testing.T) {
	_, _, err := SqrtApproxRationalUntilResidual(
		mustRat(2, 1),
		mustRat(1, 1),
		3,
		mustRat(-1, 10),
	)
	if err == nil {
		t.Fatalf("expected error for negative tolerance")
	}
}

// sqrt_until_tol_test.go v1
