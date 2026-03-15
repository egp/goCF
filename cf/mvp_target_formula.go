// mvp_target_formula.go v1
package cf

import "fmt"

// MVPTargetBounds returns a conservative inside range for:
//
//	sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))
//
// Current MVP construction:
//   - numerator is a bounded rational point approximation
//   - denominator is a certified positive inside range
//   - quotient is enclosed by dividing the numerator point by the denominator range
func MVPTargetBounds(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	sqrtPolicy SqrtPolicy2,
	angle Angle,
) (Range, error) {
	num, err := MVPNumeratorApprox(fourOverPiPrefixTerms, ePrefixTerms, sqrtPolicy)
	if err != nil {
		return Range{}, err
	}

	den, err := MVPDenominatorBounds(sqrtPolicy, angle)
	if err != nil {
		return Range{}, err
	}
	if den.Contains(intRat(0)) {
		return Range{}, fmt.Errorf("MVPTargetBounds: denominator range contains 0: %v", den)
	}
	if den.Lo.Cmp(intRat(0)) <= 0 || den.Hi.Cmp(intRat(0)) <= 0 {
		return Range{}, fmt.Errorf("MVPTargetBounds: denominator range must be strictly positive: %v", den)
	}

	numRange := NewRange(num, num, true, true)
	recipDen, err := ReciprocalRangeConservative(den)
	if err != nil {
		return Range{}, err
	}

	lo, err := numRange.Lo.Mul(recipDen.Lo)
	if err != nil {
		return Range{}, err
	}
	hi, err := numRange.Hi.Mul(recipDen.Hi)
	if err != nil {
		return Range{}, err
	}

	return NewRange(lo, hi, true, true), nil
}

func MVPTargetBoundsDefault() (Range, error) {
	return MVPTargetBounds(
		4,
		6,
		DefaultSqrtPolicy2(),
		Degrees(mustRat(69, 1)),
	)
}

// MVPTargetApprox returns a rational point only when MVPTargetBounds collapses
// to an exact point; otherwise it reports that only a bounded non-point result
// is available at the current MVP stage.
func MVPTargetApprox(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	sqrtPolicy SqrtPolicy2,
	angle Angle,
) (Rational, error) {
	r, err := MVPTargetBounds(fourOverPiPrefixTerms, ePrefixTerms, sqrtPolicy, angle)
	if err != nil {
		return Rational{}, err
	}
	if r.Lo.Cmp(r.Hi) != 0 {
		return Rational{}, fmt.Errorf("MVPTargetApprox: bounded non-point result")
	}
	return r.Lo, nil
}

func MVPTargetApproxDefault() (Rational, error) {
	return MVPTargetApprox(
		4,
		6,
		DefaultSqrtPolicy2(),
		Degrees(mustRat(69, 1)),
	)
}

// mvp_target_formula.go v1
