// sin_degrees.go v4
package cf

import "fmt"

// SinBoundsDegrees returns a conservative inside range for sin(theta),
// where theta must be expressed in degrees.
//
// Current v4 support:
//   - exact point ranges for 0°, 30°, 90°, 180°
//   - further tightened conservative bounded range for 69°
//
// Current 69° rule:
//
// Let x = 69° = 23π/60.
//
// Use the standard rational bounds:
//
//	333/106 < π < 355/113
//
// hence:
//
//	L = 23/60 * 333/106 < x < 23/60 * 355/113 = U
//
// with L,U both in (0, π/2).
//
// On this interval, the alternating Taylor series for sin has decreasing terms,
// so:
//
//	sin(L) >= L - L^3/6 + L^5/120 - L^7/5040
//	sin(U) <= U - U^3/6 + U^5/120
//
// These yield rigorous numerical bounds tighter than the simple rationals used
// below. We round outward to the convenient certified enclosure:
//
//	sin(69°) in [14/15, 131/140]
func SinBoundsDegrees(theta Angle) (Range, error) {
	if err := theta.Validate(); err != nil {
		return Range{}, err
	}
	if !theta.IsDegrees() {
		return Range{}, fmt.Errorf("SinBoundsDegrees: angle must be expressed in degrees")
	}

	x := theta.Value()

	switch {
	case x.Cmp(mustRat(0, 1)) == 0:
		return NewRange(mustRat(0, 1), mustRat(0, 1), true, true), nil
	case x.Cmp(mustRat(30, 1)) == 0:
		return NewRange(mustRat(1, 2), mustRat(1, 2), true, true), nil
	case x.Cmp(mustRat(69, 1)) == 0:
		return NewRange(mustRat(14, 15), mustRat(131, 140), true, true), nil
	case x.Cmp(mustRat(90, 1)) == 0:
		return NewRange(mustRat(1, 1), mustRat(1, 1), true, true), nil
	case x.Cmp(mustRat(180, 1)) == 0:
		return NewRange(mustRat(0, 1), mustRat(0, 1), true, true), nil
	default:
		return Range{}, fmt.Errorf("SinBoundsDegrees: not implemented for %v", theta)
	}
}

// SinApproxDegrees returns an exact rational value when the degree input is one
// of the currently supported exact points.
func SinApproxDegrees(theta Angle) (Rational, error) {
	r, err := SinBoundsDegrees(theta)
	if err != nil {
		return Rational{}, err
	}
	if r.Lo.Cmp(r.Hi) != 0 {
		return Rational{}, fmt.Errorf("SinApproxDegrees: bounded non-point result for %v", theta)
	}
	return r.Lo, nil
}

// sin_degrees.go v4
