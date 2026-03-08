// sqrt_cfsource_range_seed.go v2
package cf

import "fmt"

// NewSqrtApproxCFFromApproxRangeSeed takes a bundled CFApprox and returns a
// ContinuedFraction source for the bounded sqrt approximation under the supplied
// policy.
//
// If p.Seed is already set, it is honored and the range-derived seed is not used.
func NewSqrtApproxCFFromApproxRangeSeed(a CFApprox, p SqrtPolicy) (ContinuedFraction, error) {
	pp := p
	if pp.Seed == nil {
		seed, err := DefaultSqrtSeedFromRange(a.Range)
		if err != nil {
			return nil, err
		}
		pp.Seed = &seed
	}
	return SqrtApproxCFWithPolicy(a.Convergent, pp)
}

// NewSqrtApproxCFFromSourceRangeSeed consumes a finite prefix of src, converts
// that prefix to a CFApprox, and then returns a ContinuedFraction source for the
// bounded sqrt approximation under the supplied policy.
func NewSqrtApproxCFFromSourceRangeSeed(src ContinuedFraction, prefixTerms int, p SqrtPolicy) (ContinuedFraction, error) {
	a, err := CFApproxFromPrefix(src, prefixTerms)
	if err != nil {
		return nil, err
	}
	return NewSqrtApproxCFFromApproxRangeSeed(a, p)
}

// NewSqrtApproxCFFromSourceRangeSeedDefault is the default-policy wrapper
// around NewSqrtApproxCFFromSourceRangeSeed.
func NewSqrtApproxCFFromSourceRangeSeedDefault(src ContinuedFraction, prefixTerms int) (ContinuedFraction, error) {
	return NewSqrtApproxCFFromSourceRangeSeed(src, prefixTerms, DefaultSqrtPolicy())
}

// SqrtApproxTermsFromApproxRangeSeed returns up to digits CF terms for the
// bounded sqrt approximation produced by NewSqrtApproxCFFromApproxRangeSeed.
func SqrtApproxTermsFromApproxRangeSeed(a CFApprox, p SqrtPolicy, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromApproxRangeSeed: negative digits %d", digits)
	}
	cf, err := NewSqrtApproxCFFromApproxRangeSeed(a, p)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// SqrtApproxTermsFromSourceRangeSeed returns up to digits CF terms for the
// bounded sqrt approximation produced by NewSqrtApproxCFFromSourceRangeSeed.
func SqrtApproxTermsFromSourceRangeSeed(src ContinuedFraction, prefixTerms int, p SqrtPolicy, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromSourceRangeSeed: negative digits %d", digits)
	}
	cf, err := NewSqrtApproxCFFromSourceRangeSeed(src, prefixTerms, p)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// SqrtApproxTermsFromSourceRangeSeedDefault is the default-policy wrapper
// around SqrtApproxTermsFromSourceRangeSeed.
func SqrtApproxTermsFromSourceRangeSeedDefault(src ContinuedFraction, prefixTerms int, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromSourceRangeSeedDefault: negative digits %d", digits)
	}
	return SqrtApproxTermsFromSourceRangeSeed(src, prefixTerms, DefaultSqrtPolicy(), digits)
}

// sqrt_cfsource_range_seed.go v2
