// sqrt_stream_min_test.go v3
package cf

import (
	"strings"
	"testing"
)

type countingSliceGCF struct {
	terms [][2]int64
	i     int
	reads int
}

func newCountingSliceGCF(terms ...[2]int64) *countingSliceGCF {
	cp := append([][2]int64(nil), terms...)
	return &countingSliceGCF{terms: cp}
}

func (s *countingSliceGCF) NextPQ() (int64, int64, bool) {
	s.reads++
	if s.i >= len(s.terms) {
		return 0, 0, false
	}
	t := s.terms[s.i]
	s.i++
	return t[0], t[1], true
}

type countingRadicalCF struct {
	terms    []int64
	i        int
	reads    int
	radicand int64
}

func newCountingRadicalCF(radicand int64, terms ...int64) *countingRadicalCF {
	cp := append([]int64(nil), terms...)
	return &countingRadicalCF{terms: cp, radicand: radicand}
}

func (s *countingRadicalCF) Next() (int64, bool) {
	s.reads++
	if s.i >= len(s.terms) {
		return 0, false
	}
	v := s.terms[s.i]
	s.i++
	return v, true
}

func (s *countingRadicalCF) Radicand() (int64, bool) {
	return s.radicand, true
}

func repeatedUnitTerms(n int) [][2]int64 {
	out := make([][2]int64, n)
	for i := range out {
		out[i] = [2]int64{1, 1}
	}
	return out
}

func TestSqrtGCF_IsLazyBeforeFirstNext(t *testing.T) {
	src := newCountingSliceGCF([2]int64{4, 1})

	cf, err := SqrtGCF(src)
	if err != nil {
		t.Fatalf("SqrtGCF failed: %v", err)
	}
	if cf == nil {
		t.Fatalf("got nil cf")
	}
	if src.reads != 0 {
		t.Fatalf("reads got %d want 0 before first Next", src.reads)
	}
}

func TestSqrtGCF_FirstNextConsumesInputAndEmitsFirstDigit(t *testing.T) {
	src := newCountingSliceGCF([2]int64{4, 1})

	cf, err := SqrtGCF(src)
	if err != nil {
		t.Fatalf("SqrtGCF failed: %v", err)
	}

	d, ok := cf.Next()
	if !ok {
		t.Fatalf("expected first digit")
	}
	if d != 2 {
		t.Fatalf("digit got %d want 2", d)
	}
	if src.reads == 0 {
		t.Fatalf("expected lazy consumption on first Next")
	}
}

func TestSqrtGCF_LazyMetadataPath_SqrtFourEmitsTwoWithoutReadingInputTerms(t *testing.T) {
	src := newCountingRadicalCF(4, 2, 1, 1)

	cf, err := SqrtGCF(AdaptCFToGCF(src))
	if err != nil {
		t.Fatalf("SqrtGCF failed: %v", err)
	}
	if src.reads != 0 {
		t.Fatalf("reads got %d want 0 before first Next", src.reads)
	}

	d, ok := cf.Next()
	if !ok {
		t.Fatalf("expected first digit")
	}
	if d != 2 {
		t.Fatalf("digit got %d want 2", d)
	}
	if src.reads != 0 {
		t.Fatalf("lazy metadata path should not read source terms, got reads=%d", src.reads)
	}
}

func TestSqrtGCF_LazyMetadataPath_SqrtFourExhaustsAfterSingleDigitWithoutReadingInputTerms(t *testing.T) {
	src := newCountingRadicalCF(4, 2, 1, 1)

	cf, err := SqrtGCF(AdaptCFToGCF(src))
	if err != nil {
		t.Fatalf("SqrtGCF failed: %v", err)
	}

	d, ok := cf.Next()
	if !ok || d != 2 {
		t.Fatalf("first digit got (%d,%v) want (2,true)", d, ok)
	}

	_, ok = cf.Next()
	if ok {
		t.Fatalf("expected exhaustion after single digit")
	}
	if src.reads != 0 {
		t.Fatalf("lazy metadata path should not read source terms, got reads=%d", src.reads)
	}
}

