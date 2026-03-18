// cf/range_big_floor.go v2
package cf

import (
	"fmt"
	"math/big"
)

func (r Range) floorBigBounds() (*big.Int, *big.Int, error) {
	fLo, err := floorBigRat(r.Lo)
	if err != nil {
		return nil, nil, err
	}
	fHi, err := floorBigRat(r.Hi)
	if err != nil {
		return nil, nil, err
	}

	if fLo.Cmp(fHi) <= 0 {
		return fLo, fHi, nil
	}
	return fHi, fLo, nil
}

func floorBigRat(x Rational) (*big.Int, error) {
	num, den := x.ratNumDen()
	if den.Sign() <= 0 {
		return nil, fmt.Errorf("floorBigRat: invalid denominator %s", den.String())
	}

	quo := new(big.Int)
	rem := new(big.Int)
	quo.QuoRem(num, den, rem)

	if rem.Sign() != 0 && num.Sign() < 0 {
		quo.Sub(quo, big.NewInt(1))
	}

	return quo, nil
}

// cf/range_big_floor.go v2
