// sqrt_newton_min_test.go v2
package cf

import "testing"

func TestSqrtBootstrapStateForExactFinite_TwoStartsAtOne(t *testing.T) {
	s, err := newSqrtBootstrapState(mustRat(2, 1))
	if err != nil {
		t.Fatalf("newSqrtBootstrapState failed: %v", err)
	}

	if s.x.Cmp(mustRat(2, 1)) != 0 {
		t.Fatalf("x got %v want 2", s.x)
	}
	if s.y.Cmp(mustRat(1, 1)) != 0 {
		t.Fatalf("y got %v want 1", s.y)
	}
}

func TestSqrtBootstrapStateStep_TwoFromOne_IsThreeHalves(t *testing.T) {
	s, err := newSqrtBootstrapState(mustRat(2, 1))
	if err != nil {
		t.Fatalf("newSqrtBootstrapState failed: %v", err)
	}

	if err := s.Step(); err != nil {
		t.Fatalf("Step failed: %v", err)
	}

	want := mustRat(3, 2)
	if s.y.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", s.y, want)
	}
}

func TestSqrtBootstrapStateStepTwice_Two_IsSeventeenTwelfths(t *testing.T) {
	s, err := newSqrtBootstrapState(mustRat(2, 1))
	if err != nil {
		t.Fatalf("newSqrtBootstrapState failed: %v", err)
	}

	if err := s.Step(); err != nil {
		t.Fatalf("first Step failed: %v", err)
	}
	if err := s.Step(); err != nil {
		t.Fatalf("second Step failed: %v", err)
	}

	want := mustRat(17, 12)
	if s.y.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", s.y, want)
	}
}

func TestSqrtBootstrapStateCF_TwoAfterFourSteps_MatchesExpectedIterate(t *testing.T) {
	s, err := newSqrtBootstrapState(mustRat(2, 1))
	if err != nil {
		t.Fatalf("newSqrtBootstrapState failed: %v", err)
	}

	for i := 0; i < 4; i++ {
		if err := s.Step(); err != nil {
			t.Fatalf("Step %d failed: %v", i+1, err)
		}
	}

	cf, err := s.CF()
	if err != nil {
		t.Fatalf("CF failed: %v", err)
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

// sqrt_newton_min_test.go v2
