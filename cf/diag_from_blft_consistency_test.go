// diag_from_blft_consistency_test.go v1
package cf

import "testing"

func TestDiagFromBLFT_ApplyRatMatchesBLFTOnDiagonal(t *testing.T) {
	base := NewBLFT(2, 3, 5, 7, 11, 13, 17, 19)
	diag := DiagFromBLFT(base)

	x := mustRat(7, 3)

	want, err := base.ApplyRat(x, x)
	if err != nil {
		t.Fatalf("BLFT ApplyRat failed: %v", err)
	}

	got, err := diag.ApplyRat(x)
	if err != nil {
		t.Fatalf("DiagBLFT ApplyRat failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestDiagFromBLFT_GCFIngestionMatchesTwoSidedDiagonalSpecialization(t *testing.T) {
	base := NewBLFT(2, 3, 5, 7, 11, 13, 17, 19)

	p := int64(5)
	q := int64(7)
	xTail := mustRat(13, 1)

	// Diagonal path: specialize first, then ingest once.
	diagBase := DiagFromBLFT(base)
	diagRewritten, err := diagBase.IngestGCF(p, q)
	if err != nil {
		t.Fatalf("DiagBLFT IngestGCF failed: %v", err)
	}

	// BLFT path: ingest same term into x and y, then specialize.
	xyRewritten, err := base.IngestGCFX(p, q)
	if err != nil {
		t.Fatalf("BLFT IngestGCFX failed: %v", err)
	}
	xyRewritten, err = xyRewritten.IngestGCFY(p, q)
	if err != nil {
		t.Fatalf("BLFT IngestGCFY failed: %v", err)
	}
	diagFromXY := DiagFromBLFT(xyRewritten)

	gotDiag, err := diagRewritten.ApplyRat(xTail)
	if err != nil {
		t.Fatalf("diagRewritten ApplyRat failed: %v", err)
	}
	gotFromXY, err := diagFromXY.ApplyRat(xTail)
	if err != nil {
		t.Fatalf("diagFromXY ApplyRat failed: %v", err)
	}

	if gotDiag.Cmp(gotFromXY) != 0 {
		t.Fatalf("gotDiag=%v gotFromXY=%v", gotDiag, gotFromXY)
	}
}

func TestApplyComposedGCFDiagBLFTToTailExact_MatchesBLFTTwoSidedDiagonalCase(t *testing.T) {
	base := NewBLFT(2, 3, 5, 7, 11, 13, 17, 19)
	diag := DiagFromBLFT(base)

	srcTerms := [][2]int64{
		{3, 2},
		{4, 5},
	}
	xTail := mustRat(7, 1)

	gotDiag, diagIngested, err := ApplyComposedGCFDiagBLFTToTailExact(
		diag,
		NewSliceGCF(srcTerms...),
		xTail,
		8,
	)
	if err != nil {
		t.Fatalf("ApplyComposedGCFDiagBLFTToTailExact failed: %v", err)
	}
	if diagIngested != 2 {
		t.Fatalf("got diagIngested=%d want 2", diagIngested)
	}

	gotBLFT, xIngested, yIngested, err := ApplyComposedGCFXYBLFTToTailsExact(
		base,
		NewSliceGCF(srcTerms...), xTail, 8,
		NewSliceGCF(srcTerms...), xTail, 8,
	)
	if err != nil {
		t.Fatalf("ApplyComposedGCFXYBLFTToTailsExact failed: %v", err)
	}
	if xIngested != 2 {
		t.Fatalf("got xIngested=%d want 2", xIngested)
	}
	if yIngested != 2 {
		t.Fatalf("got yIngested=%d want 2", yIngested)
	}

	if gotDiag.Cmp(gotBLFT) != 0 {
		t.Fatalf("gotDiag=%v gotBLFT=%v", gotDiag, gotBLFT)
	}
}
