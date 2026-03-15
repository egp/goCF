// sqrt_stream_gcf_exact_tail_api_test.go v1
package cf

import (
	"strings"
	"testing"
)

func TestSqrtGCFExactTailStream_ReturnsInspectableStream(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	s, err := SqrtGCFExactTailStreamWithTail(
		NewSliceGCF([2]int64{3, 2}),
		mustRat(11, 1),
		8,
		p,
	)
	if err != nil {
		t.Fatalf("SqrtGCFExactTailStreamWithTail failed: %v", err)
	}

	snap := s.Snapshot()
	if snap.Started {
		t.Fatalf("expected Started=false before first Next")
	}
	if snap.Approx != nil {
		t.Fatalf("expected Approx=nil before start")
	}
	if s.Err() != nil {
		t.Fatalf("expected nil Err before start, got %v", s.Err())
	}
}

func TestSqrtGCFExactTailStream_AfterStartCarriesApproximation(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	s, err := SqrtGCFExactTailStreamWithTail(
		NewSliceGCF([2]int64{3, 2}),
		mustRat(11, 1),
		8,
		p,
	)
	if err != nil {
		t.Fatalf("SqrtGCFExactTailStreamWithTail failed: %v", err)
	}

	_, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit; err=%v", s.Err())
	}

	snap := s.Snapshot()
	if !snap.Started {
		t.Fatalf("expected Started=true after start")
	}
	if snap.Approx == nil {
		t.Fatalf("expected non-nil Approx after start")
	}

	want, err := SqrtApproxFromGCFWithTail2(
		NewSliceGCF([2]int64{3, 2}),
		mustRat(11, 1),
		8,
		sqrtPolicy2FromOld(p),
	)
	if err != nil {
		t.Fatalf("SqrtApproxFromGCFWithTail2 failed: %v", err)
	}
	if snap.Approx.Cmp(want) != 0 {
		t.Fatalf("got Approx=%v want %v", *snap.Approx, want)
	}
}

func TestSqrtGCFExactTailStream_MatchesCanonicalExactTailPath(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	wantCF, err := SqrtApproxCFFromGCFTailSource2(
		NewSliceGCF([2]int64{3, 2}),
		NewExactTailSource(mustRat(11, 1)),
		8,
		sqrtPolicy2FromOld(p),
	)
	if err != nil {
		t.Fatalf("SqrtApproxCFFromGCFTailSource2 failed: %v", err)
	}
	want := collectTerms(wantCF, 16)

	s, err := SqrtGCFExactTailStreamWithTail(
		NewSliceGCF([2]int64{3, 2}),
		mustRat(11, 1),
		8,
		p,
	)
	if err != nil {
		t.Fatalf("SqrtGCFExactTailStreamWithTail failed: %v", err)
	}
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

func TestSqrtGCFExactTailStream_MissingTailEvidenceFailsOnUse(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	s, err := SqrtGCFExactTailStream(
		NewSliceGCF([2]int64{3, 2}),
		NoTailSource{},
		8,
		p,
	)
	if err != nil {
		t.Fatalf("unexpected constructor error: %v", err)
	}

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
