// gcf_stream_named_sources_test.go v4
package cf

import (
	"fmt"
	"testing"
)

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
	if d != 3 {
		t.Fatalf("got second digit %d want 3", d)
	}

	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
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

func TestGCFStream_LambertInfiniteSourceEmitsFirstDigit(t *testing.T) {
	s := NewGCFStream(NewLambertPiOver4GCFSource(), GCFStreamOptions{})

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if d != 0 {
		t.Fatalf("got first digit %d want 0", d)
	}
	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

func TestGCFStream_BrounckerInfiniteSourceEmitsFirstDigit(t *testing.T) {
	s := NewGCFStream(NewBrouncker4OverPiGCFSource(), GCFStreamOptions{})

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if d != 1 {
		t.Fatalf("got first digit %d want 1", d)
	}
	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

func TestGCFStream_LambertInfinitePrefixMatchesStabilizedFinitePrefixes(t *testing.T) {
	lambertFactory := func() GCFSource { return NewLambertPiOver4GCFSource() }

	want8 := exactDigitsFromFinitePrefix(t, lambertFactory, 8, 3)
	want10 := exactDigitsFromFinitePrefix(t, lambertFactory, 10, 3)

	if len(want8) != len(want10) {
		t.Fatalf("stabilization len mismatch: want8=%v want10=%v", want8, want10)
	}
	for i := range want8 {
		if want8[i] != want10[i] {
			t.Fatalf("Lambert finite prefixes not stabilized at digit %d: p8=%v p10=%v", i, want8, want10)
		}
	}

	got := collectTerms(NewGCFStream(NewLambertPiOver4GCFSource(), GCFStreamOptions{}), 3)

	if len(got) != len(want10) {
		t.Fatalf("len mismatch: got=%v want=%v", got, want10)
	}
	for i := range want10 {
		if got[i] != want10[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, got, want10)
		}
	}
}

func TestGCFStream_BrounckerInfinitePrefixMatchesStabilizedFinitePrefixes(t *testing.T) {
	brounckerFactory := func() GCFSource { return NewBrouncker4OverPiGCFSource() }

	want8 := exactDigitsFromFinitePrefix(t, brounckerFactory, 8, 2)
	want10 := exactDigitsFromFinitePrefix(t, brounckerFactory, 10, 2)

	if len(want8) != len(want10) {
		t.Fatalf("stabilization len mismatch: want8=%v want10=%v", want8, want10)
	}
	for i := range want8 {
		if want8[i] != want10[i] {
			t.Fatalf("Brouncker finite prefixes not stabilized at digit %d: p8=%v p10=%v", i, want8, want10)
		}
	}

	got := collectTerms(NewGCFStream(NewBrouncker4OverPiGCFSource(), GCFStreamOptions{}), 2)

	if len(got) != len(want10) {
		t.Fatalf("len mismatch: got=%v want=%v", got, want10)
	}
	for i := range want10 {
		if got[i] != want10[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, got, want10)
		}
	}
}
func TestLambertPiOver4GCFSource_TailEvidenceMatchesHelperFunctions(t *testing.T) {
	src := NewLambertPiOver4GCFSource()

	check := func(wantPrefix int) {
		t.Helper()

		ev := src.TailEvidence()

		if ev.LowerBound == nil {
			t.Fatalf("prefix %d: expected non-nil lower bound", wantPrefix)
		}

		wantLB := LambertPiOver4TailLowerBoundAfterPrefix(wantPrefix)
		if ev.LowerBound.String() != wantLB.String() {
			t.Fatalf("prefix %d: lower bound mismatch: got=%v want=%v", wantPrefix, *ev.LowerBound, wantLB)
		}

		wantRange, wantOK, err := LambertPiOver4TailRangeAfterPrefix(wantPrefix)
		if err != nil {
			t.Fatalf("prefix %d: helper returned err: %v", wantPrefix, err)
		}

		if wantOK != (ev.Range != nil) {
			t.Fatalf("prefix %d: range presence mismatch: got=%v want=%v", wantPrefix, ev.Range != nil, wantOK)
		}
		if wantOK && ev.Range != nil && ev.Range.String() != wantRange.String() {
			t.Fatalf("prefix %d: range mismatch: got=%v want=%v", wantPrefix, *ev.Range, wantRange)
		}

		if ev.RangeReusable {
			t.Fatalf("prefix %d: expected non-reusable Lambert tail range", wantPrefix)
		}
		if ev.LowerBoundMinPrefix != 0 {
			t.Fatalf("prefix %d: expected LowerBoundMinPrefix=0 got %d", wantPrefix, ev.LowerBoundMinPrefix)
		}
	}

	check(0)

	_, _, ok := src.NextPQ()
	if !ok {
		t.Fatalf("expected first Lambert term")
	}
	check(1)

	_, _, ok = src.NextPQ()
	if !ok {
		t.Fatalf("expected second Lambert term")
	}
	check(2)
}

