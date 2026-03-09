// range_apply_ulft_test.go v1
package cf

import "testing"

func TestRangeApplyULFT_RejectsOutsideRange(t *testing.T) {
	r := Range{
		Lo: mustRat(2, 1),
		Hi: mustRat(1, 1),
	}

	tform := NewULFT(
		mustBig(1),
		mustBig(0),
		mustBig(0),
		mustBig(1),
	)

	_, err := r.ApplyULFT(tform)
	if err == nil {
		t.Fatalf("expected error for outside range")
	}
}

func TestRangeApplyULFT_DenomCrossesZero(t *testing.T) {
	r := NewRange(mustRat(-1, 1), mustRat(1, 1), true, true)

	tform := NewULFT(
		mustBig(1),
		mustBig(0),
		mustBig(1),
		mustBig(0),
	)

	_, err := r.ApplyULFT(tform)
	if err == nil {
		t.Fatalf("expected pole detection")
	}
}

func TestRangeApplyULFT_MonotoneImage(t *testing.T) {
	r := NewRange(mustRat(1, 1), mustRat(2, 1), true, true)

	tform := NewULFT(
		mustBig(1),
		mustBig(1),
		mustBig(0),
		mustBig(1),
	)

	out, err := r.ApplyULFT(tform)
	if err != nil {
		t.Fatalf("ApplyULFT failed: %v", err)
	}

	wantLo := mustRat(2, 1)
	wantHi := mustRat(3, 1)

	if out.Lo.Cmp(wantLo) != 0 || out.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v] want [%v,%v]", out.Lo, out.Hi, wantLo, wantHi)
	}
}
