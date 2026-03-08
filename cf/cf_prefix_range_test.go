// cf_prefix_range_test.go v1
package cf

import "testing"

func TestRangeApproxFromCFPrefix_FiniteSourceExact(t *testing.T) {
	// 355/113 = [3; 7, 16]
	got, err := RangeApproxFromCFPrefix(NewSliceCF(3, 7, 16), 10)
	if err != nil {
		t.Fatalf("RangeApproxFromCFPrefix failed: %v", err)
	}

	want := mustRat(355, 113)
	if got.Lo.Cmp(want) != 0 || got.Hi.Cmp(want) != 0 {
		t.Fatalf("got [%v,%v], want exact [%v,%v]", got.Lo, got.Hi, want, want)
	}
}

func TestRangeApproxFromCFPrefix_InfiniteSourceSingleTerm(t *testing.T) {
	// sqrt(2) = [1;2,2,...], after one term the enclosure is [1,2]
	got, err := RangeApproxFromCFPrefix(Sqrt2CF(), 1)
	if err != nil {
		t.Fatalf("RangeApproxFromCFPrefix failed: %v", err)
	}

	wantLo := mustRat(1, 1)
	wantHi := mustRat(2, 1)
	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v], want [%v,%v]", got.Lo, got.Hi, wantLo, wantHi)
	}
}

func TestRangeApproxFromCFPrefix_InfiniteSourceTwoTerms(t *testing.T) {
	// sqrt(2) prefix [1;2] => enclosure [4/3, 3/2]
	got, err := RangeApproxFromCFPrefix(Sqrt2CF(), 2)
	if err != nil {
		t.Fatalf("RangeApproxFromCFPrefix failed: %v", err)
	}

	wantLo := mustRat(4, 3)
	wantHi := mustRat(3, 2)
	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v], want [%v,%v]", got.Lo, got.Hi, wantLo, wantHi)
	}
}

func TestRangeApproxFromCFPrefix_RejectsZeroPrefixTerms(t *testing.T) {
	_, err := RangeApproxFromCFPrefix(Sqrt2CF(), 0)
	if err == nil {
		t.Fatalf("expected error for zero prefixTerms")
	}
}

func TestRangeApproxFromCFPrefix_RejectsEmptySource(t *testing.T) {
	_, err := RangeApproxFromCFPrefix(NewSliceCF(), 1)
	if err == nil {
		t.Fatalf("expected error for empty source")
	}
}

// cf_prefix_range_test.go v1
