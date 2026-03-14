// sqrt_stream2_test.go v1
package cf

import (
	"strings"
	"testing"
)

func TestNewSqrtGCFStream2_ExactTail_PerfectSquareFastPath(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(0, 1),
	}

	s := NewSqrtGCFStreamWithTail2(
		NewSliceGCF(),
		mustRat(9, 16),
		8,
		p,
	)

	got := collectTerms(s, 8)
	if err := s.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}

	want := []int64{0, 1, 3} // sqrt(9/16) = 3/4
	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d full=%v", i, got[i], want[i], got)
		}
	}
}

func TestNewSqrtGCFStream2_ExactTail_FinitePrefixMatchesCanonicalPath(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	srcA := NewSliceGCF([2]int64{3, 2})
	srcB := NewSliceGCF([2]int64{3, 2})
	tail := mustRat(11, 1)

	wantCF, err := SqrtApproxCFFromGCFTailSource2(srcA, NewExactTailSource(tail), 8, p)
	if err != nil {
		t.Fatalf("SqrtApproxCFFromGCFTailSource2 failed: %v", err)
	}
	want := collectTerms(wantCF, 16)

	s := NewSqrtGCFStream2(srcB, NewExactTailSource(tail), 8, p)
	got := collectTerms(s, 16)
	if err := s.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v wantFull=%v", len(got), len(want), got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d gotFull=%v wantFull=%v", i, got[i], want[i], got, want)
		}
	}
}

func TestNewSqrtGCFStream2_MissingTailEvidenceIsError(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	s := NewSqrtGCFStream2(
		NewSliceGCF([2]int64{3, 2}),
		NoTailSource{},
		8,
		p,
	)

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected failure, not a digit")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(s.Err().Error(), "tail evidence not implemented") {
		t.Fatalf("unexpected error: %v", s.Err())
	}
}

func TestNewSqrtGCFStream2_HonorsBoundFailure(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	s := NewSqrtGCFStream2(
		NewUnitPArithmeticQGCFSource(1, 1),
		NewExactTailSource(mustRat(1, 1)),
		3,
		p,
	)

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected failure, not a digit")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}
}

func TestSqrtGCFStream2_ExhaustedStreamStaysExhausted(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	s := NewSqrtGCFStreamWithTail2(
		NewSliceGCF(),
		mustRat(9, 16),
		8,
		p,
	)

	for {
		_, ok := s.Next()
		if !ok {
			break
		}
	}

	if _, ok := s.Next(); ok {
		t.Fatalf("expected exhausted stream to stay exhausted")
	}
	if _, ok := s.Next(); ok {
		t.Fatalf("expected exhausted stream to stay exhausted on repeated calls")
	}
}

// sqrt_stream2_test.go v1
