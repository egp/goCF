// cf_gcf_equiv_test.go v1
package cf

import "testing"

func TestFiniteCFAndAdaptedGCF_Agree_On355Over113(t *testing.T) {
	// Regular CF [3;7,16] = 355/113
	cfTerms := []int64{3, 7, 16}

	// Bounder path
	b := NewBounder()
	for _, a := range cfTerms {
		if err := b.Ingest(a); err != nil {
			t.Fatalf("Bounder.Ingest failed: %v", err)
		}
	}
	b.Finish()

	gotCF, err := b.Convergent()
	if err != nil {
		t.Fatalf("Bounder.Convergent failed: %v", err)
	}

	// Adapted GCF path
	gb, err := IngestAllGCF(AdaptCFToGCF(NewSliceCF(cfTerms...)))
	if err != nil {
		t.Fatalf("IngestAllGCF failed: %v", err)
	}
	gotGCF, err := gb.Convergent()
	if err != nil {
		t.Fatalf("GCFBounder.Convergent failed: %v", err)
	}

	want := mustRat(355, 113)

	if gotCF.Cmp(want) != 0 {
		t.Fatalf("CF got %v want %v", gotCF, want)
	}
	if gotGCF.Cmp(want) != 0 {
		t.Fatalf("GCF got %v want %v", gotGCF, want)
	}
	if gotCF.Cmp(gotGCF) != 0 {
		t.Fatalf("CF and GCF differ: cf=%v gcf=%v", gotCF, gotGCF)
	}
}

func TestFiniteCFAndAdaptedGCF_Agree_OnSingleTerm(t *testing.T) {
	cfTerms := []int64{42}

	b := NewBounder()
	for _, a := range cfTerms {
		if err := b.Ingest(a); err != nil {
			t.Fatalf("Bounder.Ingest failed: %v", err)
		}
	}
	b.Finish()

	gotCF, err := b.Convergent()
	if err != nil {
		t.Fatalf("Bounder.Convergent failed: %v", err)
	}

	gb, err := IngestAllGCF(AdaptCFToGCF(NewSliceCF(cfTerms...)))
	if err != nil {
		t.Fatalf("IngestAllGCF failed: %v", err)
	}
	gotGCF, err := gb.Convergent()
	if err != nil {
		t.Fatalf("GCFBounder.Convergent failed: %v", err)
	}

	want := mustRat(42, 1)

	if gotCF.Cmp(want) != 0 {
		t.Fatalf("CF got %v want %v", gotCF, want)
	}
	if gotGCF.Cmp(want) != 0 {
		t.Fatalf("GCF got %v want %v", gotGCF, want)
	}
	if gotCF.Cmp(gotGCF) != 0 {
		t.Fatalf("CF and GCF differ: cf=%v gcf=%v", gotCF, gotGCF)
	}
}

func TestFiniteCFAndAdaptedGCF_Agree_OnNegativeLeadingTerm(t *testing.T) {
	// [-1;2] = -1 + 1/2 = -1/2
	cfTerms := []int64{-1, 2}

	b := NewBounder()
	for _, a := range cfTerms {
		if err := b.Ingest(a); err != nil {
			t.Fatalf("Bounder.Ingest failed: %v", err)
		}
	}
	b.Finish()

	gotCF, err := b.Convergent()
	if err != nil {
		t.Fatalf("Bounder.Convergent failed: %v", err)
	}

	gb, err := IngestAllGCF(AdaptCFToGCF(NewSliceCF(cfTerms...)))
	if err != nil {
		t.Fatalf("IngestAllGCF failed: %v", err)
	}
	gotGCF, err := gb.Convergent()
	if err != nil {
		t.Fatalf("GCFBounder.Convergent failed: %v", err)
	}

	want := mustRat(-1, 2)

	if gotCF.Cmp(want) != 0 {
		t.Fatalf("CF got %v want %v", gotCF, want)
	}
	if gotGCF.Cmp(want) != 0 {
		t.Fatalf("GCF got %v want %v", gotGCF, want)
	}
	if gotCF.Cmp(gotGCF) != 0 {
		t.Fatalf("CF and GCF differ: cf=%v gcf=%v", gotCF, gotGCF)
	}
}

func TestFiniteCFAndAdaptedGCF_RangesAgree_On355Over113(t *testing.T) {
	cfTerms := []int64{3, 7, 16}

	b := NewBounder()
	for _, a := range cfTerms {
		if err := b.Ingest(a); err != nil {
			t.Fatalf("Bounder.Ingest failed: %v", err)
		}
	}
	b.Finish()

	rCF, ok, err := b.Range()
	if err != nil {
		t.Fatalf("Bounder.Range failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected CF range")
	}

	gb, err := IngestAllGCF(AdaptCFToGCF(NewSliceCF(cfTerms...)))
	if err != nil {
		t.Fatalf("IngestAllGCF failed: %v", err)
	}
	rGCF, ok, err := gb.Range()
	if err != nil {
		t.Fatalf("GCFBounder.Range failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected GCF range")
	}

	if rCF.Lo.Cmp(rGCF.Lo) != 0 || rCF.Hi.Cmp(rGCF.Hi) != 0 {
		t.Fatalf("range mismatch: cf=[%v,%v] gcf=[%v,%v]", rCF.Lo, rCF.Hi, rGCF.Lo, rGCF.Hi)
	}
}

// cf_gcf_equiv_test.go v1
