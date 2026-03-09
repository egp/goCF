// gcf_ingest_all_test.go v1
package cf

import "testing"

func TestIngestAllGCF_Empty(t *testing.T) {
	b, err := IngestAllGCF(NewSliceGCF())
	if err != nil {
		t.Fatalf("IngestAllGCF failed: %v", err)
	}
	if b == nil {
		t.Fatalf("expected non-nil bounder")
	}
	if b.HasValue() {
		t.Fatalf("expected empty bounder")
	}

	_, ok, err := b.Range()
	if err != nil {
		t.Fatalf("Range failed: %v", err)
	}
	if ok {
		t.Fatalf("expected ok=false for empty bounder")
	}
}

func TestIngestAllGCF_TwoTerms(t *testing.T) {
	b, err := IngestAllGCF(NewSliceGCF(
		[2]int64{3, 2},
		[2]int64{5, 7},
	))
	if err != nil {
		t.Fatalf("IngestAllGCF failed: %v", err)
	}

	got, err := b.Convergent()
	if err != nil {
		t.Fatalf("Convergent failed: %v", err)
	}

	want := mustRat(17, 5)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestIngestAllGCF_RangeIsExactAfterFinish(t *testing.T) {
	b, err := IngestAllGCF(NewSliceGCF(
		[2]int64{1, 1},
		[2]int64{2, 1},
		[2]int64{2, 1},
	))
	if err != nil {
		t.Fatalf("IngestAllGCF failed: %v", err)
	}

	r, ok, err := b.Range()
	if err != nil {
		t.Fatalf("Range failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	want := mustRat(7, 5)
	if r.Lo.Cmp(want) != 0 || r.Hi.Cmp(want) != 0 {
		t.Fatalf("got [%v,%v] want exact [%v,%v]", r.Lo, r.Hi, want, want)
	}
}

func TestIngestAllGCF_RejectsBadQ(t *testing.T) {
	_, err := IngestAllGCF(NewSliceGCF(
		[2]int64{3, 0},
	))
	if err == nil {
		t.Fatalf("expected error for bad q")
	}
}

// gcf_ingest_all_test.go v1
