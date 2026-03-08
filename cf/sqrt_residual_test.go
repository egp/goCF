// sqrt_residual_test.go v1
package cf

import "testing"

func TestSqrtResidual_ExactZeroForPerfectRoot(t *testing.T) {
	x := mustRat(9, 16)
	y := mustRat(3, 4)

	got, err := SqrtResidual(x, y)
	if err != nil {
		t.Fatalf("SqrtResidual failed: %v", err)
	}

	if got.Cmp(intRat(0)) != 0 {
		t.Fatalf("got %v, want 0", got)
	}
}

func TestSqrtResidual_Sqrt2_FirstNewtonStep(t *testing.T) {
	x := mustRat(2, 1)
	y := mustRat(3, 2) // first Newton iterate from seed 1

	got, err := SqrtResidual(x, y)
	if err != nil {
		t.Fatalf("SqrtResidual failed: %v", err)
	}

	// (3/2)^2 - 2 = 9/4 - 2 = 1/4
	want := mustRat(1, 4)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtResidualAbs_Sqrt2_FirstNewtonStep(t *testing.T) {
	x := mustRat(2, 1)
	y := mustRat(3, 2)

	got, err := SqrtResidualAbs(x, y)
	if err != nil {
		t.Fatalf("SqrtResidualAbs failed: %v", err)
	}

	want := mustRat(1, 4)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtResidualAbs_ImprovesAcrossNewtonSteps(t *testing.T) {
	x := mustRat(2, 1)
	seed := mustRat(1, 1)

	ys, err := NewtonSqrtIterates(x, seed, 3)
	if err != nil {
		t.Fatalf("NewtonSqrtIterates failed: %v", err)
	}

	r1, err := SqrtResidualAbs(x, ys[0])
	if err != nil {
		t.Fatalf("SqrtResidualAbs step1 failed: %v", err)
	}
	r2, err := SqrtResidualAbs(x, ys[1])
	if err != nil {
		t.Fatalf("SqrtResidualAbs step2 failed: %v", err)
	}
	r3, err := SqrtResidualAbs(x, ys[2])
	if err != nil {
		t.Fatalf("SqrtResidualAbs step3 failed: %v", err)
	}

	// Expect strict improvement for these first three sqrt(2) iterates from seed 1.
	if r2.Cmp(r1) >= 0 {
		t.Fatalf("expected r2 < r1, got r1=%v r2=%v", r1, r2)
	}
	if r3.Cmp(r2) >= 0 {
		t.Fatalf("expected r3 < r2, got r2=%v r3=%v", r2, r3)
	}
}

// sqrt_residual_test.go v1
