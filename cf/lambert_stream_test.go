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

//EOF lambert_stream_test.go
