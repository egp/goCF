// sqrt_gcf_range_seed_api2.go v1
package cf

import "fmt"

// SqrtApproxFromGCFApproxRangeSeed2 takes a GCFApprox, derives a seed from its
// range when available, and returns a bounded rational sqrt approximation under
// the supplied policy.
//
// Seed selection order:
//   - if p.Seed is set, honor it
//   - else if a.Range != nil, derive seed from range
//   - else fall back to SqrtSeedDefault(a.Convergent)
func SqrtApproxFromGCFApproxRangeSeed2(a GCFApprox, p SqrtPolicy2) (Rational, error) {
	if err := p.Validate(); err != nil {
		return Rational{}, err
	}

	pp := p
	if pp.Seed == nil {
		if a.Range != nil {
			seed, err := SqrtSeedFromRange(*a.Range)
			if err != nil {
				return Rational{}, err
			}
			pp.Seed = &seed
		} else {
			seed, err := SqrtSeedDefault(a.Convergent)
			if err != nil {
				return Rational{}, err
			}
			pp.Seed = &seed
		}
	}

	return SqrtApproxWithPolicy2(a.Convergent, pp)
}

// SqrtApproxCFFromGCFApproxRangeSeed2 returns a ContinuedFraction source for the
// bounded sqrt approximation produced by SqrtApproxFromGCFApproxRangeSeed2.
func SqrtApproxCFFromGCFApproxRangeSeed2(a GCFApprox, p SqrtPolicy2) (ContinuedFraction, error) {
	approx, err := SqrtApproxFromGCFApproxRangeSeed2(a, p)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

// SqrtApproxTermsFromGCFApproxRangeSeed2 returns up to digits CF terms for the
// bounded sqrt approximation produced by SqrtApproxCFFromGCFApproxRangeSeed2.
func SqrtApproxTermsFromGCFApproxRangeSeed2(a GCFApprox, p SqrtPolicy2, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromGCFApproxRangeSeed2: negative digits %d", digits)
	}
	cf, err := SqrtApproxCFFromGCFApproxRangeSeed2(a, p)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// SqrtApproxFromGCFSourceRangeSeed2 consumes a finite GCF prefix, forms a
// GCFApprox, then returns a bounded rational sqrt approximation under the
// supplied policy.
func SqrtApproxFromGCFSourceRangeSeed2(src GCFSource, prefixTerms int, p SqrtPolicy2) (Rational, error) {
	a, err := GCFApproxFromPrefix(src, prefixTerms)
	if err != nil {
		return Rational{}, err
	}
	return SqrtApproxFromGCFApproxRangeSeed2(a, p)
}

// SqrtApproxCFFromGCFSourceRangeSeed2 consumes a finite GCF prefix, forms a
// GCFApprox, then returns a ContinuedFraction source for the bounded sqrt
// approximation under the supplied policy.
func SqrtApproxCFFromGCFSourceRangeSeed2(src GCFSource, prefixTerms int, p SqrtPolicy2) (ContinuedFraction, error) {
	approx, err := SqrtApproxFromGCFSourceRangeSeed2(src, prefixTerms, p)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

// SqrtApproxTermsFromGCFSourceRangeSeed2 consumes a finite GCF prefix, forms a
// GCFApprox, then returns up to digits CF terms for the bounded sqrt
// approximation under the supplied policy.
func SqrtApproxTermsFromGCFSourceRangeSeed2(src GCFSource, prefixTerms int, p SqrtPolicy2, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromGCFSourceRangeSeed2: negative digits %d", digits)
	}
	cf, err := SqrtApproxCFFromGCFSourceRangeSeed2(src, prefixTerms, p)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// SqrtApproxFromGCFSourceRangeSeedDefault2 is the default-policy wrapper around
// SqrtApproxFromGCFSourceRangeSeed2.
func SqrtApproxFromGCFSourceRangeSeedDefault2(src GCFSource, prefixTerms int) (Rational, error) {
	return SqrtApproxFromGCFSourceRangeSeed2(src, prefixTerms, DefaultSqrtPolicy2())
}

// SqrtApproxCFFromGCFSourceRangeSeedDefault2 is the default-policy wrapper
// around SqrtApproxCFFromGCFSourceRangeSeed2.
func SqrtApproxCFFromGCFSourceRangeSeedDefault2(src GCFSource, prefixTerms int) (ContinuedFraction, error) {
	return SqrtApproxCFFromGCFSourceRangeSeed2(src, prefixTerms, DefaultSqrtPolicy2())
}

// SqrtApproxTermsFromGCFSourceRangeSeedDefault2 is the default-policy wrapper
// around SqrtApproxTermsFromGCFSourceRangeSeed2.
func SqrtApproxTermsFromGCFSourceRangeSeedDefault2(src GCFSource, prefixTerms int, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromGCFSourceRangeSeedDefault2: negative digits %d", digits)
	}
	return SqrtApproxTermsFromGCFSourceRangeSeed2(src, prefixTerms, DefaultSqrtPolicy2(), digits)
}

// sqrt_gcf_range_seed_api2.go v1
