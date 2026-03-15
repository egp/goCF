// reciprocal_stream_gcf_prefix_snapshot_test.go v1
package cf

import "testing"

func TestReciprocalGCFPrefixStream2_Snapshot_BeforeStart(t *testing.T) {
	s := NewReciprocalGCFPrefixStream2(AdaptCFToGCF(Sqrt2CF()), 2)
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
	if snap.GCFInputApprox != nil {
		t.Fatalf("expected GCFInputApprox=nil before start")
	}
}

func TestReciprocalGCFPrefixStream2_Snapshot_AfterStartCarriesInputApprox(t *testing.T) {
	s := NewReciprocalGCFPrefixStream2(AdaptCFToGCF(Sqrt2CF()), 2)

	_, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit; err=%v", s.Err())
	}

	snap := s.Snapshot()
	if !snap.Started {
		t.Fatalf("expected Started=true")
	}
	if snap.Approx == nil {
		t.Fatalf("expected non-nil Approx after start")
	}
	if snap.GCFInputApprox == nil {
		t.Fatalf("expected non-nil GCFInputApprox after start")
	}

	want, err := GCFApproxFromPrefix(AdaptCFToGCF(Sqrt2CF()), 2)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}
	if snap.GCFInputApprox.Convergent.Cmp(want.Convergent) != 0 {
		t.Fatalf("got convergent %v want %v", snap.GCFInputApprox.Convergent, want.Convergent)
	}
}

// reciprocal_stream_gcf_prefix_snapshot_test.go v1
