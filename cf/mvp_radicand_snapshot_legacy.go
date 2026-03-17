// mvp_radicand_snapshot_legacy.go v3
package cf

import "fmt"

// Legacy expression-specific wrapper retained temporarily while callers migrate.
func MVPThreeOverPiSquaredPlusEApproxSnapshot(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	bridgeTerms int,
) (GCFApprox, error) {
	return MVPThreeOverPiSquaredPlusERadicandSnapshot(
		fourOverPiPrefixTerms,
		ePrefixTerms,
		bridgeTerms,
	)
}

// Legacy expression-specific wrapper retained temporarily while callers migrate.
func MVPThreeOverPiSquaredPlusERadicandApproxSnapshotWithFourOverPiApprox(
	fourOverPiFn MVPFourOverPiApproxFunc,
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
) (GCFApprox, error) {
	return MVPRadicandAssembleSnapshotWithFourOverPiApprox(
		fourOverPiFn,
		fourOverPiPrefixTerms,
		ePrefixTerms,
	)
}

// Legacy expression-specific wrapper retained temporarily while callers migrate.
func MVPThreeOverPiSquaredPlusERadicandSnapshot(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	bridgeTerms int,
) (GCFApprox, error) {
	if bridgeTerms <= 0 {
		return GCFApprox{}, fmt.Errorf(
			"MVPThreeOverPiSquaredPlusERadicandSnapshot: bridgeTerms must be > 0, got %d",
			bridgeTerms,
		)
	}

	return MVPRadicandAssembleSnapshot(
		fourOverPiPrefixTerms,
		ePrefixTerms,
		bridgeTerms,
	)
}

// mvp_radicand_snapshot_legacy.go v3
