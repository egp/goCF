// gcf_stream_tail_range_refinement_test.go v1
package cf

import "testing"

type refiningTailEvidenceGCFSource struct {
	calls      int
	refinedNow bool
}

func (s *refiningTailEvidenceGCFSource) NextPQ() (int64, int64, bool) {
	s.calls++
	return 2, 1, true
}

func (s *refiningTailEvidenceGCFSource) TailEvidence() GCFTailEvidence {
	r := NewRange(mustRat(2, 1), mustRat(3, 1), true, true)
	return GCFTailEvidence{
		Range:         &r,
		RangeReusable: false,
	}
}

func (s *refiningTailEvidenceGCFSource) RefinedTailEvidence() (GCFTailEvidence, bool) {
	if s.refinedNow {
		return GCFTailEvidence{}, false
	}
	s.refinedNow = true

	r := NewRange(mustRat(2, 1), mustRat(5, 2), true, true)
	return GCFTailEvidence{
		Range:         &r,
		RangeReusable: false,
	}, true
}

func TestGCFStream_RefinedTailEvidenceCanAvoidExtraIngestion(t *testing.T) {
	factory := func() GCFSource { return &plainTwoOneGCFSource{} }

	want6 := exactDigitsFromFinitePrefix(t, factory, 6, 2)
	want8 := exactDigitsFromFinitePrefix(t, factory, 8, 2)

	if len(want6) != len(want8) {
		t.Fatalf("stabilization len mismatch: want6=%v want8=%v", want6, want8)
	}
	for i := range want6 {
		if want6[i] != want8[i] {
			t.Fatalf("oracle fixture not stabilized at digit %d: p6=%v p8=%v", i, want6, want8)
		}
	}

	src := &refiningTailEvidenceGCFSource{}
	s := NewGCFStream(src, GCFStreamOptions{})

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if d != want8[0] {
		t.Fatalf("got first digit %d want %d", d, want8[0])
	}
	callsAfterFirst := src.calls

	d, ok = s.Next()
	if !ok {
		t.Fatalf("expected second digit, err=%v", s.Err())
	}
	if d != want8[1] {
		t.Fatalf("got second digit %d want %d", d, want8[1])
	}

	if src.calls != callsAfterFirst {
		t.Fatalf(
			"expected second digit via refined tail evidence without additional ingestion, first calls=%d second calls=%d",
			callsAfterFirst,
			src.calls,
		)
	}

	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

// gcf_stream_tail_range_refinement_test.go v1
