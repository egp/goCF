// sqrt_stream_status_exact_input_test.go v1
package cf

import "testing"

func TestSqrtCFPrefixStream2_Status_ExactInputForFiniteSquare(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(0, 1),
	}
	s := NewSqrtCFPrefixStream2(NewSliceCF(4), 8, p)

	_, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit; err=%v", s.Err())
	}

	if got := s.Snapshot().Status; got != SqrtStreamStatusExactInput {
		t.Fatalf("got status %q want %q", got, SqrtStreamStatusExactInput)
	}
}

func TestSqrtGCFPrefixStream2_Status_ExactInputForFiniteSquare(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(0, 1),
	}
	s := NewSqrtGCFPrefixStream2(NewSliceGCF([2]int64{3, 1}), 8, p)

	_, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit; err=%v", s.Err())
	}

	// Under the current bounded GCF-prefix semantics, this finite exhausted
	// source yields an exact point input approximation, so the stream reports
	// exact_input.
	if got := s.Snapshot().Status; got != SqrtStreamStatusExactInput {
		t.Fatalf("got status %q want %q", got, SqrtStreamStatusExactInput)
	}
}

func TestSqrtGCFExactTailStream2_Status_ExactInputForExactSquare(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(0, 1),
	}
	s := NewSqrtGCFExactTailStreamWithTail2(
		NewSliceGCF(),
		mustRat(9, 16),
		8,
		p,
	)

	_, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit; err=%v", s.Err())
	}

	if got := s.Snapshot().Status; got != SqrtStreamStatusExactInput {
		t.Fatalf("got status %q want %q", got, SqrtStreamStatusExactInput)
	}
}
