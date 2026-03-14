// sqrt_range.go v2
package cf

// SqrtRangeExact returns the exact monotone image of a nonnegative inside range
// under sqrt(x), when both endpoints have exact rational square roots.
//
// Returns:
//   - (out, true, nil)  if exact endpoint square roots exist
//   - (_,   false, nil) if not yet supported exactly
//   - (_,   false, err) on invalid input
func SqrtRangeExact(r Range) (Range, bool, error) {
	return SqrtRangeExact2(r)
}

// SqrtRangeExactFromCFApprox applies SqrtRangeExact to the enclosure carried by CFApprox.
func SqrtRangeExactFromCFApprox(a CFApprox) (Range, bool, error) {
	return SqrtRangeExactFromCFApprox2(a)
}

// SqrtRangeHeuristic returns a bounded rational approximation to the monotone
// image of a nonnegative inside range under sqrt(x).
//
// IMPORTANT:
//   - this is a heuristic helper
//   - it is not yet a proof-safe conservative enclosure for irrational endpoints
//   - it should not be used anywhere that requires a formally guaranteed range
//
// Current policy:
//   - first try SqrtRangeExact
//   - otherwise approximate sqrt(lo) and sqrt(hi) independently using SqrtApprox
//   - preserve endpoint inclusions
func SqrtRangeHeuristic(r Range) (Range, error) {
	return SqrtRangeHeuristic2(r)
}

// SqrtRangeHeuristicFromCFApprox applies SqrtRangeHeuristic to the enclosure
// carried by CFApprox.
func SqrtRangeHeuristicFromCFApprox(a CFApprox) (Range, error) {
	return SqrtRangeHeuristicFromCFApprox2(a)
}

// sqrt_range.go v2
