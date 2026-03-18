// cf/sqrt_unary_residual_test.go v1
package cf

import "testing"

func TestSqrtUnaryResidual_FourWithIterateTwo_IsZero(t *testing.T) {
	got, err := sqrtUnaryResidualExact(mustRat(4, 1), mustRat(2, 1))
	if err != nil {
		t.Fatalf("sqrtUnaryResidualExact failed: %v", err)
	}

	want := mustRat(0, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSqrtUnaryResidual_FourWithIterateOne_IsNegativeThree(t *testing.T) {
	got, err := sqrtUnaryResidualExact(mustRat(4, 1), mustRat(1, 1))
	if err != nil {
		t.Fatalf("sqrtUnaryResidualExact failed: %v", err)
	}

	want := mustRat(-3, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSqrtUnaryResidual_FourWithIterateThree_IsPositiveFive(t *testing.T) {
	got, err := sqrtUnaryResidualExact(mustRat(4, 1), mustRat(3, 1))
	if err != nil {
		t.Fatalf("sqrtUnaryResidualExact failed: %v", err)
	}

	want := mustRat(5, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSqrtUnaryResidualSnapshot_TracksSquareAndResidual(t *testing.T) {
	s, err := newSqrtUnaryState(mustRat(4, 1), mustRat(3, 2))
	if err != nil {
		t.Fatalf("newSqrtUnaryState failed: %v", err)
	}

	snap, err := s.residualSnapshot()
	if err != nil {
		t.Fatalf("residualSnapshot failed: %v", err)
	}

	if snap.X.Cmp(mustRat(4, 1)) != 0 {
		t.Fatalf("X: got %v want %v", snap.X, mustRat(4, 1))
	}
	if snap.Y.Cmp(mustRat(3, 2)) != 0 {
		t.Fatalf("Y: got %v want %v", snap.Y, mustRat(3, 2))
	}
	if snap.YSquared.Cmp(mustRat(9, 4)) != 0 {
		t.Fatalf("YSquared: got %v want %v", snap.YSquared, mustRat(9, 4))
	}
	if snap.Residual.Cmp(mustRat(-7, 4)) != 0 {
		t.Fatalf("Residual: got %v want %v", snap.Residual, mustRat(-7, 4))
	}
}

// cf/sqrt_unary_residual_test.go v1
