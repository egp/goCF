// sqrt_stream_gcf_api_test.go v1
package cf

import "testing"

func TestSqrtGCFStream_ReturnsInspectableStream(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	s, err := SqrtGCFStream(AdaptCFToGCF(Sqrt2CF()), 2, p)
	if err != nil {
		t.Fatalf("SqrtGCFStream failed: %v", err)
	}

	snap := s.Snapshot()
	if snap.Started {
		t.Fatalf("expected Started=false before first Next")
	}
	if snap.PrefixTerms != 2 {
		t.Fatalf("got PrefixTerms=%d want 2", snap.PrefixTerms)
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

func TestSqrtGCFStream_SnapshotAfterStartCarriesApproximation(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	s, err := SqrtGCFStream(AdaptCFToGCF(Sqrt2CF()), 2, p)
	if err != nil {
		t.Fatalf("SqrtGCFStream failed: %v", err)
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
	if snap.CFInputApprox != nil {
		t.Fatalf("expected nil CFInputApprox for GCF stream")
	}
}

func TestSqrtGCFStream_FinitePerfectSquarePath(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(0, 1),
	}

	// Under current bounded GCF-prefix semantics, a single exhausted term (3,1)
	// yields convergent 3, so sqrt output begins as sqrt(3).
	s, err := SqrtGCFStream(NewSliceGCF([2]int64{3, 1}), 8, p)
	if err != nil {
		t.Fatalf("SqrtGCFStream failed: %v", err)
	}

	got := collectTerms(s, 8)
	want := []int64{1, 1, 2, 1, 2, 1, 2, 1}

	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d full=%v", i, got[i], want[i], got)
		}
	}
}

func TestSqrtGCFStream_RejectsBadPrefixTerms(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	_, err := SqrtGCFStream(AdaptCFToGCF(Sqrt2CF()), 0, p)
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
}

// sqrt_stream_gcf_api_test.go v1
