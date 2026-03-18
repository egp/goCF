// cf/sqrt_unary_digit.go v2
package cf

import "fmt"

func sqrtUnaryNextDigitIfForced(r Range) (int64, bool, error) {
	if !r.IsInside() {
		return 0, false, fmt.Errorf("sqrtUnaryNextDigitIfForced: require inside range, got %v", r)
	}

	flo, fhi, err := r.FloorBounds()
	if err != nil {
		return 0, false, err
	}
	if flo != fhi {
		return 0, false, nil
	}
	return flo, true, nil
}

func (s *sqrtUnaryOperator) nextDigitIfForced() (int64, bool, error) {
	snap := s.snapshot()
	if snap.InputApprox == nil || snap.InputApprox.Range == nil {
		return 0, false, nil
	}
	if snap.CurrentY == nil {
		return 0, false, nil
	}

	enclosure, err := sqrtUnaryRangeEnclosureFromInputRange(*snap.InputApprox.Range, *snap.CurrentY)
	if err != nil {
		return 0, false, err
	}
	return sqrtUnaryNextDigitIfForced(enclosure)
}

// cf/sqrt_unary_digit.go v2
