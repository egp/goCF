// mvp_denominator.go v2
package cf

import "fmt"

// MVPDenominatorApprox returns a bounded rational approximation for:
//
//	tanh(sqrt(5)) - sin(69°)
//
// Current status:
//   - degree-aware angle semantics are now fixed
//   - denominator evaluator shape is frozen
//   - tanh approximation kernel is not implemented yet
//   - sin approximation currently supports only a narrow exact degree subset
func MVPDenominatorApprox(
	sqrt5Policy SqrtPolicy2,
	angle Angle,
) (Rational, error) {
	if err := angle.Validate(); err != nil {
		return Rational{}, err
	}
	if !angle.IsDegrees() {
		return Rational{}, fmt.Errorf("MVPDenominatorApprox: angle must be expressed in degrees")
	}

	_, err := SinApproxDegrees(angle)
	if err != nil {
		return Rational{}, err
	}

	return Rational{}, fmt.Errorf("MVPDenominatorApprox: tanh kernel not implemented")
}

func MVPDenominatorApproxDefault() (Rational, error) {
	return MVPDenominatorApprox(
		DefaultSqrtPolicy2(),
		Degrees(mustRat(69, 1)),
	)
}

// MVP denominator shape note:
//
//	The full target formula should remain assembled in tests for now.
//	This production helper intentionally stops at the denominator:
//
//	    tanh(sqrt(5)) - sin(69°)
//
// mvp_denominator.go v2
