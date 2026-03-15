// unary_stream_taxonomy_test.go v1
package cf

import "testing"

func TestUnaryStreamTaxonomy_ReciprocalPrefix(t *testing.T) {
	s, err := ReciprocalGCFPrefixStream(AdaptCFToGCF(Sqrt2CF()), 2)
	if err != nil {
		t.Fatalf("constructor failed: %v", err)
	}

	got := classifyReciprocalSnapshot(s.Snapshot())
	want := unaryStreamClass{
		Operator: unaryOperatorReciprocal,
		Input:    unaryInputGCFPrefix,
		Progress: unaryProgressExactCollapse,
	}
	if got != want {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestUnaryStreamTaxonomy_ReciprocalExactTail(t *testing.T) {
	s, err := ReciprocalGCFExactTailStreamWithTail(
		NewSliceGCF([2]int64{3, 2}, [2]int64{5, 7}),
		mustRat(11, 1),
		8,
	)
	if err != nil {
		t.Fatalf("constructor failed: %v", err)
	}

	got := classifyReciprocalSnapshot(s.Snapshot())
	want := unaryStreamClass{
		Operator: unaryOperatorReciprocal,
		Input:    unaryInputGCFExact,
		Progress: unaryProgressExactCollapse,
	}
	if got != want {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestUnaryStreamTaxonomy_SqrtCertifiedPrefix_BeforeStart(t *testing.T) {
	s, err := NewSqrtCertifiedGCFPrefixStream(AdaptCFToGCF(Sqrt2CF()), 8)
	if err != nil {
		t.Fatalf("constructor failed: %v", err)
	}

	got := classifySqrtSnapshot(s.Snapshot())
	want := unaryStreamClass{
		Operator: unaryOperatorSqrt,
		Input:    unaryInputUnknown,
		Progress: unaryProgressUnknown,
	}
	if got != want {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestUnaryStreamTaxonomy_SqrtCertifiedPrefix_AfterStart(t *testing.T) {
	s, err := NewSqrtCertifiedGCFPrefixStream(AdaptCFToGCF(Sqrt2CF()), 8)
	if err != nil {
		t.Fatalf("constructor failed: %v", err)
	}

	_, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit; err=%v", s.Err())
	}

	got := classifySqrtSnapshot(s.Snapshot())
	want := unaryStreamClass{
		Operator: unaryOperatorSqrt,
		Input:    unaryInputGCFPrefix,
		Progress: unaryProgressProgressiveCertified,
	}
	if got != want {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

// unary_stream_taxonomy_test.go v1