func TestBrouncker4OverPiGCFSource_TailEvidenceMatchesHelperFunctions(t *testing.T) {
	src := NewBrouncker4OverPiGCFSource()

	check := func(wantPrefix int) {
		t.Helper()

		ev := src.TailEvidence()

		if ev.LowerBound == nil {
			t.Fatalf("prefix %d: expected non-nil lower bound", wantPrefix)
		}

		wantLB := Brouncker4OverPiTailLowerBoundAfterPrefix(wantPrefix)
		if ev.LowerBound.String() != wantLB.String() {
			t.Fatalf("prefix %d: lower bound mismatch: got=%v want=%v", wantPrefix, *ev.LowerBound, wantLB)
		}

		wantRange, wantOK, err := Brouncker4OverPiTailRangeAfterPrefix(wantPrefix)
		if err != nil {
			t.Fatalf("prefix %d: helper returned err: %v", wantPrefix, err)
		}

		if wantOK != (ev.Range != nil) {
			t.Fatalf("prefix %d: range presence mismatch: got=%v want=%v", wantPrefix, ev.Range != nil, wantOK)
		}
		if wantOK && ev.Range != nil && ev.Range.String() != wantRange.String() {
			t.Fatalf("prefix %d: range mismatch: got=%v want=%v", wantPrefix, *ev.Range, wantRange)
		}

		if ev.RangeReusable {
			t.Fatalf("prefix %d: expected non-reusable Brouncker tail range", wantPrefix)
		}
		if ev.LowerBoundMinPrefix != 2 {
			t.Fatalf("prefix %d: expected LowerBoundMinPrefix=2 got %d", wantPrefix, ev.LowerBoundMinPrefix)
		}
	}

	check(0)

	_, _, ok := src.NextPQ()
	if !ok {
		t.Fatalf("expected first Brouncker term")
	}
	check(1)

	_, _, ok = src.NextPQ()
	if !ok {
		t.Fatalf("expected second Brouncker term")
	}
	check(2)
}

type brounckerLowerBoundOnlyStreamSource struct {
	src *Brouncker4OverPiGCFSource
}

func newBrounckerLowerBoundOnlyStreamSource() *brounckerLowerBoundOnlyStreamSource {
	return &brounckerLowerBoundOnlyStreamSource{src: NewBrouncker4OverPiGCFSource()}
}

func (s *brounckerLowerBoundOnlyStreamSource) NextPQ() (int64, int64, bool) {
	return s.src.NextPQ()
}

func (s *brounckerLowerBoundOnlyStreamSource) TailLowerBound() Rational {
	return mustRat(1, 1)
}

func (s *brounckerLowerBoundOnlyStreamSource) LowerBoundRayMinPrefix() int {
	return 2
}

func collectDigitsAndLambertCalls(t *testing.T, s ContinuedFraction, src *LambertPiOver4GCFSource, n int) ([]int64, []int) {
	t.Helper()

	var digits []int64
	var calls []int
	for i := 0; i < n; i++ {
		d, ok := s.Next()
		if !ok {
			if gs, ok := s.(*GCFStream); ok && gs.Err() != nil {
				t.Fatalf("digit %d: unexpected stream error: %v", i, gs.Err())
			}
			t.Fatalf("digit %d: expected digit", i)
		}
		digits = append(digits, d)
		calls = append(calls, src.i)
	}
	return digits, calls
}

func collectDigitsAndBrounckerCalls(t *testing.T, s ContinuedFraction, src *Brouncker4OverPiGCFSource, n int) ([]int64, []int) {
	t.Helper()

	var digits []int64
	var calls []int
	for i := 0; i < n; i++ {
		d, ok := s.Next()
		if !ok {
			if gs, ok := s.(*GCFStream); ok && gs.Err() != nil {
				t.Fatalf("digit %d: unexpected stream error: %v", i, gs.Err())
			}
			t.Fatalf("digit %d: expected digit", i)
		}
		digits = append(digits, d)
		calls = append(calls, src.i)
	}
	return digits, calls
}

