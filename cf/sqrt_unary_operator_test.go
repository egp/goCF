// cf/sqrt_unary_operator_test.go v1
package cf

import (
	"strings"
	"testing"
)

func TestNewSqrtUnaryOperator_RejectsNilSource(t *testing.T) {
	_, err := newSqrtUnaryOperator(nil, mustRat(1, 1))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "nil src") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewSqrtUnaryOperator_RejectsNonpositiveInitialIterate(t *testing.T) {
	_, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(0, 1))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "nonpositive iterate") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewSqrtUnaryOperator_InitialSnapshotIsEmpty(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1))
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
}

func TestSqrtUnaryOperator_IngestOneAndRefine_UsesOneTermConvergent(t *testing.T) {
	// e source first term is 2, so with initial y=1:
	// y' = (1 + 2/1)/2 = 3/2
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1))
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

func TestSqrtUnaryOperator_TwoIngestsTrackChangingInputApproximation(t *testing.T) {
	// e source first two terms are (2,1), (1,1), giving convergent 3.
	// Start y=1:
	// after first ingest, y = (1 + 2)/2 = 3/2
	// after second ingest with x=3, y = (3/2 + 3/(3/2))/2 = (3/2 + 2)/2 = 7/4
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1))
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
	if snap.InputApprox == nil {
		t.Fatalf("InputApprox: got nil want non-nil")
	}
	if snap.InputApprox.Convergent.Cmp(mustRat(3, 1)) != 0 {
		t.Fatalf("input convergent: got %v want %v", snap.InputApprox.Convergent, mustRat(3, 1))
	}
	if snap.CurrentY == nil {
		t.Fatalf("CurrentY: got nil want non-nil")
	}
	if snap.CurrentY.Cmp(mustRat(7, 4)) != 0 {
		t.Fatalf("CurrentY: got %v want %v", *snap.CurrentY, mustRat(7, 4))
	}
}
