// sqrt_cfsource_range_seed.go v1
package cf

import "fmt"

// NewSqrtApproxCFFromSourceRangeSeed consumes a finite prefix of src, converts
// that prefix to:
//
//   - a rational approximation (the convergent) for x
//   - a range-derived seed for sqrt(x)
//
// and then returns a ContinuedFraction source for the bounded sqrt approximation
// under the supplied policy.
//
// If p.Seed is already set, it is honored and the range-derived seed is not used.
func NewSqrtApproxCFFromSourceRangeSeed(src ContinuedFraction, prefixTerms int, p SqrtPolicy) (ContinuedFraction, error) {
	xApprox, rng, err := ApproxFromCFPrefix(src, prefixTerms)
	if err != nil {
		return nil, err
	}

	pp := p
	if pp.Seed == nil {
		seed, err := DefaultSqrtSeedFromRange(rng)
		if err != nil {
			return nil, err
		}
		pp.Seed = &seed
	}

	return SqrtApproxCFWithPolicy(xApprox, pp)
}

// NewSqrtApproxCFFromSourceRangeSeedDefault is the default-policy wrapper
// around NewSqrtApproxCFFromSourceRangeSeed.
func NewSqrtApproxCFFromSourceRangeSeedDefault(src ContinuedFraction, prefixTerms int) (ContinuedFraction, error) {
	return NewSqrtApproxCFFromSourceRangeSeed(src, prefixTerms, DefaultSqrtPolicy())
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

// sqrt_cfsource_range_seed.go v1
