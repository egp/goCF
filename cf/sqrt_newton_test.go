// sqrt_newton_test.go v1
package cf

import "testing"

func TestRationalSqrtExact_PerfectSquareInteger(t *testing.T) {
	x := mustRat(9, 1)

	got, ok, err := RationalSqrtExact(x)
	if err != nil {
		t.Fatalf("RationalSqrtExact failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected exact square root")
	}
	want := mustRat(3, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestRationalSqrtExact_PerfectSquareRational(t *testing.T) {
	x := mustRat(9, 16)

	got, ok, err := RationalSqrtExact(x)
	if err != nil {
		t.Fatalf("RationalSqrtExact failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected exact square root")
	}
	want := mustRat(3, 4)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestRationalSqrtExact_NonSquare(t *testing.T) {
	x := mustRat(2, 1)

	_, ok, err := RationalSqrtExact(x)
	if err != nil {
		t.Fatalf("RationalSqrtExact failed: %v", err)
	}
	if ok {
		t.Fatalf("did not expect exact square root")
	}
}

func TestNewtonSqrtStep_Sqrt2FromSeed1(t *testing.T) {
	x := mustRat(2, 1)
	y := mustRat(1, 1)

	got, err := NewtonSqrtStep(x, y)
	if err != nil {
		t.Fatalf("NewtonSqrtStep failed: %v", err)
	}

	// (1 + 2/1)/2 = 3/2
	want := mustRat(3, 2)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestNewtonSqrtIterates_Sqrt2_FirstThree(t *testing.T) {
	x := mustRat(2, 1)
	seed := mustRat(1, 1)

	got, err := NewtonSqrtIterates(x, seed, 3)
	if err != nil {
		t.Fatalf("NewtonSqrtIterates failed: %v", err)
	}

	want := []Rational{
		mustRat(3, 2),   // 1.5
		mustRat(17, 12), // 1.41666...
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

func TestNewtonSqrtIterates_PerfectSquareFastPath(t *testing.T) {
	x := mustRat(9, 16)
	seed := mustRat(1, 1)

	got, err := NewtonSqrtIterates(x, seed, 3)
	if err != nil {
		t.Fatalf("NewtonSqrtIterates failed: %v", err)
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

func TestNewtonSqrtStep_RejectsNegativeInput(t *testing.T) {
	_, err := NewtonSqrtStep(mustRat(-1, 1), mustRat(1, 1))
	if err == nil {
		t.Fatalf("expected error for negative input")
	}
}

func TestNewtonSqrtStep_RejectsZeroSeed(t *testing.T) {
	_, err := NewtonSqrtStep(mustRat(2, 1), mustRat(0, 1))
	if err == nil {
		t.Fatalf("expected error for zero seed")
	}
}

// sqrt_newton_test.go v1
