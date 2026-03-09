// cf_to_gcf_adapter_test.go v1
package cf

import "testing"

func TestAdaptCFToGCF_Empty(t *testing.T) {
	g := AdaptCFToGCF(NewSliceCF())

	_, _, ok := g.NextPQ()
	if ok {
		t.Fatalf("expected empty adapted source to terminate")
	}
}

func TestAdaptCFToGCF_SliceCF(t *testing.T) {
	g := AdaptCFToGCF(NewSliceCF(3, 7, 16))

	p, q, ok := g.NextPQ()
	if !ok {
		t.Fatalf("expected first term")
	}
	if p != 3 || q != 1 {
		t.Fatalf("got (%d,%d), want (3,1)", p, q)
	}

	p, q, ok = g.NextPQ()
	if !ok {
		t.Fatalf("expected second term")
	}
	if p != 7 || q != 1 {
		t.Fatalf("got (%d,%d), want (7,1)", p, q)
	}

	p, q, ok = g.NextPQ()
	if !ok {
		t.Fatalf("expected third term")
	}
	if p != 16 || q != 1 {
		t.Fatalf("got (%d,%d), want (16,1)", p, q)
	}

	_, _, ok = g.NextPQ()
	if ok {
		t.Fatalf("expected termination after three terms")
	}
}

func TestAdaptCFToGCF_PeriodicCF(t *testing.T) {
	g := AdaptCFToGCF(Sqrt2CF())

	wantP := []int64{1, 2, 2, 2, 2}
	for i, want := range wantP {
		p, q, ok := g.NextPQ()
		if !ok {
			t.Fatalf("expected term %d", i)
		}
		if p != want || q != 1 {
			t.Fatalf("term %d: got (%d,%d), want (%d,1)", i, p, q, want)
		}
	}
}

func TestAdaptCFToGCF_IngestAllMatchesRegularFiniteValue(t *testing.T) {
	// Regular CF [3;7,16] = 355/113.
	b, err := IngestAllGCF(AdaptCFToGCF(NewSliceCF(3, 7, 16)))
	if err != nil {
		t.Fatalf("IngestAllGCF failed: %v", err)
	}

	got, err := b.Convergent()
	if err != nil {
		t.Fatalf("Convergent failed: %v", err)
	}

	want := mustRat(355, 113)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

// cf_to_gcf_adapter_test.go v1