func TestSqrtGCF_ExactFiniteTwoStillReturnsBootstrapApproximation(t *testing.T) {
	cf, err := SqrtGCF(NewSliceGCF([2]int64{2, 1}))
	if err != nil {
		t.Fatalf("SqrtGCF failed: %v", err)
	}

	got := collectTerms(cf, 8)
	want := collectTerms(NewRationalCF(mustRat(577, 408)), 8)

	if len(got) != len(want) {
		t.Fatalf("len mismatch: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, got, want)
		}
	}
}

func TestSqrtGCF_EmptySourceRecordsErrorOnFirstNext(t *testing.T) {
	cf, err := SqrtGCF(NewSliceGCF())
	if err != nil {
		t.Fatalf("SqrtGCF failed: %v", err)
	}

	_, ok := cf.Next()
	if ok {
		t.Fatalf("expected no emitted term")
	}

	s, ok := cf.(*sqrtBootstrapCFStream)
	if !ok {
		t.Fatalf("expected *sqrtBootstrapCFStream")
	}
	if s.Err() == nil {
		t.Fatalf("expected recorded error")
	}
	if !strings.Contains(s.Err().Error(), "empty source") {
		t.Fatalf("unexpected error: %v", s.Err())
	}
}

func TestSqrtGCF_NegativeFiniteInputRecordsErrorOnFirstNext(t *testing.T) {
	cf, err := SqrtGCF(NewSliceGCF([2]int64{-1, 1}))
	if err != nil {
		t.Fatalf("SqrtGCF failed: %v", err)
	}

	_, ok := cf.Next()
	if ok {
		t.Fatalf("expected no emitted term")
	}

	s, ok := cf.(*sqrtBootstrapCFStream)
	if !ok {
		t.Fatalf("expected *sqrtBootstrapCFStream")
	}
	if s.Err() == nil {
		t.Fatalf("expected recorded error")
	}
	if !strings.Contains(s.Err().Error(), "negative input") {
		t.Fatalf("unexpected error: %v", s.Err())
	}
}

func TestSqrtGCF_InputUnresolvedWithinBootstrapBudget_RecordsBudgetError(t *testing.T) {
	src := newCountingSliceGCF(repeatedUnitTerms(sqrtGCFExactBootstrapTermBudget + 20)...)

	cf, err := SqrtGCF(src)
	if err != nil {
		t.Fatalf("SqrtGCF failed: %v", err)
	}

	_, ok := cf.Next()
	if ok {
		t.Fatalf("expected no emitted term")
	}

	s, ok := cf.(*sqrtBootstrapCFStream)
	if !ok {
		t.Fatalf("expected *sqrtBootstrapCFStream")
	}
	if s.Err() == nil {
		t.Fatalf("expected recorded error")
	}
	if !strings.Contains(s.Err().Error(), "bootstrap term budget") {
		t.Fatalf("unexpected error: %v", s.Err())
	}
	if src.reads != sqrtGCFExactBootstrapTermBudget {
		t.Fatalf("reads got %d want %d", src.reads, sqrtGCFExactBootstrapTermBudget)
	}
}

func TestSqrtGCF_CurrentBootstrapCannotDistinguishLongFiniteFromInfinite(t *testing.T) {
	finiteLong := newCountingSliceGCF(repeatedUnitTerms(sqrtGCFExactBootstrapTermBudget + 20)...)

	cf, err := SqrtGCF(finiteLong)
	if err != nil {
		t.Fatalf("SqrtGCF failed: %v", err)
	}

	_, ok := cf.Next()
	if ok {
		t.Fatalf("expected no emitted term")
	}

	s, ok := cf.(*sqrtBootstrapCFStream)
	if !ok {
		t.Fatalf("expected *sqrtBootstrapCFStream")
	}
	if s.Err() == nil {
		t.Fatalf("expected recorded error")
	}
	if strings.Contains(s.Err().Error(), "non-terminating") {
		t.Fatalf("bootstrap should not claim non-terminating input: %v", s.Err())
	}
}

// sqrt_stream_min_test.go v3
