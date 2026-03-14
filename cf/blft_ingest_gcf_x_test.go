// blft_ingest_gcf_x_test.go v1
package cf

import (
	"strings"
	"testing"
)

func TestBLFTIngestGCFX_RejectsBadQ(t *testing.T) {
	base := NewBLFT(1, 2, 3, 4, 5, 6, 7, 8)

	_, err := base.IngestGCFX(3, 0)
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "require q>0") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBLFTIngestGCFX_IdentityLikeXProjection_RewriteLaw(t *testing.T) {
	// B(x,y) = x
	base := NewBLFT(0, 1, 0, 0, 0, 0, 0, 1)

	p := int64(3)
	q := int64(2)
	xTail := mustRat(11, 1)
	y := mustRat(7, 1)

	rewritten, err := base.IngestGCFX(p, q)
	if err != nil {
		t.Fatalf("IngestGCFX failed: %v", err)
	}

	x, _, err := EvalGCFWithTailExact(NewSliceGCF([2]int64{p, q}), xTail, 8)
	if err != nil {
		t.Fatalf("EvalGCFWithTailExact failed: %v", err)
	}

	want, err := base.ApplyRat(x, y)
	if err != nil {
		t.Fatalf("base ApplyRat failed: %v", err)
	}

	got, err := rewritten.ApplyRat(xTail, y)
	if err != nil {
		t.Fatalf("rewritten ApplyRat failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestBLFTIngestGCFX_GeneralRewriteLaw(t *testing.T) {
	base := NewBLFT(2, 3, 5, 7, 11, 13, 17, 19)

	p := int64(5)
	q := int64(7)
	xTail := mustRat(13, 1)
	y := mustRat(3, 2)

	rewritten, err := base.IngestGCFX(p, q)
	if err != nil {
		t.Fatalf("IngestGCFX failed: %v", err)
	}

	x, _, err := EvalGCFWithTailExact(NewSliceGCF([2]int64{p, q}), xTail, 8)
	if err != nil {
		t.Fatalf("EvalGCFWithTailExact failed: %v", err)
	}

	want, err := base.ApplyRat(x, y)
	if err != nil {
		t.Fatalf("base ApplyRat failed: %v", err)
	}

	got, err := rewritten.ApplyRat(xTail, y)
	if err != nil {
		t.Fatalf("rewritten ApplyRat failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got=%v want=%v", got, want)
	}
}
