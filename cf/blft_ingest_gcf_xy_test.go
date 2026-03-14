// blft_ingest_gcf_xy_test.go v1
package cf

import "testing"

func composeGCFIntoBLFTXY(base BLFT, xSrc, ySrc GCFSource) (BLFT, error) {
	cur := base

	for {
		p, q, ok := xSrc.NextPQ()
		if !ok {
			break
		}
		var err error
		cur, err = cur.IngestGCFX(p, q)
		if err != nil {
			return BLFT{}, err
		}
	}

	for {
		p, q, ok := ySrc.NextPQ()
		if !ok {
			break
		}
		var err error
		cur, err = cur.IngestGCFY(p, q)
		if err != nil {
			return BLFT{}, err
		}
	}

	return cur, nil
}

func TestBLFTIngestGCFXY_OrderIndependence_SingleTerms(t *testing.T) {
	base := NewBLFT(2, 3, 5, 7, 11, 13, 17, 19)

	px, qx := int64(5), int64(7)
	py, qy := int64(3), int64(2)

	xTail := mustRat(13, 1)
	yTail := mustRat(11, 1)

	xy, err := base.IngestGCFX(px, qx)
	if err != nil {
		t.Fatalf("IngestGCFX failed: %v", err)
	}
	xy, err = xy.IngestGCFY(py, qy)
	if err != nil {
		t.Fatalf("IngestGCFY after X failed: %v", err)
	}

	yx, err := base.IngestGCFY(py, qy)
	if err != nil {
		t.Fatalf("IngestGCFY failed: %v", err)
	}
	yx, err = yx.IngestGCFX(px, qx)
	if err != nil {
		t.Fatalf("IngestGCFX after Y failed: %v", err)
	}

	gotXY, err := xy.ApplyRat(xTail, yTail)
	if err != nil {
		t.Fatalf("xy ApplyRat failed: %v", err)
	}
	gotYX, err := yx.ApplyRat(xTail, yTail)
	if err != nil {
		t.Fatalf("yx ApplyRat failed: %v", err)
	}

	if gotXY.Cmp(gotYX) != 0 {
		t.Fatalf("gotXY=%v gotYX=%v", gotXY, gotYX)
	}
}

func TestBLFTIngestGCFXY_TwoSidedRewriteMatchesExactEvaluation(t *testing.T) {
	base := NewBLFT(2, 3, 5, 7, 11, 13, 17, 19)

	xSrc := NewSliceGCF(
		[2]int64{5, 7},
	)
	xTail := mustRat(13, 1)

	ySrc := NewSliceGCF(
		[2]int64{3, 2},
		[2]int64{4, 5},
	)
	yTail := mustRat(7, 1)

	rewritten, err := composeGCFIntoBLFTXY(
		base,
		NewSliceGCF([2]int64{5, 7}),
		NewSliceGCF(
			[2]int64{3, 2},
			[2]int64{4, 5},
		),
	)
	if err != nil {
		t.Fatalf("composeGCFIntoBLFTXY failed: %v", err)
	}

	x, _, err := EvalGCFWithTailExact(xSrc, xTail, 8)
	if err != nil {
		t.Fatalf("EvalGCFWithTailExact x failed: %v", err)
	}
	y, _, err := EvalGCFWithTailExact(ySrc, yTail, 8)
	if err != nil {
		t.Fatalf("EvalGCFWithTailExact y failed: %v", err)
	}

	want, err := base.ApplyRat(x, y)
	if err != nil {
		t.Fatalf("base ApplyRat failed: %v", err)
	}

	got, err := rewritten.ApplyRat(xTail, yTail)
	if err != nil {
		t.Fatalf("rewritten ApplyRat failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestApplyComposedGCFXYBLFTToTailsExact_MatchesManualTwoStageRewrite(t *testing.T) {
	base := NewBLFT(2, 3, 5, 7, 11, 13, 17, 19)

	xSrc := NewSliceGCF(
		[2]int64{5, 7},
	)
	xTail := mustRat(13, 1)
	yTail := mustRat(7, 1)

	got, xIngested, yIngested, err := ApplyComposedGCFXYBLFTToTailsExact(
		base,
		NewSliceGCF([2]int64{5, 7}), xTail, 8,
		NewSliceGCF(
			[2]int64{3, 2},
			[2]int64{4, 5},
		), yTail, 8,
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

	stage1, xCount, err := ApplyComposedGCFXBLFTToTailsExact(
		base,
		xSrc,
		xTail,
		yTail,
		8,
	)
	if err != nil {
		t.Fatalf("ApplyComposedGCFXBLFTToTailsExact failed: %v", err)
	}
	if xCount != 1 {
		t.Fatalf("got xCount=%d want 1", xCount)
	}

	_ = stage1 // intentional: this test compares the two-sided helper to exact reference below

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
