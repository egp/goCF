// gcf_stream_tail_range_basic_test.go v1
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

// gcf_stream_tail_range_basic_test.go v1
