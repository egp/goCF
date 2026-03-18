// cf/sqrt_unary_range_enclosure.go v3
package cf

import "fmt"

func sqrtUnaryRangeEnclosureFromInputRange(input Range, y Rational) (Range, error) {
	if !input.IsInside() {
		return Range{}, fmt.Errorf("sqrtUnaryRangeEnclosureFromInputRange: require inside range, got %v", input)
	}
	if input.Lo.Cmp(intRat(0)) <= 0 {
		return Range{}, fmt.Errorf(
			"sqrtUnaryRangeEnclosureFromInputRange: require positive lower bound, got %v",
			input.Lo,
		)
	}
	if y.Cmp(intRat(0)) <= 0 {
		return Range{}, fmt.Errorf(
			"sqrtUnaryRangeEnclosureFromInputRange: require positive iterate, got %v",
			y,
		)
	}

	loEnclosure, err := sqrtUnaryPointEnclosureExact(input.Lo, y)
	if err != nil {
		return Range{}, err
	}
	hiEnclosure, err := sqrtUnaryPointEnclosureExact(input.Hi, y)
	if err != nil {
		return Range{}, err
	}

	lo := loEnclosure.Lo
	if hiEnclosure.Lo.Cmp(lo) < 0 {
		lo = hiEnclosure.Lo
	}

	hi := loEnclosure.Hi
	if hiEnclosure.Hi.Cmp(hi) > 0 {
		hi = hiEnclosure.Hi
	}

	return NewRange(lo, hi, true, true), nil
}

// cf/sqrt_unary_range_enclosure.go v3
