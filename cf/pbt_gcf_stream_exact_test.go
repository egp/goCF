// pbt_gcf_stream_exact_test.go v1
package cf

import (
	"testing"

	"pgregory.net/rapid"
)

func genSmallFiniteGCF() *rapid.Generator[*SliceGCF] {
	return rapid.Custom(func(t *rapid.T) *SliceGCF {
		n := rapid.IntRange(0, 4).Draw(t, "n")
		terms := make([][2]int64, 0, n)
		for i := 0; i < n; i++ {
			p := rapid.Int64Range(-5, 5).Draw(t, "p")
			q := rapid.Int64Range(1, 9).Draw(t, "q")
			terms = append(terms, [2]int64{p, q})
		}
		return NewSliceGCF(terms...)
	})
}

func genSmallTailRat() *rapid.Generator[Rational] {
	return rapid.Custom(func(t *rapid.T) Rational {
		p := rapid.Int64Range(1, 9).Draw(t, "p")
		q := rapid.Int64Range(1, 9).Draw(t, "q")
		return mustRat(p, q)
	})
}

func genSmallULFTForStreams() *rapid.Generator[ULFT] {
	return rapid.Custom(func(t *rapid.T) ULFT {
		for {
			u := NewULFT(
				mustBig(rapid.Int64Range(-5, 5).Draw(t, "A")),
				mustBig(rapid.Int64Range(-5, 5).Draw(t, "B")),
				mustBig(rapid.Int64Range(-5, 5).Draw(t, "C")),
				mustBig(rapid.Int64Range(-5, 5).Draw(t, "D")),
			)
			if err := u.Validate(); err != nil {
				continue
			}
			return u
		}
	})
}

func genSmallDiagBLFTForStreams() *rapid.Generator[DiagBLFT] {
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
func TestPBT_GCFULFTStreamWithTail_MatchesExactRationalImage(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		u := genSmallULFTForStreams().Draw(t, "u")
		srcTerms := genSmallFiniteGCF().Draw(t, "srcTerms")
		tail := genSmallTailRat().Draw(t, "tail")

		wantRat, _, err := ApplyComposedGCFULFTToTailExact(
			u,
			NewSliceGCF(srcTerms.terms...),
			tail,
			16,
		)
		if err != nil {
			t.Skip()
		}
		want := collectTermsBounded(NewRationalCF(wantRat), 64)

		s := NewGCFULFTStreamWithTail(
			u,
			NewSliceGCF(srcTerms.terms...),
			tail,
			GCFULFTStreamOptions{
				MaxIngestTerms: 16,
			},
		)

		got := collectTermsBounded(s, 64)
		if err := s.Err(); err != nil {
			t.Fatalf("GCFULFTStream error: %v", err)
		}

		if len(got) != len(want) {
			t.Fatalf("len mismatch: got=%v want=%v u=%v tail=%v terms=%v", got, want, u, tail, srcTerms.terms)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("digit mismatch at %d: got=%v want=%v u=%v tail=%v terms=%v", i, got, want, u, tail, srcTerms.terms)
			}
		}
	})
}
func TestPBT_GCFDiagStreamWithTail_MatchesExactRationalImage(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		d := genSmallDiagBLFTForStreams().Draw(t, "d")
		srcTerms := genSmallFiniteGCF().Draw(t, "srcTerms")
		tail := genSmallTailRat().Draw(t, "tail")

		wantRat, _, err := ApplyComposedGCFDiagBLFTToTailExact(
			d,
			NewSliceGCF(srcTerms.terms...),
			tail,
			16,
		)
		if err != nil {
			t.Skip()
		}
		want := collectTermsBounded(NewRationalCF(wantRat), 64)

		s := NewGCFDiagStreamWithTail(
			d,
			NewSliceGCF(srcTerms.terms...),
			tail,
			GCFULFTStreamOptions{
				MaxIngestTerms: 16,
			},
		)

		got := collectTermsBounded(s, 64)
		if err := s.Err(); err != nil {
			t.Fatalf("GCFDiagStream error: %v", err)
		}

		if len(got) != len(want) {
			t.Fatalf("len mismatch: got=%v want=%v d=%v tail=%v terms=%v", got, want, d, tail, srcTerms.terms)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("digit mismatch at %d: got=%v want=%v d=%v tail=%v terms=%v", i, got, want, d, tail, srcTerms.terms)
			}
		}
	})
}

func collectTermsBounded(cf ContinuedFraction, max int) []int64 {
	var out []int64
	for i := 0; i < max; i++ {
		d, ok := cf.Next()
		if !ok {
			return out
		}
		out = append(out, d)
	}
	return out
}

// pbt_gcf_stream_exact_test.go v1
