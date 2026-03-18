// cf/sqrt_unary_newton_state_test.go v2
package cf

import (
	"strings"
	"testing"
)

func TestNewSqrtUnaryNewtonState_StoresInitialState(t *testing.T) {
	s, err := newSqrtUnaryNewtonState(mustRat(4, 1), mustRat(1, 1))
	if err != nil {
		t.Fatalf("newSqrtUnaryNewtonState failed: %v", err)
	}

	if s.xValue().Cmp(mustRat(4, 1)) != 0 {
		t.Fatalf("xValue: got %v want %v", s.xValue(), mustRat(4, 1))
	}
	if s.yValue().Cmp(mustRat(1, 1)) != 0 {
		t.Fatalf("yValue: got %v want %v", s.yValue(), mustRat(1, 1))
	}
}

func TestNewSqrtUnaryNewtonState_RejectsZeroInput(t *testing.T) {
	_, err := newSqrtUnaryNewtonState(mustRat(0, 1), mustRat(1, 1))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "nonpositive input") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewSqrtUnaryNewtonState_RejectsNegativeInput(t *testing.T) {
	_, err := newSqrtUnaryNewtonState(mustRat(-1, 1), mustRat(1, 1))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "nonpositive input") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewSqrtUnaryNewtonState_RejectsZeroIterate(t *testing.T) {
	_, err := newSqrtUnaryNewtonState(mustRat(4, 1), mustRat(0, 1))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "nonpositive iterate") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewSqrtUnaryNewtonState_RejectsNegativeIterate(t *testing.T) {
	_, err := newSqrtUnaryNewtonState(mustRat(4, 1), mustRat(-1, 1))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "nonpositive iterate") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSqrtUnaryNewtonState_Step_UpdatesOnlyIterate(t *testing.T) {
	s, err := newSqrtUnaryNewtonState(mustRat(4, 1), mustRat(1, 1))
	if err != nil {
		t.Fatalf("newSqrtUnaryNewtonState failed: %v", err)
	}

	if err := s.step(); err != nil {
		t.Fatalf("step failed: %v", err)
	}

	if s.xValue().Cmp(mustRat(4, 1)) != 0 {
		t.Fatalf("x changed: got %v want %v", s.xValue(), mustRat(4, 1))
	}
	if s.yValue().Cmp(mustRat(5, 2)) != 0 {
		t.Fatalf("y after one step: got %v want %v", s.yValue(), mustRat(5, 2))
	}
}

func TestSqrtUnaryNewtonState_TwoStepsOnFourFromOne(t *testing.T) {
	s, err := newSqrtUnaryNewtonState(mustRat(4, 1), mustRat(1, 1))
	if err != nil {
		t.Fatalf("newSqrtUnaryNewtonState failed: %v", err)
	}

	if err := s.step(); err != nil {
		t.Fatalf("first step failed: %v", err)
	}
	if err := s.step(); err != nil {
		t.Fatalf("second step failed: %v", err)
	}

	if s.yValue().Cmp(mustRat(41, 20)) != 0 {
		t.Fatalf("y after two steps: got %v want %v", s.yValue(), mustRat(41, 20))
	}
}

func TestSqrtUnaryNewtonState_FixedPointRemainsFixed(t *testing.T) {
	s, err := newSqrtUnaryNewtonState(mustRat(4, 1), mustRat(2, 1))
	if err != nil {
		t.Fatalf("newSqrtUnaryNewtonState failed: %v", err)
	}

	if err := s.step(); err != nil {
		t.Fatalf("step failed: %v", err)
	}

	if s.yValue().Cmp(mustRat(2, 1)) != 0 {
		t.Fatalf("y after fixed-point step: got %v want %v", s.yValue(), mustRat(2, 1))
	}
}

// cf/sqrt_unary_newton_state_test.go v2
