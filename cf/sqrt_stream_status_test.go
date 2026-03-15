// sqrt_stream_status_test.go v1
package cf

import "testing"

func TestSqrtCFPrefixStream2_Status_UnstartedThenBoundedCollapse(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}
	s := NewSqrtCFPrefixStream2(Sqrt2CF(), 2, p)

	if got := s.Snapshot().Status; got != SqrtStreamStatusUnstarted {
		t.Fatalf("before start got status %q want %q", got, SqrtStreamStatusUnstarted)
	}

	_, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit; err=%v", s.Err())
	}

	if got := s.Snapshot().Status; got != SqrtStreamStatusBoundedCollapse {
		t.Fatalf("after start got status %q want %q", got, SqrtStreamStatusBoundedCollapse)
	}
}

func TestSqrtGCFPrefixStream2_Status_FailedOnBadPrefixTerms(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}
	s := NewSqrtGCFPrefixStream2(NewSliceGCF([2]int64{3, 2}), 0, p)

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected failure")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}

	if got := s.Snapshot().Status; got != SqrtStreamStatusFailed {
		t.Fatalf("after failure got status %q want %q", got, SqrtStreamStatusFailed)
	}
}

func TestSqrtGCFExactTailStream2_Status_UnstartedThenBoundedCollapse(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}
	s := NewSqrtGCFExactTailStreamWithTail2(
		NewSliceGCF([2]int64{3, 2}),
		mustRat(11, 1),
		8,
		p,
	)

	if got := s.Snapshot().Status; got != SqrtStreamStatusUnstarted {
		t.Fatalf("before start got status %q want %q", got, SqrtStreamStatusUnstarted)
	}

	_, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit; err=%v", s.Err())
	}

	if got := s.Snapshot().Status; got != SqrtStreamStatusBoundedCollapse {
		t.Fatalf("after start got status %q want %q", got, SqrtStreamStatusBoundedCollapse)
	}
}
