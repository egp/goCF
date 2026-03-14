// sqrt_cf.go v3
package cf

import "fmt"

// collectTerms reads up to n CF terms from src.
// It stops early if the source terminates.
func collectTerms(src ContinuedFraction, n int) []int64 {
	if n <= 0 {
		return []int64{}
	}
	out := make([]int64, 0, n)
	for i := 0; i < n; i++ {
		a, ok := src.Next()
		if !ok {
			break
		}
		out = append(out, a)
	}
	return out
}

// DefaultSqrtSeed returns a simple exact positive seed for Newton sqrt.
//
// Policy:
//   - x < 0  => error
//   - x = 0  => 1
//   - 0 < x <= 1 => 1
//   - x > 1  => x
//
// This is intentionally simple and exact. It is a seed policy, not a proof.
func DefaultSqrtSeed(x Rational) (Rational, error) {
	if x.Cmp(intRat(0)) < 0 {
		return Rational{}, fmt.Errorf("sqrt of negative rational: %v", x)
	}
	if x.Cmp(intRat(1)) <= 0 {
		return intRat(1), nil
	}
	return x, nil
}

// NewSqrtApproxCF returns a ContinuedFraction source for the bounded rational
// Newton approximation to sqrt(x) produced after the requested number of steps
// from the given seed.
//
// This is a convenience adapter:
//
//	sqrt target -> exact rational approximation -> RationalCF
//
// It is not yet a true streaming sqrt operator.
func NewSqrtApproxCF(x, seed Rational, steps int) (ContinuedFraction, error) {
	approx, err := SqrtApproxRational(x, seed, steps)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

// NewSqrtApproxCFDefault is a convenience wrapper that uses DefaultSqrtSeed.
func NewSqrtApproxCFDefault(x Rational, steps int) (ContinuedFraction, error) {
	seed, err := DefaultSqrtSeed(x)
	if err != nil {
		return nil, err
	}
	return NewSqrtApproxCF(x, seed, steps)
}

// SqrtApproxTerms computes a bounded rational Newton approximation to sqrt(x),
// converts that exact rational approximation to a finite CF, and returns up to
// digits terms of that CF.
//
// This is a convenience bridge from exact-rational sqrt approximation to CF
// output terms. It is not yet a true streaming sqrt operator.
func SqrtApproxTerms(x, seed Rational, steps, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTerms: negative digits %d", digits)
	}

	approx, err := SqrtApproxRational(x, seed, steps)
	if err != nil {
		return nil, err
	}

	return collectTerms(NewRationalCF(approx), digits), nil
}

// SqrtApproxTermsDefault computes CF terms for a bounded rational Newton
// approximation to sqrt(x), using DefaultSqrtSeed(x).
func SqrtApproxTermsDefault(x Rational, steps, digits int) ([]int64, error) {
	seed, err := DefaultSqrtSeed(x)
	if err != nil {
		return nil, err
	}
	return SqrtApproxTerms(x, seed, steps, digits)
}

