// sqrt_stream_api_test.go v1
package cf

import "testing"

func TestSqrtStream_FinitePerfectSquare(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(0, 1),
	}

	cf, err := SqrtStream(NewSliceCF(4), 8, p)
	if err != nil {
		t.Fatalf("SqrtStream failed: %v", err)
	}

	got := collectTerms(cf, 8)
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

func TestSqrtStream_Sqrt2PrefixMatchesBoundedCFPrefixStream(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	wantStream := NewSqrtCFPrefixStream2(Sqrt2CF(), 2, sqrtPolicy2FromOld(p))
	want := collectTerms(wantStream, 16)
	if err := wantStream.Err(); err != nil {
		t.Fatalf("want stream error: %v", err)
	}

	cf, err := SqrtStream(Sqrt2CF(), 2, p)
	if err != nil {
		t.Fatalf("SqrtStream failed: %v", err)
	}
	got := collectTerms(cf, 16)

	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v wantFull=%v", len(got), len(want), got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d gotFull=%v wantFull=%v", i, got[i], want[i], got, want)
		}
	}
}

func TestSqrtStream_RejectsBadPrefixTerms(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	_, err := SqrtStream(Sqrt2CF(), 0, p)
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
}

// sqrt_stream_api_test.go v1
