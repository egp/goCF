// gcf_diag_exact_test.go v1
package cf

import (
	"strings"
	"testing"
)

// func composeGCFIntoDiagBLFT(base DiagBLFT, src GCFSource) (DiagBLFT, error) {
// 	cur := base
// 	for {
// 		p, q, ok := src.NextPQ()
// 		if !ok {
// 			return cur, nil
// 		}
// 		var err error
// 		cur, err = cur.IngestGCF(p, q)
// 		if err != nil {
// 			return DiagBLFT{}, err
// 		}
// 	}
// }

func TestApplyComposedGCFDiagBLFTToTailExact_IdentityLikeProjectionMatchesExactEvaluation(t *testing.T) {
	// D(x)=x
	base := NewDiagBLFT(
		mustBig(0), mustBig(1), mustBig(0),
		mustBig(0), mustBig(0), mustBig(1),
	)

	src := NewSliceGCF(
		[2]int64{3, 2},
		[2]int64{5, 7},
	)
	xTail := mustRat(11, 1)

	got, ingested, err := ApplyComposedGCFDiagBLFTToTailExact(base, src, xTail, 8)
	if err != nil {
		t.Fatalf("ApplyComposedGCFDiagBLFTToTailExact failed: %v", err)
	}
	if ingested != 2 {
		t.Fatalf("got ingested=%d want 2", ingested)
	}

	wantX, _, err := EvalGCFWithTailExact(
		NewSliceGCF(
			[2]int64{3, 2},
			[2]int64{5, 7},
		),
		xTail,
		8,
	)
	if err != nil {
		t.Fatalf("EvalGCFWithTailExact failed: %v", err)
	}

	want, err := base.ApplyRat(wantX)
	if err != nil {
		t.Fatalf("base ApplyRat failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestApplyComposedGCFDiagBLFTToTailExact_GeneralMatchesExactEvaluation(t *testing.T) {
	base := NewDiagBLFT(
		mustBig(2), mustBig(3), mustBig(5),
		mustBig(7), mustBig(11), mustBig(13),
	)

	src := NewSliceGCF([2]int64{5, 7})
	xTail := mustRat(13, 1)

	got, ingested, err := ApplyComposedGCFDiagBLFTToTailExact(base, src, xTail, 8)
	if err != nil {
		t.Fatalf("ApplyComposedGCFDiagBLFTToTailExact failed: %v", err)
	}
	if ingested != 1 {
		t.Fatalf("got ingested=%d want 1", ingested)
	}

	x, _, err := EvalGCFWithTailExact(NewSliceGCF([2]int64{5, 7}), xTail, 8)
	if err != nil {
		t.Fatalf("EvalGCFWithTailExact failed: %v", err)
	}
	want, err := base.ApplyRat(x)
	if err != nil {
		t.Fatalf("base ApplyRat failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestApplyComposedGCFDiagBLFTToTailExact_RequiresExhaustionWithinBound(t *testing.T) {
	base := NewDiagBLFT(
		mustBig(0), mustBig(1), mustBig(0),
		mustBig(0), mustBig(0), mustBig(1),
	)

	_, ingested, err := ApplyComposedGCFDiagBLFTToTailExact(
		base,
		NewUnitPArithmeticQGCFSource(1, 1),
		mustRat(1, 1),
		3,
	)
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "exceeded MaxIngestTerms=3") {
		t.Fatalf("unexpected error: %v", err)
	}
	if ingested != 3 {
		t.Fatalf("got ingested=%d want 3", ingested)
	}
}

func TestApplyComposedGCFDiagBLFTToTailExact_InvalidTermPropagatesError(t *testing.T) {
	base := NewDiagBLFT(
		mustBig(0), mustBig(1), mustBig(0),
		mustBig(0), mustBig(0), mustBig(1),
	)

	_, ingested, err := ApplyComposedGCFDiagBLFTToTailExact(
		base,
		NewSliceGCF([2]int64{7, 0}),
		mustRat(1, 1),
		8,
	)
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "require q>0") {
		t.Fatalf("unexpected error: %v", err)
	}
	if ingested != 0 {
		t.Fatalf("got ingested=%d want 0", ingested)
	}
}

// gcf_diag_exact_test.go v1
