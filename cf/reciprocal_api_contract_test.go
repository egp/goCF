// reciprocal_api_contract_test.go v1
package cf

import "testing"

func TestReciprocalGCFPrefixStream_APIContract(t *testing.T) {
	s, err := ReciprocalGCFPrefixStream(AdaptCFToGCF(Sqrt2CF()), 2)
	if err != nil {
		t.Fatalf("constructor failed: %v", err)
	}

	snap := s.Snapshot()
	if snap.Started {
		t.Fatalf("expected Started=false")
	}
	if snap.PrefixTerms != 2 {
		t.Fatalf("got PrefixTerms=%d want 2", snap.PrefixTerms)
	}
	if snap.MaxIngestTerms != 0 {
		t.Fatalf("got MaxIngestTerms=%d want 0 for prefix stream", snap.MaxIngestTerms)
	}
	if snap.ConsumedTerms != 0 {
		t.Fatalf("got ConsumedTerms=%d want 0 before start", snap.ConsumedTerms)
	}
}

func TestReciprocalGCFExactTailStream_APIContract(t *testing.T) {
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
	if snap.MaxIngestTerms != 8 {
		t.Fatalf("got MaxIngestTerms=%d want 8", snap.MaxIngestTerms)
	}
	if snap.PrefixTerms != 0 {
		t.Fatalf("got PrefixTerms=%d want 0 for exact-tail stream", snap.PrefixTerms)
	}
	if snap.ConsumedTerms != 0 {
		t.Fatalf("got ConsumedTerms=%d want 0 before start", snap.ConsumedTerms)
	}
}

func TestReciprocalAndSqrtSnapshotContract_CommonBeforeStartShape(t *testing.T) {
	recip, err := ReciprocalGCFPrefixStream(AdaptCFToGCF(Sqrt2CF()), 2)
	if err != nil {
		t.Fatalf("reciprocal constructor failed: %v", err)
	}
	rs := recip.Snapshot()
	if rs.Started {
		t.Fatalf("expected reciprocal Started=false")
	}
	if rs.Approx != nil {
		t.Fatalf("expected reciprocal Approx=nil before start")
	}

	sqrt, err := NewSqrtCertifiedGCFPrefixStream(AdaptCFToGCF(Sqrt2CF()), 8)
	if err != nil {
		t.Fatalf("sqrt constructor failed: %v", err)
	}
	ss := sqrt.Snapshot()
	if ss.Started {
		t.Fatalf("expected sqrt Started=false")
	}
	if ss.Approx != nil {
		t.Fatalf("expected sqrt Approx=nil before start")
	}
	if ss.Status != SqrtStreamStatusUnstarted {
		t.Fatalf("got sqrt Status=%q want %q", ss.Status, SqrtStreamStatusUnstarted)
	}
}
