// blft_ingest_gcf_y_test.go v1
package cf

import (
	"strings"
	"testing"
)

func TestBLFTIngestGCFY_RejectsBadQ(t *testing.T) {
	base := NewBLFT(1, 2, 3, 4, 5, 6, 7, 8)

	_, err := base.IngestGCFY(3, 0)
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "require q>0") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBLFTIngestGCFY_IdentityLikeYProjection_RewriteLaw(t *testing.T) {
	// B(x,y) = y
	base := NewBLFT(0, 0, 1, 0, 0, 0, 0, 1)

	p := int64(3)
	q := int64(2)
	yTail := mustRat(11, 1)
	x := mustRat(7, 1)

	rewritten, err := base.IngestGCFY(p, q)
	if err != nil {
		t.Fatalf("IngestGCFY failed: %v", err)
	}

	y, _, err := EvalGCFWithTailExact(NewSliceGCF([2]int64{p, q}), yTail, 8)
	if err != nil {
		t.Fatalf("EvalGCFWithTailExact failed: %v", err)
	}

	want, err := base.ApplyRat(x, y)
	if err != nil {
		t.Fatalf("base ApplyRat failed: %v", err)
	}

	got, err := rewritten.ApplyRat(x, yTail)
	if err != nil {
		t.Fatalf("rewritten ApplyRat failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestBLFTIngestGCFY_GeneralRewriteLaw(t *testing.T) {
	base := NewBLFT(2, 3, 5, 7, 11, 13, 17, 19)

	p := int64(5)
	q := int64(7)
	yTail := mustRat(13, 1)
	x := mustRat(3, 2)

	rewritten, err := base.IngestGCFY(p, q)
	if err != nil {
		t.Fatalf("IngestGCFY failed: %v", err)
	}

	y, _, err := EvalGCFWithTailExact(NewSliceGCF([2]int64{p, q}), yTail, 8)
	if err != nil {
		t.Fatalf("EvalGCFWithTailExact failed: %v", err)
	}

	want, err := base.ApplyRat(x, y)
	if err != nil {
		t.Fatalf("base ApplyRat failed: %v", err)
	}

	got, err := rewritten.ApplyRat(x, yTail)
	if err != nil {
		t.Fatalf("rewritten ApplyRat failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got=%v want=%v", got, want)
	}
}
