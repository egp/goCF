// diag_blft_ingest_gcf_test.go v1
package cf

import (
	"strings"
	"testing"
)

func TestDiagBLFTIngestGCF_RejectsBadQ(t *testing.T) {
	base := NewDiagBLFT(
		mustBig(1), mustBig(2), mustBig(3),
		mustBig(4), mustBig(5), mustBig(6),
	)

	_, err := base.IngestGCF(3, 0)
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "require q>0") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDiagBLFTIngestGCF_IdentityLikeProjection_RewriteLaw(t *testing.T) {
	// D(x) = x  =>  (0*x^2 + 1*x + 0) / (0*x^2 + 0*x + 1)
	base := NewDiagBLFT(
		mustBig(0), mustBig(1), mustBig(0),
		mustBig(0), mustBig(0), mustBig(1),
	)

	p := int64(3)
	q := int64(2)
	xTail := mustRat(11, 1)

	rewritten, err := base.IngestGCF(p, q)
	if err != nil {
		t.Fatalf("IngestGCF failed: %v", err)
	}

	x, _, err := EvalGCFWithTailExact(NewSliceGCF([2]int64{p, q}), xTail, 8)
	if err != nil {
		t.Fatalf("EvalGCFWithTailExact failed: %v", err)
	}

	want, err := base.ApplyRat(x)
	if err != nil {
		t.Fatalf("base ApplyRat failed: %v", err)
	}

	got, err := rewritten.ApplyRat(xTail)
	if err != nil {
		t.Fatalf("rewritten ApplyRat failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestDiagBLFTIngestGCF_GeneralRewriteLaw(t *testing.T) {
	base := NewDiagBLFT(
		mustBig(2), mustBig(3), mustBig(5),
		mustBig(7), mustBig(11), mustBig(13),
	)

	p := int64(5)
	q := int64(7)
	xTail := mustRat(13, 1)

	rewritten, err := base.IngestGCF(p, q)
	if err != nil {
		t.Fatalf("IngestGCF failed: %v", err)
	}

	x, _, err := EvalGCFWithTailExact(NewSliceGCF([2]int64{p, q}), xTail, 8)
	if err != nil {
		t.Fatalf("EvalGCFWithTailExact failed: %v", err)
	}

	want, err := base.ApplyRat(x)
	if err != nil {
		t.Fatalf("base ApplyRat failed: %v", err)
	}

	got, err := rewritten.ApplyRat(xTail)
	if err != nil {
		t.Fatalf("rewritten ApplyRat failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got=%v want=%v", got, want)
	}
}
