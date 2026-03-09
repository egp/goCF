// gcf_specialized_prefix_test.go v1
package cf

import "testing"

func TestSpecializedGCFApproxFromPrefix_RejectsNegativePrefixTerms(t *testing.T) {
	_, err := specializedGCFApproxFromPrefix(
		-1,
		func() GCFSource { return NewSliceGCF([2]int64{1, 1}) },
		func(prefixTerms int) (Range, bool, error) { return Range{}, false, nil },
		func(prefixTerms int) Rational { return mustRat(1, 1) },
	)
	if err == nil {
		t.Fatalf("expected error for negative prefixTerms")
	}
}

func TestSpecializedGCFApproxFromPrefix_UsesSpecializedRangeWhenAvailable(t *testing.T) {
	got, err := specializedGCFApproxFromPrefix(
		1,
		func() GCFSource { return NewSliceGCF([2]int64{1, 1}) },
		func(prefixTerms int) (Range, bool, error) {
			return NewRange(mustRat(2, 1), mustRat(3, 1), true, true), true, nil
		},
		func(prefixTerms int) Rational { return mustRat(1, 1) },
	)
	if err != nil {
		t.Fatalf("specializedGCFApproxFromPrefix failed: %v", err)
	}
	if got.Range == nil {
		t.Fatalf("expected non-nil range")
	}
	wantLo := mustRat(4, 3) // x = 1 + 1/tail with tail in [2,3]
	wantHi := mustRat(3, 2)
	if got.Range.Lo.Cmp(wantLo) != 0 || got.Range.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got range [%v,%v] want [%v,%v]", got.Range.Lo, got.Range.Hi, wantLo, wantHi)
	}
}

func TestSpecializedGCFApproxFromPrefix_FallsBackToLowerBound(t *testing.T) {
	got, err := specializedGCFApproxFromPrefix(
		1,
		func() GCFSource { return NewSliceGCF([2]int64{1, 1}) },
		func(prefixTerms int) (Range, bool, error) { return Range{}, false, nil },
		func(prefixTerms int) Rational { return mustRat(1, 1) },
	)
	if err != nil {
		t.Fatalf("specializedGCFApproxFromPrefix failed: %v", err)
	}
	if got.Range == nil {
		t.Fatalf("expected non-nil range")
	}
	wantLo := mustRat(1, 1) // x = 1 + 1/tail, tail >= 1
	wantHi := mustRat(2, 1)
	if got.Range.Lo.Cmp(wantLo) != 0 || got.Range.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got range [%v,%v] want [%v,%v]", got.Range.Lo, got.Range.Hi, wantLo, wantHi)
	}
}

// gcf_specialized_prefix_test.go v1
