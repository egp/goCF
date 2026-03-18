// cf/sqrt_gcf_metadata_boundary_test.go v3
package cf

import (
	"testing"
)

type testNativeQuadraticRadicalGCFSource struct {
	terms    [][2]int64
	i        int
	radicand int64
}

func newTestNativeQuadraticRadicalGCFSource(radicand int64, terms ...[2]int64) *testNativeQuadraticRadicalGCFSource {
	cp := append([][2]int64(nil), terms...)
	return &testNativeQuadraticRadicalGCFSource{
		terms:    cp,
		radicand: radicand,
	}
}

func (s *testNativeQuadraticRadicalGCFSource) NextPQ() (int64, int64, bool) {
	if s.i >= len(s.terms) {
		return 0, 0, false
	}
	t := s.terms[s.i]
	s.i++
	return t[0], t[1], true
}

func (s *testNativeQuadraticRadicalGCFSource) Radicand() (int64, bool) {
	return s.radicand, true
}

func TestSqrtGCFMetadataBoundary_TanhAlreadyAcceptsNativeGCFRadicalMetadata(t *testing.T) {
	src := newTestNativeQuadraticRadicalGCFSource(
		5,
		[2]int64{2, 1},
		[2]int64{4, 1},
	)

	got, err := TanhBoundsSpecialFromGCF2(src)
	if err != nil {
		t.Fatalf("TanhBoundsSpecialFromGCF2 error: %v", err)
	}

	want := TanhBoundsSqrt5()

	if got.Lo.Cmp(want.Lo) != 0 || got.Hi.Cmp(want.Hi) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
	if got.IncLo != want.IncLo || got.IncHi != want.IncHi {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSqrtGCFMetadataBoundary_NonSquareNativeMetadataDoesNotForceExactRoot(t *testing.T) {
	src := newTestNativeQuadraticRadicalGCFSource(
		5,
		[2]int64{2, 1},
		[2]int64{4, 1},
	)

	_, ok, err := sqrtExactRootCFViaSourceMetadata(src)
	if err != nil {
		t.Fatalf("sqrtExactRootCFViaSourceMetadata error: %v", err)
	}
	if ok {
		t.Fatalf("did not expect exact-root fast path for non-square radicand metadata")
	}
}

// cf/sqrt_gcf_metadata_boundary_test.go v3
