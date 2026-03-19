// cf/sqrt_unary_digit.go v5
package cf

import (
	"fmt"
	"math/big"
)

func sqrtUnaryNextDigitIfForced(r Range) (*big.Int, bool, error) {
	if !r.IsInside() {
		return nil, false, fmt.Errorf("sqrtUnaryNextDigitIfForced: require inside range, got %v", r)
	}

	flo, fhi, err := r.floorBigBounds()
	if err != nil {
		return nil, false, err
	}
	if flo.Cmp(fhi) != 0 {
		return nil, false, nil
	}
	return new(big.Int).Set(flo), true, nil
}

func (s *sqrtUnaryOperator) nextDigitIfForced() (*big.Int, bool, error) {
	if s.currentEnclosure != nil {
		return sqrtUnaryNextDigitIfForced(*s.currentEnclosure)
	}

	snap := s.snapshot()
	if snap.InputApprox == nil || snap.InputApprox.Range == nil {
		return nil, false, nil
	}
	if snap.CurrentY == nil {
		return nil, false, nil
	}

	enclosure, err := sqrtUnaryRangeEnclosureFromInputRange(*snap.InputApprox.Range, *snap.CurrentY)
	if err != nil {
		return nil, false, err
	}
	return sqrtUnaryNextDigitIfForced(enclosure)
}

// cf/sqrt_unary_digit.go v5
