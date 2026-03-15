// reciprocal_stream_gcf_exact_tail_test.go v1
package cf

import (
	"strings"
	"testing"
)

func TestReciprocalGCFExactTailStream_ReturnsInspectableStream(t *testing.T) {
	s, err := ReciprocalGCFExactTailStreamWithTail(
		NewSliceGCF([2]int64{3, 2}),
		mustRat(11, 1),
		8,
	)
	if err != nil {
		t.Fatalf("ReciprocalGCFExactTailStreamWithTail failed: %v", err)
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

func TestReciprocalGCFExactTailStream_AfterStartCarriesApproximation(t *testing.T) {
	s, err := ReciprocalGCFExactTailStreamWithTail(
		NewSliceGCF([2]int64{3, 2}),
		mustRat(11, 1),
		8,
	)
	if err != nil {
		t.Fatalf("ReciprocalGCFExactTailStreamWithTail failed: %v", err)
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

	// x = 3 + 2/11 = 35/11, reciprocal = 11/35 = [0; 3, 5, 2]
	want := mustRat(11, 35)
	if snap.Approx.Cmp(want) != 0 {
		t.Fatalf("got Approx=%v want %v", *snap.Approx, want)
	}
}

func TestReciprocalGCFExactTailStream_EmitsExpectedCF(t *testing.T) {
	s, err := ReciprocalGCFExactTailStreamWithTail(
		NewSliceGCF([2]int64{3, 2}),
		mustRat(11, 1),
		8,
	)
	if err != nil {
		t.Fatalf("ReciprocalGCFExactTailStreamWithTail failed: %v", err)
	}

	got := collectTerms(s, 8)
	if err := s.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}

	want := []int64{0, 3, 5, 2} // 11/35
	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d gotFull=%v", i, got[i], want[i], got)
		}
	}
}

func TestReciprocalGCFExactTailStream_MissingTailEvidenceFailsOnUse(t *testing.T) {
	s, err := ReciprocalGCFExactTailStream(
		NewSliceGCF([2]int64{3, 2}),
		NoTailSource{},
		8,
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

func TestReciprocalGCFExactTailStream_ReciprocalOfZeroFails(t *testing.T) {
	s, err := ReciprocalGCFExactTailStreamWithTail(
		NewSliceGCF(),
		mustRat(0, 1),
		8,
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
	if !strings.Contains(s.Err().Error(), "reciprocal of zero") {
		t.Fatalf("unexpected error: %v", s.Err())
	}
}

func TestReciprocalGCFExactTailStream_RejectsZeroBound(t *testing.T) {
	_, err := ReciprocalGCFExactTailStreamWithTail(
		NewSliceGCF([2]int64{3, 2}),
		mustRat(11, 1),
		0,
	)
	if err == nil {
		t.Fatalf("expected constructor error")
	}
}

// reciprocal_stream_gcf_exact_tail_test.go v1
