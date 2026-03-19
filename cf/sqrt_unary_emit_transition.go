// cf/sqrt_unary_emit_transition.go v2
package cf

import (
	"fmt"
	"math/big"
)

func sqrtUnaryEmitForcedDigitTransition(r Range) (*big.Int, *Range, bool, error) {
	d, ok, err := sqrtUnaryNextDigitIfForced(r)
	if err != nil {
		return nil, nil, false, err
	}
	if !ok {
		return nil, nil, false, fmt.Errorf("sqrtUnaryEmitForcedDigitTransition: digit not forced for %v", r)
	}

	// Exact integer point: emit digit and terminate.
	if r.Lo.Cmp(r.Hi) == 0 {
		dRat, err := newRationalBig(new(big.Int).Set(d), big.NewInt(1))
		if err != nil {
			return nil, nil, false, err
		}
		if r.Lo.Cmp(dRat) == 0 {
			return new(big.Int).Set(d), nil, true, nil
		}
	}

	dRat, err := newRationalBig(new(big.Int).Set(d), big.NewInt(1))
	if err != nil {
		return nil, nil, false, err
	}

	loShift, err := r.Lo.Sub(dRat)
	if err != nil {
		return nil, nil, false, err
	}
	hiShift, err := r.Hi.Sub(dRat)
	if err != nil {
		return nil, nil, false, err
	}

	// Nonterminal reciprocal step requires the shifted interval to stay strictly positive.
	if loShift.Cmp(intRat(0)) <= 0 || hiShift.Cmp(intRat(0)) <= 0 {
		return nil, nil, false, fmt.Errorf(
			"sqrtUnaryEmitForcedDigitTransition: nonterminal remainder not strictly positive for %v after digit %s",
			r,
			d.String(),
		)
	}

	remLo, err := intRat(1).Div(hiShift)
	if err != nil {
		return nil, nil, false, err
	}
	remHi, err := intRat(1).Div(loShift)
	if err != nil {
		return nil, nil, false, err
	}

	rem := NewRange(remLo, remHi, true, true)
	return new(big.Int).Set(d), &rem, false, nil
}

// cf/sqrt_unary_emit_transition.go v2
