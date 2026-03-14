// sqrt_stream_cf2_test.go v1
package cf

import "testing"

func TestNewSqrtCFPrefixStream2_FiniteSourceExact(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(0, 1),
	}

	s := NewSqrtCFPrefixStream2(NewSliceCF(4), 8, p)

	got := collectTerms(s, 8)
	if err := s.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}

	want := []int64{2}
	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d full=%v", i, got[i], want[i], got)
		}
	}
}

func TestNewSqrtCFPrefixStream2_BoundedInfiniteSourceMatchesCanonicalPath(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	srcA := Sqrt2CF()
	srcB := Sqrt2CF()

	wantCF, err := SqrtApproxCFFromSourceRangeSeed2(srcA, 2, p)
	if err != nil {
		t.Fatalf("SqrtApproxCFFromSourceRangeSeed2 failed: %v", err)
	}
	want := collectTerms(wantCF, 16)

	s := NewSqrtCFPrefixStream2(srcB, 2, p)
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

func TestNewSqrtCFPrefixStream2_RejectsBadPrefixTerms(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	s := NewSqrtCFPrefixStream2(Sqrt2CF(), 0, p)

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected failure, not a digit")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}
}

func TestSqrtCFPrefixStream2_ExhaustedStreamStaysExhausted(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(0, 1),
	}

	s := NewSqrtCFPrefixStream2(NewSliceCF(4), 8, p)

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

// sqrt_stream_cf2_test.go v1
