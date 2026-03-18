// cf/range_big_floor_test.go v1
package cf

import (
	"math/big"
	"testing"
)

func TestRangeFloorBigBounds_PointTwoToTwo(t *testing.T) {
	r := NewRange(mustRat(2, 1), mustRat(2, 1), true, true)

	flo, fhi, err := r.floorBigBounds()
	if err != nil {
		t.Fatalf("floorBigBounds failed: %v", err)
	}

	if flo.Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("flo: got %v want 2", flo)
	}
	if fhi.Cmp(big.NewInt(2)) != 0 {
		t.Fatalf("fhi: got %v want 2", fhi)
	}
}

func TestRangeFloorBigBounds_RangeFourThirdsToThreeHalves(t *testing.T) {
	r := NewRange(mustRat(4, 3), mustRat(3, 2), true, true)

	flo, fhi, err := r.floorBigBounds()
	if err != nil {
		t.Fatalf("floorBigBounds failed: %v", err)
	}

	if flo.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("flo: got %v want 1", flo)
	}
	if fhi.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("fhi: got %v want 1", fhi)
	}
}

func TestRangeFloorBigBounds_RangeNineTenthsToElevenTenths(t *testing.T) {
	r := NewRange(mustRat(9, 10), mustRat(11, 10), true, true)

	flo, fhi, err := r.floorBigBounds()
	if err != nil {
		t.Fatalf("floorBigBounds failed: %v", err)
	}

	if flo.Cmp(big.NewInt(0)) != 0 {
		t.Fatalf("flo: got %v want 0", flo)
	}
	if fhi.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("fhi: got %v want 1", fhi)
	}
}

func TestRangeFloorBigBounds_HugeExactIntegerPoint(t *testing.T) {
	n, ok := new(big.Int).SetString("10000000000000000000000000", 10)
	if !ok {
		t.Fatalf("failed to parse big integer")
	}
	x, err := newRationalBig(n, big.NewInt(1))
	if err != nil {
		t.Fatalf("newRationalBig failed: %v", err)
	}

	r := NewRange(x, x, true, true)

	flo, fhi, err := r.floorBigBounds()
	if err != nil {
		t.Fatalf("floorBigBounds failed: %v", err)
	}

	if flo.Cmp(n) != 0 {
		t.Fatalf("flo: got %v want %v", flo, n)
	}
	if fhi.Cmp(n) != 0 {
		t.Fatalf("fhi: got %v want %v", fhi, n)
	}
}
