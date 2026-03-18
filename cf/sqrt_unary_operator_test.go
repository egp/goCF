// cf/sqrt_unary_operator_test.go v5
package cf

import (
	"strings"
	"testing"
)

func TestNewSqrtUnaryOperator_RejectsNilSource(t *testing.T) {
	_, err := newSqrtUnaryOperator(nil, mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "nil src") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewSqrtUnaryOperator_RejectsNonpositiveInitialIterate(t *testing.T) {
	_, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(0, 1), defaultSqrtUnaryRefinementPolicy())
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "nonpositive iterate") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewSqrtUnaryOperator_RejectsInvalidPolicy(t *testing.T) {
	_, err := newSqrtUnaryOperator(
		NewECFGSource(),
		mustRat(1, 1),
		sqrtUnaryRefinementPolicy{StepsPerInput: 0},
	)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "StepsPerInput must be > 0") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewSqrtUnaryOperator_InitialSnapshotIsEmpty(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	snap := op.snapshot()
	if snap.HasInputApprox {
		t.Fatalf("HasInputApprox: got true want false")
	}
	if snap.InputApprox != nil {
		t.Fatalf("InputApprox: got non-nil want nil")
	}
	if snap.CurrentY == nil {
		t.Fatalf("CurrentY: got nil want non-nil")
	}
	if snap.CurrentY.Cmp(mustRat(1, 1)) != 0 {
		t.Fatalf("CurrentY: got %v want %v", *snap.CurrentY, mustRat(1, 1))
	}
	if snap.Residual != nil {
		t.Fatalf("Residual: got non-nil want nil")
	}
}

func TestSqrtUnaryOperator_DefaultPolicy_UsesOneStepPerInput(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	if err := op.ingestOneAndRefine(); err != nil {
		t.Fatalf("ingestOneAndRefine failed: %v", err)
	}

	snap := op.snapshot()
	if !snap.HasInputApprox {
		t.Fatalf("HasInputApprox: got false want true")
	}
	if snap.InputApprox == nil {
		t.Fatalf("InputApprox: got nil want non-nil")
	}
	if snap.InputApprox.Convergent.Cmp(mustRat(2, 1)) != 0 {
		t.Fatalf("input convergent: got %v want %v", snap.InputApprox.Convergent, mustRat(2, 1))
	}
	if snap.CurrentY == nil {
		t.Fatalf("CurrentY: got nil want non-nil")
	}
	if snap.CurrentY.Cmp(mustRat(3, 2)) != 0 {
		t.Fatalf("CurrentY: got %v want %v", *snap.CurrentY, mustRat(3, 2))
	}
}

func TestSqrtUnaryOperator_TwoStepPolicy_RefinesTwicePerInput(t *testing.T) {
	op, err := newSqrtUnaryOperator(
		NewECFGSource(),
		mustRat(1, 1),
		sqrtUnaryRefinementPolicy{StepsPerInput: 2},
	)
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	if err := op.ingestOneAndRefine(); err != nil {
		t.Fatalf("ingestOneAndRefine failed: %v", err)
	}

	snap := op.snapshot()
	if snap.InputApprox == nil {
		t.Fatalf("InputApprox: got nil want non-nil")
	}
	if snap.InputApprox.Convergent.Cmp(mustRat(2, 1)) != 0 {
		t.Fatalf("input convergent: got %v want %v", snap.InputApprox.Convergent, mustRat(2, 1))
	}

	if snap.CurrentY == nil {
		t.Fatalf("CurrentY: got nil want non-nil")
	}
	if snap.CurrentY.Cmp(mustRat(17, 12)) != 0 {
		t.Fatalf("CurrentY: got %v want %v", *snap.CurrentY, mustRat(17, 12))
	}
}

func TestSqrtUnaryOperator_UsesPrefixStateForInputApproximation(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	if err := op.ingestOneAndRefine(); err != nil {
		t.Fatalf("first ingestOneAndRefine failed: %v", err)
	}

	if op.prefixState == nil {
		t.Fatalf("prefixState: got nil want non-nil")
	}
	if !op.prefixState.hasApprox() {
		t.Fatalf("prefixState.hasApprox: got false want true")
	}

	got := op.snapshot()
	want := op.prefixState.approx()

	if got.InputApprox == nil {
		t.Fatalf("snapshot InputApprox: got nil want non-nil")
	}
	if got.InputApprox.Convergent.Cmp(want.Convergent) != 0 {
		t.Fatalf("snapshot convergent: got %v want %v", got.InputApprox.Convergent, want.Convergent)
	}
}

