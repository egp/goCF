// gcf_blft_exact_xy_test.go v1
package cf

import (
	"strings"
	"testing"
)

// func composeGCFIntoBLFTXY(base BLFT, xSrc, ySrc GCFSource) (BLFT, error) {
// 	cur := base

// 	for {
// 		p, q, ok := xSrc.NextPQ()
// 		if !ok {
// 			break
// 		}
// 		var err error
// 		cur, err = cur.IngestGCFX(p, q)
// 		if err != nil {
// 			return BLFT{}, err
// 		}
// 	}

// 	for {
// 		p, q, ok := ySrc.NextPQ()
// 		if !ok {
// 			break
// 		}
// 		var err error
// 		cur, err = cur.IngestGCFY(p, q)
// 		if err != nil {
// 			return BLFT{}, err
// 		}
// 	}

// 	return cur, nil
// }

func TestApplyComposedGCFXYBLFTToTailsExact_ProjectionXMatchesExactEvaluation(t *testing.T) {
	// B(x,y)=x
	base := NewBLFT(0, 1, 0, 0, 0, 0, 0, 1)

	xSrc := NewSliceGCF(
		[2]int64{3, 2},
		[2]int64{5, 7},
	)
	xTail := mustRat(11, 1)

	ySrc := NewSliceGCF(
		[2]int64{2, 3},
	)
	yTail := mustRat(13, 1)

	got, xIngested, yIngested, err := ApplyComposedGCFXYBLFTToTailsExact(
		base,
		xSrc, xTail, 8,
		ySrc, yTail, 8,
	)
	if err != nil {
		t.Fatalf("ApplyComposedGCFXYBLFTToTailsExact failed: %v", err)
	}
	if xIngested != 2 {
		t.Fatalf("got xIngested=%d want 2", xIngested)
	}
	if yIngested != 1 {
		t.Fatalf("got yIngested=%d want 1", yIngested)
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
		t.Fatalf("EvalGCFWithTailExact x failed: %v", err)
	}
	wantY, _, err := EvalGCFWithTailExact(
		NewSliceGCF(
			[2]int64{2, 3},
		),
		yTail,
		8,
	)
	if err != nil {
		t.Fatalf("EvalGCFWithTailExact y failed: %v", err)
	}

	want, err := base.ApplyRat(wantX, wantY)
	if err != nil {
		t.Fatalf("base ApplyRat failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestApplyComposedGCFXYBLFTToTailsExact_GeneralMatchesExactEvaluation(t *testing.T) {
	base := NewBLFT(2, 3, 5, 7, 11, 13, 17, 19)

	xSrc := NewSliceGCF([2]int64{5, 7})
	xTail := mustRat(13, 1)

	ySrc := NewSliceGCF(
		[2]int64{3, 2},
		[2]int64{4, 5},
	)
	yTail := mustRat(7, 1)

	got, xIngested, yIngested, err := ApplyComposedGCFXYBLFTToTailsExact(
		base,
		xSrc, xTail, 8,
		ySrc, yTail, 8,
	)
	if err != nil {
		t.Fatalf("ApplyComposedGCFXYBLFTToTailsExact failed: %v", err)
	}
	if xIngested != 1 {
		t.Fatalf("got xIngested=%d want 1", xIngested)
	}
	if yIngested != 2 {
		t.Fatalf("got yIngested=%d want 2", yIngested)
	}

	x, _, err := EvalGCFWithTailExact(NewSliceGCF([2]int64{5, 7}), xTail, 8)
	if err != nil {
		t.Fatalf("EvalGCFWithTailExact x failed: %v", err)
	}
	y, _, err := EvalGCFWithTailExact(
		NewSliceGCF(
			[2]int64{3, 2},
			[2]int64{4, 5},
		),
		yTail,
		8,
	)
	if err != nil {
		t.Fatalf("EvalGCFWithTailExact y failed: %v", err)
	}

	want, err := base.ApplyRat(x, y)
	if err != nil {
		t.Fatalf("base ApplyRat failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestApplyComposedGCFXYBLFTToTailsExact_XRequiresExhaustionWithinBound(t *testing.T) {
	base := NewBLFT(0, 1, 0, 0, 0, 0, 0, 1)

	_, xIngested, yIngested, err := ApplyComposedGCFXYBLFTToTailsExact(
		base,
		NewUnitPArithmeticQGCFSource(1, 1), mustRat(1, 1), 3,
		NewSliceGCF([2]int64{2, 3}), mustRat(5, 1), 8,
	)
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "exceeded MaxIngestTerms=3") {
		t.Fatalf("unexpected error: %v", err)
	}
	if xIngested != 3 {
		t.Fatalf("got xIngested=%d want 3", xIngested)
	}
	if yIngested != 0 {
		t.Fatalf("got yIngested=%d want 0", yIngested)
	}
}

func TestApplyComposedGCFXYBLFTToTailsExact_YRequiresExhaustionWithinBound(t *testing.T) {
	base := NewBLFT(0, 0, 1, 0, 0, 0, 0, 1)

	_, xIngested, yIngested, err := ApplyComposedGCFXYBLFTToTailsExact(
		base,
		NewSliceGCF([2]int64{2, 3}), mustRat(5, 1), 8,
		NewUnitPArithmeticQGCFSource(1, 1), mustRat(1, 1), 3,
	)
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "exceeded MaxIngestTerms=3") {
		t.Fatalf("unexpected error: %v", err)
	}
	if xIngested != 1 {
		t.Fatalf("got xIngested=%d want 1", xIngested)
	}
	if yIngested != 3 {
		t.Fatalf("got yIngested=%d want 3", yIngested)
	}
}

func TestApplyComposedGCFXYBLFTToTailsExact_InvalidXTermPropagatesError(t *testing.T) {
	base := NewBLFT(0, 1, 0, 0, 0, 0, 0, 1)

	_, xIngested, yIngested, err := ApplyComposedGCFXYBLFTToTailsExact(
		base,
		NewSliceGCF([2]int64{7, 0}), mustRat(1, 1), 8,
		NewSliceGCF([2]int64{2, 3}), mustRat(5, 1), 8,
	)
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "require q>0") {
		t.Fatalf("unexpected error: %v", err)
	}
	if xIngested != 0 {
		t.Fatalf("got xIngested=%d want 0", xIngested)
	}
	if yIngested != 0 {
		t.Fatalf("got yIngested=%d want 0", yIngested)
	}
}

func TestApplyComposedGCFXYBLFTToTailsExact_InvalidYTermPropagatesError(t *testing.T) {
	base := NewBLFT(0, 0, 1, 0, 0, 0, 0, 1)

	_, xIngested, yIngested, err := ApplyComposedGCFXYBLFTToTailsExact(
		base,
		NewSliceGCF([2]int64{2, 3}), mustRat(5, 1), 8,
		NewSliceGCF([2]int64{7, 0}), mustRat(1, 1), 8,
	)
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "require q>0") {
		t.Fatalf("unexpected error: %v", err)
	}
	if xIngested != 1 {
		t.Fatalf("got xIngested=%d want 1", xIngested)
	}
	if yIngested != 0 {
		t.Fatalf("got yIngested=%d want 0", yIngested)
	}
}

// gcf_blft_exact_xy_test.go v1
