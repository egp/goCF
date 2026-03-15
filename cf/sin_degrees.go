// sin_degrees.go v2
package cf

import "fmt"

// SinBoundsDegrees returns a conservative inside range for sin(theta),
// where theta must be expressed in degrees.
//
// Current v2 support:
//   - exact point ranges for 0°, 30°, 90°, 180°
//   - conservative bounded range for 69°
//
// Current 69° rule:
//   - 69° is in the first quadrant
//   - sin is increasing on [0°,90°]
//   - therefore sin(69°) is strictly between sin(30°)=1/2 and sin(90°)=1
//
// This is intentionally conservative and should be tightened later.
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
		return NewRange(mustRat(1, 2), mustRat(1, 1), true, true), nil
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

// sin_degrees.go v2
