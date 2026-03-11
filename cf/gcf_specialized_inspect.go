// gcf_specialized_inspect.go v1
package cf

// specializedInspectGCFSource is the inspection analogue of
// specializedGCFApproxFromPrefix. It forms a specialized bounded-prefix GCF
// snapshot and returns that snapshot together with up to digits regular CF terms
// of the exact rational convergent.
func specializedInspectGCFSource(
	prefixTerms int,
	digits int,
	approxFn func(prefixTerms int) (GCFApprox, error),
) (GCFInspect, error) {
	a, err := approxFn(prefixTerms)
	if err != nil {
		return GCFInspect{}, err
	}
	return gcfInspectFromApprox(a, digits, "specializedInspectGCFSource:")
}

// InspectLambertPiOver4Prefix returns a specialized bounded-prefix Lambert pi/4
// snapshot together with up to digits regular CF terms of its exact rational
// convergent.
func InspectLambertPiOver4Prefix(prefixTerms int, digits int) (GCFInspect, error) {
	return specializedInspectGCFSource(prefixTerms, digits, LambertPiOver4ApproxFromPrefix)
}

// InspectBrouncker4OverPiPrefix returns a specialized bounded-prefix Brouncker
// 4/pi snapshot together with up to digits regular CF terms of its exact
// rational convergent.
func InspectBrouncker4OverPiPrefix(prefixTerms int, digits int) (GCFInspect, error) {
	return specializedInspectGCFSource(prefixTerms, digits, Brouncker4OverPiApproxFromPrefix)
}

// gcf_specialized_inspect.go v1
