// gcf_compose_test.go v1
package cf

import (
	"math/big"
	"strings"
	"testing"
)

func TestComposeGCFToULFT_EmptyIsIdentity(t *testing.T) {
	got, err := ComposeGCFToULFT(NewSliceGCF())
	if err != nil {
		t.Fatalf("ComposeGCFToULFT failed: %v", err)
	}

	want := NewULFT(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1))
	if got.A.Cmp(want.A) != 0 || got.B.Cmp(want.B) != 0 || got.C.Cmp(want.C) != 0 || got.D.Cmp(want.D) != 0 {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestComposeGCFToULFT_SingleTerm(t *testing.T) {
	got, err := ComposeGCFToULFT(NewSliceGCF([2]int64{3, 2}))
	if err != nil {
		t.Fatalf("ComposeGCFToULFT failed: %v", err)
	}

	// x = 3 + 2/x' = (3x' + 2)/x'
	want := NewULFT(big.NewInt(3), big.NewInt(2), big.NewInt(1), big.NewInt(0))
	if got.A.Cmp(want.A) != 0 || got.B.Cmp(want.B) != 0 || got.C.Cmp(want.C) != 0 || got.D.Cmp(want.D) != 0 {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestComposeGCFToULFT_TwoTermsByEvaluation(t *testing.T) {
	// x = 3 + 2/(5 + 7/x_tail)
	got, err := ComposeGCFToULFT(NewSliceGCF(
		[2]int64{3, 2},
		[2]int64{5, 7},
	))
	if err != nil {
		t.Fatalf("ComposeGCFToULFT failed: %v", err)
	}

	// Let x_tail = 11.
	tail := mustRat(11, 1)

	// x = 3 + 2/(5 + 7/11) = 3 + 2/(62/11) = 104/31
	want := mustRat(104, 31)

	gotVal, err := got.ApplyRat(tail)
	if err != nil {
		t.Fatalf("ApplyRat failed: %v", err)
	}
	if gotVal.Cmp(want) != 0 {
		t.Fatalf("got=%v want=%v", gotVal, want)
	}
}

func TestComposeGCFToULFT_RejectsBadQ(t *testing.T) {
	_, err := ComposeGCFToULFT(NewSliceGCF([2]int64{3, 0}))
	if err == nil {
		t.Fatalf("expected error for q=0")
	}
}

func TestComposeGCFIntoULFTBounded_FiniteSourceExhausts(t *testing.T) {
	base := NewULFT(mustBig(1), mustBig(0), mustBig(0), mustBig(1))
	src := NewSliceGCF(
		[2]int64{1, 2},
		[2]int64{3, 4},
	)

	got, ingested, exhausted, err := ComposeGCFIntoULFTBounded(base, src, 8)
	if err != nil {
		t.Fatalf("ComposeGCFIntoULFTBounded failed: %v", err)
	}
	if !exhausted {
		t.Fatalf("expected exhausted=true")
	}
	if ingested != 2 {
		t.Fatalf("got ingested=%d want 2", ingested)
	}

	want, err := composeGCFIntoULFT(base, NewSliceGCF(
		[2]int64{1, 2},
		[2]int64{3, 4},
	))
	if err != nil {
		t.Fatalf("composeGCFIntoULFT failed: %v", err)
	}

	if got.A.Cmp(want.A) != 0 || got.B.Cmp(want.B) != 0 || got.C.Cmp(want.C) != 0 || got.D.Cmp(want.D) != 0 {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestComposeGCFIntoULFTBounded_ZeroBoundFailsImmediately(t *testing.T) {
	base := NewULFT(mustBig(1), mustBig(0), mustBig(0), mustBig(1))
	src := NewSliceGCF([2]int64{1, 1})

	_, ingested, exhausted, err := ComposeGCFIntoULFTBounded(base, src, 0)
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "exceeded MaxIngestTerms=0") {
		t.Fatalf("unexpected error: %v", err)
	}
	if ingested != 0 {
		t.Fatalf("got ingested=%d want 0", ingested)
	}
	if exhausted {
		t.Fatalf("expected exhausted=false")
	}
}

func TestComposeGCFIntoULFTBounded_InfiniteSourceHitsBound(t *testing.T) {
	base := NewULFT(mustBig(1), mustBig(0), mustBig(0), mustBig(1))
	src := NewUnitPArithmeticQGCFSource(1, 1)

	_, ingested, exhausted, err := ComposeGCFIntoULFTBounded(base, src, 3)
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "exceeded MaxIngestTerms=3") {
		t.Fatalf("unexpected error: %v", err)
	}
	if ingested != 3 {
		t.Fatalf("got ingested=%d want 3", ingested)
	}
	if exhausted {
		t.Fatalf("expected exhausted=false")
	}
}

func TestComposeGCFIntoULFTBounded_InvalidTermPropagatesError(t *testing.T) {
	base := NewULFT(mustBig(1), mustBig(0), mustBig(0), mustBig(1))
	src := NewSliceGCF([2]int64{7, 0})

	_, ingested, exhausted, err := ComposeGCFIntoULFTBounded(base, src, 8)
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "require q>0") {
		t.Fatalf("unexpected error: %v", err)
	}
	if ingested != 0 {
		t.Fatalf("got ingested=%d want 0", ingested)
	}
	if exhausted {
		t.Fatalf("expected exhausted=false")
	}
}

// EOF gcf_compose_test.go v1
