// cf/sqrt_unary_range_enclosure_test.go v1
package cf

import (
	"strings"
	"testing"
)

func TestSqrtUnaryRangeEnclosureFromInputRange_PointFourWithIterateOne_MatchesPointEnclosure(t *testing.T) {
	in := NewRange(mustRat(4, 1), mustRat(4, 1), true, true)

	got, err := sqrtUnaryRangeEnclosureFromInputRange(in, mustRat(1, 1))
	if err != nil {
		t.Fatalf("sqrtUnaryRangeEnclosureFromInputRange failed: %v", err)
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
}

func TestSqrtUnaryRangeEnclosureFromInputRange_PointFourWithExactIterate_CollapsesToPoint(t *testing.T) {
	in := NewRange(mustRat(4, 1), mustRat(4, 1), true, true)

	got, err := sqrtUnaryRangeEnclosureFromInputRange(in, mustRat(2, 1))
	if err != nil {
		t.Fatalf("sqrtUnaryRangeEnclosureFromInputRange failed: %v", err)
	}

	if got.Lo.Cmp(mustRat(2, 1)) != 0 || got.Hi.Cmp(mustRat(2, 1)) != 0 {
		t.Fatalf("got %v want point [2,2]", got)
	}
}

func TestSqrtUnaryRangeEnclosureFromInputRange_RangeOneToFour_EnclosesOneAndTwo(t *testing.T) {
	in := NewRange(mustRat(1, 1), mustRat(4, 1), true, true)

	got, err := sqrtUnaryRangeEnclosureFromInputRange(in, mustRat(2, 1))
	if err != nil {
		t.Fatalf("sqrtUnaryRangeEnclosureFromInputRange failed: %v", err)
	}

	if !got.Contains(mustRat(1, 1)) {
		t.Fatalf("expected enclosure %v to contain 1", got)
	}
	if !got.Contains(mustRat(2, 1)) {
		t.Fatalf("expected enclosure %v to contain 2", got)
	}
	if got.Lo.Cmp(got.Hi) > 0 {
		t.Fatalf("expected inside range, got %v", got)
	}
}

func TestSqrtUnaryRangeEnclosureFromInputRange_RejectsOutsideRange(t *testing.T) {
	in := NewRange(mustRat(4, 1), mustRat(1, 1), true, true)

	_, err := sqrtUnaryRangeEnclosureFromInputRange(in, mustRat(1, 1))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "inside range") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSqrtUnaryRangeEnclosureFromInputRange_RejectsNonpositiveLowerBound(t *testing.T) {
	in := NewRange(mustRat(0, 1), mustRat(4, 1), true, true)

	_, err := sqrtUnaryRangeEnclosureFromInputRange(in, mustRat(1, 1))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "positive lower bound") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSqrtUnaryRangeEnclosureFromInputRange_OperatorSnapshotRangeProducesPositiveEnclosure(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	if err := op.ingestOneAndRefine(); err != nil {
		t.Fatalf("ingestOneAndRefine failed: %v", err)
	}

	snap := op.snapshot()
	if snap.InputApprox == nil || snap.InputApprox.Range == nil {
		t.Fatalf("expected input range in snapshot")
	}
	if snap.CurrentY == nil {
		t.Fatalf("expected current iterate in snapshot")
	}

	got, err := sqrtUnaryRangeEnclosureFromInputRange(*snap.InputApprox.Range, *snap.CurrentY)
	if err != nil {
		t.Fatalf("sqrtUnaryRangeEnclosureFromInputRange failed: %v", err)
	}

	if got.Lo.Cmp(intRat(0)) <= 0 {
		t.Fatalf("expected positive lower bound, got %v", got.Lo)
	}
	if got.Lo.Cmp(got.Hi) > 0 {
		t.Fatalf("expected inside range, got %v", got)
	}
}

// cf/sqrt_unary_range_enclosure_test.go v1
