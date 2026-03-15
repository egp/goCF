// sqrt_certified_gcf_stream_test.go v1
package cf

import "testing"

func TestSqrtCertifiedGCFPrefixStream_RejectsBadBound(t *testing.T) {
	_, err := NewSqrtCertifiedGCFPrefixStream(AdaptCFToGCF(Sqrt2CF()), 0)
	if err == nil {
		t.Fatalf("expected constructor error")
	}
}

func TestSqrtCertifiedGCFPrefixStream_EmptySourceFails(t *testing.T) {
	s, err := NewSqrtCertifiedGCFPrefixStream(NewSliceGCF(), 4)
	if err != nil {
		t.Fatalf("constructor failed: %v", err)
	}

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected failure, not a digit")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}
	if snap := s.Snapshot(); snap.Status != SqrtStreamStatusFailed {
		t.Fatalf("got status %q want %q", snap.Status, SqrtStreamStatusFailed)
	}
}

func TestSqrtCertifiedGCFPrefixStream_Sqrt2InputYieldsFourthRootOf2(t *testing.T) {
	s, err := NewSqrtCertifiedGCFPrefixStream(AdaptCFToGCF(Sqrt2CF()), 12)
	if err != nil {
		t.Fatalf("constructor failed: %v", err)
	}

	got := collectTerms(s, 5)
	if err := s.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}

	// Input is sqrt(2), so the stream computes sqrt(sqrt(2)) = 2^(1/4).
	wantPrefix := []int64{1, 5, 3, 1, 1}
	if len(got) < len(wantPrefix) {
		t.Fatalf("len(got)=%d want at least %d got=%v", len(got), len(wantPrefix), got)
	}
	for i := range wantPrefix {
		if got[i] != wantPrefix[i] {
			t.Fatalf("got[%d]=%d want=%d full=%v", i, got[i], wantPrefix[i], got)
		}
	}

	snap := s.Snapshot()
	if snap.Status != SqrtStreamStatusCertifiedProgressive {
		t.Fatalf("got status %q want %q", snap.Status, SqrtStreamStatusCertifiedProgressive)
	}
	if snap.PrefixTerms > 12 {
		t.Fatalf("used too many terms: got %d want <= 12", snap.PrefixTerms)
	}
	if snap.GCFInputApprox == nil {
		t.Fatalf("expected non-nil GCFInputApprox")
	}
}

func TestSqrtCertifiedGCFStream_PublicConstructor(t *testing.T) {
	s, err := SqrtCertifiedGCFStream(AdaptCFToGCF(Sqrt2CF()), 8)
	if err != nil {
		t.Fatalf("SqrtCertifiedGCFStream failed: %v", err)
	}

	snap := s.Snapshot()
	if snap.Started {
		t.Fatalf("expected Started=false before first Next")
	}
	if snap.PrefixTerms != 0 {
		t.Fatalf("got PrefixTerms=%d want 0 before start", snap.PrefixTerms)
	}
}

// sqrt_certified_gcf_stream_test.go v1
