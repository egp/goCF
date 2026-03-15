// sqrt_stream_snapshot_input_test.go v1
package cf

import "testing"

func TestSqrtCFPrefixStream2_Snapshot_CarriesCFInputApprox(t *testing.T) {
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
	if snap.CFInputApprox == nil {
		t.Fatalf("expected non-nil CFInputApprox")
	}
	if snap.GCFInputApprox != nil {
		t.Fatalf("expected nil GCFInputApprox for CF stream")
	}

	want, err := CFApproxFromPrefix(Sqrt2CF(), 2)
	if err != nil {
		t.Fatalf("CFApproxFromPrefix failed: %v", err)
	}
	if snap.CFInputApprox.Convergent.Cmp(want.Convergent) != 0 {
		t.Fatalf("got convergent %v want %v", snap.CFInputApprox.Convergent, want.Convergent)
	}
}

func TestSqrtGCFPrefixStream2_Snapshot_CarriesGCFInputApprox(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	s := NewSqrtGCFPrefixStream2(AdaptCFToGCF(Sqrt2CF()), 2, p)

	_, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit; err=%v", s.Err())
	}

	snap := s.Snapshot()
	if snap.GCFInputApprox == nil {
		t.Fatalf("expected non-nil GCFInputApprox")
	}
	if snap.CFInputApprox != nil {
		t.Fatalf("expected nil CFInputApprox for GCF stream")
	}

	want, err := GCFApproxFromPrefix(AdaptCFToGCF(Sqrt2CF()), 2)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}
	if snap.GCFInputApprox.Convergent.Cmp(want.Convergent) != 0 {
		t.Fatalf("got convergent %v want %v", snap.GCFInputApprox.Convergent, want.Convergent)
	}
}
