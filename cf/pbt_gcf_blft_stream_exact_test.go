// pbt_gcf_blft_stream_exact_test.go v1
package cf

import (
	"testing"

	"pgregory.net/rapid"
)

func genSmallBLFTForStreams() *rapid.Generator[BLFT] {
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

func cloneSliceGCF(src *SliceGCF) *SliceGCF {
	return NewSliceGCF(src.terms...)
}

func TestPBT_GCFBLFTStreamWithTails_MatchesExactRationalImage(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		b := genSmallBLFTForStreams().Draw(t, "b")
		xTerms := genSmallFiniteGCF().Draw(t, "xTerms")
		yTerms := genSmallFiniteGCF().Draw(t, "yTerms")
		xTail := genSmallTailRat().Draw(t, "xTail")
		yTail := genSmallTailRat().Draw(t, "yTail")

		wantRat, _, _, err := ApplyComposedGCFXYBLFTToTailsExact(
			b,
			cloneSliceGCF(xTerms),
			xTail,
			16,
			cloneSliceGCF(yTerms),
			yTail,
			16,
		)
		if err != nil {
			t.Skip()
		}
		want := collectTermsBounded(NewRationalCF(wantRat), 64)

		s := NewGCFBLFTStreamWithTails(
			b,
			cloneSliceGCF(xTerms),
			xTail,
			cloneSliceGCF(yTerms),
			yTail,
			GCFULFTStreamOptions{
				MaxIngestTerms: 16,
			},
		)

		got := collectTermsBounded(s, 64)
		if err := s.Err(); err != nil {
			t.Fatalf("GCFBLFTStream error: %v", err)
		}

		if len(got) != len(want) {
			t.Fatalf("len mismatch: got=%v want=%v b=%v xTail=%v yTail=%v xTerms=%v yTerms=%v",
				got, want, b, xTail, yTail, xTerms.terms, yTerms.terms)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("digit mismatch at %d: got=%v want=%v b=%v xTail=%v yTail=%v xTerms=%v yTerms=%v",
					i, got, want, b, xTail, yTail, xTerms.terms, yTerms.terms)
			}
		}
	})
}
