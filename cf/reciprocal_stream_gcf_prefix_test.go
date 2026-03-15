// reciprocal_stream_gcf_prefix_test.go v1
package cf

import "testing"

func TestReciprocalGCFPrefixStream_ReturnsInspectableStream(t *testing.T) {
	s, err := ReciprocalGCFPrefixStream(AdaptCFToGCF(Sqrt2CF()), 2)
	if err != nil {
		t.Fatalf("ReciprocalGCFPrefixStream failed: %v", err)
	}

	snap := s.Snapshot()
	if snap.Started {
		t.Fatalf("expected Started=false before first Next")
	}
	if snap.Approx != nil {
		t.Fatalf("expected Approx=nil before start")
	}
	if snap.GCFInputApprox != nil {
		t.Fatalf("expected GCFInputApprox=nil before start")
	}
	if s.Err() != nil {
		t.Fatalf("expected nil Err before start, got %v", s.Err())
	}
}

func TestReciprocalGCFPrefixStream_AfterStartCarriesApproximation(t *testing.T) {
	s, err := ReciprocalGCFPrefixStream(AdaptCFToGCF(Sqrt2CF()), 2)
	if err != nil {
		t.Fatalf("ReciprocalGCFPrefixStream failed: %v", err)
	}

	_, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit; err=%v", s.Err())
	}

	snap := s.Snapshot()
	if !snap.Started {
		t.Fatalf("expected Started=true after start")
	}
	if snap.Approx == nil {
		t.Fatalf("expected non-nil Approx after start")
	}
	if snap.GCFInputApprox == nil {
		t.Fatalf("expected non-nil GCFInputApprox after start")
	}

	// For sqrt(2) first two terms as GCF-adapted CF, convergent is 3/2, reciprocal is 2/3.
	want := mustRat(2, 3)
	if snap.Approx.Cmp(want) != 0 {
		t.Fatalf("got Approx=%v want %v", *snap.Approx, want)
	}
}

func TestReciprocalGCFPrefixStream_EmitsExpectedCF(t *testing.T) {
	s, err := ReciprocalGCFPrefixStream(AdaptCFToGCF(Sqrt2CF()), 2)
	if err != nil {
		t.Fatalf("ReciprocalGCFPrefixStream failed: %v", err)
	}

	got := collectTerms(s, 8)
	if err := s.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}

	want := []int64{0, 1, 2} // 2/3
	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d gotFull=%v", i, got[i], want[i], got)
		}
	}
}

func TestReciprocalGCFPrefixStream_FiniteSourceCarriesInputApprox(t *testing.T) {
	s, err := ReciprocalGCFPrefixStream(NewSliceGCF(
		[2]int64{3, 2},
		[2]int64{5, 7},
	), 10)
	if err != nil {
		t.Fatalf("ReciprocalGCFPrefixStream failed: %v", err)
	}

	_, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit; err=%v", s.Err())
	}

	snap := s.Snapshot()
	if snap.GCFInputApprox == nil {
		t.Fatalf("expected non-nil GCFInputApprox")
	}

	// Existing GCFApprox finite-source semantics in this codebase give convergent 17/5.
	wantConv := mustRat(17, 5)
	if snap.GCFInputApprox.Convergent.Cmp(wantConv) != 0 {
		t.Fatalf("got convergent %v want %v", snap.GCFInputApprox.Convergent, wantConv)
	}

	wantRecip := mustRat(5, 17)
	if snap.Approx.Cmp(wantRecip) != 0 {
		t.Fatalf("got reciprocal approx %v want %v", *snap.Approx, wantRecip)
	}
}

func TestReciprocalGCFPrefixStream_ReciprocalOfZeroFails(t *testing.T) {
	s, err := ReciprocalGCFPrefixStream(NewSliceGCF([2]int64{0, 1}), 1)
	if err != nil {
		t.Fatalf("unexpected constructor error: %v", err)
	}

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected failure, not a digit")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}
}

func TestReciprocalGCFPrefixStream_RejectsBadPrefixTerms(t *testing.T) {
	_, err := ReciprocalGCFPrefixStream(AdaptCFToGCF(Sqrt2CF()), 0)
	if err == nil {
		t.Fatalf("expected constructor error")
	}
}

// reciprocal_stream_gcf_prefix_test.go v1
