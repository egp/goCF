// mvp_radicand_snapshot_legacy.go v1
package cf

import "fmt"

// MVPThreeOverPiSquaredPlusEApproxSnapshot is the legacy name for the current
// radicand snapshot for:
//
//	3/pi^2 + e
//
// Deprecated MVP note:
//   - prefer MVPThreeOverPiSquaredPlusERadicandApproxSnapshotWithFourOverPiApprox
//     or MVPRadicandSnapshot for new code
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

// Production radicand snapshot path: build a direct point snapshot instead
// of crossing the former finite bridge boundary.

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
	)
}

// mvp_radicand_snapshot_legacy.go v1
