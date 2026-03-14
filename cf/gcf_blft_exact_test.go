// gcf_blft_exact_test.go v1
package cf

import (
	"strings"
	"testing"
)

// func composeGCFIntoBLFTX(base BLFT, src GCFSource) (BLFT, error) {
// 	cur := base
// 	for {
// 		p, q, ok := src.NextPQ()
// 		if !ok {
// 			return cur, nil
// 		}
// 		var err error
// 		cur, err = cur.IngestGCFX(p, q)
// 		if err != nil {
// 			return BLFT{}, err
// 		}
// 	}
// }

func TestApplyComposedGCFXBLFTToTailsExact_ProjectionXMatchesExactEvaluation(t *testing.T) {
	// B(x,y)=x
	base := NewBLFT(0, 1, 0, 0, 0, 0, 0, 1)
	xSrc := NewSliceGCF(
		[2]int64{3, 2},
		[2]int64{5, 7},
	)
	xTail := mustRat(11, 1)
	y := mustRat(7, 1)

	got, ingested, err := ApplyComposedGCFXBLFTToTailsExact(base, xSrc, xTail, y, 8)
	if err != nil {
		t.Fatalf("ApplyComposedGCFXBLFTToTailsExact failed: %v", err)
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

	want, err := base.ApplyRat(wantX, y)
	if err != nil {
		t.Fatalf("base ApplyRat failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestApplyComposedGCFXBLFTToTailsExact_GeneralMatchesExactEvaluation(t *testing.T) {
	base := NewBLFT(2, 3, 5, 7, 11, 13, 17, 19)
	xSrc := NewSliceGCF(
		[2]int64{5, 7},
	)
	xTail := mustRat(13, 1)
	y := mustRat(3, 2)

	got, ingested, err := ApplyComposedGCFXBLFTToTailsExact(base, xSrc, xTail, y, 8)
	if err != nil {
		t.Fatalf("ApplyComposedGCFXBLFTToTailsExact failed: %v", err)
	}
	if ingested != 1 {
		t.Fatalf("got ingested=%d want 1", ingested)
	}

	x, _, err := EvalGCFWithTailExact(NewSliceGCF([2]int64{5, 7}), xTail, 8)
	if err != nil {
		t.Fatalf("EvalGCFWithTailExact failed: %v", err)
	}
	want, err := base.ApplyRat(x, y)
	if err != nil {
		t.Fatalf("base ApplyRat failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestApplyComposedGCFXBLFTToTailsExact_RequiresExhaustionWithinBound(t *testing.T) {
	base := NewBLFT(0, 1, 0, 0, 0, 0, 0, 1)

	_, ingested, err := ApplyComposedGCFXBLFTToTailsExact(
		base,
		NewUnitPArithmeticQGCFSource(1, 1),
		mustRat(1, 1),
		mustRat(2, 1),
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

func TestApplyComposedGCFXBLFTToTailsExact_InvalidTermPropagatesError(t *testing.T) {
	base := NewBLFT(0, 1, 0, 0, 0, 0, 0, 1)

	_, ingested, err := ApplyComposedGCFXBLFTToTailsExact(
		base,
		NewSliceGCF([2]int64{7, 0}),
		mustRat(1, 1),
		mustRat(2, 1),
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

// gcf_blft_exact_test.go v1
