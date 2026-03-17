// mvp_numerator_gcf_bridge.go v4
package cf

import "fmt"

// MVPThreeOverPiSquaredPlusEApproxSnapshot is the legacy name for the current
// radicand snapshot for:
//
//	3/pi^2 + e
//
// Deprecated MVP note:
//   - prefer MVPThreeOverPiSquaredPlusERadicandApproxSnapshotWithFourOverPiApprox
//     or MVPNumeratorRadicandApproxSnapshot for new code
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

// New production radicand snapshot path: build a direct point snapshot instead
// of crossing the temporary finite bridge boundary.
func MVPThreeOverPiSquaredPlusERadicandApproxSnapshotWithFourOverPiApprox(
	fourOverPiFn MVPFourOverPiApproxFunc,
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
) (GCFApprox, error) {
	fourOverPi, err := MVPApproxSnapshotFromApproxFunc(
		fourOverPiFn,
		fourOverPiPrefixTerms,
	)
	if err != nil {
		return GCFApprox{}, err
	}

	eApprox, err := MVPRadicandDefaultEApproxSnapshot(ePrefixTerms)
	if err != nil {
		return GCFApprox{}, err
	}

	return MVPRadicandAssembleFromSnapshots(fourOverPi, eApprox)
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

	return MVPThreeOverPiSquaredPlusERadicandApproxSnapshotWithFourOverPiApprox(
		MVPDefaultFourOverPiApproxFunc(),
		fourOverPiPrefixTerms,
		ePrefixTerms,
	)
}

// mvp_numerator_gcf_bridge.go v4
