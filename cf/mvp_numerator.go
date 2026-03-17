// mvp_numerator.go v5
package cf

import "fmt"

// Default bounded-prefix choices for the current MVP numerator path.
//
// These are intentionally centralized so target-level code and tests lock the
// same chosen approximation budget.
const (
	MVPDefaultFourOverPiPrefixTerms = 6
	MVPDefaultEPrefixTerms          = 8

	// Temporary compatibility budget retained for snapshot-based callers.
	// The live numerator production path consumes a direct radicand snapshot.
	MVPNumeratorBridgePrefixTerms = 64
)

// MVPNumeratorRadicandApprox returns the bounded rational subexpression
//
//	3/pi^2 + e
//
// for the numerator path.
func MVPNumeratorRadicandApprox(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
) (Rational, error) {
	return MVPThreeOverPiSquaredPlusEApprox(fourOverPiPrefixTerms, ePrefixTerms)
}

// MVPNumeratorApprox returns a bounded rational approximation for:
//
//	sqrt(3/pi^2 + e)
//
// Current MVP construction:
//   - assemble a direct radicand snapshot for 3/pi^2 + e
//   - route the final sqrt from that snapshot through the GCF-aware unary path
func MVPNumeratorApprox(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	sqrtPolicy SqrtPolicy2,
) (Rational, error) {
	return MVPNumeratorApproxWithBridgeTerms(
		fourOverPiPrefixTerms,
		ePrefixTerms,
		sqrtPolicy,
		MVPNumeratorBridgePrefixTerms,
	)
}

func MVPNumeratorApproxWithBridgeTerms(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	sqrtPolicy SqrtPolicy2,
	bridgeTerms int,
) (Rational, error) {
	a, err := MVPNumeratorRadicandApproxSnapshot(
		fourOverPiPrefixTerms,
		ePrefixTerms,
		bridgeTerms,
	)
	if err != nil {
		return Rational{}, err
	}

	return MVPNumeratorApproxFromRadicandApprox(a, sqrtPolicy)
}

func MVPNumeratorRadicandApproxSnapshot(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	bridgeTerms int,
) (GCFApprox, error) {
	if bridgeTerms <= 0 {
		return GCFApprox{}, fmt.Errorf(
			"MVPNumeratorRadicandApproxSnapshot: bridgeTerms must be > 0, got %d",
			bridgeTerms,
		)
	}

	return MVPThreeOverPiSquaredPlusERadicandSnapshot(
		fourOverPiPrefixTerms,
		ePrefixTerms,
		bridgeTerms,
	)
}

func MVPNumeratorApproxFromRadicandApprox(
	a GCFApprox,
	sqrtPolicy SqrtPolicy2,
) (Rational, error) {
	return SqrtApproxFromGCFApproxRangeSeed2(a, sqrtPolicy)
}

// MVPNumeratorApproxDefault uses the default sqrt policy and the current chosen
// bounded-prefix budgets for the MVP numerator.
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

// MVPNumeratorApproxCurrentDefault returns the current chosen default numerator
// approximation for the MVP target.
func MVPNumeratorApproxCurrentDefault() (Rational, error) {
	return MVPNumeratorApprox(
		MVPDefaultFourOverPiPrefixTerms,
		MVPDefaultEPrefixTerms,
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
// mvp_numerator.go v5
