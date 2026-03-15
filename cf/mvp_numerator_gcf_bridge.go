// mvp_numerator_gcf_bridge.go v1
package cf

// MVPThreeOverPiSquaredPlusEAsGCFSource adapts the current bounded rational
// approximation of
//
//	3/pi^2 + e
//
// into a regular CF and then into a GCF source, so unary GCF-ingesting entry
// points can operate on it.
func MVPThreeOverPiSquaredPlusEAsGCFSource(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
) (GCFSource, error) {
	x, err := MVPThreeOverPiSquaredPlusEApprox(fourOverPiPrefixTerms, ePrefixTerms)
	if err != nil {
		return nil, err
	}
	return AdaptCFToGCF(NewRationalCF(x)), nil
}

// mvp_numerator_gcf_bridge.go v1
