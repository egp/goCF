// sqrt_seed_range.go v1
package cf

import "fmt"

// SqrtSeedDefault returns a simple exact positive seed for Newton sqrt.
//
// Policy:
//   - x < 0       => error
//   - x = 0       => 1
//   - 0 < x <= 1  => 1
//   - x > 1       => x
func SqrtSeedDefault(x Rational) (Rational, error) {
	if x.Cmp(intRat(0)) < 0 {
		return Rational{}, fmt.Errorf("sqrt of negative rational: %v", x)
	}
	if x.Cmp(intRat(1)) <= 0 {
		return intRat(1), nil
	}
	return x, nil
}

// SqrtRangeMidpoint returns (Lo + Hi) / 2 for an inside range.
func SqrtRangeMidpoint(r Range) (Rational, error) {
	if !r.IsInside() {
		return Rational{}, fmt.Errorf("SqrtRangeMidpoint requires inside range; got %v", r)
	}
	sum, err := r.Lo.Add(r.Hi)
	if err != nil {
		return Rational{}, err
	}
	return sum.Div(intRat(2))
}

// SqrtSeedFromRange derives a better exact seed for sqrt(x) from an enclosure of x.
//
// Current policy:
//   - range must be inside and nonnegative
//   - let m = midpoint(range)
//   - if m has exact rational square root, return it
//   - otherwise:
//     base = SqrtSeedDefault(m)
//     seed = one Newton step toward sqrt(m) from base
func SqrtSeedFromRange(r Range) (Rational, error) {
	if !r.IsInside() {
		return Rational{}, fmt.Errorf("SqrtSeedFromRange: requires inside range; got %v", r)
	}
	if r.Lo.Cmp(intRat(0)) < 0 {
		return Rational{}, fmt.Errorf("SqrtSeedFromRange: negative range %v", r)
	}

	m, err := SqrtRangeMidpoint(r)
	if err != nil {
		return Rational{}, err
	}

	if root, ok, err := SqrtCoreRationalExact(m); err != nil {
		return Rational{}, err
	} else if ok {
		return root, nil
	}

	base, err := SqrtSeedDefault(m)
	if err != nil {
		return Rational{}, err
	}
	return SqrtCoreNewtonStep(m, base)
}

// SqrtSeedFromCFPrefix derives a sqrt seed from the enclosure implied by
// a finite CF prefix.
func SqrtSeedFromCFPrefix(src ContinuedFraction, prefixTerms int) (Rational, error) {
	rng, err := RangeApproxFromCFPrefix(src, prefixTerms)
	if err != nil {
		return Rational{}, err
	}
	return SqrtSeedFromRange(rng)
}

// SqrtRangeExact2 returns the exact monotone image of a nonnegative inside range
// under sqrt(x), when both endpoints have exact rational square roots.
func SqrtRangeExact2(r Range) (Range, bool, error) {
	if !r.IsInside() {
		return Range{}, false, fmt.Errorf("SqrtRangeExact2: requires inside range; got %v", r)
	}
	if r.Lo.Cmp(intRat(0)) < 0 {
		return Range{}, false, fmt.Errorf("SqrtRangeExact2: negative range %v", r)
	}

	lo, okLo, err := SqrtCoreRationalExact(r.Lo)
	if err != nil {
		return Range{}, false, err
	}
	if !okLo {
		return Range{}, false, nil
	}

	hi, okHi, err := SqrtCoreRationalExact(r.Hi)
	if err != nil {
		return Range{}, false, err
	}
	if !okHi {
		return Range{}, false, nil
	}

	return NewRange(lo, hi, r.IncLo, r.IncHi), true, nil
}

// SqrtRangeExactFromCFApprox2 applies SqrtRangeExact2 to the enclosure carried by CFApprox.
func SqrtRangeExactFromCFApprox2(a CFApprox) (Range, bool, error) {
	return SqrtRangeExact2(a.Range)
}

// SqrtRangeHeuristic2 returns a bounded rational approximation to the monotone
// image of a nonnegative inside range under sqrt(x).
//
// IMPORTANT:
//   - this is a heuristic helper
//   - it is not yet a proof-safe conservative enclosure for irrational endpoints
func SqrtRangeHeuristic2(r Range) (Range, error) {
	if !r.IsInside() {
		return Range{}, fmt.Errorf("SqrtRangeHeuristic2: requires inside range; got %v", r)
	}
	if r.Lo.Cmp(intRat(0)) < 0 {
		return Range{}, fmt.Errorf("SqrtRangeHeuristic2: negative range %v", r)
	}

	if exact, ok, err := SqrtRangeExact2(r); err != nil {
		return Range{}, err
	} else if ok {
		return exact, nil
	}

	lo, err := SqrtApprox(r.Lo)
	if err != nil {
		return Range{}, err
	}
	hi, err := SqrtApprox(r.Hi)
	if err != nil {
		return Range{}, err
	}

	if lo.Cmp(hi) > 0 {
		lo, hi = hi, lo
	}
	return NewRange(lo, hi, r.IncLo, r.IncHi), nil
}

// SqrtRangeHeuristicFromCFApprox2 applies SqrtRangeHeuristic2 to the enclosure
// carried by CFApprox.
func SqrtRangeHeuristicFromCFApprox2(a CFApprox) (Range, error) {
	return SqrtRangeHeuristic2(a.Range)
}

// sqrt_seed_range.go v1