// NewSqrtApproxCFUntilResidual returns a ContinuedFraction source for the bounded
// Newton approximation to sqrt(x), stopping early once |y^2 - x| <= tol.
func NewSqrtApproxCFUntilResidual(x, seed Rational, maxSteps int, tol Rational) (ContinuedFraction, error) {
	approx, _, err := SqrtApproxRationalUntilResidual(x, seed, maxSteps, tol)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

// NewSqrtApproxCFUntilResidualDefault is the default-seed wrapper around
// NewSqrtApproxCFUntilResidual.
func NewSqrtApproxCFUntilResidualDefault(x Rational, maxSteps int, tol Rational) (ContinuedFraction, error) {
	approx, _, err := SqrtApproxRationalUntilResidualDefault(x, maxSteps, tol)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

// SqrtApproxTermsUntilResidual computes CF terms for the bounded Newton
// approximation to sqrt(x), stopping early once |y^2 - x| <= tol.
func SqrtApproxTermsUntilResidual(x, seed Rational, maxSteps int, tol Rational, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsUntilResidual: negative digits %d", digits)
	}
	cf, err := NewSqrtApproxCFUntilResidual(x, seed, maxSteps, tol)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// SqrtApproxTermsUntilResidualDefault is the default-seed wrapper around
// SqrtApproxTermsUntilResidual.
func SqrtApproxTermsUntilResidualDefault(x Rational, maxSteps int, tol Rational, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsUntilResidualDefault: negative digits %d", digits)
	}
	cf, err := NewSqrtApproxCFUntilResidualDefault(x, maxSteps, tol)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

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

// RationalApproxFromCFPrefix ingests up to prefixTerms terms from src and returns
// the resulting convergent as an exact Rational.
//
// Behavior:
//   - prefixTerms < 0 => error
//   - prefixTerms == 0 => error
//   - if src terminates early, the bounder is finished and the exact rational
//     represented by the finite source is returned
//   - if src does not terminate within prefixTerms, the convergent of the prefix
//     is returned
func RationalApproxFromCFPrefix(src ContinuedFraction, prefixTerms int) (Rational, error) {
	a, err := CFApproxFromPrefix(src, prefixTerms)
	if err != nil {
		return Rational{}, err
	}
	return a.Convergent, nil
}

// NewSqrtApproxCFFromSource consumes a finite prefix of src, converts that prefix
// to a rational approximation, then returns a ContinuedFraction source for a
// bounded sqrt approximation under the supplied policy.
//
// This is a bridge from CF input to the existing rational sqrt machinery.
// It is still bounded/approximate, not a true streaming sqrt operator.
func NewSqrtApproxCFFromSource(src ContinuedFraction, prefixTerms int, p SqrtPolicy) (ContinuedFraction, error) {
	xApprox, err := RationalApproxFromCFPrefix(src, prefixTerms)
	if err != nil {
		return nil, err
	}
	return SqrtApproxCFWithPolicy(xApprox, p)
}

// NewSqrtApproxCFFromSourceDefault is the default-policy wrapper around
// NewSqrtApproxCFFromSource.
func NewSqrtApproxCFFromSourceDefault(src ContinuedFraction, prefixTerms int) (ContinuedFraction, error) {
	return NewSqrtApproxCFFromSource(src, prefixTerms, DefaultSqrtPolicy())
}

// SqrtApproxTermsFromSource consumes a finite prefix of src, converts that prefix
// to a rational approximation, then returns up to digits CF terms for the bounded
// sqrt approximation under the supplied policy.
func SqrtApproxTermsFromSource(src ContinuedFraction, prefixTerms int, p SqrtPolicy, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromSource: negative digits %d", digits)
	}
	cf, err := NewSqrtApproxCFFromSource(src, prefixTerms, p)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// SqrtApproxTermsFromSourceDefault is the default-policy wrapper around
// SqrtApproxTermsFromSource.
func SqrtApproxTermsFromSourceDefault(src ContinuedFraction, prefixTerms int, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromSourceDefault: negative digits %d", digits)
	}
	return SqrtApproxTermsFromSource(src, prefixTerms, DefaultSqrtPolicy(), digits)
}

// NewSqrtApproxCFFromApproxRangeSeed takes a bundled CFApprox and returns a
// ContinuedFraction source for the bounded sqrt approximation under the supplied
// policy.
//
// If p.Seed is already set, it is honored and the range-derived seed is not used.
func NewSqrtApproxCFFromApproxRangeSeed(a CFApprox, p SqrtPolicy) (ContinuedFraction, error) {
	return SqrtApproxCFFromApproxRangeSeed2(
		a,
		sqrtPolicy2FromOld(p),
	)
}

// NewSqrtApproxCFFromSourceRangeSeed consumes a finite prefix of src, converts
// that prefix to a CFApprox, and then returns a ContinuedFraction source for the
// bounded sqrt approximation under the supplied policy.
func NewSqrtApproxCFFromSourceRangeSeed(src ContinuedFraction, prefixTerms int, p SqrtPolicy) (ContinuedFraction, error) {
	return SqrtApproxCFFromSourceRangeSeed2(
		src,
		prefixTerms,
		sqrtPolicy2FromOld(p),
	)
}

// NewSqrtApproxCFFromSourceRangeSeedDefault is the default-policy wrapper
// around NewSqrtApproxCFFromSourceRangeSeed.
func NewSqrtApproxCFFromSourceRangeSeedDefault(src ContinuedFraction, prefixTerms int) (ContinuedFraction, error) {
	return SqrtApproxCFFromSourceRangeSeedDefault2(src, prefixTerms)
}

// SqrtApproxTermsFromApproxRangeSeed returns up to digits CF terms for the
// bounded sqrt approximation produced by NewSqrtApproxCFFromApproxRangeSeed.
func SqrtApproxTermsFromApproxRangeSeed(a CFApprox, p SqrtPolicy, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromApproxRangeSeed: negative digits %d", digits)
	}
	return SqrtApproxTermsFromApproxRangeSeed2(
		a,
		sqrtPolicy2FromOld(p),
		digits,
	)
}

