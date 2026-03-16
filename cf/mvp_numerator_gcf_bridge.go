// mvp_numerator_gcf_bridge.go v2
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
// finite bridge snapshot for:
//
//	3/pi^2 + e
//
// Deprecated MVP note:
//   - prefer MVPThreeOverPiSquaredPlusEFiniteBridgeSnapshot for new code
func MVPThreeOverPiSquaredPlusEApproxSnapshot(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	bridgeTerms int,
) (GCFApprox, error) {
	return MVPThreeOverPiSquaredPlusEFiniteBridgeSnapshot(
		fourOverPiPrefixTerms,
		ePrefixTerms,
		bridgeTerms,
	)
}

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

func MVPThreeOverPiSquaredPlusEFiniteBridgeSnapshot(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	bridgeTerms int,
) (GCFApprox, error) {
	return MVPThreeOverPiSquaredPlusEFiniteBridgeSnapshotWithFourOverPiApprox(
		MVPDefaultFourOverPiApproxFunc(),
		fourOverPiPrefixTerms,
		ePrefixTerms,
		bridgeTerms,
	)
}

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
	return MVPThreeOverPiSquaredPlusEFiniteBridgeSnapshot(
		fourOverPiPrefixTerms,
		ePrefixTerms,
		bridgeTerms,
	)
}

// mvp_numerator_gcf_bridge.go v2
