// sqrt_stream_semantics_test.go v1
package cf

import "testing"

func TestSqrtStream_CurrentSemantics_MatchesBoundedCFPrefixStream(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	wantStream := NewSqrtCFPrefixStream2(
		Sqrt2CF(),
		2,
		sqrtPolicy2FromOld(p),
	)
	want := collectTerms(wantStream, 16)
	if err := wantStream.Err(); err != nil {
		t.Fatalf("want stream error: %v", err)
	}

	gotStream, err := SqrtStream(Sqrt2CF(), 2, p)
	if err != nil {
		t.Fatalf("SqrtStream failed: %v", err)
	}
	got := collectTerms(gotStream, 16)

	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v wantFull=%v", len(got), len(want), got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d gotFull=%v wantFull=%v", i, got[i], want[i], got, want)
		}
	}
}
