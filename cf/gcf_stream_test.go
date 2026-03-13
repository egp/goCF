// gcf_stream_test.go v1
package cf

import (
	"fmt"
	"testing"
)

func TestGCFStream_FiniteSliceGCF_MatchesExactRational(t *testing.T) {
	src := NewSliceGCF(
		[2]int64{1, 1},
		[2]int64{2, 1},
		[2]int64{3, 1},
	)

	gotStream := NewGCFStream(src, GCFStreamOptions{})

	got := collectTerms(gotStream, 32)

	wantRat, err := EvaluateFiniteGCF(NewSliceGCF(
		[2]int64{1, 1},
		[2]int64{2, 1},
		[2]int64{3, 1},
	))
	if err != nil {
		t.Fatalf("EvaluateFiniteGCF failed: %v", err)
	}
	want := collectTerms(NewRationalCF(wantRat), 32)

	if len(got) != len(want) {
		t.Fatalf("len mismatch: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, got, want)
		}
	}
	if err := gotStream.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

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

	// This is the real test: with a generic stable TailRange() contract, the
	// stream should be able to emit after exactly two ingested GCF terms.
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

	// Lower-bound-only logic would need more than two ingested terms here.
	// Stable TailRange() should make the stronger generic proof available.
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

func TestGCFStream_AdaptedRegularCFRoundTrip(t *testing.T) {
	orig := NewSliceCF(1, 2, 3, 4)
	src := AdaptCFToGCF(NewSliceCF(1, 2, 3, 4))

	gotStream := NewGCFStream(src, GCFStreamOptions{})
	got := collectTerms(gotStream, 32)
	want := collectTerms(orig, 32)

	if len(got) != len(want) {
		t.Fatalf("len mismatch: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, got, want)
		}
	}
	if err := gotStream.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

func TestGCFStream_EmptySourceIsError(t *testing.T) {
	s := NewGCFStream(NewSliceGCF(), GCFStreamOptions{})

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected no digit")
	}
	if err := s.Err(); err == nil {
		t.Fatalf("expected non-nil error")
	}
}

func TestGCFStream_FiniteSingleTermTerminatesCleanly(t *testing.T) {
	s := NewGCFStream(NewSliceGCF([2]int64{5, 1}), GCFStreamOptions{})

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if d != 5 {
		t.Fatalf("got %d want 5", d)
	}

	_, ok = s.Next()
	if ok {
		t.Fatalf("expected termination")
	}
	if err := s.Err(); err != nil {
		t.Fatalf("expected clean termination, got %v", err)
	}
}

func TestGCFStream_FiniteIngestMatchesEvaluateFiniteGCF_OnSeveralFixtures(t *testing.T) {
	cases := []struct {
		name  string
		terms [][2]int64
	}{
		{"single", [][2]int64{{5, 1}}},
		{"simple", [][2]int64{{1, 1}, {2, 1}, {3, 1}}},
		{"mixed-q", [][2]int64{{1, 2}, {3, 4}, {5, 6}}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src1 := NewSliceGCF(tc.terms...)
			src2 := NewSliceGCF(tc.terms...)

			got := collectTerms(NewGCFStream(src1, GCFStreamOptions{}), 64)

			wantRat, err := EvaluateFiniteGCF(src2)
			if err != nil {
				t.Fatalf("EvaluateFiniteGCF failed: %v", err)
			}
			want := collectTerms(NewRationalCF(wantRat), 64)

			if len(got) != len(want) {
				t.Fatalf("len mismatch: got=%v want=%v", got, want)
			}
			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("digit %d: got=%v want=%v", i, got, want)
				}
			}
		})
	}
}

func TestGCFStream_AdaptedRegularCFInfinitePrefixCanEmitEarly(t *testing.T) {
	s := NewGCFStream(AdaptCFToGCF(NewSliceCF(1, 2, 3)), GCFStreamOptions{})

	got := collectTerms(s, 32)
	want := []int64{1, 2, 3}

	if len(got) != len(want) {
		t.Fatalf("len mismatch: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, got, want)
		}
	}
	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

func TestGCFStream_LambertCanEmitAtLeastOneDigit(t *testing.T) {
	s := NewGCFStream(NewLambertPiOver4GCFSource(), GCFStreamOptions{})

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected at least one digit, err=%v", s.Err())
	}

	// pi/4 is in (0,1), so first ordinary CF digit should be 0.
	if d != 0 {
		t.Fatalf("got %d want 0", d)
	}
}

