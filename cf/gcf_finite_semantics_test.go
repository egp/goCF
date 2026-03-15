// gcf_finite_semantics_test.go v1
package cf

import "testing"

func TestFiniteGCFSemanticContract_AllCanonicalPathsAgree(t *testing.T) {
	cases := []struct {
		name  string
		terms [][2]int64
	}{
		{
			name:  "single",
			terms: [][2]int64{{5, 1}},
		},
		{
			name:  "simple-q1",
			terms: [][2]int64{{1, 1}, {2, 1}, {3, 1}},
		},
		{
			name:  "mixed-q",
			terms: [][2]int64{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			name:  "two-terms",
			terms: [][2]int64{{3, 2}, {5, 7}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Path 1: direct exact finite evaluation.
			exact, err := EvaluateFiniteGCF(NewSliceGCF(tc.terms...))
			if err != nil {
				t.Fatalf("EvaluateFiniteGCF failed: %v", err)
			}

			// Path 2: ingest-all bounder convergent.
			b, err := IngestAllGCF(NewSliceGCF(tc.terms...))
			if err != nil {
				t.Fatalf("IngestAllGCF failed: %v", err)
			}

			gotConv, err := b.Convergent()
			if err != nil {
				t.Fatalf("Convergent failed: %v", err)
			}
			if gotConv.Cmp(exact) != 0 {
				t.Fatalf("Convergent mismatch: got %v want %v", gotConv, exact)
			}

			// Path 3: finished bounder exact range.
			r, ok, err := b.Range()
			if err != nil {
				t.Fatalf("Range failed: %v", err)
			}
			if !ok {
				t.Fatalf("expected exact finished range")
			}
			if r.Lo.Cmp(exact) != 0 || r.Hi.Cmp(exact) != 0 {
				t.Fatalf("finished range mismatch: got [%v,%v] want exact [%v,%v]", r.Lo, r.Hi, exact, exact)
			}

			// Path 4: finite GCF stream emitted regular CF terms.
			gotTerms := collectTerms(NewGCFStream(NewSliceGCF(tc.terms...), GCFStreamOptions{}), 64)
			wantTerms := collectTerms(NewRationalCF(exact), 64)

			if len(gotTerms) != len(wantTerms) {
				t.Fatalf("CF term length mismatch: got=%v want=%v", gotTerms, wantTerms)
			}
			for i := range wantTerms {
				if gotTerms[i] != wantTerms[i] {
					t.Fatalf("CF term mismatch at %d: got=%v want=%v", i, gotTerms, wantTerms)
				}
			}
		})
	}
}

func TestFiniteGCFSemanticContract_IngestPrefixThenFinishMatchesExactWhenSourceEndsEarly(t *testing.T) {
	terms := [][2]int64{{3, 2}, {5, 7}}

	b, err := IngestGCFPrefix(NewSliceGCF(terms...), 10)
	if err != nil {
		t.Fatalf("IngestGCFPrefix failed: %v", err)
	}

	exact, err := EvaluateFiniteGCF(NewSliceGCF(terms...))
	if err != nil {
		t.Fatalf("EvaluateFiniteGCF failed: %v", err)
	}

	gotConv, err := b.Convergent()
	if err != nil {
		t.Fatalf("Convergent failed: %v", err)
	}
	if gotConv.Cmp(exact) != 0 {
		t.Fatalf("Convergent mismatch: got %v want %v", gotConv, exact)
	}

	r, ok, err := b.Range()
	if err != nil {
		t.Fatalf("Range failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected exact finished range")
	}
	if r.Lo.Cmp(exact) != 0 || r.Hi.Cmp(exact) != 0 {
		t.Fatalf("finished range mismatch: got [%v,%v] want exact [%v,%v]", r.Lo, r.Hi, exact, exact)
	}
}

func TestFiniteGCFSemanticContract_AdaptedRegularFiniteCFAgreesWithRoundTrip(t *testing.T) {
	orig := NewSliceCF(1, 2, 3, 4)
	got := collectTerms(NewGCFStream(AdaptCFToGCF(NewSliceCF(1, 2, 3, 4)), GCFStreamOptions{}), 64)
	want := collectTerms(orig, 64)

	if len(got) != len(want) {
		t.Fatalf("len mismatch: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, got, want)
		}
	}
}

// gcf_finite_semantics_test.go v1
