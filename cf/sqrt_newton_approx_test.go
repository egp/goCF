// sqrt_newton_approx_test.go v1
package cf

import "testing"

func TestSqrtApproxRational_Sqrt2_ThreeSteps(t *testing.T) {
	x := mustRat(2, 1)
	seed := mustRat(1, 1)

	got, err := SqrtApproxRational(x, seed, 3)
	if err != nil {
		t.Fatalf("SqrtApproxRational failed: %v", err)
	}

	want := mustRat(577, 408)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtApproxRational_PerfectSquareFastPath(t *testing.T) {
	x := mustRat(9, 16)
	seed := mustRat(1, 1)

	got, err := SqrtApproxRational(x, seed, 5)
	if err != nil {
		t.Fatalf("SqrtApproxRational failed: %v", err)
	}

	want := mustRat(3, 4)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtApproxRational_ZeroStepsReturnsSeed(t *testing.T) {
	x := mustRat(2, 1)
	seed := mustRat(7, 5)

	got, err := SqrtApproxRational(x, seed, 0)
	if err != nil {
		t.Fatalf("SqrtApproxRational failed: %v", err)
	}

	if got.Cmp(seed) != 0 {
		t.Fatalf("got %v, want seed %v", got, seed)
	}
}

func TestSqrtApproxRational_RejectsNegativeInput(t *testing.T) {
	_, err := SqrtApproxRational(mustRat(-2, 1), mustRat(1, 1), 3)
	if err == nil {
		t.Fatalf("expected error for negative input")
	}
}

func TestSqrtApproxRational_RejectsZeroSeed(t *testing.T) {
	_, err := SqrtApproxRational(mustRat(2, 1), mustRat(0, 1), 3)
	if err == nil {
		t.Fatalf("expected error for zero seed")
	}
}

// sqrt_newton_approx_test.go v1
