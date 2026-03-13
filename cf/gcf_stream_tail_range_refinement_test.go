// gcf_stream_tail_range_refinement_test.go v2
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

type multiRefiningTailEvidenceGCFSource struct {
	calls          int
	refinementStep int
}

func (s *multiRefiningTailEvidenceGCFSource) NextPQ() (int64, int64, bool) {
	s.calls++
	return 2, 1, true
}

func (s *multiRefiningTailEvidenceGCFSource) TailEvidence() GCFTailEvidence {
	// Initial evidence certifies the first digit 2 after one ingested term,
	// but is too coarse to certify the next digit.
	r := NewRange(mustRat(2, 1), mustRat(3, 1), true, true)
	return GCFTailEvidence{
		Range:         &r,
		RangeReusable: false,
	}
}

func (s *multiRefiningTailEvidenceGCFSource) RefinedTailEvidence() (GCFTailEvidence, bool) {
	s.refinementStep++

	switch s.refinementStep {
	case 1:
		// Still too coarse for the second digit.
		r := NewRange(mustRat(2, 1), mustRat(3, 1), true, true)
		return GCFTailEvidence{
			Range:         &r,
			RangeReusable: false,
		}, true
	case 2:
		// Now tight enough to certify second digit 2 without additional ingestion.
		r := NewRange(mustRat(2, 1), mustRat(5, 2), true, true)
		return GCFTailEvidence{
			Range:         &r,
			RangeReusable: false,
		}, true
	default:
		return GCFTailEvidence{}, false
	}
}

func TestGCFStream_MultipleRefinementsCanAvoidExtraIngestion(t *testing.T) {
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

	src := &multiRefiningTailEvidenceGCFSource{}
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
			"expected second digit via multiple refinements without additional ingestion, first calls=%d second calls=%d",
			callsAfterFirst,
			src.calls,
		)
	}

	if src.refinementStep < 2 {
		t.Fatalf("expected at least 2 refinement attempts, got %d", src.refinementStep)
	}

	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

type candidateTailEvidenceGCFSource struct {
	calls int
}

func (s *candidateTailEvidenceGCFSource) NextPQ() (int64, int64, bool) {
	s.calls++
	return 2, 1, true
}

func (s *candidateTailEvidenceGCFSource) TailEvidence() GCFTailEvidence {
	// Base evidence is intentionally too coarse for the second digit.
	r := NewRange(mustRat(2, 1), mustRat(3, 1), true, true)
	return GCFTailEvidence{
		Range:         &r,
		RangeReusable: false,
	}
}

func (s *candidateTailEvidenceGCFSource) CandidateTailEvidence() []GCFTailEvidence {
	// First candidate is still too coarse.
	r1 := NewRange(mustRat(2, 1), mustRat(3, 1), true, true)

	// Second candidate is tight enough to certify second digit 2
	// without another ingest.
	r2 := NewRange(mustRat(2, 1), mustRat(5, 2), true, true)

	return []GCFTailEvidence{
		{Range: &r1, RangeReusable: false},
		{Range: &r2, RangeReusable: false},
	}
}

func TestGCFStream_CandidateTailEvidenceCanAvoidExtraIngestion(t *testing.T) {
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

	src := &candidateTailEvidenceGCFSource{}
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
			"expected second digit via candidate tail evidence without additional ingestion, first calls=%d second calls=%d",
			callsAfterFirst,
			src.calls,
		)
	}

	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

// gcf_stream_tail_range_refinement_test.go v2
