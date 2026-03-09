// gcf_tail_ray_test.go v1
package cf

import (
	"math/big"
	"testing"
)

func TestApplyULFTToTailRay_Reciprocal(t *testing.T) {
	// T(x) = 1/x, tail in [1,+∞) => image [0,1]
	tform := NewULFT(big.NewInt(0), big.NewInt(1), big.NewInt(1), big.NewInt(0))

	got, err := ApplyULFTToTailRay(tform, mustRat(1, 1))
	if err != nil {
		t.Fatalf("ApplyULFTToTailRay failed: %v", err)
	}

	wantLo := mustRat(0, 1)
	wantHi := mustRat(1, 1)
	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v] want [%v,%v]", got.Lo, got.Hi, wantLo, wantHi)
	}
}

func TestApplyULFTToTailRay_Constant(t *testing.T) {
	// T(x) = 2
	tform := NewULFT(big.NewInt(0), big.NewInt(2), big.NewInt(0), big.NewInt(1))

	got, err := ApplyULFTToTailRay(tform, mustRat(1, 1))
	if err != nil {
		t.Fatalf("ApplyULFTToTailRay failed: %v", err)
	}

	want := mustRat(2, 1)
	if got.Lo.Cmp(want) != 0 || got.Hi.Cmp(want) != 0 {
		t.Fatalf("got [%v,%v] want exact [%v,%v]", got.Lo, got.Hi, want, want)
	}
}

func TestApplyULFTToTailRay_RejectsPoleOnRay(t *testing.T) {
	// T(x) = 1/(x-1), pole at x=1
	tform := NewULFT(big.NewInt(0), big.NewInt(1), big.NewInt(1), big.NewInt(-1))

	_, err := ApplyULFTToTailRay(tform, mustRat(1, 1))
	if err == nil {
		t.Fatalf("expected pole-crossing error")
	}
}

func TestApplyULFTToTailRay_RejectsUnboundedAffineCase(t *testing.T) {
	// T(x) = x + 1
	tform := NewULFT(big.NewInt(1), big.NewInt(1), big.NewInt(0), big.NewInt(1))

	_, err := ApplyULFTToTailRay(tform, mustRat(1, 1))
	if err == nil {
		t.Fatalf("expected unbounded affine image error")
	}
}

func TestGCFBounder_RangeUsesTailLowerBound(t *testing.T) {
	// Prefix terms (1,1), (2,1) define:
	// x = 1 + 1/(2 + 1/tail)
	//
	// With tail in [1,+∞), image is [4/3, 3/2].
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

	got, ok, err := b.Range()
	if err != nil {
		t.Fatalf("Range failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	wantLo := mustRat(4, 3)
	wantHi := mustRat(3, 2)
	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v] want [%v,%v]", got.Lo, got.Hi, wantLo, wantHi)
	}
}

func TestGCFBounder_RangeFallsBackToPointPlaceholderWithoutTailBound(t *testing.T) {
	b := NewGCFBounder()
	if err := b.IngestPQ(1, 1); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}
	if err := b.IngestPQ(2, 1); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}

	got, ok, err := b.Range()
	if err != nil {
		t.Fatalf("Range failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	want := mustRat(3, 2)
	if got.Lo.Cmp(want) != 0 || got.Hi.Cmp(want) != 0 {
		t.Fatalf("got [%v,%v] want placeholder [%v,%v]", got.Lo, got.Hi, want, want)
	}
}

func TestGCFBounder_FinishStillReturnsExactPoint(t *testing.T) {
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

	b.Finish()

	got, ok, err := b.Range()
	if err != nil {
		t.Fatalf("Range failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	want := mustRat(3, 2)
	if got.Lo.Cmp(want) != 0 || got.Hi.Cmp(want) != 0 {
		t.Fatalf("got [%v,%v] want exact [%v,%v]", got.Lo, got.Hi, want, want)
	}
}

func TestGCFBounder_RejectsNonPositiveTailLowerBound(t *testing.T) {
	b := NewGCFBounder()

	if err := b.SetTailLowerBound(mustRat(0, 1)); err == nil {
		t.Fatalf("expected error for zero lower bound")
	}
}

// gcf_tail_ray_test.go v1
