// gcf_stream_named_sources_test.go v1
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

// gcf_stream_named_sources_test.go v1
