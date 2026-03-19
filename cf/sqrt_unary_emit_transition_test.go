// cf/sqrt_unary_emit_transition_test.go v1
package cf

import (
	"math/big"
	"strings"
	"testing"
)

func TestSqrtUnaryEmitForcedDigitTransition_ExactIntegerPoint_Terminates(t *testing.T) {
	r := NewRange(mustRat(2, 1), mustRat(2, 1), true, true)

	d, rem, done, err := sqrtUnaryEmitForcedDigitTransition(r)
	if err != nil {
		t.Fatalf("sqrtUnaryEmitForcedDigitTransition failed: %v", err)
	}
	if d == nil {
		t.Fatalf("digit: got nil want non-nil")
	}
	if d.Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("digit: got %v want 2", d)
	}
	if !done {
		t.Fatalf("done: got false want true")
	}
	if rem != nil {
		t.Fatalf("remainder: got %v want nil", rem)
	}
}

func TestSqrtUnaryEmitForcedDigitTransition_NonterminalInterval_ProducesReciprocalRemainder(t *testing.T) {
	r := NewRange(mustRat(4, 3), mustRat(3, 2), true, true)

	d, rem, done, err := sqrtUnaryEmitForcedDigitTransition(r)
	if err != nil {
		t.Fatalf("sqrtUnaryEmitForcedDigitTransition failed: %v", err)
	}
	if d == nil {
		t.Fatalf("digit: got nil want non-nil")
	}
	if d.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("digit: got %v want 1", d)
	}
	if done {
		t.Fatalf("done: got true want false")
	}
	if rem == nil {
		t.Fatalf("remainder: got nil want non-nil")
	}

	// [4/3, 3/2] -> subtract 1 => [1/3, 1/2] -> reciprocal => [2, 3]
	if rem.Lo.Cmp(mustRat(2, 1)) != 0 {
		t.Fatalf("remainder.Lo: got %v want %v", rem.Lo, mustRat(2, 1))
	}
	if rem.Hi.Cmp(mustRat(3, 1)) != 0 {
		t.Fatalf("remainder.Hi: got %v want %v", rem.Hi, mustRat(3, 1))
	}
	if !rem.IncLo || !rem.IncHi {
		t.Fatalf("expected closed remainder interval, got %v", *rem)
	}
}

func TestSqrtUnaryEmitForcedDigitTransition_NotForced_Rejects(t *testing.T) {
	r := NewRange(mustRat(9, 10), mustRat(11, 10), true, true)

	_, _, _, err := sqrtUnaryEmitForcedDigitTransition(r)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "digit not forced") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSqrtUnaryEmitForcedDigitTransition_OperatorOneIngestEnclosure_ProducesRemainder(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	if err := op.ingestOneAndRefine(); err != nil {
		t.Fatalf("ingestOneAndRefine failed: %v", err)
	}

	snap := op.snapshot()
	if snap.SqrtEnclosure == nil {
		t.Fatalf("SqrtEnclosure: got nil want non-nil")
	}

	d, rem, done, err := sqrtUnaryEmitForcedDigitTransition(*snap.SqrtEnclosure)
	if err != nil {
		t.Fatalf("sqrtUnaryEmitForcedDigitTransition failed: %v", err)
	}
	if d == nil || d.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("digit: got %v want 1", d)
	}
	if done {
		t.Fatalf("done: got true want false")
	}
	if rem == nil {
		t.Fatalf("remainder: got nil want non-nil")
	}
	if rem.Lo.Cmp(mustRat(2, 1)) != 0 {
		t.Fatalf("remainder.Lo: got %v want %v", rem.Lo, mustRat(2, 1))
	}
	if rem.Hi.Cmp(mustRat(3, 1)) != 0 {
		t.Fatalf("remainder.Hi: got %v want %v", rem.Hi, mustRat(3, 1))
	}
}

// cf/sqrt_unary_emit_transition_test.go v1
