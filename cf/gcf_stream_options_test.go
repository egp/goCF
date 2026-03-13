// gcf_stream_options_test.go v1
package cf

import "testing"

type cappedRefiningTailEvidenceGCFSource struct {
	calls          int
	refinementStep int
}

func (s *cappedRefiningTailEvidenceGCFSource) NextPQ() (int64, int64, bool) {
	s.calls++
	return 2, 1, true
}

func (s *cappedRefiningTailEvidenceGCFSource) TailEvidence() GCFTailEvidence {
	r := NewRange(mustRat(2, 1), mustRat(3, 1), true, true)
	return GCFTailEvidence{
		Range:         &r,
		RangeReusable: false,
	}
}

func (s *cappedRefiningTailEvidenceGCFSource) RefinedTailEvidence() (GCFTailEvidence, bool) {
	s.refinementStep++

	switch s.refinementStep {
	case 1:
		r := NewRange(mustRat(2, 1), mustRat(3, 1), true, true)
		return GCFTailEvidence{Range: &r, RangeReusable: false}, true
	case 2:
		r := NewRange(mustRat(2, 1), mustRat(5, 2), true, true)
		return GCFTailEvidence{Range: &r, RangeReusable: false}, true
	default:
		return GCFTailEvidence{}, false
	}
}

func TestGCFStream_MaxRefinementSteps_LimitsRefinementBeforeIngest(t *testing.T) {
	src := &cappedRefiningTailEvidenceGCFSource{}
	s := NewGCFStream(src, GCFStreamOptions{MaxRefinementSteps: 1})

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if d != 2 {
		t.Fatalf("got first digit %d want 2", d)
	}
	callsAfterFirst := src.calls

	d, ok = s.Next()
	if !ok {
		t.Fatalf("expected second digit, err=%v", s.Err())
	}
	if d != 2 {
		t.Fatalf("got second digit %d want 2", d)
	}

	if src.calls <= callsAfterFirst {
		t.Fatalf("expected extra ingestion when refinement cap is too small, first calls=%d second calls=%d", callsAfterFirst, src.calls)
	}

	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

func TestGCFStream_NegativeMaxRefinementSteps_IsError(t *testing.T) {
	s := NewGCFStream(&cappedRefiningTailEvidenceGCFSource{}, GCFStreamOptions{MaxRefinementSteps: -1})

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected no digit")
	}
	if err := s.Err(); err == nil {
		t.Fatalf("expected non-nil error")
	}
}
