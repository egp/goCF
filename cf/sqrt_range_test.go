// sqrt_range_test.go v1
package cf

import "testing"

func TestSqrtRangeExact_ExactPointPerfectSquare(t *testing.T) {
	r := NewRange(mustRat(9, 16), mustRat(9, 16), true, true)

	got, ok, err := SqrtRangeExact(r)
	if err != nil {
		t.Fatalf("SqrtRangeExact failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	want := mustRat(3, 4)
	if got.Lo.Cmp(want) != 0 || got.Hi.Cmp(want) != 0 {
		t.Fatalf("got [%v,%v], want exact [%v,%v]", got.Lo, got.Hi, want, want)
	}
}

func TestSqrtRangeExact_PerfectSquareEndpoints(t *testing.T) {
	r := NewRange(mustRat(1, 4), mustRat(9, 16), true, true)

	got, ok, err := SqrtRangeExact(r)
	if err != nil {
		t.Fatalf("SqrtRangeExact failed: %v", err)
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

func TestSqrtRangeExact_NonSquareEndpointNotYetSupported(t *testing.T) {
	r := NewRange(mustRat(4, 3), mustRat(3, 2), true, true)

	_, ok, err := SqrtRangeExact(r)
	if err != nil {
		t.Fatalf("SqrtRangeExact failed: %v", err)
	}
	if ok {
		t.Fatalf("did not expect exact sqrt range support here")
	}
}

func TestSqrtRangeExact_RejectsNegativeRange(t *testing.T) {
	r := NewRange(mustRat(-1, 1), mustRat(1, 1), true, true)

	_, _, err := SqrtRangeExact(r)
	if err == nil {
		t.Fatalf("expected error for negative range")
	}
}

func TestSqrtRangeExactFromCFApprox_FiniteExact(t *testing.T) {
	a, err := CFApproxFromPrefix(NewSliceCF(4), 2)
	if err != nil {
		t.Fatalf("CFApproxFromPrefix failed: %v", err)
	}

	got, ok, err := SqrtRangeExactFromCFApprox(a)
	if err != nil {
		t.Fatalf("SqrtRangeExactFromCFApprox failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	want := mustRat(2, 1)
	if got.Lo.Cmp(want) != 0 || got.Hi.Cmp(want) != 0 {
		t.Fatalf("got [%v,%v], want exact [%v,%v]", got.Lo, got.Hi, want, want)
	}
}

func TestSqrtRangeHeuristic_ExactFallsBackToExact(t *testing.T) {
	r := NewRange(mustRat(1, 4), mustRat(9, 16), true, true)

	got, err := SqrtRangeHeuristic(r)
	if err != nil {
		t.Fatalf("SqrtRangeHeuristic failed: %v", err)
	}

	wantLo := mustRat(1, 2)
	wantHi := mustRat(3, 4)
	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v], want [%v,%v]", got.Lo, got.Hi, wantLo, wantHi)
	}
}

func TestSqrtRangeHeuristic_NonExactEndpoints(t *testing.T) {
	r := NewRange(mustRat(4, 3), mustRat(3, 2), true, true)

	got, err := SqrtRangeHeuristic(r)
	if err != nil {
		t.Fatalf("SqrtRangeHeuristic failed: %v", err)
	}

	if !got.IsInside() {
		t.Fatalf("expected inside range, got %v", got)
	}
	if got.Lo.Cmp(intRat(0)) < 0 {
		t.Fatalf("expected nonnegative lower bound, got %v", got.Lo)
	}
	if got.Lo.Cmp(got.Hi) > 0 {
		t.Fatalf("expected ordered range, got [%v,%v]", got.Lo, got.Hi)
	}
}

func TestSqrtRangeHeuristicFromCFApprox_Sqrt2Prefix2(t *testing.T) {
	a, err := CFApproxFromPrefix(Sqrt2CF(), 2)
	if err != nil {
		t.Fatalf("CFApproxFromPrefix failed: %v", err)
	}

	got, err := SqrtRangeHeuristicFromCFApprox(a)
	if err != nil {
		t.Fatalf("SqrtRangeHeuristicFromCFApprox failed: %v", err)
	}

	if !got.IsInside() {
		t.Fatalf("expected inside range, got %v", got)
	}
	if got.Lo.Cmp(intRat(0)) < 0 {
		t.Fatalf("expected nonnegative lower bound, got %v", got.Lo)
	}
}

// sqrt_range_test.go v1
