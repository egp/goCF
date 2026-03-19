// cf/sqrt_unary_multi_emit_test.go v1
package cf

import (
	"math/big"
	"testing"
)

func TestSqrtUnaryOperator_NextDigitIfForced_UsesActiveEnclosureAfterEmit(t *testing.T) {
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
	if !ok || done || d == nil || d.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("first emit: got digit=%v ok=%v done=%v want digit=1 ok=true done=false", d, ok, done)
	}

	// After [4/3, 3/2] emits 1, active enclosure becomes [2,3], which does not force a digit.
	d2, ok2, err := op.nextDigitIfForced()
	if err != nil {
		t.Fatalf("nextDigitIfForced failed: %v", err)
	}
	if ok2 {
		t.Fatalf("unexpected forced digit after emit: %v", d2)
	}
}

func TestSqrtUnaryOperator_EmitForcedDigitsUpTo_InitialStateEmitsNone(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	got, done, err := op.emitForcedDigitsUpTo(5)
	if err != nil {
		t.Fatalf("emitForcedDigitsUpTo failed: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("got %v want empty", got)
	}
	if done {
		t.Fatalf("unexpected done=true")
	}
}

func TestSqrtUnaryOperator_EmitForcedDigitsUpTo_AfterOneIngestEmitsExactlyOneDigit(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	if err := op.ingestOneAndRefine(); err != nil {
		t.Fatalf("ingestOneAndRefine failed: %v", err)
	}

	got, done, err := op.emitForcedDigitsUpTo(5)
	if err != nil {
		t.Fatalf("emitForcedDigitsUpTo failed: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(got): got %d want 1", len(got))
	}
	if got[0].Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("got[0]: got %v want 1", got[0])
	}
	if done {
		t.Fatalf("unexpected done=true")
	}

	snap := op.snapshot()
	if snap.SqrtEnclosure == nil {
		t.Fatalf("SqrtEnclosure: got nil want non-nil")
	}
	if snap.SqrtEnclosure.Lo.Cmp(mustRat(2, 1)) != 0 {
		t.Fatalf("SqrtEnclosure.Lo: got %v want %v", snap.SqrtEnclosure.Lo, mustRat(2, 1))
	}
	if snap.SqrtEnclosure.Hi.Cmp(mustRat(3, 1)) != 0 {
		t.Fatalf("SqrtEnclosure.Hi: got %v want %v", snap.SqrtEnclosure.Hi, mustRat(3, 1))
	}
	if len(snap.EmittedDigits) != 1 {
		t.Fatalf("EmittedDigits len: got %d want 1", len(snap.EmittedDigits))
	}
	if snap.EmittedDigits[0].Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("EmittedDigits[0]: got %v want 1", snap.EmittedDigits[0])
	}
}

func TestSqrtUnaryOperator_EmitForcedDigitsUpTo_ExactIntegerEnclosureTerminates(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	// Seed a terminal active enclosure directly.
	r := NewRange(mustRat(2, 1), mustRat(2, 1), true, true)
	op.currentEnclosure = &r
	op.currentForcedDigit = big.NewInt(2)

	got, done, err := op.emitForcedDigitsUpTo(5)
	if err != nil {
		t.Fatalf("emitForcedDigitsUpTo failed: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len(got): got %d want 1", len(got))
	}
	if got[0].Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("got[0]: got %v want 2", got[0])
	}
	if !done {
		t.Fatalf("done: got false want true")
	}

	snap := op.snapshot()
	if len(snap.EmittedDigits) != 1 {
		t.Fatalf("EmittedDigits len: got %d want 1", len(snap.EmittedDigits))
	}
	if snap.EmittedDigits[0].Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("EmittedDigits[0]: got %v want 2", snap.EmittedDigits[0])
	}
}
