// mvp_numerator_test.go v3
package cf

import "testing"

func TestMVPNumeratorApprox_RejectsBadBounds(t *testing.T) {
	if _, err := MVPNumeratorApproxDefault(0, 6); err == nil {
		t.Fatalf("expected error for fourOverPiPrefixTerms=0")
	}
	if _, err := MVPNumeratorApproxDefault(4, 0); err == nil {
		t.Fatalf("expected error for ePrefixTerms=0")
	}
}

func TestMVPThreeOverPiSquaredPlusEAsGCFSource_IsUsableByGCFUnaryPath(t *testing.T) {
	src, err := MVPThreeOverPiSquaredPlusEAsGCFSource(4, 6)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEAsGCFSource failed: %v", err)
	}

	got, err := SqrtApproxFromGCFSourceRangeSeed2(src, 64, DefaultSqrtPolicy2())
	if err != nil {
		t.Fatalf("SqrtApproxFromGCFSourceRangeSeed2 failed: %v", err)
	}

	want, err := MVPNumeratorApproxDefault(4, 6)
	if err != nil {
		t.Fatalf("MVPNumeratorApproxDefault failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPNumeratorApprox_UsesGCFIngestingUnarySqrtPath(t *testing.T) {
	got, err := MVPNumeratorApproxDefault(4, 6)
	if err != nil {
		t.Fatalf("MVPNumeratorApproxDefault failed: %v", err)
	}

	x, err := MVPThreeOverPiSquaredPlusEApprox(4, 6)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEApprox failed: %v", err)
	}

	src := AdaptCFToGCF(NewRationalCF(x))
	want, err := SqrtApproxFromGCFSourceRangeSeed2(src, 64, DefaultSqrtPolicy2())
	if err != nil {
		t.Fatalf("SqrtApproxFromGCFSourceRangeSeed2 failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPNumeratorApprox_CurrentDefaultUsesSharperBudgets(t *testing.T) {
	got, err := MVPNumeratorApproxCurrentDefault()
	if err != nil {
		t.Fatalf("MVPNumeratorApproxCurrentDefault failed: %v", err)
	}

	want, err := MVPNumeratorApproxDefault(
		MVPDefaultFourOverPiPrefixTerms,
		MVPDefaultEPrefixTerms,
	)
	if err != nil {
		t.Fatalf("MVPNumeratorApproxDefault failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPNumeratorApproxCF_MatchesApprox(t *testing.T) {
	gotTerms, err := MVPNumeratorApproxTermsDefault(4, 6, 12)
	if err != nil {
		t.Fatalf("MVPNumeratorApproxTermsDefault failed: %v", err)
	}

	gotApprox, err := MVPNumeratorApproxDefault(4, 6)
	if err != nil {
		t.Fatalf("MVPNumeratorApproxDefault failed: %v", err)
	}
	wantTerms := collectTerms(NewRationalCF(gotApprox), 12)

	if len(gotTerms) != len(wantTerms) {
		t.Fatalf("len mismatch: got=%v want=%v", gotTerms, wantTerms)
	}
	for i := range wantTerms {
		if gotTerms[i] != wantTerms[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, gotTerms, wantTerms)
		}
	}
}

func TestMVPNumeratorApprox_IsPositive(t *testing.T) {
	got, err := MVPNumeratorApproxCurrentDefault()
	if err != nil {
		t.Fatalf("MVPNumeratorApproxCurrentDefault failed: %v", err)
	}
	if got.Cmp(intRat(0)) <= 0 {
		t.Fatalf("got %v want positive", got)
	}
}

func TestMVPNumeratorApprox_ExceedsOne(t *testing.T) {
	got, err := MVPNumeratorApproxCurrentDefault()
	if err != nil {
		t.Fatalf("MVPNumeratorApproxCurrentDefault failed: %v", err)
	}
	if got.Cmp(intRat(1)) <= 0 {
		t.Fatalf("got %v want > 1", got)
	}
}

// Full target formula intentionally stays in test code for now.
// This test fixes only the numerator shape:
//
//	sqrt(3/pi^2 + e)
//
// mvp_numerator_test.go v3
