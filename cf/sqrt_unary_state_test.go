// cf/sqrt_unary_state_test.go v1
package cf

import (
	"strings"
	"testing"
)

func TestNewSqrtUnaryState_StoresInitialState(t *testing.T) {
	s, err := newSqrtUnaryState(mustRat(4, 1), mustRat(1, 1))
	if err != nil {
		t.Fatalf("newSqrtUnaryState failed: %v", err)
	}

	if s.xValue().Cmp(mustRat(4, 1)) != 0 {
		t.Fatalf("xValue: got %v want %v", s.xValue(), mustRat(4, 1))
	}
	if s.yValue().Cmp(mustRat(1, 1)) != 0 {
		t.Fatalf("yValue: got %v want %v", s.yValue(), mustRat(1, 1))
	}
}

func TestNewSqrtUnaryState_RejectsZeroInput(t *testing.T) {
	_, err := newSqrtUnaryState(mustRat(0, 1), mustRat(1, 1))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "nonpositive input") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewSqrtUnaryState_RejectsZeroIterate(t *testing.T) {
	_, err := newSqrtUnaryState(mustRat(4, 1), mustRat(0, 1))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "nonpositive iterate") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSqrtUnaryState_Step_UsesCurrentInputApproximation(t *testing.T) {
	s, err := newSqrtUnaryState(mustRat(4, 1), mustRat(1, 1))
	if err != nil {
		t.Fatalf("newSqrtUnaryState failed: %v", err)
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

func TestSqrtUnaryState_UpdateInput_ReplacesCurrentInput(t *testing.T) {
	s, err := newSqrtUnaryState(mustRat(4, 1), mustRat(1, 1))
	if err != nil {
		t.Fatalf("newSqrtUnaryState failed: %v", err)
	}

	if err := s.updateInput(mustRat(9, 1)); err != nil {
		t.Fatalf("updateInput failed: %v", err)
	}

	if s.xValue().Cmp(mustRat(9, 1)) != 0 {
		t.Fatalf("x after update: got %v want %v", s.xValue(), mustRat(9, 1))
	}
	if s.yValue().Cmp(mustRat(1, 1)) != 0 {
		t.Fatalf("y changed on input update: got %v want %v", s.yValue(), mustRat(1, 1))
	}
}

func TestSqrtUnaryState_UpdateInputThenStep_UsesUpdatedInput(t *testing.T) {
	s, err := newSqrtUnaryState(mustRat(4, 1), mustRat(1, 1))
	if err != nil {
		t.Fatalf("newSqrtUnaryState failed: %v", err)
	}

	if err := s.updateInput(mustRat(9, 1)); err != nil {
		t.Fatalf("updateInput failed: %v", err)
	}
	if err := s.step(); err != nil {
		t.Fatalf("step failed: %v", err)
	}

	// Using x=9 and y=1:
	// y' = (1 + 9/1) / 2 = 5
	if s.yValue().Cmp(mustRat(5, 1)) != 0 {
		t.Fatalf("y after update+step: got %v want %v", s.yValue(), mustRat(5, 1))
	}
}

func TestSqrtUnaryState_UpdateInput_RejectsZeroInput(t *testing.T) {
	s, err := newSqrtUnaryState(mustRat(4, 1), mustRat(1, 1))
	if err != nil {
		t.Fatalf("newSqrtUnaryState failed: %v", err)
	}

	err = s.updateInput(mustRat(0, 1))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "nonpositive input") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSqrtUnaryState_UpdateInput_RejectsNegativeInput(t *testing.T) {
	s, err := newSqrtUnaryState(mustRat(4, 1), mustRat(1, 1))
	if err != nil {
		t.Fatalf("newSqrtUnaryState failed: %v", err)
	}

	err = s.updateInput(mustRat(-1, 1))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "nonpositive input") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// cf/sqrt_unary_state_test.go v1