func TestSqrtUnaryOperator_SnapshotReportsPositiveStateAfterEachRefinement(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	for i := 0; i < 3; i++ {
		if err := op.ingestOneAndRefine(); err != nil {
			t.Fatalf("ingestOneAndRefine #%d failed: %v", i+1, err)
		}

		snap := op.snapshot()
		if !snap.HasInputApprox {
			t.Fatalf("after step %d HasInputApprox: got false want true", i+1)
		}
		if snap.InputApprox == nil {
			t.Fatalf("after step %d InputApprox: got nil want non-nil", i+1)
		}
		if snap.CurrentY == nil {
			t.Fatalf("after step %d CurrentY: got nil want non-nil", i+1)
		}
		if snap.InputApprox.Convergent.Cmp(intRat(0)) <= 0 {
			t.Fatalf("after step %d input convergent: got %v want positive", i+1, snap.InputApprox.Convergent)
		}
		if snap.CurrentY.Cmp(intRat(0)) <= 0 {
			t.Fatalf("after step %d CurrentY: got %v want positive", i+1, *snap.CurrentY)
		}
	}
}

func TestSqrtUnaryOperator_SnapshotCarriesInputRangeWhenAvailable(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	if err := op.ingestOneAndRefine(); err != nil {
		t.Fatalf("ingestOneAndRefine failed: %v", err)
	}

	snap := op.snapshot()
	if snap.InputApprox == nil {
		t.Fatalf("InputApprox: got nil want non-nil")
	}
	if snap.InputApprox.Range == nil {
		t.Fatalf("InputApprox.Range: got nil want non-nil")
	}
}

func TestSqrtUnaryOperator_SnapshotCarriesResidualAfterOneIngest(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	if err := op.ingestOneAndRefine(); err != nil {
		t.Fatalf("ingestOneAndRefine failed: %v", err)
	}

	snap := op.snapshot()
	if snap.Residual == nil {
		t.Fatalf("Residual: got nil want non-nil")
	}
	if snap.Residual.X.Cmp(mustRat(2, 1)) != 0 {
		t.Fatalf("Residual.X: got %v want %v", snap.Residual.X, mustRat(2, 1))
	}
	if snap.Residual.Y.Cmp(mustRat(3, 2)) != 0 {
		t.Fatalf("Residual.Y: got %v want %v", snap.Residual.Y, mustRat(3, 2))
	}
	if snap.Residual.YSquared.Cmp(mustRat(9, 4)) != 0 {
		t.Fatalf("Residual.YSquared: got %v want %v", snap.Residual.YSquared, mustRat(9, 4))
	}
	if snap.Residual.Residual.Cmp(mustRat(1, 4)) != 0 {
		t.Fatalf("Residual.Residual: got %v want %v", snap.Residual.Residual, mustRat(1, 4))
	}
}

func TestSqrtUnaryOperator_SnapshotCarriesResidualAfterTwoIngests(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	if err := op.ingestOneAndRefine(); err != nil {
		t.Fatalf("first ingestOneAndRefine failed: %v", err)
	}
	if err := op.ingestOneAndRefine(); err != nil {
		t.Fatalf("second ingestOneAndRefine failed: %v", err)
	}

	snap := op.snapshot()
	if snap.Residual == nil {
		t.Fatalf("Residual: got nil want non-nil")
	}
	if snap.Residual.X.Cmp(mustRat(3, 1)) != 0 {
		t.Fatalf("Residual.X: got %v want %v", snap.Residual.X, mustRat(3, 1))
	}
	if snap.Residual.Y.Cmp(mustRat(7, 4)) != 0 {
		t.Fatalf("Residual.Y: got %v want %v", snap.Residual.Y, mustRat(7, 4))
	}
	if snap.Residual.YSquared.Cmp(mustRat(49, 16)) != 0 {
		t.Fatalf("Residual.YSquared: got %v want %v", snap.Residual.YSquared, mustRat(49, 16))
	}
	if snap.Residual.Residual.Cmp(mustRat(1, 16)) != 0 {
		t.Fatalf("Residual.Residual: got %v want %v", snap.Residual.Residual, mustRat(1, 16))
	}
}

// cf/sqrt_unary_operator_test.go v5
