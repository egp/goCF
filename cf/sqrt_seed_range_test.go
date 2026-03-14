// sqrt_seed_range_test.go v1
package cf

import "testing"

func TestSqrtSeedDefault_NegativeRejected(t *testing.T) {
	_, err := SqrtSeedDefault(mustRat(-2, 1))
	if err == nil {
		t.Fatalf("expected error for negative input")
	}
}

func TestSqrtSeedDefault_ZeroReturnsOne(t *testing.T) {
	got, err := SqrtSeedDefault(mustRat(0, 1))
	if err != nil {
		t.Fatalf("SqrtSeedDefault failed: %v", err)
	}
	want := mustRat(1, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSqrtSeedDefault_UnitIntervalReturnsOne(t *testing.T) {
	got, err := SqrtSeedDefault(mustRat(3, 4))
	if err != nil {
		t.Fatalf("SqrtSeedDefault failed: %v", err)
	}
	want := mustRat(1, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSqrtSeedDefault_GreaterThanOneReturnsX(t *testing.T) {
	x := mustRat(5, 2)
	got, err := SqrtSeedDefault(x)
	if err != nil {
		t.Fatalf("SqrtSeedDefault failed: %v", err)
	}
	if got.Cmp(x) != 0 {
		t.Fatalf("got %v want %v", got, x)
	}
}

func TestSqrtRangeMidpoint_InsideRange(t *testing.T) {
	r := NewRange(mustRat(1, 2), mustRat(3, 2), true, true)
	got, err := SqrtRangeMidpoint(r)
	if err != nil {
		t.Fatalf("SqrtRangeMidpoint failed: %v", err)
	}
	want := mustRat(1, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSqrtSeedFromRange_ExactMidpointSquare(t *testing.T) {
	r := NewRange(mustRat(1, 2), mustRat(3, 2), true, true)

	got, err := SqrtSeedFromRange(r)
	if err != nil {
		t.Fatalf("SqrtSeedFromRange failed: %v", err)
	}

	want := mustRat(1, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestSqrtRangeExact2_ExactPointPerfectSquare(t *testing.T) {
	r := NewRange(mustRat(9, 16), mustRat(9, 16), true, true)

	got, ok, err := SqrtRangeExact2(r)
	if err != nil {
		t.Fatalf("SqrtRangeExact2 failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	want := mustRat(3, 4)
	if got.Lo.Cmp(want) != 0 || got.Hi.Cmp(want) != 0 {
		t.Fatalf("got [%v,%v], want exact [%v,%v]", got.Lo, got.Hi, want, want)
	}
}

func TestSqrtRangeExact2_PerfectSquareEndpoints(t *testing.T) {
	r := NewRange(mustRat(1, 4), mustRat(9, 16), true, true)

	got, ok, err := SqrtRangeExact2(r)
	if err != nil {
		t.Fatalf("SqrtRangeExact2 failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	wantLo := mustRat(1, 2)
	wantHi := mustRat(3, 4)
	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v], want [%v,%v]", got.Lo, got.Hi, wantLo, wantHi)
	}
}

func TestSqrtRangeExact2_NonSquareEndpointNotYetSupported(t *testing.T) {
	r := NewRange(mustRat(4, 3), mustRat(3, 2), true, true)

	_, ok, err := SqrtRangeExact2(r)
	if err != nil {
		t.Fatalf("SqrtRangeExact2 failed: %v", err)
	}
	if ok {
		t.Fatalf("did not expect exact sqrt range support here")
	}
}

func TestSqrtRangeExact2_RejectsNegativeRange(t *testing.T) {
	r := NewRange(mustRat(-1, 1), mustRat(1, 1), true, true)

	_, _, err := SqrtRangeExact2(r)
	if err == nil {
		t.Fatalf("expected error for negative range")
	}
}

func TestSqrtRangeHeuristic2_ExactFallsBackToExact(t *testing.T) {
	r := NewRange(mustRat(1, 4), mustRat(9, 16), true, true)

	got, err := SqrtRangeHeuristic2(r)
	if err != nil {
		t.Fatalf("SqrtRangeHeuristic2 failed: %v", err)
	}

	wantLo := mustRat(1, 2)
	wantHi := mustRat(3, 4)
	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v], want [%v,%v]", got.Lo, got.Hi, wantLo, wantHi)
	}
}

func TestSqrtRangeHeuristic2_PreservesEndpointInclusions(t *testing.T) {
	r := NewRange(mustRat(1, 4), mustRat(9, 16), false, false)

	got, err := SqrtRangeHeuristic2(r)
	if err != nil {
		t.Fatalf("SqrtRangeHeuristic2 failed: %v", err)
	}

	if got.IncLo || got.IncHi {
		t.Fatalf("expected inclusions preserved as false,false; got %v %v", got.IncLo, got.IncHi)
	}
}

func TestSqrtRangeHeuristic2_RejectsOutsideRange(t *testing.T) {
	r := NewRange(mustRat(2, 1), mustRat(1, 1), true, true)

	_, err := SqrtRangeHeuristic2(r)
	if err == nil {
		t.Fatalf("expected error for outside range")
	}
}

// sqrt_seed_range_test.go v1
