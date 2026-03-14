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
	p := SqrtPolicy2{
		MaxSteps: steps,
		Tol:      mustRat(0, 1),
		Seed:     &seed,
	}
	approx, err := SqrtApproxWithPolicy2(x, p)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

// NewSqrtApproxCFDefault is a convenience wrapper that uses DefaultSqrtSeed.
func NewSqrtApproxCFDefault(x Rational, steps int) (ContinuedFraction, error) {
	p := SqrtPolicy2{
		MaxSteps: steps,
		Tol:      mustRat(0, 1),
		Seed:     nil,
	}
	approx, err := SqrtApproxWithPolicy2(x, p)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
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
	cf, err := NewSqrtApproxCF(x, seed, steps)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// SqrtApproxTermsDefault computes CF terms for a bounded rational Newton
// approximation to sqrt(x), using DefaultSqrtSeed(x).
func SqrtApproxTermsDefault(x Rational, steps, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsDefault: negative digits %d", digits)
	}
	cf, err := NewSqrtApproxCFDefault(x, steps)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// NewSqrtApproxCFUntilResidual returns a ContinuedFraction source for the bounded
// Newton approximation to sqrt(x), stopping early once |y^2 - x| <= tol.
func NewSqrtApproxCFUntilResidual(x, seed Rational, maxSteps int, tol Rational) (ContinuedFraction, error) {
	p := SqrtPolicy2{
		MaxSteps: maxSteps,
		Tol:      tol,
		Seed:     &seed,
	}
	approx, err := SqrtApproxWithPolicy2(x, p)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

// NewSqrtApproxCFUntilResidualDefault is the default-seed wrapper around
// NewSqrtApproxCFUntilResidual.
func NewSqrtApproxCFUntilResidualDefault(x Rational, maxSteps int, tol Rational) (ContinuedFraction, error) {
	p := SqrtPolicy2{
		MaxSteps: maxSteps,
		Tol:      tol,
		Seed:     nil,
	}
	approx, err := SqrtApproxWithPolicy2(x, p)
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
	return SqrtApproxCFFromApproxRangeMidpoint2(a, sqrtPolicy2FromOld(p))
}

// NewSqrtApproxCFFromSourceRangeMidpoint consumes a finite prefix of src,
// converts that prefix to a CFApprox, uses the midpoint of the enclosure as
// the sqrt target, and returns a ContinuedFraction source for the bounded
// sqrt approximation under the supplied policy.
func NewSqrtApproxCFFromSourceRangeMidpoint(src ContinuedFraction, prefixTerms int, p SqrtPolicy) (ContinuedFraction, error) {
	return SqrtApproxCFFromSourceRangeMidpoint2(src, prefixTerms, sqrtPolicy2FromOld(p))
}

// NewSqrtApproxCFFromSourceRangeMidpointDefault is the default-policy wrapper
// around NewSqrtApproxCFFromSourceRangeMidpoint.
func NewSqrtApproxCFFromSourceRangeMidpointDefault(src ContinuedFraction, prefixTerms int) (ContinuedFraction, error) {
	return SqrtApproxCFFromSourceRangeMidpointDefault2(src, prefixTerms)
}

// SqrtApproxTermsFromApproxRangeMidpoint returns up to digits CF terms for the
// bounded sqrt approximation produced by NewSqrtApproxCFFromApproxRangeMidpoint.
func SqrtApproxTermsFromApproxRangeMidpoint(a CFApprox, p SqrtPolicy, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromApproxRangeMidpoint: negative digits %d", digits)
	}
	return SqrtApproxTermsFromApproxRangeMidpoint2(a, sqrtPolicy2FromOld(p), digits)
}

// SqrtApproxTermsFromSourceRangeMidpoint returns up to digits CF terms for the
// bounded sqrt approximation produced by NewSqrtApproxCFFromSourceRangeMidpoint.
func SqrtApproxTermsFromSourceRangeMidpoint(src ContinuedFraction, prefixTerms int, p SqrtPolicy, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromSourceRangeMidpoint: negative digits %d", digits)
	}
	return SqrtApproxTermsFromSourceRangeMidpoint2(src, prefixTerms, sqrtPolicy2FromOld(p), digits)
}

// SqrtApproxTermsFromSourceRangeMidpointDefault is the default-policy wrapper
// around SqrtApproxTermsFromSourceRangeMidpoint.
func SqrtApproxTermsFromSourceRangeMidpointDefault(src ContinuedFraction, prefixTerms int, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromSourceRangeMidpointDefault: negative digits %d", digits)
	}
	return SqrtApproxTermsFromSourceRangeMidpointDefault2(src, prefixTerms, digits)
}

// sqrt_cf.go v