func TestGCFStream_BrounckerCanEmitAtLeastOneDigit(t *testing.T) {
	s := NewGCFStream(NewBrouncker4OverPiGCFSource(), GCFStreamOptions{})

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected at least one digit, err=%v", s.Err())
	}

	// 4/pi is > 1, so first ordinary CF digit should be 1.
	if d != 1 {
		t.Fatalf("got %d want 1", d)
	}
}
func TestGCFStream_LambertFirstTwoDigits(t *testing.T) {
	s := NewGCFStream(NewLambertPiOver4GCFSource(), GCFStreamOptions{})

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if d != 0 {
		t.Fatalf("got first digit %d want 0", d)
	}

	d, ok = s.Next()
	if !ok {
		t.Fatalf("expected second digit, err=%v", s.Err())
	}

	// pi/4 ≈ 0.785..., so ordinary CF begins [0;1,...]
	if d != 1 {
		t.Fatalf("got second digit %d want 1", d)
	}

	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

func TestGCFStream_BrounckerFirstTwoDigits(t *testing.T) {
	s := NewGCFStream(NewBrouncker4OverPiGCFSource(), GCFStreamOptions{})

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if d != 1 {
		t.Fatalf("got first digit %d want 1", d)
	}

	d, ok = s.Next()
	if !ok {
		t.Fatalf("expected second digit, err=%v", s.Err())
	}

	// 4/pi ≈ 1.273..., so ordinary CF begins [1;3,...]
	if d != 3 {
		t.Fatalf("got second digit %d want 3", d)
	}

	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

func TestGCFStream_FiniteSourceStaysExhausted(t *testing.T) {
	s := NewGCFStream(NewSliceGCF([2]int64{3, 1}), GCFStreamOptions{})

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if d != 3 {
		t.Fatalf("got %d want 3", d)
	}

	_, ok = s.Next()
	if ok {
		t.Fatalf("expected first exhaustion")
	}
	_, ok = s.Next()
	if ok {
		t.Fatalf("expected repeated exhaustion to stay false")
	}
	if err := s.Err(); err != nil {
		t.Fatalf("expected clean termination, got %v", err)
	}
}

func TestGCFStream_BrounckerDoesNotEmitWrongSecondDigit(t *testing.T) {
	s := NewGCFStream(NewBrouncker4OverPiGCFSource(), GCFStreamOptions{})

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if d != 1 {
		t.Fatalf("got first digit %d want 1", d)
	}

	d, ok = s.Next()
	if !ok {
		t.Fatalf("expected second digit, err=%v", s.Err())
	}
	if d == 2 {
		t.Fatalf("unsound second digit 2 for Brouncker 4/pi")
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

	// After one more ingested GCF term, the same generic TailRange() contract
	// should support another sound emission.
	if src.calls != 3 {
		t.Fatalf("expected second digit after exactly 3 ingested terms, got %d", src.calls)
	}

	// We do not hard-code the exact second digit here; this test is about the
	// ingestion/emission cadence enabled by the generic TailRange() contract.
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
	case 0:
		// Before any ingestion, too wide to prove anything useful.
		return NewRange(mustRat(1, 1), mustRat(3, 1), true, true)
	case 1:
		// Still too wide.
		return NewRange(mustRat(1, 1), mustRat(3, 1), true, true)
	default:
		// Tightens after enough ingestion; now should permit first digit 1.
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

	// If GCFStream were reusing stale constructor-time TailRange evidence,
	// it would not see the tightened [1,2] range and this would either take
	// longer or fail to emit here.
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
	// Wider than [1,2], so it should not prove the first digit as early.
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
	// Stronger explicit enclosure.
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

	// With the explicit TailRange() contract, the stream should emit after
	// exactly two ingested terms for this source family. If it were behaving
	// like a weaker lower-bound-ray-only source, it would need at least as
	// much evidence and possibly more.
	if src.calls != 2 {
		t.Fatalf("expected TailRange-driven emission after exactly 2 ingested terms, got %d", src.calls)
	}

	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

type finitePrefixGCFSource struct {
	src   GCFSource
	limit int
	n     int
}

func newFinitePrefixGCFSource(src GCFSource, limit int) *finitePrefixGCFSource {
	return &finitePrefixGCFSource{
		src:   src,
		limit: limit,
	}
}

func (s *finitePrefixGCFSource) NextPQ() (int64, int64, bool) {
	if s.n >= s.limit {
		return 0, 0, false
	}
	p, q, ok := s.src.NextPQ()
	if !ok {
		return 0, 0, false
	}
	s.n++
	return p, q, true
}

func collectFinitePrefixTerms(src GCFSource, n int) [][2]int64 {
	var out [][2]int64
	for i := 0; i < n; i++ {
		p, q, ok := src.NextPQ()
		if !ok {
			break
		}
		out = append(out, [2]int64{p, q})
	}
	return out
}

func TestGCFStream_FiniteLambertPrefixMatchesEvaluateFiniteGCF(t *testing.T) {
	for _, prefixLen := range []int{1, 2, 3, 4, 5, 6} {
		t.Run(fmt.Sprintf("prefix_%d", prefixLen), func(t *testing.T) {
			streamSrc := newFinitePrefixGCFSource(NewLambertPiOver4GCFSource(), prefixLen)
			evalTerms := collectFinitePrefixTerms(NewLambertPiOver4GCFSource(), prefixLen)

			got := collectTerms(NewGCFStream(streamSrc, GCFStreamOptions{}), 64)

			wantRat, err := EvaluateFiniteGCF(NewSliceGCF(evalTerms...))
			if err != nil {
				t.Fatalf("EvaluateFiniteGCF failed: %v", err)
			}
			want := collectTerms(NewRationalCF(wantRat), 64)

			if len(got) != len(want) {
				t.Fatalf("len mismatch: got=%v want=%v", got, want)
			}
			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("digit %d: got=%v want=%v", i, got, want)
				}
			}
		})
	}
}

func TestGCFStream_FiniteBrounckerPrefixMatchesEvaluateFiniteGCF(t *testing.T) {
	for _, prefixLen := range []int{1, 2, 3, 4, 5, 6} {
		t.Run(fmt.Sprintf("prefix_%d", prefixLen), func(t *testing.T) {
			streamSrc := newFinitePrefixGCFSource(NewBrouncker4OverPiGCFSource(), prefixLen)
			evalTerms := collectFinitePrefixTerms(NewBrouncker4OverPiGCFSource(), prefixLen)

			got := collectTerms(NewGCFStream(streamSrc, GCFStreamOptions{}), 64)

			wantRat, err := EvaluateFiniteGCF(NewSliceGCF(evalTerms...))
			if err != nil {
				t.Fatalf("EvaluateFiniteGCF failed: %v", err)
			}
			want := collectTerms(NewRationalCF(wantRat), 64)

			if len(got) != len(want) {
				t.Fatalf("len mismatch: got=%v want=%v", got, want)
			}
			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("digit %d: got=%v want=%v", i, got, want)
				}
			}
		})
	}
}

// gcf_stream_test.go v1
