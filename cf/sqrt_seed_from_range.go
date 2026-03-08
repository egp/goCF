// sqrt_seed_from_range.go v1
package cf

import "fmt"

// rangeMidpoint returns (Lo + Hi) / 2 for an inside range.
func rangeMidpoint(r Range) (Rational, error) {
	if !r.IsInside() {
		return Rational{}, fmt.Errorf("rangeMidpoint requires inside range; got %v", r)
	}
	sum, err := r.Lo.Add(r.Hi)
	if err != nil {
		return Rational{}, err
	}
	return sum.Div(intRat(2))
}

// DefaultSqrtSeedFromRange derives a better exact seed for sqrt(x) from an
// enclosure of x.
//
// Current policy:
//   - range must be inside and nonnegative
//   - let m = midpoint(range)
//   - if m has exact rational square root, return it
//   - otherwise:
//     base = DefaultSqrtSeed(m)
//     seed = one Newton step toward sqrt(m) from base
//
// This is intended as a better seed than using only the convergent or the
// generic DefaultSqrtSeed on a single rational approximation of x.
func DefaultSqrtSeedFromRange(r Range) (Rational, error) {
	if !r.IsInside() {
		return Rational{}, fmt.Errorf("DefaultSqrtSeedFromRange: requires inside range; got %v", r)
	}
	if r.Lo.Cmp(intRat(0)) < 0 {
		return Rational{}, fmt.Errorf("DefaultSqrtSeedFromRange: negative range %v", r)
	}

	m, err := rangeMidpoint(r)
	if err != nil {
		return Rational{}, err
	}

	if root, ok, err := RationalSqrtExact(m); err != nil {
		return Rational{}, err
	} else if ok {
		return root, nil
	}

	base, err := DefaultSqrtSeed(m)
	if err != nil {
		return Rational{}, err
	}
	return NewtonSqrtStep(m, base)
}

// DefaultSqrtSeedFromCFPrefix derives a sqrt seed from the enclosure implied by
// a finite CF prefix.
func DefaultSqrtSeedFromCFPrefix(src ContinuedFraction, prefixTerms int) (Rational, error) {
	rng, err := RangeApproxFromCFPrefix(src, prefixTerms)
	if err != nil {
		return Rational{}, err
	}
	return DefaultSqrtSeedFromRange(rng)
}

// sqrt_seed_from_range.go v1
