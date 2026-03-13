// gcf_stream_tail_evidence_test.go v3
package cf

import "testing"

type countingStableTailRangeGCFSource struct {
	i     int
	calls int
}

func (s *countingStableTailRangeGCFSource) NextPQ() (int64, int64, bool) {
	s.calls++
	s.i++

	// Emit (1,1) forever.
	//
	// The represented GCF is:
	//   x = 1 + 1/(1 + 1/(1 + 1/(...)))
	//
	// With the stable explicit unfinished-tail range [1,2], the stream can prove:
	//   after 2 ingested terms: x in [3/2, 5/3] => first CF digit 1
	//   after 3 ingested terms, the next ordinary CF digit is still 1
	//
	// This matches the regular continued fraction of the golden ratio:
	//   [1; 1, 1, 1, ...]

	return 1, 1, true
}

func (s *countingStableTailRangeGCFSource) TailLowerBound() Rational {
	return mustRat(1, 1)
}

func (s *countingStableTailRangeGCFSource) TailRange() Range {
	return NewRange(mustRat(1, 1), mustRat(2, 1), true, true)
}

func TestGCFStream_UsesGenericStableTailRangeForEarlierEmission(t *testing.T) {
	src := &countingStableTailRangeGCFSource{}
	s := NewGCFStream(src, GCFStreamOptions{})

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if d != 1 {
		t.Fatalf("got first digit %d want 1", d)
	}

	if src.calls != 2 {
		t.Fatalf("expected first digit after exactly 2 ingested GCF terms, got %d", src.calls)
	}

	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

func TestGCFStream_StableTailRangeBeatsLowerBoundRay(t *testing.T) {
	src := &countingStableTailRangeGCFSource{}
	s := NewGCFStream(src, GCFStreamOptions{})

	_, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}

	if src.calls > 2 {
		t.Fatalf("expected TailRange-driven early emission, but stream needed %d ingested terms", src.calls)
	}
}

