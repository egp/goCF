// gcf_stream_tail_range_unified_test.go v1
package cf

import "testing"

type plainTwoOneGCFSource struct {
	calls int
}

func (s *plainTwoOneGCFSource) NextPQ() (int64, int64, bool) {
	s.calls++
	return 2, 1, true
}

type nonReusableOracleTailRangeGCFSource struct {
	calls int
}

func (s *nonReusableOracleTailRangeGCFSource) NextPQ() (int64, int64, bool) {
	s.calls++
	return 2, 1, true
}

func (s *nonReusableOracleTailRangeGCFSource) TailEvidence() GCFTailEvidence {
	r := NewRange(mustRat(2, 1), mustRat(5, 2), true, true)
	return GCFTailEvidence{
		Range:         &r,
		RangeReusable: false,
	}
}

type reusableOracleTailRangeGCFSource struct {
	calls int
}

func (s *reusableOracleTailRangeGCFSource) NextPQ() (int64, int64, bool) {
	s.calls++
	return 2, 1, true
}

func (s *reusableOracleTailRangeGCFSource) TailEvidence() GCFTailEvidence {
	r := NewRange(mustRat(2, 1), mustRat(5, 2), true, true)
	return GCFTailEvidence{
		Range:         &r,
		RangeReusable: true,
	}
}

func TestGCFStream_ReusableTailRangePolicyBeatsNonReusablePolicy_OnOracleBackedFixture(t *testing.T) {
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

	reusableSrc := &reusableOracleTailRangeGCFSource{}
	nonReusableSrc := &nonReusableOracleTailRangeGCFSource{}

	reusable := NewGCFStream(reusableSrc, GCFStreamOptions{})
	nonReusable := NewGCFStream(nonReusableSrc, GCFStreamOptions{})

	d, ok := reusable.Next()
	if !ok {
		t.Fatalf("reusable source: expected first digit, err=%v", reusable.Err())
	}
	if d != want8[0] {
		t.Fatalf("reusable source: got first digit %d want %d", d, want8[0])
	}
	reusableCallsAfterFirst := reusableSrc.calls

	d, ok = nonReusable.Next()
	if !ok {
		t.Fatalf("non-reusable source: expected first digit, err=%v", nonReusable.Err())
	}
	if d != want8[0] {
		t.Fatalf("non-reusable source: got first digit %d want %d", d, want8[0])
	}
	nonReusableCallsAfterFirst := nonReusableSrc.calls

	d, ok = reusable.Next()
	if !ok {
		t.Fatalf("reusable source: expected second digit, err=%v", reusable.Err())
	}
	if d != want8[1] {
		t.Fatalf("reusable source: got second digit %d want %d", d, want8[1])
	}

	d, ok = nonReusable.Next()
	if !ok {
		t.Fatalf("non-reusable source: expected second digit, err=%v", nonReusable.Err())
	}
	if d != want8[1] {
		t.Fatalf("non-reusable source: got second digit %d want %d", d, want8[1])
	}

	if reusableSrc.calls > reusableCallsAfterFirst {
		t.Fatalf(
			"reusable source: expected second digit without additional ingestion, first calls=%d second calls=%d",
			reusableCallsAfterFirst,
			reusableSrc.calls,
		)
	}

	if nonReusableSrc.calls <= nonReusableCallsAfterFirst {
		t.Fatalf(
			"non-reusable source: expected additional ingestion before second digit, first calls=%d second calls=%d",
			nonReusableCallsAfterFirst,
			nonReusableSrc.calls,
		)
	}

	if err := reusable.Err(); err != nil {
		t.Fatalf("reusable source: unexpected stream err: %v", err)
	}
	if err := nonReusable.Err(); err != nil {
		t.Fatalf("non-reusable source: unexpected stream err: %v", err)
	}
}

type unifiedTailEvidenceGCFSource struct {
	calls int
}

func (s *unifiedTailEvidenceGCFSource) NextPQ() (int64, int64, bool) {
	s.calls++
	return 2, 1, true
}

func (s *unifiedTailEvidenceGCFSource) TailEvidence() GCFTailEvidence {
	r := NewRange(mustRat(2, 1), mustRat(5, 2), true, true)
	return GCFTailEvidence{
		Range:         &r,
		RangeReusable: true,
	}
}

func TestGCFStream_UnifiedTailEvidenceMatchesReusableSplitContracts(t *testing.T) {
	unifiedSrc := &unifiedTailEvidenceGCFSource{}
	splitSrc := &reusableOracleTailRangeGCFSource{}

	unified := NewGCFStream(unifiedSrc, GCFStreamOptions{})
	split := NewGCFStream(splitSrc, GCFStreamOptions{})

	gotUnified := collectTerms(unified, 2)
	gotSplit := collectTerms(split, 2)

	if len(gotUnified) != len(gotSplit) {
		t.Fatalf("len mismatch: unified=%v split=%v", gotUnified, gotSplit)
	}
	for i := range gotSplit {
		if gotUnified[i] != gotSplit[i] {
			t.Fatalf("digit %d: unified=%v split=%v", i, gotUnified, gotSplit)
		}
	}

	if unifiedSrc.calls != splitSrc.calls {
		t.Fatalf("expected matching ingest cadence, unified=%d split=%d", unifiedSrc.calls, splitSrc.calls)
	}

	if err := unified.Err(); err != nil {
		t.Fatalf("unified source: unexpected stream err: %v", err)
	}
	if err := split.Err(); err != nil {
		t.Fatalf("split source: unexpected stream err: %v", err)
	}
}

