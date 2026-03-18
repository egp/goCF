// cf/sqrt_unary_enclosure.go v2
package cf

func sqrtUnaryPointEnclosureExact(x Rational, y Rational) (Range, error) {
	xy, err := x.Div(y)
	if err != nil {
		return Range{}, err
	}

	lo := y
	hi := xy
	if lo.Cmp(hi) > 0 {
		lo, hi = hi, lo
	}

	return NewRange(lo, hi, true, true), nil
}

// cf/sqrt_unary_enclosure.go v2
