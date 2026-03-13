// lambert_stream_test.go
package cf

import "testing"

func TestGCFStream_LambertInfinite_FirstTwoDigitsAndCadence(t *testing.T) {
	src := NewLambertPiOver4GCFSource()
	s := NewGCFStream(src, GCFStreamOptions{})

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if d != 0 {
		t.Fatalf("got first digit %d want 0", d)
	}
	callsAfterFirst := src.i

	d, ok = s.Next()
	if !ok {
		t.Fatalf("expected second digit, err=%v", s.Err())
	}
	if d != 1 {
		t.Fatalf("got second digit %d want 1", d)
	}
	callsAfterSecond := src.i

	if callsAfterSecond <= callsAfterFirst {
		t.Fatalf("expected additional Lambert ingestion for second digit, first=%d second=%d", callsAfterFirst, callsAfterSecond)
	}

	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}
func TestGCFStream_LambertInfinite_FirstThreeDigitsMatchStabilizedFinitePrefixes(t *testing.T) {
	factory := func() GCFSource { return NewLambertPiOver4GCFSource() }

	want8 := exactDigitsFromFinitePrefix(t, factory, 8, 3)
	want10 := exactDigitsFromFinitePrefix(t, factory, 10, 3)

	if len(want8) != len(want10) {
		t.Fatalf("stabilization len mismatch: want8=%v want10=%v", want8, want10)
	}
	for i := range want8 {
		if want8[i] != want10[i] {
			t.Fatalf("Lambert finite prefixes not stabilized at digit %d: p8=%v p10=%v", i, want8, want10)
		}
	}

	s := NewGCFStream(NewLambertPiOver4GCFSource(), GCFStreamOptions{})
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

type lambertLowerBoundOnlyStreamSource struct {
	src *LambertPiOver4GCFSource
}

func newLambertLowerBoundOnlyStreamSource() *lambertLowerBoundOnlyStreamSource {
	return &lambertLowerBoundOnlyStreamSource{src: NewLambertPiOver4GCFSource()}
}

func (s *lambertLowerBoundOnlyStreamSource) NextPQ() (int64, int64, bool) {
	return s.src.NextPQ()
}

func (s *lambertLowerBoundOnlyStreamSource) TailLowerBound() Rational {
	return mustRat(1, 1)
}

func TestGCFStream_LambertInfinite_SpecializedEvidenceBeatsLowerBoundOnlyCadence(t *testing.T) {
	specSrc := NewLambertPiOver4GCFSource()
	genSrc := newLambertLowerBoundOnlyStreamSource()

	spec := NewGCFStream(specSrc, GCFStreamOptions{})
	gen := NewGCFStream(genSrc, GCFStreamOptions{})

	want := []int64{0, 1, 3}

	for i, w := range want {
		d, ok := spec.Next()
		if !ok {
			t.Fatalf("specialized stream: expected digit %d, err=%v", i, spec.Err())
		}
		if d != w {
			t.Fatalf("specialized stream digit %d: got %d want %d", i, d, w)
		}

		d, ok = gen.Next()
		if !ok {
			t.Fatalf("generic stream: expected digit %d, err=%v", i, gen.Err())
		}
		if d != w {
			t.Fatalf("generic stream digit %d: got %d want %d", i, d, w)
		}
	}

	specCalls := specSrc.i
	genCalls := genSrc.src.i

	if specCalls > genCalls {
		t.Fatalf("expected specialized Lambert evidence to use no more ingestion than lower-bound-only baseline, specialized=%d generic=%d", specCalls, genCalls)
	}

	if err := spec.Err(); err != nil {
		t.Fatalf("specialized stream: unexpected err=%v", err)
	}
	if err := gen.Err(); err != nil {
		t.Fatalf("generic stream: unexpected err=%v", err)
	}
}

//EOF lambert_stream_test.go