func TestGCFStream_LambertInfinite_SpecializedEvidenceCadencePayoff(t *testing.T) {
	lambertFactory := func() GCFSource { return NewLambertPiOver4GCFSource() }

	want10 := exactDigitsFromFinitePrefix(t, lambertFactory, 10, 3)
	want12 := exactDigitsFromFinitePrefix(t, lambertFactory, 12, 3)

	if len(want10) != len(want12) {
		t.Fatalf("stabilization len mismatch: want10=%v want12=%v", want10, want12)
	}
	for i := range want10 {
		if want10[i] != want12[i] {
			t.Fatalf("Lambert finite prefixes not stabilized at digit %d: p10=%v p12=%v", i, want10, want12)
		}
	}

	specSrc := NewLambertPiOver4GCFSource()
	genWrap := newLambertLowerBoundOnlyStreamSource()

	specStream := NewGCFStream(specSrc, GCFStreamOptions{})
	genStream := NewGCFStream(genWrap, GCFStreamOptions{})

	gotSpec, specCalls := collectDigitsAndLambertCalls(t, specStream, specSrc, 3)
	gotGen, genCalls := collectDigitsAndLambertCalls(t, genStream, genWrap.src, 3)

	for i := range want12 {
		if gotSpec[i] != want12[i] {
			t.Fatalf("specialized digit %d: got=%v want=%v", i, gotSpec, want12)
		}
		if gotGen[i] != want12[i] {
			t.Fatalf("generic digit %d: got=%v want=%v", i, gotGen, want12)
		}
	}

	for i := range specCalls {
		if specCalls[i] > genCalls[i] {
			t.Fatalf("Lambert digit %d: specialized used more ingestion: specCalls=%v genCalls=%v", i, specCalls, genCalls)
		}
	}

	strictGain := false
	for i := range specCalls {
		if specCalls[i] < genCalls[i] {
			strictGain = true
			break
		}
	}
	if !strictGain {
		t.Fatalf("expected Lambert specialized evidence to improve cadence for at least one digit: specCalls=%v genCalls=%v", specCalls, genCalls)
	}

	if err := specStream.Err(); err != nil {
		t.Fatalf("specialized stream: unexpected err=%v", err)
	}
	if err := genStream.Err(); err != nil {
		t.Fatalf("generic stream: unexpected err=%v", err)
	}
}

func TestGCFStream_BrounckerInfinite_SpecializedEvidenceCadencePayoff(t *testing.T) {
	brounckerFactory := func() GCFSource { return NewBrouncker4OverPiGCFSource() }

	want10 := exactDigitsFromFinitePrefix(t, brounckerFactory, 10, 2)
	want12 := exactDigitsFromFinitePrefix(t, brounckerFactory, 12, 2)

	if len(want10) != len(want12) {
		t.Fatalf("stabilization len mismatch: want10=%v want12=%v", want10, want12)
	}
	for i := range want10 {
		if want10[i] != want12[i] {
			t.Fatalf("Brouncker finite prefixes not stabilized at digit %d: p10=%v p12=%v", i, want10, want12)
		}
	}

	specSrc := NewBrouncker4OverPiGCFSource()
	genWrap := newBrounckerLowerBoundOnlyStreamSource()

	specStream := NewGCFStream(specSrc, GCFStreamOptions{})
	genStream := NewGCFStream(genWrap, GCFStreamOptions{})

	gotSpec, specCalls := collectDigitsAndBrounckerCalls(t, specStream, specSrc, 2)
	gotGen, genCalls := collectDigitsAndBrounckerCalls(t, genStream, genWrap.src, 2)

	for i := range want12 {
		if gotSpec[i] != want12[i] {
			t.Fatalf("specialized digit %d: got=%v want=%v", i, gotSpec, want12)
		}
		if gotGen[i] != want12[i] {
			t.Fatalf("generic digit %d: got=%v want=%v", i, gotGen, want12)
		}
	}

	for i := range specCalls {
		if specCalls[i] > genCalls[i] {
			t.Fatalf("Brouncker digit %d: specialized used more ingestion: specCalls=%v genCalls=%v", i, specCalls, genCalls)
		}
	}

	strictGain := false
	for i := range specCalls {
		if specCalls[i] < genCalls[i] {
			strictGain = true
			break
		}
	}
	if !strictGain {
		t.Fatalf("expected Brouncker specialized evidence to improve cadence for at least one digit: specCalls=%v genCalls=%v", specCalls, genCalls)
	}

	if err := specStream.Err(); err != nil {
		t.Fatalf("specialized stream: unexpected err=%v", err)
	}
	if err := genStream.Err(); err != nil {
		t.Fatalf("generic stream: unexpected err=%v", err)
	}
}

// gcf_stream_named_sources_test.go v4
