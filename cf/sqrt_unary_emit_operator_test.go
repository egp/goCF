// cf/sqrt_unary_emit_operator_test.go v1
package cf

import (
	"math/big"
	"testing"
)

func TestSqrtUnaryOperator_EmitFirstDigitIfForced_InitialStateNotForced(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	d, ok, done, err := op.emitFirstDigitIfForced()
	if err != nil {
		t.Fatalf("emitFirstDigitIfForced failed: %v", err)
	}
	if ok {
		t.Fatalf("unexpected emitted digit %v", d)
	}
	if done {
		t.Fatalf("unexpected done=true")
	}

	snap := op.snapshot()
	if len(snap.EmittedDigits) != 0 {
		t.Fatalf("EmittedDigits: got %v want empty", snap.EmittedDigits)
	}
}

func TestSqrtUnaryOperator_EmitFirstDigitIfForced_AfterOneIngest_EmitsOneAndAdvancesEnclosure(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	if err := op.ingestOneAndRefine(); err != nil {
		t.Fatalf("ingestOneAndRefine failed: %v", err)
	}

	d, ok, done, err := op.emitFirstDigitIfForced()
	if err != nil {
		t.Fatalf("emitFirstDigitIfForced failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected emitted digit")
	}
	if done {
		t.Fatalf("done: got true want false")
	}
	if d == nil || d.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("digit: got %v want 1", d)
	}

	snap := op.snapshot()
	if len(snap.EmittedDigits) != 1 {
		t.Fatalf("EmittedDigits len: got %d want 1", len(snap.EmittedDigits))
	}
	if snap.EmittedDigits[0].Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("EmittedDigits[0]: got %v want 1", snap.EmittedDigits[0])
	}
	if snap.SqrtEnclosure == nil {
		t.Fatalf("SqrtEnclosure: got nil want non-nil")
	}
	if snap.SqrtEnclosure.Lo.Cmp(mustRat(2, 1)) != 0 {
		t.Fatalf("SqrtEnclosure.Lo: got %v want %v", snap.SqrtEnclosure.Lo, mustRat(2, 1))
	}
	if snap.SqrtEnclosure.Hi.Cmp(mustRat(3, 1)) != 0 {
		t.Fatalf("SqrtEnclosure.Hi: got %v want %v", snap.SqrtEnclosure.Hi, mustRat(3, 1))
	}
	if snap.ForcedDigit != nil {
		t.Fatalf("ForcedDigit: got %v want nil", snap.ForcedDigit)
	}
}

func TestSqrtUnaryOperator_EmitFirstDigitIfForced_SecondCallWithoutMoreInfoDoesNotEmitAgain(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	if err := op.ingestOneAndRefine(); err != nil {
		t.Fatalf("ingestOneAndRefine failed: %v", err)
	}

	_, ok, done, err := op.emitFirstDigitIfForced()
	if err != nil {
		t.Fatalf("first emitFirstDigitIfForced failed: %v", err)
	}
	if !ok || done {
		t.Fatalf("first emit: got ok=%v done=%v want ok=true done=false", ok, done)
	}

	d, ok, done, err := op.emitFirstDigitIfForced()
	if err != nil {
		t.Fatalf("second emitFirstDigitIfForced failed: %v", err)
	}
	if ok {
		t.Fatalf("unexpected second emitted digit %v", d)
	}
	if done {
		t.Fatalf("unexpected done=true on second emit attempt")
	}

	snap := op.snapshot()
	if len(snap.EmittedDigits) != 1 {
		t.Fatalf("EmittedDigits len: got %d want 1", len(snap.EmittedDigits))
	}
	if snap.EmittedDigits[0].Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("EmittedDigits[0]: got %v want 1", snap.EmittedDigits[0])
	}
}
