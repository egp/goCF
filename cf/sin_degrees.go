// sin_degrees.go v1
package cf

import "fmt"

// SinApproxDegrees returns a bounded rational approximation for sin(theta),
// where theta must be expressed in degrees.
//
// Current v1 support is intentionally narrow but mathematically exact:
//   - 0°   -> 0
//   - 30°  -> 1/2
//   - 90°  -> 1
//   - 180° -> 0
//
// Future work:
//   - bounded approximation for arbitrary degree inputs such as 69°
//   - radian-input support
func SinApproxDegrees(theta Angle) (Rational, error) {
	if err := theta.Validate(); err != nil {
		return Rational{}, err
	}
	if !theta.IsDegrees() {
		return Rational{}, fmt.Errorf("SinApproxDegrees: angle must be expressed in degrees")
	}

	x := theta.Value()

	switch {
	case x.Cmp(mustRat(0, 1)) == 0:
		return mustRat(0, 1), nil
	case x.Cmp(mustRat(30, 1)) == 0:
		return mustRat(1, 2), nil
	case x.Cmp(mustRat(90, 1)) == 0:
		return mustRat(1, 1), nil
	case x.Cmp(mustRat(180, 1)) == 0:
		return mustRat(0, 1), nil
	default:
		return Rational{}, fmt.Errorf("SinApproxDegrees: not implemented for %v", theta)
	}
}

// sin_degrees.go v1
