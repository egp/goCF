// cf/sqrt_unary_demand_test.go v1
package cf

import (
	"math/big"
	"strings"
	"testing"
)

func TestSqrtUnaryOperator_ForceFirstDigitWithin_RejectsNegativeBudget(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	_, _, err = op.forceFirstDigitWithin(-1)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "negative maxIngests") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSqrtUnaryOperator_ForceFirstDigitWithin_ZeroBudgetInitialStateNotForced(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	d, ok, err := op.forceFirstDigitWithin(0)
	if err != nil {
		t.Fatalf("forceFirstDigitWithin failed: %v", err)
	}
	if ok {
		t.Fatalf("unexpected forced digit %v", d)
	}

	snap := op.snapshot()
	if snap.HasInputApprox {
		t.Fatalf("HasInputApprox: got true want false")
	}
}

func TestSqrtUnaryOperator_ForceFirstDigitWithin_OneBudgetForESourceForcesDigitOne(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	d, ok, err := op.forceFirstDigitWithin(1)
	if err != nil {
		t.Fatalf("forceFirstDigitWithin failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected forced digit")
	}
	if d == nil || d.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("got %v want 1", d)
	}
}

func TestSqrtUnaryOperator_ForceFirstDigitWithin_AlreadyForcedDoesNotIngestAgain(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	if err := op.ingestOneAndRefine(); err != nil {
		t.Fatalf("ingestOneAndRefine failed: %v", err)
	}

	before := op.snapshot()

	d, ok, err := op.forceFirstDigitWithin(5)
	if err != nil {
		t.Fatalf("forceFirstDigitWithin failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected forced digit")
	}
	if d == nil || d.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("got %v want 1", d)
	}

	after := op.snapshot()

	if before.InputApprox == nil || after.InputApprox == nil {
		t.Fatalf("expected input approximations before and after")
	}
	if before.InputApprox.Convergent.Cmp(after.InputApprox.Convergent) != 0 {
		t.Fatalf("input convergent changed: before=%v after=%v", before.InputApprox.Convergent, after.InputApprox.Convergent)
	}
	if before.CurrentY == nil || after.CurrentY == nil {
		t.Fatalf("expected current iterates before and after")
	}
	if before.CurrentY.Cmp(*after.CurrentY) != 0 {
		t.Fatalf("current iterate changed: before=%v after=%v", *before.CurrentY, *after.CurrentY)
	}
}
