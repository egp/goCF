// pbt_gcf_blft_diag_test.go v1
package cf

import (
	"testing"

	"pgregory.net/rapid"
)

func genSmallPositiveRat() *rapid.Generator[Rational] {
	return rapid.Custom(func(t *rapid.T) Rational {
		p := rapid.Int64Range(1, 9).Draw(t, "p")
		q := rapid.Int64Range(1, 9).Draw(t, "q")
		return mustRat(p, q)
	})
}

func genSmallBLFT() *rapid.Generator[BLFT] {
	return rapid.Custom(func(t *rapid.T) BLFT {
		return NewBLFT(
			rapid.Int64Range(-5, 5).Draw(t, "A"),
			rapid.Int64Range(-5, 5).Draw(t, "B"),
			rapid.Int64Range(-5, 5).Draw(t, "C"),
			rapid.Int64Range(-5, 5).Draw(t, "D"),
			rapid.Int64Range(-5, 5).Draw(t, "E"),
			rapid.Int64Range(-5, 5).Draw(t, "F"),
			rapid.Int64Range(-5, 5).Draw(t, "G"),
			rapid.Int64Range(-5, 5).Draw(t, "H"),
		)
	})
}

func genSmallDiagBLFT() *rapid.Generator[DiagBLFT] {
	return rapid.Custom(func(t *rapid.T) DiagBLFT {
		return NewDiagBLFT(
			mustBig(rapid.Int64Range(-5, 5).Draw(t, "A")),
			mustBig(rapid.Int64Range(-5, 5).Draw(t, "B")),
			mustBig(rapid.Int64Range(-5, 5).Draw(t, "C")),
			mustBig(rapid.Int64Range(-5, 5).Draw(t, "D")),
			mustBig(rapid.Int64Range(-5, 5).Draw(t, "E")),
			mustBig(rapid.Int64Range(-5, 5).Draw(t, "F")),
		)
	})
}

func TestPBT_BLFTIngestGCFX_RewriteLaw(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		base := genSmallBLFT().Draw(t, "base")
		p := rapid.Int64Range(-5, 5).Draw(t, "p")
		q := rapid.Int64Range(1, 9).Draw(t, "q")
		xTail := genSmallPositiveRat().Draw(t, "xTail")
		y := genSmallPositiveRat().Draw(t, "y")

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
			t.Skip()
		}

		got, err := rewritten.ApplyRat(xTail, y)
		if err != nil {
			t.Skip()
		}

		if got.Cmp(want) != 0 {
			t.Fatalf("got=%v want=%v base=%v p=%d q=%d xTail=%v y=%v x=%v",
				got, want, base, p, q, xTail, y, x)
		}
	})
}

func TestPBT_BLFTIngestGCFY_RewriteLaw(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		base := genSmallBLFT().Draw(t, "base")
		p := rapid.Int64Range(-5, 5).Draw(t, "p")
		q := rapid.Int64Range(1, 9).Draw(t, "q")
		x := genSmallPositiveRat().Draw(t, "x")
		yTail := genSmallPositiveRat().Draw(t, "yTail")

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
			t.Skip()
		}

		got, err := rewritten.ApplyRat(x, yTail)
		if err != nil {
			t.Skip()
		}

		if got.Cmp(want) != 0 {
			t.Fatalf("got=%v want=%v base=%v p=%d q=%d x=%v yTail=%v y=%v",
				got, want, base, p, q, x, yTail, y)
		}
	})
}

func TestPBT_DiagBLFTIngestGCF_RewriteLaw(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		base := genSmallDiagBLFT().Draw(t, "base")
		p := rapid.Int64Range(-5, 5).Draw(t, "p")
		q := rapid.Int64Range(1, 9).Draw(t, "q")
		xTail := genSmallPositiveRat().Draw(t, "xTail")

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
			t.Skip()
		}

		got, err := rewritten.ApplyRat(xTail)
		if err != nil {
			t.Skip()
		}

		if got.Cmp(want) != 0 {
			t.Fatalf("got=%v want=%v base=%v p=%d q=%d xTail=%v x=%v",
				got, want, base, p, q, xTail, x)
		}
	})
}

// pbt_gcf_blft_diag_test.go v1
