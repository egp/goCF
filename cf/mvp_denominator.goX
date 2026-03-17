// mvp_denominator.go v7
package cf

import "fmt"

// MVPDenominatorBounds returns a conservative inside range for:
//
//	tanh(sqrt(5)) - sin(69°)
//
// Current status:
//   - degree-aware angle semantics are fixed
//   - sin is routed through a GCF+exact-tail unary entry point
//   - tanh(sqrt(5)) is routed through a GCF-ingesting metadata-driven unary entry poin

func MVPDenominatorBounds(
	sqrt5Policy SqrtPolicy2,
	angle Angle,
) (Range, error) {
	_ = sqrt5Policy // reserved for later tighter tanh(sqrt(5)) work

	if err := angle.Validate(); err != nil {
		return Range{}, err
	}
	if !angle.IsDegrees() {
		return Range{}, fmt.Errorf("MVPDenominatorBounds: angle must be expressed in degrees")
	}

	tanhR, err := TanhBoundsSpecialFromGCF2(AdaptCFToGCF(Sqrt5CF()))
	if err != nil {
		return Range{}, err
	}

	sinR, err := SinBoundsDegreesFromGCFPrefix2(
		MVP69DegreeGCFSource(),
		2,
	)
	if err != nil {
		return Range{}, err
	}

	lo, err := tanhR.Lo.Sub(sinR.Hi)
	if err != nil {
		return Range{}, err
	}
	hi, err := tanhR.Hi.Sub(sinR.Lo)
	if err != nil {
		return Range{}, err
	}

	return NewRange(lo, hi, true, true), nil
}

// MVPDenominatorApprox returns a bounded rational approximation for:
//
//	tanh(sqrt(5)) - sin(69°)
//
// Current MVP behavior:
//   - derive a certified denominator range
//   - if it is an exact point, return that point
//   - otherwise report that only a bounded non-point result is available
func MVPDenominatorApprox(
	sqrt5Policy SqrtPolicy2,
	angle Angle,
) (Rational, error) {
	r, err := MVPDenominatorBounds(sqrt5Policy, angle)
	if err != nil {
		return Rational{}, err
	}
	if r.Lo.Cmp(r.Hi) != 0 {
		return Rational{}, fmt.Errorf("MVPDenominatorApprox: bounded non-point result for %v", angle)
	}
	return r.Lo, nil
}

func MVPDenominatorBoundsDefault() (Range, error) {
	return MVPDenominatorBounds(
		DefaultSqrtPolicy2(),
		Degrees(mustRat(69, 1)),
	)
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
// mvp_denominator.go v7
