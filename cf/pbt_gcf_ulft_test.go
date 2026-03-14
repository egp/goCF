package cf

import (
	"testing"

	"pgregory.net/rapid"
)

func genSmallPositiveInt64(min, max int64) *rapid.Generator[int64] {
	return rapid.Int64Range(min, max)
}

func genSmallRatNonzeroPositive() *rapid.Generator[Rational] {
	return rapid.Custom(func(t *rapid.T) Rational {
		p := rapid.Int64Range(1, 9).Draw(t, "p")
		q := rapid.Int64Range(1, 9).Draw(t, "q")
		return mustRat(p, q)
	})
}

func genSmallValidULFT() *rapid.Generator[ULFT] {
	return rapid.Custom(func(t *rapid.T) ULFT {
		for {
			a := rapid.Int64Range(-5, 5).Draw(t, "a")
			b := rapid.Int64Range(-5, 5).Draw(t, "b")
			c := rapid.Int64Range(-5, 5).Draw(t, "c")
			d := rapid.Int64Range(-5, 5).Draw(t, "d")

			u := NewULFT(mustBig(a), mustBig(b), mustBig(c), mustBig(d))
			if err := u.Validate(); err != nil {
				continue
			}
			return u
		}
	})
}

func TestPBT_ULFTIngestGCF_RewriteLaw(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		base := genSmallValidULFT().Draw(t, "base")
		p := rapid.Int64Range(-5, 5).Draw(t, "p")
		q := genSmallPositiveInt64(1, 9).Draw(t, "q")
		xTail := genSmallRatNonzeroPositive().Draw(t, "xTail")

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
			t.Fatalf("got=%v want=%v base=%v p=%d q=%d xTail=%v x=%v", got, want, base, p, q, xTail, x)
		}
	})
}
