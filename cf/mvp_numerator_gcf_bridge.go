// mvp_numerator_gcf_bridge.go v3
package cf

import "fmt"

// MVPThreeOverPiSquaredPlusEAsGCFSource is the legacy name for the current
// finite bridge source for:
//
//	3/pi^2 + e
//
// Deprecated MVP note:
//   - this helper adapts a bounded rational approximation into a regular CF and
//     then into a GCF source
//   - prefer MVPThreeOverPiSquaredPlusEFiniteBridgeSource for new code
func MVPThreeOverPiSquaredPlusEAsGCFSource(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
) (GCFSource, error) {
	return MVPThreeOverPiSquaredPlusEFiniteBridgeSource(
		fourOverPiPrefixTerms,
		ePrefixTerms,
	)
}

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

// Legacy finite bridge source retained for compatibility/tests.
func MVPThreeOverPiSquaredPlusEFiniteBridgeSource(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
) (GCFSource, error) {
	return MVPThreeOverPiSquaredPlusEFiniteBridgeSourceWithFourOverPiApprox(
		MVPDefaultFourOverPiApproxFunc(),
		fourOverPiPrefixTerms,
		ePrefixTerms,
	)
}

// Legacy finite bridge snapshot retained for compatibility/tests.
func MVPThreeOverPiSquaredPlusEFiniteBridgeSnapshot(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	bridgeTerms int,
) (GCFApprox, error) {
	if bridgeTerms <= 0 {
		return GCFApprox{}, fmt.Errorf(
			"MVPThreeOverPiSquaredPlusEFiniteBridgeSnapshot: bridgeTerms must be > 0, got %d",
			bridgeTerms,
		)
	}

	src, err := MVPThreeOverPiSquaredPlusEFiniteBridgeSource(
		fourOverPiPrefixTerms,
		ePrefixTerms,
	)
	if err != nil {
		return GCFApprox{}, err
	}

	return GCFApproxFromPrefix(src, bridgeTerms)
}

// Legacy family-parameterized finite bridge retained for compatibility/tests.
func MVPThreeOverPiSquaredPlusEFiniteBridgeSourceWithFourOverPiApprox(
	fourOverPiFn MVPFourOverPiApproxFunc,
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
) (GCFSource, error) {
	x, err := MVPThreeOverPiSquaredPlusEApproxWithFourOverPiApprox(
		fourOverPiFn,
		fourOverPiPrefixTerms,
		ePrefixTerms,
	)
	if err != nil {
		return nil, err
	}
	return AdaptCFToGCF(NewRationalCF(x)), nil
}

// Legacy family-parameterized finite bridge snapshot retained for compatibility/tests.
func MVPThreeOverPiSquaredPlusEFiniteBridgeSnapshotWithFourOverPiApprox(
	fourOverPiFn MVPFourOverPiApproxFunc,
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	bridgeTerms int,
) (GCFApprox, error) {
	if bridgeTerms <= 0 {
		return GCFApprox{}, fmt.Errorf(
			"MVPThreeOverPiSquaredPlusEFiniteBridgeSnapshotWithFourOverPiApprox: bridgeTerms must be > 0, got %d",
			bridgeTerms,
		)
	}

	src, err := MVPThreeOverPiSquaredPlusEFiniteBridgeSourceWithFourOverPiApprox(
		fourOverPiFn,
		fourOverPiPrefixTerms,
		ePrefixTerms,
	)
	if err != nil {
		return GCFApprox{}, err
	}

	return GCFApproxFromPrefix(src, bridgeTerms)
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

func MVPThreeOverPiSquaredPlusERadicandSource(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
) (GCFSource, error) {
	return MVPThreeOverPiSquaredPlusEFiniteBridgeSource(
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

	return MVPThreeOverPiSquaredPlusERadicandApproxSnapshotWithFourOverPiApprox(
		MVPDefaultFourOverPiApproxFunc(),
		fourOverPiPrefixTerms,
		ePrefixTerms,
	)
}

// mvp_numerator_gcf_bridge.go v3
