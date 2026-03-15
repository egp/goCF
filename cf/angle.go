// angle.go v1
package cf

import "fmt"

// Angle represents an angle with explicit unit semantics.
//
// MVP rule:
//   - production code must not accept a bare Rational where angle semantics matter
//   - tests and callers must choose degrees or radians explicitly
type Angle struct {
	value   Rational
	degrees bool
}

func Degrees(x Rational) Angle {
	return Angle{value: x, degrees: true}
}

func Radians(x Rational) Angle {
	return Angle{value: x, degrees: false}
}

func (a Angle) Value() Rational {
	return a.value
}

func (a Angle) IsDegrees() bool {
	return a.degrees
}

func (a Angle) Validate() error {
	// Current MVP: all rational angles are acceptable.
	return nil
}

func (a Angle) String() string {
	if a.degrees {
		return fmt.Sprintf("%v°", a.value)
	}
	return fmt.Sprintf("%v rad", a.value)
}

// angle.go v1
