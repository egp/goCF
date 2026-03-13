// gcf_stream_tail_range_post_emit_test.go v1
package cf

import "testing"

type postEmitTailEvidenceGCFSource struct {
	calls int
}

func (s *postEmitTailEvidenceGCFSource) NextPQ() (int64, int64, bool) {
	s.calls++
	return 2, 1, true
}

func (s *postEmitTailEvidenceGCFSource) TailEvidence() GCFTailEvidence {
	r := NewRange(mustRat(2, 1), mustRat(5, 2), true, true)
	return GCFTailEvidence{
		Range:         &r,
		RangeReusable: false,
	}
}

func (s *postEmitTailEvidenceGCFSource) PostEmitTailEvidence(emittedDigit int64) (GCFTailEvidence, bool) {
	if emittedDigit != 2 {
		return GCFTailEvidence{}, false
	}

	r := NewRange(mustRat(2, 1), mustRat(5, 2), true, true)
	return GCFTailEvidence{
		Range:         &r,
		RangeReusable: false,
	}, true
}

func TestGCFStream_PostEmitTailEvidenceCanContinueWithoutReusableStaticRange(t *testing.T) {
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

	src := &postEmitTailEvidenceGCFSource{}
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
			"expected second digit via post-emit evidence without additional ingestion, first calls=%d second calls=%d",
			callsAfterFirst,
			src.calls,
		)
	}

	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

type chainedPostEmitTailEvidenceGCFSource struct {
	calls int
}

func (s *chainedPostEmitTailEvidenceGCFSource) NextPQ() (int64, int64, bool) {
	s.calls++
	return 2, 1, true
}

func (s *chainedPostEmitTailEvidenceGCFSource) TailEvidence() GCFTailEvidence {
	r := NewRange(mustRat(2, 1), mustRat(5, 2), true, true)
	return GCFTailEvidence{
		Range:         &r,
		RangeReusable: false,
	}
}

func (s *chainedPostEmitTailEvidenceGCFSource) PostEmitTailEvidence(emittedDigit int64) (GCFTailEvidence, bool) {
	if emittedDigit != 2 {
		return GCFTailEvidence{}, false
	}

	r := NewRange(mustRat(2, 1), mustRat(5, 2), true, true)
	return GCFTailEvidence{
		Range:         &r,
		RangeReusable: false,
	}, true
}

func TestGCFStream_ChainedPostEmitTailEvidenceImprovesCadenceOverStaticEvidence(t *testing.T) {
	factory := func() GCFSource { return &plainTwoOneGCFSource{} }

	want8 := exactDigitsFromFinitePrefix(t, factory, 8, 3)
	want10 := exactDigitsFromFinitePrefix(t, factory, 10, 3)

	if len(want8) != len(want10) {
		t.Fatalf("stabilization len mismatch: want8=%v want10=%v", want8, want10)
	}
	for i := range want8 {
		if want8[i] != want10[i] {
			t.Fatalf("oracle fixture not stabilized at digit %d: p8=%v p10=%v", i, want8, want10)
		}
	}

	chainedSrc := &chainedPostEmitTailEvidenceGCFSource{}
	staticSrc := &reusableOracleTailRangeGCFSource{}

	chained := NewGCFStream(chainedSrc, GCFStreamOptions{})
	static := NewGCFStream(staticSrc, GCFStreamOptions{})

	gotChained := collectTerms(chained, 3)
	gotStatic := collectTerms(static, 3)

	if len(gotChained) != 3 {
		t.Fatalf("expected 3 digits from chained source, got=%v err=%v", gotChained, chained.Err())
	}
	if len(gotStatic) != 3 {
		t.Fatalf("expected 3 digits from static source, got=%v err=%v", gotStatic, static.Err())
	}

	for i := range want10 {
		if gotChained[i] != want10[i] {
			t.Fatalf("chained digit %d: got=%v want=%v", i, gotChained, want10)
		}
		if gotStatic[i] != want10[i] {
			t.Fatalf("static digit %d: got=%v want=%v", i, gotStatic, want10)
		}
	}

	if chainedSrc.calls > staticSrc.calls {
		t.Fatalf("expected chained post-emit evidence to use no more ingestion than static evidence, chained=%d static=%d", chainedSrc.calls, staticSrc.calls)
	}

	if err := chained.Err(); err != nil {
		t.Fatalf("chained source: unexpected err=%v", err)
	}
	if err := static.Err(); err != nil {
		t.Fatalf("static source: unexpected err=%v", err)
	}
}

// gcf_stream_tail_range_post_emit_test.go v1
