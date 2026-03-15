// sqrt_stream_iface_test.go v1
package cf

import "testing"

func TestSqrtStream_ReturnsInspectableStream(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	s, err := SqrtStream(Sqrt2CF(), 2, p)
	if err != nil {
		t.Fatalf("SqrtStream failed: %v", err)
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
	if s.Err() != nil {
		t.Fatalf("expected nil Err before start, got %v", s.Err())
	}
}

func TestSqrtStream_SnapshotAfterStartCarriesApproximation(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	s, err := SqrtStream(Sqrt2CF(), 2, p)
	if err != nil {
		t.Fatalf("SqrtStream failed: %v", err)
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
	if snap.CFInputApprox == nil {
		t.Fatalf("expected non-nil CFInputApprox after start")
	}
}
