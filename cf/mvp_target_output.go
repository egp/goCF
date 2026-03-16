// mvp_target_output.go v1
package cf

import "fmt"

func mvpRangeMidpoint(r Range) (Rational, error) {
	sum, err := r.Lo.Add(r.Hi)
	if err != nil {
		return Rational{}, err
	}
	return sum.Div(mustRat(2, 1))
}

// MVPTargetMidpointApprox returns the midpoint rational of the current bounded
// MVP target range:
//
//	sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))
//
// This is an explicit approximation helper, not a proof that the target range
// has collapsed to an exact point.
func MVPTargetMidpointApprox(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	sqrtPolicy SqrtPolicy2,
	angle Angle,
) (Rational, error) {
	r, err := MVPTargetBounds(
		fourOverPiPrefixTerms,
		ePrefixTerms,
		sqrtPolicy,
		angle,
	)
	if err != nil {
		return Rational{}, err
	}
	return mvpRangeMidpoint(r)
}

func MVPTargetMidpointApproxDefault() (Rational, error) {
	return MVPTargetMidpointApprox(
		MVPDefaultFourOverPiPrefixTerms,
		MVPDefaultEPrefixTerms,
		DefaultSqrtPolicy2(),
		Degrees(mustRat(69, 1)),
	)
}

// MVPTargetMidpointApproxCF returns a regular continued fraction for the target
// midpoint approximation.
func MVPTargetMidpointApproxCF(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	sqrtPolicy SqrtPolicy2,
	angle Angle,
) (ContinuedFraction, error) {
	approx, err := MVPTargetMidpointApprox(
		fourOverPiPrefixTerms,
		ePrefixTerms,
		sqrtPolicy,
		angle,
	)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

func MVPTargetMidpointApproxCFDefault() (ContinuedFraction, error) {
	return MVPTargetMidpointApproxCF(
		MVPDefaultFourOverPiPrefixTerms,
		MVPDefaultEPrefixTerms,
		DefaultSqrtPolicy2(),
		Degrees(mustRat(69, 1)),
	)
}

// MVPTargetMidpointApproxTerms returns up to digits regular CF terms for the
// target midpoint approximation.
func MVPTargetMidpointApproxTerms(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	sqrtPolicy SqrtPolicy2,
	angle Angle,
	digits int,
) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("MVPTargetMidpointApproxTerms: negative digits %d", digits)
	}

	cf, err := MVPTargetMidpointApproxCF(
		fourOverPiPrefixTerms,
		ePrefixTerms,
		sqrtPolicy,
		angle,
	)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

func MVPTargetMidpointApproxTermsDefault(digits int) ([]int64, error) {
	return MVPTargetMidpointApproxTerms(
		MVPDefaultFourOverPiPrefixTerms,
		MVPDefaultEPrefixTerms,
		DefaultSqrtPolicy2(),
		Degrees(mustRat(69, 1)),
		digits,
	)
}

// mvp_target_output.go v1
