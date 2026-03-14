// sqrt_stream_cf2_snapshot_test.go v1
package cf

import "testing"

func TestSqrtCFPrefixStream2_Snapshot_BeforeStart(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	s := NewSqrtCFPrefixStream2(Sqrt2CF(), 2, p)
	snap := s.Snapshot()

	if snap.Started {
		t.Fatalf("expected Started=false")
	}
	if snap.PrefixTerms != 2 {
		t.Fatalf("got PrefixTerms=%d want 2", snap.PrefixTerms)
	}
	if snap.Approx != nil {
		t.Fatalf("expected Approx=nil before start")
	}
}

func TestSqrtCFPrefixStream2_Snapshot_AfterStartCarriesApproximation(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	s := NewSqrtCFPrefixStream2(Sqrt2CF(), 2, p)

	_, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit; err=%v", s.Err())
	}

	snap := s.Snapshot()
	if !snap.Started {
		t.Fatalf("expected Started=true")
	}
	if snap.PrefixTerms != 2 {
		t.Fatalf("got PrefixTerms=%d want 2", snap.PrefixTerms)
	}
	if snap.Approx == nil {
		t.Fatalf("expected non-nil Approx after start")
	}

	want, err := SqrtApproxFromSourceRangeSeed2(Sqrt2CF(), 2, p)
	if err != nil {
		t.Fatalf("SqrtApproxFromSourceRangeSeed2 failed: %v", err)
	}
	if snap.Approx.Cmp(want) != 0 {
		t.Fatalf("got Approx=%v want %v", *snap.Approx, want)
	}
}

func TestSqrtCFPrefixStream2_Snapshot_StartFailureLeavesApproxNil(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	s := NewSqrtCFPrefixStream2(Sqrt2CF(), 0, p)

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected failure")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}

	snap := s.Snapshot()
	if !snap.Started {
		t.Fatalf("expected Started=true after attempted start")
	}
	if snap.Approx != nil {
		t.Fatalf("expected Approx=nil on failed start")
	}
}
