// infinite_gcf_contract_test.go v4
package cf

import (
	"strings"
	"testing"
)

func TestInfiniteGCFContract_CanonicalAlgorithmicSourcesDoNotExhaustEarly(t *testing.T) {
	cases := []struct {
		name string
		src  GCFSource
		n    int
	}{
		{"brouncker_4_over_pi", MVPReciprocalPiGCFSource(), 8},
		{"e_source", MVPEGCFSource(), 8},
		{"sqrt5_adapted_cf", AdaptCFToGCF(Sqrt5CF()), 8},
		{"sqrt2_adapted_cf", AdaptCFToGCF(Sqrt2CF()), 8},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_ = mustReadNPQWithoutExhaustion(t, tc.src, tc.n)
		})
	}
}

func TestInfiniteGCFContract_SinPrefixEntryRejectsNonExactInfiniteAngle(t *testing.T) {
	// Prefix 68 + 1/x with an infinite tail is not an exact finite angle.
	src := NewPeriodicGCF(
		[][2]int64{{68, 1}},
		[][2]int64{{1, 1}},
	)

	_, err := SinBoundsDegreesFromGCFPrefix2(src, 1)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "angle not exact") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInfiniteGCFContract_CurrentAngleConstant_IsExactFiniteGCFPrefix(t *testing.T) {
	got, err := GCFSourceConvergent(MVP69DegreeGCFSource(), 2)
	if err != nil {
		t.Fatalf("GCFSourceConvergent failed: %v", err)
	}

	want := mustRat(69, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestInfiniteGCFContract_CurrentNumeratorBridge_IsFiniteByDesign(t *testing.T) {
	src, err := MVPThreeOverPiSquaredPlusEAsGCFSource(
		MVPDefaultFourOverPiPrefixTerms,
		MVPDefaultEPrefixTerms,
	)
	if err != nil {
		t.Fatalf("MVPThreeOverPiSquaredPlusEAsGCFSource failed: %v", err)
	}

	_, exhausted := readUntilExhaustionPQ(src, 128)
	if !exhausted {
		t.Fatalf("expected current numerator bridge to be finite within 128 terms")
	}
}

func TestInfiniteGCFContract_CurrentNumeratorBridge_BudgetIsExplicit(t *testing.T) {
	if MVPNumeratorBridgePrefixTerms <= 0 {
		t.Fatalf("got MVPNumeratorBridgePrefixTerms=%d want > 0", MVPNumeratorBridgePrefixTerms)
	}
}

func TestInfiniteGCFContract_CurrentMVPTargetStillWorksDespiteExceptions(t *testing.T) {
	got, err := MVPTargetBoundsDefault()
	if err != nil {
		t.Fatalf("MVPTargetBoundsDefault failed: %v", err)
	}
	if !got.IsInside() {
		t.Fatalf("got %v want inside range", got)
	}
	if got.Contains(intRat(0)) {
		t.Fatalf("got %v want zero excluded", got)
	}
}

// This file documents the current state:
//
// 1. Canonical mathematical sources are infinite/algorithmic.
// 2. Some MVP helpers still rely on explicit finite/exact-tail exceptions.
// 3. The numerator path currently uses an explicit temporary finite bridge and
//    an explicit bridge-prefix budget.
// 4. Post-MVP goal: retire those exceptions and move prod ingestion toward
//    infinite-GCF-only operator plumbing.
//
// infinite_gcf_contract_test.go v4
