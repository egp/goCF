// gcf_blft_exact_y_test.go v1
package cf

import (
	"strings"
	"testing"
)

// func composeGCFIntoBLFTY(base BLFT, src GCFSource) (BLFT, error) {
// 	cur := base
// 	for {
// 		p, q, ok := src.NextPQ()
// 		if !ok {
// 			return cur, nil
// 		}
// 		var err error
// 		cur, err = cur.IngestGCFY(p, q)
// 		if err != nil {
// 			return BLFT{}, err
// 		}
// 	}
// }

func TestApplyComposedGCFYBLFTToTailsExact_ProjectionYMatchesExactEvaluation(t *testing.T) {
	// B(x,y)=y
	base := NewBLFT(0, 0, 1, 0, 0, 0, 0, 1)
	x := mustRat(7, 1)
	ySrc := NewSliceGCF(
		[2]int64{3, 2},
		[2]int64{5, 7},
	)
	yTail := mustRat(11, 1)

	got, ingested, err := ApplyComposedGCFYBLFTToTailsExact(base, x, ySrc, yTail, 8)
	if err != nil {
		t.Fatalf("ApplyComposedGCFYBLFTToTailsExact failed: %v", err)
	}
	if ingested != 2 {
		t.Fatalf("got ingested=%d want 2", ingested)
	}

	wantY, _, err := EvalGCFWithTailExact(
		NewSliceGCF(
			[2]int64{3, 2},
			[2]int64{5, 7},
		),
		yTail,
		8,
	)
	if err != nil {
		t.Fatalf("EvalGCFWithTailExact failed: %v", err)
	}

	want, err := base.ApplyRat(x, wantY)
	if err != nil {
		t.Fatalf("base ApplyRat failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestApplyComposedGCFYBLFTToTailsExact_GeneralMatchesExactEvaluation(t *testing.T) {
	base := NewBLFT(2, 3, 5, 7, 11, 13, 17, 19)
	x := mustRat(3, 2)
	ySrc := NewSliceGCF(
		[2]int64{5, 7},
	)
	yTail := mustRat(13, 1)

	got, ingested, err := ApplyComposedGCFYBLFTToTailsExact(base, x, ySrc, yTail, 8)
	if err != nil {
		t.Fatalf("ApplyComposedGCFYBLFTToTailsExact failed: %v", err)
	}
	if ingested != 1 {
		t.Fatalf("got ingested=%d want 1", ingested)
	}

	y, _, err := EvalGCFWithTailExact(NewSliceGCF([2]int64{5, 7}), yTail, 8)
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

func TestApplyComposedGCFYBLFTToTailsExact_RequiresExhaustionWithinBound(t *testing.T) {
	base := NewBLFT(0, 0, 1, 0, 0, 0, 0, 1)

	_, ingested, err := ApplyComposedGCFYBLFTToTailsExact(
		base,
		mustRat(2, 1),
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

func TestApplyComposedGCFYBLFTToTailsExact_InvalidTermPropagatesError(t *testing.T) {
	base := NewBLFT(0, 0, 1, 0, 0, 0, 0, 1)

	_, ingested, err := ApplyComposedGCFYBLFTToTailsExact(
		base,
		mustRat(2, 1),
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

// gcf_blft_exact_y_test.go v1