// SqrtApproxTermsFromSourceRangeSeed returns up to digits CF terms for the
// bounded sqrt approximation produced by NewSqrtApproxCFFromSourceRangeSeed.
func SqrtApproxTermsFromSourceRangeSeed(src ContinuedFraction, prefixTerms int, p SqrtPolicy, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromSourceRangeSeed: negative digits %d", digits)
	}
	return SqrtApproxTermsFromSourceRangeSeed2(
		src,
		prefixTerms,
		sqrtPolicy2FromOld(p),
		digits,
	)
}

// SqrtApproxTermsFromSourceRangeSeedDefault is the default-policy wrapper
// around SqrtApproxTermsFromSourceRangeSeed.
func SqrtApproxTermsFromSourceRangeSeedDefault(src ContinuedFraction, prefixTerms int, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromSourceRangeSeedDefault: negative digits %d", digits)
	}
	return SqrtApproxTermsFromSourceRangeSeedDefault2(src, prefixTerms, digits)
}

// NewSqrtApproxCFFromApproxRangeMidpoint takes a bundled CFApprox, uses the
// midpoint of its enclosure as the sqrt target, and returns a ContinuedFraction
// source for the bounded sqrt approximation under the supplied policy.
//
// IMPORTANT:
//   - this is an experimental heuristic path
//   - it is not yet a proof-safe conservative sqrt operator
func NewSqrtApproxCFFromApproxRangeMidpoint(a CFApprox, p SqrtPolicy) (ContinuedFraction, error) {
	m, err := rangeMidpoint(a.Range)
	if err != nil {
		return nil, err
	}
	return SqrtApproxCFWithPolicy(m, p)
}

// NewSqrtApproxCFFromSourceRangeMidpoint consumes a finite prefix of src,
// converts that prefix to a CFApprox, uses the midpoint of the enclosure as
// the sqrt target, and returns a ContinuedFraction source for the bounded
// sqrt approximation under the supplied policy.
func NewSqrtApproxCFFromSourceRangeMidpoint(src ContinuedFraction, prefixTerms int, p SqrtPolicy) (ContinuedFraction, error) {
	a, err := CFApproxFromPrefix(src, prefixTerms)
	if err != nil {
		return nil, err
	}
	return NewSqrtApproxCFFromApproxRangeMidpoint(a, p)
}

// NewSqrtApproxCFFromSourceRangeMidpointDefault is the default-policy wrapper
// around NewSqrtApproxCFFromSourceRangeMidpoint.
func NewSqrtApproxCFFromSourceRangeMidpointDefault(src ContinuedFraction, prefixTerms int) (ContinuedFraction, error) {
	return NewSqrtApproxCFFromSourceRangeMidpoint(src, prefixTerms, DefaultSqrtPolicy())
}

// SqrtApproxTermsFromApproxRangeMidpoint returns up to digits CF terms for the
// bounded sqrt approximation produced by NewSqrtApproxCFFromApproxRangeMidpoint.
func SqrtApproxTermsFromApproxRangeMidpoint(a CFApprox, p SqrtPolicy, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromApproxRangeMidpoint: negative digits %d", digits)
	}
	cf, err := NewSqrtApproxCFFromApproxRangeMidpoint(a, p)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// SqrtApproxTermsFromSourceRangeMidpoint returns up to digits CF terms for the
// bounded sqrt approximation produced by NewSqrtApproxCFFromSourceRangeMidpoint.
func SqrtApproxTermsFromSourceRangeMidpoint(src ContinuedFraction, prefixTerms int, p SqrtPolicy, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromSourceRangeMidpoint: negative digits %d", digits)
	}
	cf, err := NewSqrtApproxCFFromSourceRangeMidpoint(src, prefixTerms, p)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// SqrtApproxTermsFromSourceRangeMidpointDefault is the default-policy wrapper
// around SqrtApproxTermsFromSourceRangeMidpoint.
func SqrtApproxTermsFromSourceRangeMidpointDefault(src ContinuedFraction, prefixTerms int, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromSourceRangeMidpointDefault: negative digits %d", digits)
	}
	return SqrtApproxTermsFromSourceRangeMidpoint(src, prefixTerms, DefaultSqrtPolicy(), digits)
}

// sqrt_cf.go v
