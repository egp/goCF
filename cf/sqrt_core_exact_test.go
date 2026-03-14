// sqrt_core_exact_test.go v1
package cf

import "testing"

func TestSqrtCoreRationalExact_PerfectSquareInteger(t *testing.T) {
	x := mustRat(9, 1)

	got, ok, err := SqrtCoreRationalExact(x)
	if err != nil {
		t.Fatalf("SqrtCoreRationalExact failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected exact square root")
	}
	want := mustRat(3, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtCoreRationalExact_PerfectSquareRational(t *testing.T) {
	x := mustRat(9, 16)

	got, ok, err := SqrtCoreRationalExact(x)
	if err != nil {
		t.Fatalf("SqrtCoreRationalExact failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected exact square root")
	}
	want := mustRat(3, 4)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtCoreRationalExact_NonSquare(t *testing.T) {
	x := mustRat(2, 1)

	_, ok, err := SqrtCoreRationalExact(x)
	if err != nil {
		t.Fatalf("SqrtCoreRationalExact failed: %v", err)
	}
	if ok {
		t.Fatalf("did not expect exact square root")
	}
}

func TestSqrtCoreNewtonStep_Sqrt2FromSeed1(t *testing.T) {
	x := mustRat(2, 1)
	y := mustRat(1, 1)

	got, err := SqrtCoreNewtonStep(x, y)
	if err != nil {
		t.Fatalf("SqrtCoreNewtonStep failed: %v", err)
	}

	want := mustRat(3, 2)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtCoreNewtonIterates_Sqrt2_FirstThree(t *testing.T) {
	x := mustRat(2, 1)
	seed := mustRat(1, 1)

	got, err := SqrtCoreNewtonIterates(x, seed, 3)
	if err != nil {
		t.Fatalf("SqrtCoreNewtonIterates failed: %v", err)
	}

	want := []Rational{
		mustRat(3, 2),
		mustRat(17, 12),
		mustRat(577, 408),
	}

	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d", len(got), len(want))
	}
	for i := range want {
		if got[i].Cmp(want[i]) != 0 {
			t.Fatalf("got[%d]=%v want=%v", i, got[i], want[i])
		}
	}
}

func TestSqrtCoreNewtonIterates_PerfectSquareFastPath(t *testing.T) {
	x := mustRat(9, 16)
	seed := mustRat(1, 1)

	got, err := SqrtCoreNewtonIterates(x, seed, 3)
	if err != nil {
		t.Fatalf("SqrtCoreNewtonIterates failed: %v", err)
	}

	want := mustRat(3, 4)
	if len(got) != 3 {
		t.Fatalf("len(got)=%d want=3", len(got))
	}
	for i := range got {
		if got[i].Cmp(want) != 0 {
			t.Fatalf("got[%d]=%v want=%v", i, got[i], want)
		}
	}
}

func TestSqrtCoreNewtonStep_RejectsNegativeInput(t *testing.T) {
	_, err := SqrtCoreNewtonStep(mustRat(-1, 1), mustRat(1, 1))
	if err == nil {
		t.Fatalf("expected error for negative input")
	}
}

func TestSqrtCoreNewtonStep_RejectsZeroSeed(t *testing.T) {
	_, err := SqrtCoreNewtonStep(mustRat(2, 1), mustRat(0, 1))
	if err == nil {
		t.Fatalf("expected error for zero seed")
	}
}

func TestSqrtCoreApproxRational_StepsZeroReturnsSeed(t *testing.T) {
	x := mustRat(2, 1)
	seed := mustRat(3, 2)

	got, err := SqrtCoreApproxRational(x, seed, 0)
	if err != nil {
		t.Fatalf("SqrtCoreApproxRational failed: %v", err)
	}
	if got.Cmp(seed) != 0 {
		t.Fatalf("got %v, want %v", got, seed)
	}
}

func TestSqrtCoreResidual_Sqrt2Approx(t *testing.T) {
	x := mustRat(2, 1)
	y := mustRat(3, 2)

	got, err := SqrtCoreResidual(x, y)
	if err != nil {
		t.Fatalf("SqrtCoreResidual failed: %v", err)
	}

	want := mustRat(1, 4)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtCoreResidualAbs_NegativeResidual(t *testing.T) {
	x := mustRat(2, 1)
	y := mustRat(1, 1)

	got, err := SqrtCoreResidualAbs(x, y)
	if err != nil {
		t.Fatalf("SqrtCoreResidualAbs failed: %v", err)
	}

	want := mustRat(1, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtCoreApproxRationalUntilExact_PerfectSquareFastPath(t *testing.T) {
	x := mustRat(9, 16)
	seed := mustRat(1, 1)

	got, exact, err := SqrtCoreApproxRationalUntilExact(x, seed, 5)
	if err != nil {
		t.Fatalf("SqrtCoreApproxRationalUntilExact failed: %v", err)
	}
	if !exact {
		t.Fatalf("expected exact=true")
	}

	want := mustRat(3, 4)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtCoreApproxRationalUntilResidual_PerfectSquareFastPath(t *testing.T) {
	x := mustRat(9, 16)
	seed := mustRat(1, 1)

	got, satisfied, err := SqrtCoreApproxRationalUntilResidual(x, seed, 5, mustRat(0, 1))
	if err != nil {
		t.Fatalf("SqrtCoreApproxRationalUntilResidual failed: %v", err)
	}
	if !satisfied {
		t.Fatalf("expected satisfied=true")
	}

	want := mustRat(3, 4)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// sqrt_core_exact_test.go v1