func TestGCFStream_StableTailRangeFirstTwoDigitsAndCadence(t *testing.T) {
	src := &countingStableTailRangeGCFSource{}
	s := NewGCFStream(src, GCFStreamOptions{})

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if d != 1 {
		t.Fatalf("got first digit %d want 1", d)
	}
	if src.calls != 2 {
		t.Fatalf("expected first digit after exactly 2 ingested terms, got %d", src.calls)
	}

	d, ok = s.Next()
	if !ok {
		t.Fatalf("expected second digit, err=%v", s.Err())
	}
	if d != 1 {
		t.Fatalf("got second digit %d want 1", d)
	}
	if src.calls != 3 {
		t.Fatalf("expected second digit after exactly 3 ingested terms, got %d", src.calls)
	}

	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

func TestGCFStream_StableTailRangeSecondDigitAfterOneMoreIngest(t *testing.T) {
	src := &countingStableTailRangeGCFSource{}
	s := NewGCFStream(src, GCFStreamOptions{})

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if d != 1 {
		t.Fatalf("got first digit %d want 1", d)
	}
	if src.calls != 2 {
		t.Fatalf("expected first digit after 2 ingested terms, got %d", src.calls)
	}

	_, ok = s.Next()
	if !ok {
		t.Fatalf("expected second digit, err=%v", s.Err())
	}

	if src.calls != 3 {
		t.Fatalf("expected second digit after exactly 3 ingested terms, got %d", src.calls)
	}

	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

type prefixSensitiveTailRangeGCFSource struct {
	i     int
	calls int
}

func (s *prefixSensitiveTailRangeGCFSource) NextPQ() (int64, int64, bool) {
	s.calls++
	s.i++
	return 1, 1, true
}

func (s *prefixSensitiveTailRangeGCFSource) TailLowerBound() Rational {
	return mustRat(1, 1)
}

func (s *prefixSensitiveTailRangeGCFSource) TailRange() Range {
	switch s.i {
	case 0, 1:
		return NewRange(mustRat(1, 1), mustRat(3, 1), true, true)
	default:
		return NewRange(mustRat(1, 1), mustRat(2, 1), true, true)
	}
}

func TestGCFStream_UsesCurrentPrefixSensitiveTailRange(t *testing.T) {
	src := &prefixSensitiveTailRangeGCFSource{}
	s := NewGCFStream(src, GCFStreamOptions{})

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if d != 1 {
		t.Fatalf("got first digit %d want 1", d)
	}

	if src.calls != 2 {
		t.Fatalf("expected first digit after exactly 2 ingested terms, got %d", src.calls)
	}

	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

type weakTailRangeGCFSource struct {
	calls int
}

func (s *weakTailRangeGCFSource) NextPQ() (int64, int64, bool) {
	s.calls++
	return 1, 1, true
}

func (s *weakTailRangeGCFSource) TailLowerBound() Rational {
	return mustRat(1, 1)
}

func (s *weakTailRangeGCFSource) TailRange() Range {
	return NewRange(mustRat(1, 1), mustRat(3, 1), true, true)
}

type strongTailRangeGCFSource struct {
	calls int
}

func (s *strongTailRangeGCFSource) NextPQ() (int64, int64, bool) {
	s.calls++
	return 1, 1, true
}

func (s *strongTailRangeGCFSource) TailLowerBound() Rational {
	return mustRat(1, 1)
}

func (s *strongTailRangeGCFSource) TailRange() Range {
	return NewRange(mustRat(1, 1), mustRat(2, 1), true, true)
}

func TestGCFStream_StrongerTailRangeEmitsNoLaterThanWeakerTailRange(t *testing.T) {
	strongSrc := &strongTailRangeGCFSource{}
	weakSrc := &weakTailRangeGCFSource{}

	strong := NewGCFStream(strongSrc, GCFStreamOptions{})
	weak := NewGCFStream(weakSrc, GCFStreamOptions{})

	dStrong, ok := strong.Next()
	if !ok {
		t.Fatalf("strong source: expected first digit, err=%v", strong.Err())
	}
	if dStrong != 1 {
		t.Fatalf("strong source: got first digit %d want 1", dStrong)
	}

	dWeak, ok := weak.Next()
	if !ok {
		t.Fatalf("weak source: expected first digit, err=%v", weak.Err())
	}
	if dWeak != 1 {
		t.Fatalf("weak source: got first digit %d want 1", dWeak)
	}

	if strongSrc.calls > weakSrc.calls {
		t.Fatalf("expected stronger TailRange to emit no later than weaker one, strong=%d weak=%d", strongSrc.calls, weakSrc.calls)
	}

	if err := strong.Err(); err != nil {
		t.Fatalf("strong source: unexpected stream err: %v", err)
	}
	if err := weak.Err(); err != nil {
		t.Fatalf("weak source: unexpected stream err: %v", err)
	}
}

type competingTailEvidenceGCFSource struct {
	calls int
}

func (s *competingTailEvidenceGCFSource) NextPQ() (int64, int64, bool) {
	s.calls++
	return 1, 1, true
}

func (s *competingTailEvidenceGCFSource) TailLowerBound() Rational {
	return mustRat(1, 1)
}

func (s *competingTailEvidenceGCFSource) TailRange() Range {
	return NewRange(mustRat(1, 1), mustRat(2, 1), true, true)
}

func TestGCFStream_PrefersExplicitTailRangeOverLowerBoundRay(t *testing.T) {
	src := &competingTailEvidenceGCFSource{}
	s := NewGCFStream(src, GCFStreamOptions{})

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if d != 1 {
		t.Fatalf("got first digit %d want 1", d)
	}

	if src.calls != 2 {
		t.Fatalf("expected TailRange-driven emission after exactly 2 ingested terms, got %d", src.calls)
	}

	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

type delayedLowerBoundRayGCFSource struct {
	calls int
}

func (s *delayedLowerBoundRayGCFSource) NextPQ() (int64, int64, bool) {
	s.calls++
	return 1, 1, true
}

func (s *delayedLowerBoundRayGCFSource) TailLowerBound() Rational {
	return mustRat(1, 1)
}

func (s *delayedLowerBoundRayGCFSource) LowerBoundRayMinPrefix() int {
	return 2
}

func TestGCFStream_UsesGenericLowerBoundRayMinPrefixPolicy(t *testing.T) {
	src := &delayedLowerBoundRayGCFSource{}
	s := NewGCFStream(src, GCFStreamOptions{})

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if d != 1 {
		t.Fatalf("got first digit %d want 1", d)
	}

	if src.calls < 2 {
		t.Fatalf("expected at least 2 ingested terms before first digit, got %d", src.calls)
	}

	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

func TestGCFStream_LowerBoundOnlySourceCanStillEmit(t *testing.T) {
	src := &delayedLowerBoundRayGCFSource{}
	s := NewGCFStream(src, GCFStreamOptions{})

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if d != 1 {
		t.Fatalf("got first digit %d want 1", d)
	}

	if src.calls < 2 {
		t.Fatalf("expected at least 2 ingested terms before first digit, got %d", src.calls)
	}

	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

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

func (s *nonReusableOracleTailRangeGCFSource) TailRange() Range {
	// After one ingested (2,1) term, x = 2 + 1/y with y in [2, 5/2].
	// So x in [12/5, 5/2], certifying first digit 2.
	//
	// After emitting 2, the remainder is exactly y, and [2, 5/2] still has
	// unique floor 2. So a reusable tail-range contract can certify the second
	// digit without another ingest.
	return NewRange(mustRat(2, 1), mustRat(5, 2), true, true)
}

func (s *nonReusableOracleTailRangeGCFSource) TailRangeReusable() bool {
	return false
}

type reusableOracleTailRangeGCFSource struct {
	calls int
}

func (s *reusableOracleTailRangeGCFSource) NextPQ() (int64, int64, bool) {
	s.calls++
	return 2, 1, true
}

func (s *reusableOracleTailRangeGCFSource) TailRange() Range {
	return NewRange(mustRat(2, 1), mustRat(5, 2), true, true)
}

func (s *reusableOracleTailRangeGCFSource) TailRangeReusable() bool {
	return true
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

// gcf_stream_tail_evidence_test.go v3
