// sqrt_midpoint_api2.go v1
package cf

import "fmt"

// SqrtApproxFromApproxRangeMidpoint2 takes a bundled CFApprox, uses the
// midpoint of its enclosure as the sqrt target, and returns a bounded rational
// sqrt approximation under the supplied policy.
//
// IMPORTANT:
//   - this is an experimental heuristic path
//   - it is not yet a proof-safe conservative sqrt operator
func SqrtApproxFromApproxRangeMidpoint2(a CFApprox, p SqrtPolicy2) (Rational, error) {
	if err := p.Validate(); err != nil {
		return Rational{}, err
	}

	m, err := SqrtRangeMidpoint(a.Range)
	if err != nil {
		return Rational{}, err
	}
	return SqrtApproxWithPolicy2(m, p)
}

func SqrtApproxCFFromApproxRangeMidpoint2(a CFApprox, p SqrtPolicy2) (ContinuedFraction, error) {
	approx, err := SqrtApproxFromApproxRangeMidpoint2(a, p)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

func SqrtApproxTermsFromApproxRangeMidpoint2(a CFApprox, p SqrtPolicy2, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromApproxRangeMidpoint2: negative digits %d", digits)
	}
	cf, err := SqrtApproxCFFromApproxRangeMidpoint2(a, p)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

func SqrtApproxFromSourceRangeMidpoint2(src ContinuedFraction, prefixTerms int, p SqrtPolicy2) (Rational, error) {
	a, err := CFApproxFromPrefix(src, prefixTerms)
	if err != nil {
		return Rational{}, err
	}
	return SqrtApproxFromApproxRangeMidpoint2(a, p)
}

func SqrtApproxCFFromSourceRangeMidpoint2(src ContinuedFraction, prefixTerms int, p SqrtPolicy2) (ContinuedFraction, error) {
	approx, err := SqrtApproxFromSourceRangeMidpoint2(src, prefixTerms, p)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

func SqrtApproxTermsFromSourceRangeMidpoint2(src ContinuedFraction, prefixTerms int, p SqrtPolicy2, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromSourceRangeMidpoint2: negative digits %d", digits)
	}
	cf, err := SqrtApproxCFFromSourceRangeMidpoint2(src, prefixTerms, p)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

func SqrtApproxCFFromSourceRangeMidpointDefault2(src ContinuedFraction, prefixTerms int) (ContinuedFraction, error) {
	return SqrtApproxCFFromSourceRangeMidpoint2(src, prefixTerms, DefaultSqrtPolicy2())
}

func SqrtApproxTermsFromSourceRangeMidpointDefault2(src ContinuedFraction, prefixTerms int, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromSourceRangeMidpointDefault2: negative digits %d", digits)
	}
	return SqrtApproxTermsFromSourceRangeMidpoint2(src, prefixTerms, DefaultSqrtPolicy2(), digits)
}

// sqrt_midpoint_api2.go v1
