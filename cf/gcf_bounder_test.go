// gcf_bounder_test.go v1
package cf

import "testing"

func TestGCFBounder_Empty(t *testing.T) {
	b := NewGCFBounder()

	if b.HasValue() {
		t.Fatalf("expected HasValue=false")
	}

	_, ok, err := b.Range()
	if err != nil {
		t.Fatalf("Range failed: %v", err)
	}
	if ok {
		t.Fatalf("expected ok=false for empty range")
	}

	_, err = b.Convergent()
	if err == nil {
		t.Fatalf("expected error for empty convergent")
	}
}

func TestGCFBounder_SingleTerm(t *testing.T) {
	b := NewGCFBounder()
	if err := b.IngestPQ(3, 2); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}

	got, err := b.Convergent()
	if err != nil {
		t.Fatalf("Convergent failed: %v", err)
	}

	want := mustRat(3, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestGCFBounder_TwoTerms(t *testing.T) {
	b := NewGCFBounder()
	if err := b.IngestPQ(3, 2); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}
	if err := b.IngestPQ(5, 7); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}

	got, err := b.Convergent()
	if err != nil {
		t.Fatalf("Convergent failed: %v", err)
	}

	// 3 + 2/5 = 17/5
	want := mustRat(17, 5)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestGCFBounder_ThreeTerms(t *testing.T) {
	b := NewGCFBounder()
	if err := b.IngestPQ(1, 1); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}
	if err := b.IngestPQ(2, 1); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}
	if err := b.IngestPQ(2, 1); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}

	got, err := b.Convergent()
	if err != nil {
		t.Fatalf("Convergent failed: %v", err)
	}

	// 1 + 1/(2 + 1/2) = 7/5
	want := mustRat(7, 5)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestGCFBounder_RangeIsPoint(t *testing.T) {
	b := NewGCFBounder()
	if err := b.IngestPQ(3, 2); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}
	if err := b.IngestPQ(5, 7); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}

	r, ok, err := b.Range()
	if err != nil {
		t.Fatalf("Range failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	want := mustRat(17, 5)
	if r.Lo.Cmp(want) != 0 || r.Hi.Cmp(want) != 0 {
		t.Fatalf("got [%v,%v] want exact [%v,%v]", r.Lo, r.Hi, want, want)
	}
}

func TestGCFBounder_FinishAndIngestAfterFinish(t *testing.T) {
	b := NewGCFBounder()
	if err := b.IngestPQ(3, 2); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}

	b.Finish()

	if err := b.IngestPQ(5, 7); err == nil {
		t.Fatalf("expected error ingesting after Finish")
	}

	r, ok, err := b.Range()
	if err != nil {
		t.Fatalf("Range failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}
	want := mustRat(3, 1)
	if r.Lo.Cmp(want) != 0 || r.Hi.Cmp(want) != 0 {
		t.Fatalf("got [%v,%v] want exact [%v,%v]", r.Lo, r.Hi, want, want)
	}
}

func TestGCFBounder_RejectsBadQ(t *testing.T) {
	b := NewGCFBounder()

	if err := b.IngestPQ(3, 0); err == nil {
		t.Fatalf("expected error for q=0")
	}
	if err := b.IngestPQ(3, -1); err == nil {
		t.Fatalf("expected error for q<0")
	}
}

// gcf_bounder_test.go v1
