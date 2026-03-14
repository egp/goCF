// sqrt_gcf_prefix_stream2_test.go v1
package cf

import "testing"

func TestNewSqrtGCFPrefixStream2_FiniteSource_ExactRangeFastPath(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(0, 1),
	}

	// Under the current GCFApprox/GCFBounder finite-prefix semantics,
	// a finite exhausted source with single term (3,1) yields convergent 3.
	// So this stream approximates sqrt(3), whose CF begins [1; 1,2,1,2,...].
	s := NewSqrtGCFPrefixStream2(
		NewSliceGCF([2]int64{3, 1}),
		8,
		p,
	)

	got := collectTerms(s, 8)
	if err := s.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}

	want := []int64{1, 1, 2, 1, 2, 1, 2, 1}
	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d full=%v", i, got[i], want[i], got)
		}
	}
}

func TestNewSqrtGCFPrefixStream2_BoundedInfiniteSource_MatchesCanonicalPath(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	srcA := AdaptCFToGCF(Sqrt2CF())
	srcB := AdaptCFToGCF(Sqrt2CF())

	wantApprox, err := SqrtApproxFromGCFSourceRangeSeed2(srcA, 2, p)
	if err != nil {
		t.Fatalf("SqrtApproxFromGCFSourceRangeSeed2 failed: %v", err)
	}
	want := collectTerms(NewRationalCF(wantApprox), 16)

	s := NewSqrtGCFPrefixStream2(srcB, 2, p)
	got := collectTerms(s, 16)
	if err := s.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v wantFull=%v", len(got), len(want), got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d gotFull=%v wantFull=%v", i, got[i], want[i], got, want)
		}
	}
}

func TestNewSqrtGCFPrefixStream2_RejectsBadPrefixTerms(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	s := NewSqrtGCFPrefixStream2(
		NewSliceGCF([2]int64{3, 2}),
		0,
		p,
	)

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected failure, not a digit")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}
}

func TestSqrtGCFPrefixStream2_ExhaustedStreamStaysExhausted(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(0, 1),
	}

	s := NewSqrtGCFPrefixStream2(
		NewSliceGCF([2]int64{3, 1}),
		8,
		p,
	)

	for {
		_, ok := s.Next()
		if !ok {
			break
		}
	}

	if _, ok := s.Next(); ok {
		t.Fatalf("expected exhausted stream to stay exhausted")
	}
	if _, ok := s.Next(); ok {
		t.Fatalf("expected exhausted stream to stay exhausted on repeated calls")
	}
}

// sqrt_gcf_prefix_stream2.go v1
