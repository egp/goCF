// unary_stream_contract_test.go v1
package cf

import "testing"

func TestUnaryStreamContract_ReciprocalPrefix_BeforeStart(t *testing.T) {
	s, err := ReciprocalGCFPrefixStream(AdaptCFToGCF(Sqrt2CF()), 2)
	if err != nil {
		t.Fatalf("constructor failed: %v", err)
	}

	snap := s.Snapshot()
	if snap.Started {
		t.Fatalf("expected Started=false")
	}
	if snap.Approx != nil {
		t.Fatalf("expected Approx=nil before start")
	}
	if snap.GCFInputApprox != nil {
		t.Fatalf("expected GCFInputApprox=nil before start")
	}
	if snap.PrefixTerms != 2 {
		t.Fatalf("got PrefixTerms=%d want 2", snap.PrefixTerms)
	}
	if snap.ConsumedTerms != 0 {
		t.Fatalf("got ConsumedTerms=%d want 0", snap.ConsumedTerms)
	}
	if s.Err() != nil {
		t.Fatalf("expected nil Err before start, got %v", s.Err())
	}
}

func TestUnaryStreamContract_ReciprocalPrefix_AfterStart(t *testing.T) {
	s, err := ReciprocalGCFPrefixStream(AdaptCFToGCF(Sqrt2CF()), 2)
	if err != nil {
		t.Fatalf("constructor failed: %v", err)
	}

	_, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit; err=%v", s.Err())
	}

	snap := s.Snapshot()
	if !snap.Started {
		t.Fatalf("expected Started=true")
	}
	if snap.Approx == nil {
		t.Fatalf("expected non-nil Approx")
	}
	if snap.GCFInputApprox == nil {
		t.Fatalf("expected non-nil GCFInputApprox")
	}
	if snap.ConsumedTerms != 2 {
		t.Fatalf("got ConsumedTerms=%d want 2", snap.ConsumedTerms)
	}
}

func TestUnaryStreamContract_ReciprocalExactTail_BeforeStart(t *testing.T) {
	s, err := ReciprocalGCFExactTailStreamWithTail(
		NewSliceGCF([2]int64{3, 2}, [2]int64{5, 7}),
		mustRat(11, 1),
		8,
	)
	if err != nil {
		t.Fatalf("constructor failed: %v", err)
	}

	snap := s.Snapshot()
	if snap.Started {
		t.Fatalf("expected Started=false")
	}
	if snap.Approx != nil {
		t.Fatalf("expected Approx=nil before start")
	}
	if snap.MaxIngestTerms != 8 {
		t.Fatalf("got MaxIngestTerms=%d want 8", snap.MaxIngestTerms)
	}
	if snap.ConsumedTerms != 0 {
		t.Fatalf("got ConsumedTerms=%d want 0", snap.ConsumedTerms)
	}
	if s.Err() != nil {
		t.Fatalf("expected nil Err before start, got %v", s.Err())
	}
}

func TestUnaryStreamContract_ReciprocalExactTail_AfterStart(t *testing.T) {
	s, err := ReciprocalGCFExactTailStreamWithTail(
		NewSliceGCF([2]int64{3, 2}, [2]int64{5, 7}),
		mustRat(11, 1),
		8,
	)
	if err != nil {
		t.Fatalf("constructor failed: %v", err)
	}

	_, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit; err=%v", s.Err())
	}

	snap := s.Snapshot()
	if !snap.Started {
		t.Fatalf("expected Started=true")
	}
	if snap.Approx == nil {
		t.Fatalf("expected non-nil Approx")
	}
	if snap.ConsumedTerms != 2 {
		t.Fatalf("got ConsumedTerms=%d want 2", snap.ConsumedTerms)
	}
}

func TestUnaryStreamContract_SqrtCertifiedPrefix_BeforeStart(t *testing.T) {
	s, err := NewSqrtCertifiedGCFPrefixStream(AdaptCFToGCF(Sqrt2CF()), 8)
	if err != nil {
		t.Fatalf("constructor failed: %v", err)
	}

	snap := s.Snapshot()
	if snap.Started {
		t.Fatalf("expected Started=false")
	}
	if snap.Status != SqrtStreamStatusUnstarted {
		t.Fatalf("got Status=%q want %q", snap.Status, SqrtStreamStatusUnstarted)
	}
	if snap.PrefixTerms != 0 {
		t.Fatalf("got PrefixTerms=%d want 0", snap.PrefixTerms)
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

func TestUnaryStreamContract_SqrtCertifiedPrefix_AfterStart(t *testing.T) {
	s, err := NewSqrtCertifiedGCFPrefixStream(AdaptCFToGCF(Sqrt2CF()), 8)
	if err != nil {
		t.Fatalf("constructor failed: %v", err)
	}

	_, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit; err=%v", s.Err())
	}

	snap := s.Snapshot()
	if !snap.Started {
		t.Fatalf("expected Started=true")
	}
	if snap.Status != SqrtStreamStatusCertifiedProgressive {
		t.Fatalf("got Status=%q want %q", snap.Status, SqrtStreamStatusCertifiedProgressive)
	}
	if snap.PrefixTerms == 0 {
		t.Fatalf("expected PrefixTerms>0 after start")
	}
	if snap.GCFInputApprox == nil {
		t.Fatalf("expected non-nil GCFInputApprox after start")
	}
}

func TestUnaryStreamContract_ExhaustionIsStable(t *testing.T) {
	recip, err := ReciprocalGCFPrefixStream(NewSliceGCF([2]int64{3, 2}, [2]int64{5, 7}), 10)
	if err != nil {
		t.Fatalf("constructor failed: %v", err)
	}
	_ = collectTerms(recip, 32)
	if recip.Err() != nil {
		t.Fatalf("unexpected reciprocal err: %v", recip.Err())
	}
	if _, ok := recip.Next(); ok {
		t.Fatalf("expected reciprocal exhaustion")
	}
	if _, ok := recip.Next(); ok {
		t.Fatalf("expected reciprocal repeated exhaustion")
	}

	sqrt, err := NewSqrtCertifiedGCFPrefixStream(NewSliceGCF([2]int64{4, 1}), 4)
	if err != nil {
		t.Fatalf("constructor failed: %v", err)
	}
	_ = collectTerms(sqrt, 32)
	if sqrt.Err() != nil {
		t.Fatalf("unexpected sqrt err: %v", sqrt.Err())
	}
	if _, ok := sqrt.Next(); ok {
		t.Fatalf("expected sqrt exhaustion")
	}
	if _, ok := sqrt.Next(); ok {
		t.Fatalf("expected sqrt repeated exhaustion")
	}
}

// unary_stream_contract_test.go v1
