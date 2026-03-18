// sqrt_newton_min_test.go v1
package cf

import "testing"

func TestSqrtNewtonStep_TwoFromOne_IsThreeHalves(t *testing.T) {
	got, err := sqrtNewtonStep(mustRat(2, 1), mustRat(1, 1))
	if err != nil {
		t.Fatalf("sqrtNewtonStep failed: %v", err)
	}

	want := mustRat(3, 2)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSqrtNewtonApprox_TwoTwoSteps_IsSeventeenTwelfths(t *testing.T) {
	got, err := sqrtNewtonApprox(mustRat(2, 1), 2)
	if err != nil {
		t.Fatalf("sqrtNewtonApprox failed: %v", err)
	}

	want := mustRat(17, 12)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSqrtNewtonApprox_FourOneStep_IsFiveHalves(t *testing.T) {
	got, err := sqrtNewtonApprox(mustRat(4, 1), 1)
	if err != nil {
		t.Fatalf("sqrtNewtonApprox failed: %v", err)
	}

	want := mustRat(5, 2)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSqrtNewtonApprox_FourManySteps_ImprovesTowardTwo(t *testing.T) {
	got, err := sqrtNewtonApprox(mustRat(4, 1), 4)
	if err != nil {
		t.Fatalf("sqrtNewtonApprox failed: %v", err)
	}

	two := mustRat(2, 1)
	if got.Cmp(two) == 0 {
		t.Fatalf("expected finite Newton iterate started from 1 to remain inexact after 4 steps")
	}

	diff, err := got.Sub(two)
	if err != nil {
		t.Fatalf("Sub failed: %v", err)
	}
	if diff.Cmp(intRat(0)) <= 0 {
		t.Fatalf("got %v want iterate above 2 from initial seed 1", got)
	}

	want := mustRat(21523361, 10761680)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSqrtGCF_ExactFiniteTwo_ReturnsNewtonApproxCF(t *testing.T) {
	cf, err := SqrtGCF(NewSliceGCF([2]int64{2, 1}))
	if err != nil {
		t.Fatalf("SqrtGCF failed: %v", err)
	}

	got := collectTerms(cf, 8)
	want := collectTerms(NewRationalCF(mustRat(577, 408)), 8)

	if len(got) != len(want) {
		t.Fatalf("len mismatch: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, got, want)
		}
	}
}

// sqrt_newton_min_test.go v1
