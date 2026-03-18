// cf/sqrt_unary_newton_test.go v2
package cf

import (
	"strings"
	"testing"
)

func TestSqrtUnaryNewtonStep_FourFromOne_GivesFiveHalves(t *testing.T) {
	got, err := SqrtUnaryNewtonStepExact(mustRat(4, 1), mustRat(1, 1))
	if err != nil {
		t.Fatalf("SqrtUnaryNewtonStepExact failed: %v", err)
	}

	want := mustRat(5, 2)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSqrtUnaryNewtonStep_FourSecondStep_GivesFortyOneTwentieths(t *testing.T) {
	got, err := SqrtUnaryNewtonStepExact(mustRat(4, 1), mustRat(5, 2))
	if err != nil {
		t.Fatalf("SqrtUnaryNewtonStepExact failed: %v", err)
	}

	want := mustRat(41, 20)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSqrtUnaryNewtonStep_FourAtFixedPoint_StaysTwo(t *testing.T) {
	got, err := SqrtUnaryNewtonStepExact(mustRat(4, 1), mustRat(2, 1))
	if err != nil {
		t.Fatalf("SqrtUnaryNewtonStepExact failed: %v", err)
	}

	want := mustRat(2, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSqrtUnaryNewtonStep_RejectsZeroInput(t *testing.T) {
	_, err := SqrtUnaryNewtonStepExact(mustRat(0, 1), mustRat(1, 1))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "nonpositive input") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSqrtUnaryNewtonStep_RejectsNegativeInput(t *testing.T) {
	_, err := SqrtUnaryNewtonStepExact(mustRat(-1, 1), mustRat(1, 1))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "nonpositive input") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSqrtUnaryNewtonStep_RejectsZeroIterate(t *testing.T) {
	_, err := SqrtUnaryNewtonStepExact(mustRat(4, 1), mustRat(0, 1))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "nonpositive iterate") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSqrtUnaryNewtonStep_RejectsNegativeIterate(t *testing.T) {
	_, err := SqrtUnaryNewtonStepExact(mustRat(4, 1), mustRat(-1, 1))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "nonpositive iterate") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// cf/sqrt_unary_newton_test.go v2
