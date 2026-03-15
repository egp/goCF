// reciprocal_finite_semantics_test.go v2
package cf

import "testing"

func TestReciprocalFiniteSemanticContract_PrefixAgreesWithFiniteConventionOnFiniteSource(t *testing.T) {
	terms := [][2]int64{
		{3, 2},
		{5, 7},
	}

	wantRat := mustRat(5, 17) // reciprocal of finite convention value 17/5
	wantTerms := collectTerms(NewRationalCF(wantRat), 32)

	s, err := ReciprocalGCFPrefixStream(NewSliceGCF(terms...), 10)
	if err != nil {
		t.Fatalf("ReciprocalGCFPrefixStream failed: %v", err)
	}

	got := collectTerms(s, 32)
	if err := s.Err(); err != nil {
		t.Fatalf("prefix stream error: %v", err)
	}

	if len(got) != len(wantTerms) {
		t.Fatalf("len mismatch: got=%v want=%v", got, wantTerms)
	}
	for i := range wantTerms {
		if got[i] != wantTerms[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, got, wantTerms)
		}
	}
}

func TestReciprocalExactTailSemanticContract_ExactTailAgreesWithDirectTailEvaluation(t *testing.T) {
	s, err := ReciprocalGCFExactTailStreamWithTail(
		NewSliceGCF(
			[2]int64{3, 2},
			[2]int64{5, 7},
		),
		mustRat(1, 1),
		10,
	)
	if err != nil {
		t.Fatalf("ReciprocalGCFExactTailStreamWithTail failed: %v", err)
	}

	got := collectTerms(s, 32)
	if err := s.Err(); err != nil {
		t.Fatalf("exact-tail stream error: %v", err)
	}

	// x = 3 + 2/(5 + 7/1) = 19/6, so reciprocal = 6/19 = [0;3,6].
	want := []int64{0, 3, 6}

	if len(got) != len(want) {
		t.Fatalf("len mismatch: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, got, want)
		}
	}
}
func TestReciprocalGCFExactTailStream2_Snapshot_BeforeStartSeparatesConfiguredAndConsumedTerms(t *testing.T) {
	s := NewReciprocalGCFExactTailStreamWithTail2(
		NewSliceGCF([2]int64{3, 2}),
		mustRat(11, 1),
		8,
	)

	snap := s.Snapshot()
	if snap.Started {
		t.Fatalf("expected Started=false")
	}
	if snap.MaxIngestTerms != 8 {
		t.Fatalf("got MaxIngestTerms=%d want 8", snap.MaxIngestTerms)
	}
	if snap.ConsumedTerms != 0 {
		t.Fatalf("got ConsumedTerms=%d want 0 before start", snap.ConsumedTerms)
	}
}
