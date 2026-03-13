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

// gcf_stream_tail_evidence_test.go v3
