// infinite_gcf_contract_test.go v6
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

func TestInfiniteGCFContract_CurrentNumeratorSnapshotPath_IsAvailable(t *testing.T) {
	got, err := MVPNumeratorRadicandApproxSnapshot(
		MVPDefaultFourOverPiPrefixTerms,
		MVPDefaultEPrefixTerms,
		MVPNumeratorBridgePrefixTerms,
	)
	if err != nil {
		t.Fatalf("MVPNumeratorRadicandApproxSnapshot failed: %v", err)
	}
	if got.Convergent.Cmp(intRat(0)) <= 0 {
		t.Fatalf("got %v want positive convergent", got.Convergent)
	}
}

func TestInfiniteGCFContract_CurrentNumeratorBridge_BudgetIsExplicit(t *testing.T) {
	if MVPNumeratorBridgePrefixTerms <= 0 {
		t.Fatalf("got MVPNumeratorBridgePrefixTerms=%d want > 0", MVPNumeratorBridgePrefixTerms)
	}
}

func TestInfiniteGCFContract_CurrentMVPTargetStillWorksDespiteExceptions(t *testing.T) {
	got, err := mvpTestTargetBoundsDefault()
	if err != nil {
		t.Fatalf("mvpTestTargetBoundsDefault failed: %v", err)
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
// 3. The numerator live path now uses snapshot assembly rather than a finite bridge source.
// 4. Post-MVP goal: retire the remaining compatibility wrappers and bridge-budget naming.
//
// infinite_gcf_contract_test.go v6
