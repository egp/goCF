// sqrt_certified_remainder.go v2
package cf

import "fmt"

// ShiftRangeByInt translates an inside range by subtracting d.
//
//	out = r - d
func ShiftRangeByInt(r Range, d int64) (Range, error) {
	if !r.IsInside() {
		return Range{}, fmt.Errorf("ShiftRangeByInt: requires inside range; got %v", r)
	}

	dd := intRat(d)

	lo, err := r.Lo.Sub(dd)
	if err != nil {
		return Range{}, err
	}
	hi, err := r.Hi.Sub(dd)
	if err != nil {
		return Range{}, err
	}

	return NewRange(lo, hi, r.IncLo, r.IncHi), nil
}

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
		return Range{}, fmt.Errorf(
			"CertifiedRemainderRange: range does not certify digit %d; got floor bounds (%d,%d)",
			d, loFloor, hiFloor,
		)
	}

	shifted, err := ShiftRangeByInt(r, d)
	if err != nil {
		return Range{}, err
	}

	if shifted.Lo.Cmp(intRat(0)) <= 0 {
		return Range{}, fmt.Errorf(
			"CertifiedRemainderRange: lower shifted endpoint not strictly positive: %v",
			shifted.Lo,
		)
	}
	if shifted.Hi.Cmp(intRat(0)) <= 0 {
		return Range{}, fmt.Errorf(
			"CertifiedRemainderRange: upper shifted endpoint not strictly positive: %v",
			shifted.Hi,
		)
	}

	return ReciprocalRangeConservative(shifted)
}

// sqrt_certified_remainder.go v2
