// sqrt_canonical_source_api.go v1
package cf

import "fmt"

func sqrtApproxFromApproxRangeSeedCanonical(a CFApprox, p SqrtPolicy2) (Rational, error) {
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

	return sqrtApproxWithPolicyCanonical(a.Convergent, pp)
}

func sqrtApproxCFFromApproxRangeSeedCanonical(a CFApprox, p SqrtPolicy2) (ContinuedFraction, error) {
	approx, err := sqrtApproxFromApproxRangeSeedCanonical(a, p)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

func sqrtApproxTermsFromApproxRangeSeedCanonical(a CFApprox, p SqrtPolicy2, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("sqrtApproxTermsFromApproxRangeSeedCanonical: negative digits %d", digits)
	}
	cf, err := sqrtApproxCFFromApproxRangeSeedCanonical(a, p)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

func sqrtApproxFromSourceRangeSeedCanonical(src ContinuedFraction, prefixTerms int, p SqrtPolicy2) (Rational, error) {
	a, err := CFApproxFromPrefix(src, prefixTerms)
	if err != nil {
		return Rational{}, err
	}
	return sqrtApproxFromApproxRangeSeedCanonical(a, p)
}

func sqrtApproxCFFromSourceRangeSeedCanonical(src ContinuedFraction, prefixTerms int, p SqrtPolicy2) (ContinuedFraction, error) {
	approx, err := sqrtApproxFromSourceRangeSeedCanonical(src, prefixTerms, p)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

func sqrtApproxTermsFromSourceRangeSeedCanonical(src ContinuedFraction, prefixTerms int, p SqrtPolicy2, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("sqrtApproxTermsFromSourceRangeSeedCanonical: negative digits %d", digits)
	}
	cf, err := sqrtApproxCFFromSourceRangeSeedCanonical(src, prefixTerms, p)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// sqrt_canonical_source_api.go v1
