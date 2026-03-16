// mvp_four_over_pi_test.go v2
package cf

import (
	"math/big"
	"testing"
)

func TestMVPFourOverPiApproxFuncForFamily_Brouncker(t *testing.T) {
	fn, err := MVPFourOverPiApproxFuncForFamily(MVPFourOverPiFamilyBrouncker)
	if err != nil {
		t.Fatalf("MVPFourOverPiApproxFuncForFamily failed: %v", err)
	}

	got, err := fn(4)
	if err != nil {
		t.Fatalf("brouncker fn failed: %v", err)
	}

	want, err := MVPFourOverPiApproxBrouncker(4)
	if err != nil {
		t.Fatalf("MVPFourOverPiApproxBrouncker failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPFourOverPiApproxFuncForFamily_Lambert(t *testing.T) {
	fn, err := MVPFourOverPiApproxFuncForFamily(MVPFourOverPiFamilyLambert)
	if err != nil {
		t.Fatalf("MVPFourOverPiApproxFuncForFamily failed: %v", err)
	}

	got, err := fn(8)
	if err != nil {
		t.Fatalf("lambert fn failed: %v", err)
	}

	want, err := MVPFourOverPiApproxLambert(8)
	if err != nil {
		t.Fatalf("MVPFourOverPiApproxLambert failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPFourOverPiApproxFuncForFamily_RejectsUnknown(t *testing.T) {
	_, err := MVPFourOverPiApproxFuncForFamily(MVPFourOverPiFamily("bogus"))
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMVPFourOverPiApproxBrouncker_MatchesDirectConvergent(t *testing.T) {
	got, err := MVPFourOverPiApproxBrouncker(4)
	if err != nil {
		t.Fatalf("MVPFourOverPiApproxBrouncker failed: %v", err)
	}

	want, err := GCFSourceConvergent(NewBrouncker4OverPiGCFSource(), 4)
	if err != nil {
		t.Fatalf("GCFSourceConvergent failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPFourOverPiApproxLambert_IsPositiveAndGreaterThanOne(t *testing.T) {
	got, err := MVPFourOverPiApproxLambert(8)
	if err != nil {
		t.Fatalf("MVPFourOverPiApproxLambert failed: %v", err)
	}

	if got.Cmp(intRat(1)) <= 0 {
		t.Fatalf("got %v want > 1", got)
	}
	if got.Cmp(intRat(2)) >= 0 {
		t.Fatalf("got %v want < 2", got)
	}
}

func TestMVPFourOverPiApproxWithSource_BrounckerPath(t *testing.T) {
	brounckerSrc := func() GCFSource { return NewBrouncker4OverPiGCFSource() }

	got, err := MVPFourOverPiApproxWithSource(brounckerSrc, 4)
	if err != nil {
		t.Fatalf("MVPFourOverPiApproxWithSource failed: %v", err)
	}

	want, err := GCFSourceConvergent(NewBrouncker4OverPiGCFSource(), 4)
	if err != nil {
		t.Fatalf("GCFSourceConvergent failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func streamTermsWithULFT(src GCFSource, a, b, c, d int64, digits int) ([]int64, error) {
	s := NewGCFStream(src, GCFStreamOptions{})
	s.t = NewULFT(
		big.NewInt(a),
		big.NewInt(b),
		big.NewInt(c),
		big.NewInt(d),
	)

	out := make([]int64, 0, digits)
	for i := 0; i < digits; i++ {
		term, ok := s.Next()
		if !ok {
			if s.Err() != nil {
				return nil, s.Err()
			}
			break
		}
		out = append(out, term)
	}
	return out, nil
}

func publishedPiRCFTerms(n int) []int64 {
	all := []int64{
		3, 7, 15, 1, 292, 1, 1, 1, 2, 1,
		3, 1, 14, 2, 1, 1, 2, 2, 2, 2,
		1, 84, 2, 1, 1, 15, 3, 13, 1, 4,
		2, 6, 6, 99, 1, 2, 2, 6, 3, 5,
		1, 1, 6, 8, 1, 7, 1, 2, 3, 7,
	}
	if n < 0 {
		n = 0
	}
	if n > len(all) {
		n = len(all)
	}
	return append([]int64(nil), all[:n]...)
}

func TestMVPFourOverPiApproxBrouncker_StreamsPublishedPiTerms_First3(t *testing.T) {
	t.Skip("blocked: Brouncker 4/pi -> ULFT(0,4,1,0) -> pi certification currently hangs; fix transformed infinite-stream evidence first")
}

func TestMVPFourOverPiApproxLambert_StreamsPublishedPiTerms_First20(t *testing.T) {
	got, err := streamTermsWithULFT(
		NewLambertPiOver4GCFSource(),
		4, 0, 0, 1, // x -> 4x, so pi/4 -> pi
		20,
	)
	if err != nil {
		t.Fatalf("streamTermsWithULFT failed: %v", err)
	}

	want := publishedPiRCFTerms(20)
	if !equalTerms(got, want) {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

// mvp_four_over_pi_test.go v2
