// sqrt_source_api2.go v1
package cf

import "fmt"

// SqrtApproxFromApproxRangeSeed2 takes a bundled CFApprox, derives a seed from
// its range when needed, and returns a bounded rational sqrt approximation
// under the supplied policy.
//
// If p.Seed is already set, it is honored and the range-derived seed is not used.
func SqrtApproxFromApproxRangeSeed2(a CFApprox, p SqrtPolicy2) (Rational, error) {
	if err := p.Validate(); err != nil {
		return Rational{}, err
	}

	pp := p
	if pp.Seed == nil {
		seed, err := SqrtSeedFromRange(a.Range)
		if err != nil {
			return Rational{}, err
		}
		pp.Seed = &seed
	}

	return SqrtApproxWithPolicy2(a.Convergent, pp)
}

// SqrtApproxCFFromApproxRangeSeed2 returns a ContinuedFraction source for the
// bounded sqrt approximation produced by SqrtApproxFromApproxRangeSeed2.
func SqrtApproxCFFromApproxRangeSeed2(a CFApprox, p SqrtPolicy2) (ContinuedFraction, error) {
	approx, err := SqrtApproxFromApproxRangeSeed2(a, p)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

// SqrtApproxTermsFromApproxRangeSeed2 returns up to digits CF terms for the
// bounded sqrt approximation produced by SqrtApproxCFFromApproxRangeSeed2.
func SqrtApproxTermsFromApproxRangeSeed2(a CFApprox, p SqrtPolicy2, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromApproxRangeSeed2: negative digits %d", digits)
	}
	cf, err := SqrtApproxCFFromApproxRangeSeed2(a, p)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// sqrt_source_api2.go v1
