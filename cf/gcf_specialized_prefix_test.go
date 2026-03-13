// gcf_specialized_prefix_test.go v2
package cf

import "testing"

type sourceEvidenceRangeGCFSource struct {
	i int
}

func (s *sourceEvidenceRangeGCFSource) NextPQ() (int64, int64, bool) {
	if s.i >= 1 {
		return 0, 0, false
	}
	s.i++
	return 1, 1, true
}

func (s *sourceEvidenceRangeGCFSource) TailEvidence() GCFTailEvidence {
	r := NewRange(mustRat(2, 1), mustRat(3, 1), true, true)
	return GCFTailEvidence{
		Range:         &r,
		RangeReusable: false,
	}
}

type sourceEvidenceLowerBoundGCFSource struct {
	i int
}

func (s *sourceEvidenceLowerBoundGCFSource) NextPQ() (int64, int64, bool) {
	if s.i >= 1 {
		return 0, 0, false
	}
	s.i++
	return 1, 1, true
}

func (s *sourceEvidenceLowerBoundGCFSource) TailEvidence() GCFTailEvidence {
	lb := mustRat(1, 1)
	return GCFTailEvidence{
		LowerBound:          &lb,
		LowerBoundMinPrefix: 0,
	}
}

func TestSpecializedGCFApproxFromPrefix_LegacyHelperRejectsNegativePrefixTerms(t *testing.T) {
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

func TestSpecializedGCFApproxFromPrefixUsingSourceEvidence_UsesRangeWhenAvailable(t *testing.T) {
	got, err := specializedGCFApproxFromPrefixUsingSourceEvidence(
		1,
		func() GCFSource { return &sourceEvidenceRangeGCFSource{} },
	)
	if err != nil {
		t.Fatalf("specializedGCFApproxFromPrefixUsingSourceEvidence failed: %v", err)
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

func TestSpecializedGCFApproxFromPrefixUsingSourceEvidence_FallsBackToLowerBound(t *testing.T) {
	got, err := specializedGCFApproxFromPrefixUsingSourceEvidence(
		1,
		func() GCFSource { return &sourceEvidenceLowerBoundGCFSource{} },
	)
	if err != nil {
		t.Fatalf("specializedGCFApproxFromPrefixUsingSourceEvidence failed: %v", err)
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

// gcf_specialized_prefix_test.go v2
