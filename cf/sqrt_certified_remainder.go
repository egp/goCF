// sqrt_certified_remainder.go v1
package cf

import "fmt"

// CertifiedRemainderRange maps a certified-output range for z into the range for
// the continued-fraction remainder z' = 1 / (z - d).
//
// Preconditions:
//   - r must be an inside range
//   - every value in r must satisfy floor(z) = d
//   - z-d must stay strictly positive across the range
func CertifiedRemainderRange(r Range, d int64) (Range, error) {
	if !r.IsInside() {
		return Range{}, fmt.Errorf("CertifiedRemainderRange: requires inside range; got %v", r)
	}

	loFloor, hiFloor, err := r.FloorBounds()
	if err != nil {
		return Range{}, err
	}
	if loFloor != d || hiFloor != d {
		return Range{}, fmt.Errorf("CertifiedRemainderRange: range does not certify digit %d; got floor bounds (%d,%d)", d, loFloor, hiFloor)
	}

	dd := intRat(d)

	loShift, err := r.Lo.Sub(dd)
	if err != nil {
		return Range{}, err
	}
	hiShift, err := r.Hi.Sub(dd)
	if err != nil {
		return Range{}, err
	}

	if loShift.Cmp(intRat(0)) <= 0 {
		return Range{}, fmt.Errorf("CertifiedRemainderRange: lower shifted endpoint not strictly positive: %v", loShift)
	}
	if hiShift.Cmp(intRat(0)) <= 0 {
		return Range{}, fmt.Errorf("CertifiedRemainderRange: upper shifted endpoint not strictly positive: %v", hiShift)
	}

	// Reciprocal reverses order on positive reals.
	hiRecip, err := intRat(1).Div(hiShift)
	if err != nil {
		return Range{}, err
	}
	loRecip, err := intRat(1).Div(loShift)
	if err != nil {
		return Range{}, err
	}

	return NewRange(hiRecip, loRecip, r.IncHi, r.IncLo), nil
}

// sqrt_certified_remainder.go v1
