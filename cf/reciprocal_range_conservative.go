// reciprocal_range_conservative.go v1
package cf

import "fmt"

// ReciprocalRangeConservative returns a proof-safe enclosure for 1/x over a
// strictly positive inside range r.
//
// Preconditions:
//   - r must be inside
//   - every value in r must be strictly positive
//
// For positive ranges, reciprocal reverses order:
//
//	x in [lo, hi]  =>  1/x in [1/hi, 1/lo]
func ReciprocalRangeConservative(r Range) (Range, error) {
	if !r.IsInside() {
		return Range{}, fmt.Errorf("ReciprocalRangeConservative: requires inside range; got %v", r)
	}
	if r.Lo.Cmp(intRat(0)) <= 0 {
		return Range{}, fmt.Errorf("ReciprocalRangeConservative: lower endpoint must be > 0; got %v", r.Lo)
	}
	if r.Hi.Cmp(intRat(0)) <= 0 {
		return Range{}, fmt.Errorf("ReciprocalRangeConservative: upper endpoint must be > 0; got %v", r.Hi)
	}

	lo, err := intRat(1).Div(r.Hi)
	if err != nil {
		return Range{}, err
	}
	hi, err := intRat(1).Div(r.Lo)
	if err != nil {
		return Range{}, err
	}

	return NewRange(lo, hi, r.IncHi, r.IncLo), nil
}

// reciprocal_range_conservative.go v1
