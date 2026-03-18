// cf/sqrt_unary_enclosure_test.go v1
package cf

import "testing"

func TestSqrtUnaryPointEnclosure_FourWithIterateOne_BracketsTwo(t *testing.T) {
	got, err := sqrtUnaryPointEnclosureExact(mustRat(4, 1), mustRat(1, 1))
	if err != nil {
		t.Fatalf("sqrtUnaryPointEnclosureExact failed: %v", err)
	}

	if got.Lo.Cmp(mustRat(1, 1)) != 0 {
		t.Fatalf("Lo: got %v want %v", got.Lo, mustRat(1, 1))
	}
	if got.Hi.Cmp(mustRat(4, 1)) != 0 {
		t.Fatalf("Hi: got %v want %v", got.Hi, mustRat(4, 1))
	}
	if !got.IncLo || !got.IncHi {
		t.Fatalf("expected closed interval, got %v", got)
	}
	if !got.Contains(mustRat(2, 1)) {
		t.Fatalf("expected enclosure %v to contain 2", got)
	}
}

func TestSqrtUnaryPointEnclosure_FourWithIterateThree_BracketsTwo(t *testing.T) {
	got, err := sqrtUnaryPointEnclosureExact(mustRat(4, 1), mustRat(3, 1))
	if err != nil {
		t.Fatalf("sqrtUnaryPointEnclosureExact failed: %v", err)
	}

	if got.Lo.Cmp(mustRat(4, 3)) != 0 {
		t.Fatalf("Lo: got %v want %v", got.Lo, mustRat(4, 3))
	}
	if got.Hi.Cmp(mustRat(3, 1)) != 0 {
		t.Fatalf("Hi: got %v want %v", got.Hi, mustRat(3, 1))
	}
	if !got.Contains(mustRat(2, 1)) {
		t.Fatalf("expected enclosure %v to contain 2", got)
	}
}

func TestSqrtUnaryPointEnclosure_FourWithExactIterateCollapsesToPoint(t *testing.T) {
	got, err := sqrtUnaryPointEnclosureExact(mustRat(4, 1), mustRat(2, 1))
	if err != nil {
		t.Fatalf("sqrtUnaryPointEnclosureExact failed: %v", err)
	}

	if got.Lo.Cmp(mustRat(2, 1)) != 0 || got.Hi.Cmp(mustRat(2, 1)) != 0 {
		t.Fatalf("got %v want point [2,2]", got)
	}
	if !got.IncLo || !got.IncHi {
		t.Fatalf("expected closed point interval, got %v", got)
	}
}

func TestSqrtUnaryPointEnclosure_UsesIterateAndReciprocalSymmetrically(t *testing.T) {
	a, err := sqrtUnaryPointEnclosureExact(mustRat(4, 1), mustRat(1, 1))
	if err != nil {
		t.Fatalf("sqrtUnaryPointEnclosureExact(a) failed: %v", err)
	}
	b, err := sqrtUnaryPointEnclosureExact(mustRat(4, 1), mustRat(4, 1))
	if err != nil {
		t.Fatalf("sqrtUnaryPointEnclosureExact(b) failed: %v", err)
	}

	if a.Lo.Cmp(b.Lo) != 0 || a.Hi.Cmp(b.Hi) != 0 {
		t.Fatalf("expected same enclosure, got a=%v b=%v", a, b)
	}
}

func TestSqrtUnaryOperator_PointInputResidualSupportsPointEnclosure(t *testing.T) {
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
	if snap.Residual == nil {
		t.Fatalf("Residual: got nil want non-nil")
	}

	got, err := sqrtUnaryPointEnclosureExact(snap.Residual.X, snap.Residual.Y)
	if err != nil {
		t.Fatalf("sqrtUnaryPointEnclosureExact failed: %v", err)
	}

	// x = 2, y = 3/2 gives enclosure [4/3, 3/2], which contains sqrt(2).
	if got.Lo.Cmp(mustRat(4, 3)) != 0 {
		t.Fatalf("Lo: got %v want %v", got.Lo, mustRat(4, 3))
	}
	if got.Hi.Cmp(mustRat(3, 2)) != 0 {
		t.Fatalf("Hi: got %v want %v", got.Hi, mustRat(3, 2))
	}
}
