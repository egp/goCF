// mvp_numerator.go v1
package cf

import "fmt"

// MVPNumeratorApprox returns a bounded rational approximation for:
//
//	sqrt(3/pi^2 + e)
//
// Current MVP construction:
//   - choose canonical reciprocal-pi source via MVPReciprocalPiGCFSource()
//   - choose canonical e source via MVPEGCFSource()
//   - form bounded rational approximation to 3/pi^2 + e
//   - apply bounded sqrt approximation under the supplied policy
func MVPNumeratorApprox(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	sqrtPolicy SqrtPolicy2,
) (Rational, error) {
	x, err := MVPThreeOverPiSquaredPlusEApprox(fourOverPiPrefixTerms, ePrefixTerms)
	if err != nil {
		return Rational{}, err
	}
	return SqrtApproxWithPolicy2(x, sqrtPolicy)
}

// MVPNumeratorApproxDefault uses the default sqrt policy.
func MVPNumeratorApproxDefault(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
) (Rational, error) {
	return MVPNumeratorApprox(
		fourOverPiPrefixTerms,
		ePrefixTerms,
		DefaultSqrtPolicy2(),
	)
}

// MVPNumeratorApproxCF returns a ContinuedFraction for the bounded numerator
// approximation.
func MVPNumeratorApproxCF(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	sqrtPolicy SqrtPolicy2,
) (ContinuedFraction, error) {
	approx, err := MVPNumeratorApprox(fourOverPiPrefixTerms, ePrefixTerms, sqrtPolicy)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

// MVPNumeratorApproxCFDefault uses the default sqrt policy.
func MVPNumeratorApproxCFDefault(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
) (ContinuedFraction, error) {
	return MVPNumeratorApproxCF(
		fourOverPiPrefixTerms,
		ePrefixTerms,
		DefaultSqrtPolicy2(),
	)
}

// MVPNumeratorApproxTerms returns up to digits CF terms for the bounded numerator
// approximation.
func MVPNumeratorApproxTerms(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	sqrtPolicy SqrtPolicy2,
	digits int,
) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("MVPNumeratorApproxTerms: negative digits %d", digits)
	}

	cf, err := MVPNumeratorApproxCF(fourOverPiPrefixTerms, ePrefixTerms, sqrtPolicy)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// MVPNumeratorApproxTermsDefault uses the default sqrt policy.
func MVPNumeratorApproxTermsDefault(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	digits int,
) ([]int64, error) {
	return MVPNumeratorApproxTerms(
		fourOverPiPrefixTerms,
		ePrefixTerms,
		DefaultSqrtPolicy2(),
		digits,
	)
}

// MVP numerator shape note:
//
//	The full target formula should remain assembled in tests for now.
//	This production helper intentionally stops at the numerator:
//
//	    sqrt(3/pi^2 + e)
//
// mvp_numerator.go v1
