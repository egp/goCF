// gcf_stream_test.go v1
package cf

import "testing"

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

// gcf_stream_test.go v1
