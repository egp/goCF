// sin_degrees.go v3
package cf

import "fmt"

// SinBoundsDegrees returns a conservative inside range for sin(theta),
// where theta must be expressed in degrees.
//
// Current v3 support:
//   - exact point ranges for 0°, 30°, 90°, 180°
//   - tightened conservative bounded range for 69°
//
// Current 69° rule:
//
// Lower bound:
//   - 69° is in the first quadrant
//   - sin is increasing on [0°,90°]
//   - so sin(69°) > sin(60°) = √3/2 > 6/7
//
// Upper bound:
//
//   - on [0,π], sin is concave, so it lies below its tangent lines
//
//   - at 60° = π/3, with x = 69° = 23π/60:
//
//     sin(x) <= sin(π/3) + cos(π/3)(x - π/3)
//     = √3/2 + (1/2)(π/20)
//     = √3/2 + π/40
//
//   - using √3/2 < 7/8 and π < 22/7:
//
//     sin(69°) < 7/8 + 11/140 = 267/280
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
		return NewRange(mustRat(6, 7), mustRat(267, 280), true, true), nil
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

// sin_degrees.go v3