type invalidReusableTailEvidenceGCFSource struct{}

func (s *invalidReusableTailEvidenceGCFSource) NextPQ() (int64, int64, bool) {
	return 1, 1, true
}

func (s *invalidReusableTailEvidenceGCFSource) TailEvidence() GCFTailEvidence {
	return GCFTailEvidence{
		Range:         nil,
		RangeReusable: true,
	}
}

func TestGCFStream_InvalidUnifiedTailEvidence_IsError(t *testing.T) {
	s := NewGCFStream(&invalidReusableTailEvidenceGCFSource{}, GCFStreamOptions{})

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected no digit")
	}
	if err := s.Err(); err == nil {
		t.Fatalf("expected non-nil error")
	}
}

func TestGCFStream_ReusableTailRangePoleBoundaryFallsThroughToFurtherIngestion(t *testing.T) {
	src := &reusableOracleTailRangeGCFSource{}
	s := NewGCFStream(src, GCFStreamOptions{})

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
	callsAfterSecond := src.calls

	d, ok = s.Next()
	if !ok {
		t.Fatalf("expected third digit after pole-boundary fallback, err=%v", s.Err())
	}

	if src.calls <= callsAfterSecond {
		t.Fatalf(
			"expected additional ingestion after reusable tail-range pole boundary, second calls=%d third calls=%d",
			callsAfterSecond,
			src.calls,
		)
	}

	if err := s.Err(); err != nil {
		t.Fatalf("expected no hard error after pole-boundary fallback, got %v", err)
	}

	t.Logf("digits so far: 2, 2, %d; calls after first=%d second=%d third=%d",
		d, callsAfterFirst, callsAfterSecond, src.calls)
}

func TestGCFStream_ReusableTailRangePoleBoundaryIsNotImmediateHardError(t *testing.T) {
	src := &reusableOracleTailRangeGCFSource{}
	s := NewGCFStream(src, GCFStreamOptions{})

	d, ok := s.Next()
	if !ok || d != 2 {
		t.Fatalf("expected first digit 2, got ok=%v d=%d err=%v", ok, d, s.Err())
	}

	d, ok = s.Next()
	if !ok || d != 2 {
		t.Fatalf("expected second digit 2, got ok=%v d=%d err=%v", ok, d, s.Err())
	}

	_, _ = s.Next()
	if err := s.Err(); err != nil {
		t.Fatalf("expected no hard error after reusable tail-range pole boundary, got %v", err)
	}
}

func TestGCFStream_ReusableTailRangePoleBoundary_ThirdDigitMatchesExactFiniteOracle(t *testing.T) {
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

	src := &reusableOracleTailRangeGCFSource{}
	s := NewGCFStream(src, GCFStreamOptions{})

	got := collectTerms(s, 3)
	if len(got) != 3 {
		t.Fatalf("expected 3 digits, got=%v err=%v", got, s.Err())
	}

	for i := range want10 {
		if got[i] != want10[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, got, want10)
		}
	}

	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

type exactPointPoleTailEvidenceGCFSource struct {
	calls int
}

func (s *exactPointPoleTailEvidenceGCFSource) NextPQ() (int64, int64, bool) {
	s.calls++
	return 2, 1, true
}

func (s *exactPointPoleTailEvidenceGCFSource) TailEvidence() GCFTailEvidence {
	// With one ingested (2,1) term and exact tail y=-1:
	//   x = 2 + 1/y = 1
	// so the first ordinary CF digit 1 is emitted successfully.
	//
	// After emitting 1, the remainder transform is singular at exact point x=1,
	// so this should remain a hard exact-point error.
	r := NewRange(mustRat(-1, 1), mustRat(-1, 1), true, true)
	return GCFTailEvidence{
		Range:         &r,
		RangeReusable: true,
	}
}

func TestGCFStream_ExactPointExplicitTailRangePoleRemainsHardError(t *testing.T) {
	src := &exactPointPoleTailEvidenceGCFSource{}
	s := NewGCFStream(src, GCFStreamOptions{})

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if d != 1 {
		t.Fatalf("got first digit %d want 1", d)
	}

	_, ok = s.Next()
	if ok {
		t.Fatalf("expected second step to fail on exact-point pole")
	}
	if err := s.Err(); err == nil {
		t.Fatalf("expected non-nil hard error on exact-point explicit tail-range pole")
	}
}

func TestGCFStream_NonPointExplicitTailRangePoleFallsThroughWithoutHardError(t *testing.T) {
	src := &reusableOracleTailRangeGCFSource{}
	s := NewGCFStream(src, GCFStreamOptions{})

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if d != 2 {
		t.Fatalf("got first digit %d want 2", d)
	}

	d, ok = s.Next()
	if !ok {
		t.Fatalf("expected second digit, err=%v", s.Err())
	}
	if d != 2 {
		t.Fatalf("got second digit %d want 2", d)
	}

	_, ok = s.Next()
	if !ok {
		t.Fatalf("expected third digit after non-point pole fallback, err=%v", s.Err())
	}
	if err := s.Err(); err != nil {
		t.Fatalf("expected no hard error after non-point explicit tail-range pole, got %v", err)
	}
}

// gcf_stream_tail_range_unified_test.go v1
