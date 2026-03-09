// gcf_tail_ray.go v1
package cf

import (
	"fmt"
	"math/big"
)

// ApplyULFTToTailRay maps the positive ray [lower, +∞) conservatively through
// a ULFT.
//
// This is the first enclosure primitive for unfinished GCF prefixes, where a
// finite ingested prefix defines a transform x = T(tail) and the unknown tail is
// known to satisfy tail >= lower > 0.
//
// Current support:
//   - requires lower > 0
//   - requires a valid ULFT
//   - requires that the denominator does not cross zero anywhere on [lower, +∞)
//   - supports bounded image cases only
//
// Bounded image cases:
//   - C != 0: image is bounded by T(lower) and the limit at infinity A/C
//   - C == 0 and A == 0: constant transform
//
// Unbounded affine cases (C == 0 and A != 0) currently return an error because
// Range does not represent infinite endpoints.
func ApplyULFTToTailRay(t ULFT, lower Rational) (Range, error) {
	if lower.Cmp(intRat(0)) <= 0 {
		return Range{}, fmt.Errorf("ApplyULFTToTailRay: require lower > 0, got %v", lower)
	}
	if err := t.Validate(); err != nil {
		return Range{}, err
	}

	// T(lower)
	y0, err := t.ApplyRat(lower)
	if err != nil {
		return Range{}, err
	}

	// C == 0 => denominator constant
	if t.C.Sign() == 0 {
		// Validated ULFT guarantees D != 0 here.
		if t.A.Sign() == 0 {
			// Constant transform B/D
			return NewRange(y0, y0, true, true), nil
		}
		return Range{}, fmt.Errorf("ApplyULFTToTailRay: unbounded affine image not representable")
	}

	// Pole location root = -D/C. If root >= lower, the ray intersects the pole.
	var negD big.Int
	negD.Neg(t.D)
	root, err := newRationalBig(&negD, t.C)
	if err != nil {
		return Range{}, err
	}
	if root.Cmp(lower) >= 0 {
		return Range{}, fmt.Errorf("ApplyULFTToTailRay: denominator crosses 0 on [%v,+∞)", lower)
	}

	// Finite limit at infinity = A/C
	yInf, err := newRationalBig(new(big.Int).Set(t.A), new(big.Int).Set(t.C))
	if err != nil {
		return Range{}, err
	}

	return NewRange(minRat(y0, yInf), maxRat(y0, yInf), true, true), nil
}

// gcf_tail_ray.go v1
