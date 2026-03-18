// sqrt_stream_min_test.go v1
package cf

import "testing"

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

// sqrt_stream_min_test.go v1
