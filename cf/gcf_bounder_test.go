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

func TestGCFBounderRange_NoValue(t *testing.T) {
	b := NewGCFBounder()

	r, ok, err := b.Range()
	if err != nil {
		t.Fatalf("Range failed: %v", err)
	}
	if ok {
		t.Fatalf("expected ok=false, got range %v", r)
	}
}

func TestGCFBounderRange_FinishedIsExactPoint(t *testing.T) {
	b := NewGCFBounder()
	if err := b.IngestPQ(1, 1); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}
	if err := b.IngestPQ(2, 1); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}
	b.Finish()

	r, ok, err := b.Range()
	if err != nil {
		t.Fatalf("Range failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	want := mustRat(3, 2)
	if r.Lo.Cmp(want) != 0 || r.Hi.Cmp(want) != 0 {
		t.Fatalf("got [%v,%v] want exact [%v,%v]", r.Lo, r.Hi, want, want)
	}
}

func TestGCFBounderRange_UnfinishedWithoutTailMetadataIsPlaceholderPoint(t *testing.T) {
	b := NewGCFBounder()
	if err := b.IngestPQ(1, 1); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}
	if err := b.IngestPQ(2, 1); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}

	r, ok, err := b.Range()
	if err != nil {
		t.Fatalf("Range failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	want := mustRat(3, 2)
	if r.Lo.Cmp(want) != 0 || r.Hi.Cmp(want) != 0 {
		t.Fatalf("got [%v,%v] want placeholder [%v,%v]", r.Lo, r.Hi, want, want)
	}
}

func TestGCFBounderRange_UnfinishedWithTailLowerBoundUsesRayEnclosure(t *testing.T) {
	// Prefix terms (1,1), (2,1):
	// x = 1 + 1/(2 + 1/tail), tail >= 1
	// so range is [4/3, 3/2].
	b := NewGCFBounder()
	if err := b.IngestPQ(1, 1); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}
	if err := b.IngestPQ(2, 1); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}
	if err := b.SetTailLowerBound(mustRat(1, 1)); err != nil {
		t.Fatalf("SetTailLowerBound failed: %v", err)
	}

	r, ok, err := b.Range()
	if err != nil {
		t.Fatalf("Range failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	wantLo := mustRat(4, 3)
	wantHi := mustRat(3, 2)
	if r.Lo.Cmp(wantLo) != 0 || r.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v] want [%v,%v]", r.Lo, r.Hi, wantLo, wantHi)
	}
}

func TestGCFBounderRange_UnfinishedWithExplicitTailRangeUsesIntervalEnclosure(t *testing.T) {
	// Prefix terms (1,1), (2,1):
	// x = 1 + 1/(2 + 1/tail), tail in [1,2]
	// tail=1 -> 4/3
	// tail=2 -> 7/5
	b := NewGCFBounder()
	if err := b.IngestPQ(1, 1); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}
	if err := b.IngestPQ(2, 1); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}
	if err := b.SetTailRange(NewRange(mustRat(1, 1), mustRat(2, 1), true, true)); err != nil {
		t.Fatalf("SetTailRange failed: %v", err)
	}

	r, ok, err := b.Range()
	if err != nil {
		t.Fatalf("Range failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	wantLo := mustRat(4, 3)
	wantHi := mustRat(7, 5)
	if r.Lo.Cmp(wantLo) != 0 || r.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v] want [%v,%v]", r.Lo, r.Hi, wantLo, wantHi)
	}
}

func TestGCFBounderSetTailRange_RejectsOutsideRange(t *testing.T) {
	b := NewGCFBounder()

	err := b.SetTailRange(NewRange(mustRat(2, 1), mustRat(1, 1), true, true))
	if err == nil {
		t.Fatalf("expected error for outside tail range")
	}
}

func TestGCFBounderSetTailRange_RejectsNonPositiveRange(t *testing.T) {
	b := NewGCFBounder()

	err := b.SetTailRange(NewRange(mustRat(0, 1), mustRat(2, 1), true, true))
	if err == nil {
		t.Fatalf("expected error for non-positive tail range")
	}
}

// gcf_bounder_test.go v1
