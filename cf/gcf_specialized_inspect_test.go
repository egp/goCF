// gcf_specialized_inspect_test.go v2
package cf

import "testing"

func TestSpecializedInspectGCFSource_RejectsNegativeDigits(t *testing.T) {
	_, err := specializedInspectGCFSource(
		1,
		-1,
		func(prefixTerms int) (GCFApprox, error) {
			return GCFApprox{
				Convergent:  mustRat(1, 1),
				Range:       nil,
				PrefixTerms: 1,
			}, nil
		},
	)
	if err == nil {
		t.Fatalf("expected error for negative digits")
	}
}

func TestInspectLambertPiOver4Prefix(t *testing.T) {
	got, err := InspectLambertPiOver4Prefix(2, 8)
	if err != nil {
		t.Fatalf("InspectLambertPiOver4Prefix failed: %v", err)
	}

	wantConv := mustRat(1, 1)
	if got.Approx.Convergent.Cmp(wantConv) != 0 {
		t.Fatalf("got convergent %v want %v", got.Approx.Convergent, wantConv)
	}
	if got.Approx.Range == nil {
		t.Fatalf("expected non-nil range")
	}

	wantLo := mustRat(1, 2)
	wantHi := mustRat(1, 1)
	if got.Approx.Range.Lo.Cmp(wantLo) != 0 || got.Approx.Range.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got range [%v,%v] want [%v,%v]", got.Approx.Range.Lo, got.Approx.Range.Hi, wantLo, wantHi)
	}

	wantTerms := []int64{1}
	if len(got.Terms) != len(wantTerms) {
		t.Fatalf("len(got)=%d want=%d got=%v", len(got.Terms), len(wantTerms), got.Terms)
	}
	for i := range wantTerms {
		if got.Terms[i] != wantTerms[i] {
			t.Fatalf("got[%d]=%d want=%d full=%v", i, got.Terms[i], wantTerms[i], got.Terms)
		}
	}
}

func TestInspectBrouncker4OverPiPrefix(t *testing.T) {
	got, err := InspectBrouncker4OverPiPrefix(2, 8)
	if err != nil {
		t.Fatalf("InspectBrouncker4OverPiPrefix failed: %v", err)
	}

	wantConv := mustRat(3, 2)
	if got.Approx.Convergent.Cmp(wantConv) != 0 {
		t.Fatalf("got convergent %v want %v", got.Approx.Convergent, wantConv)
	}
	if got.Approx.Range == nil {
		t.Fatalf("expected non-nil range")
	}

	wantLo := mustRat(12, 11)
	wantHi := mustRat(3, 2)
	if got.Approx.Range.Lo.Cmp(wantLo) != 0 || got.Approx.Range.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got range [%v,%v] want [%v,%v]", got.Approx.Range.Lo, got.Approx.Range.Hi, wantLo, wantHi)
	}

	wantTerms := []int64{1, 2}
	if len(got.Terms) != len(wantTerms) {
		t.Fatalf("len(got)=%d want=%d got=%v", len(got.Terms), len(wantTerms), got.Terms)
	}
	for i := range wantTerms {
		if got.Terms[i] != wantTerms[i] {
			t.Fatalf("got[%d]=%d want=%d full=%v", i, got.Terms[i], wantTerms[i], got.Terms)
		}
	}
}

// gcf_specialized_inspect_test.go v2
