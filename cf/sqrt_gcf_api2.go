// sqrt_gcf_api2.go v1
package cf

import "fmt"

// SqrtApproxFromGCFWithTail2 evaluates a bounded finite GCF prefix with an
// explicit exact tail, then computes a bounded sqrt approximation under the
// supplied policy.
func SqrtApproxFromGCFWithTail2(src GCFSource, tail Rational, maxIngestTerms int, p SqrtPolicy2) (Rational, error) {
	x, _, err := EvalGCFWithTailExact(src, tail, maxIngestTerms)
	if err != nil {
		return Rational{}, err
	}
	return SqrtApproxWithPolicy2(x, p)
}

// SqrtApproxCFFromGCFWithTail2 evaluates a bounded finite GCF prefix with an
// explicit exact tail, then returns a ContinuedFraction source for the bounded
// sqrt approximation under the supplied policy.
func SqrtApproxCFFromGCFWithTail2(src GCFSource, tail Rational, maxIngestTerms int, p SqrtPolicy2) (ContinuedFraction, error) {
	approx, err := SqrtApproxFromGCFWithTail2(src, tail, maxIngestTerms, p)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

// SqrtApproxTermsFromGCFWithTail2 evaluates a bounded finite GCF prefix with an
// explicit exact tail, then returns up to digits CF terms for the bounded sqrt
// approximation under the supplied policy.
func SqrtApproxTermsFromGCFWithTail2(src GCFSource, tail Rational, maxIngestTerms int, p SqrtPolicy2, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromGCFWithTail2: negative digits %d", digits)
	}
	cf, err := SqrtApproxCFFromGCFWithTail2(src, tail, maxIngestTerms, p)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// SqrtApproxFromGCFWithTailDefault2 is the default-policy wrapper around
// SqrtApproxFromGCFWithTail2.
func SqrtApproxFromGCFWithTailDefault2(src GCFSource, tail Rational, maxIngestTerms int) (Rational, error) {
	return SqrtApproxFromGCFWithTail2(src, tail, maxIngestTerms, DefaultSqrtPolicy2())
}

// SqrtApproxCFFromGCFWithTailDefault2 is the default-policy wrapper around
// SqrtApproxCFFromGCFWithTail2.
func SqrtApproxCFFromGCFWithTailDefault2(src GCFSource, tail Rational, maxIngestTerms int) (ContinuedFraction, error) {
	return SqrtApproxCFFromGCFWithTail2(src, tail, maxIngestTerms, DefaultSqrtPolicy2())
}

// SqrtApproxTermsFromGCFWithTailDefault2 is the default-policy wrapper around
// SqrtApproxTermsFromGCFWithTail2.
func SqrtApproxTermsFromGCFWithTailDefault2(src GCFSource, tail Rational, maxIngestTerms int, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromGCFWithTailDefault2: negative digits %d", digits)
	}
	return SqrtApproxTermsFromGCFWithTail2(src, tail, maxIngestTerms, DefaultSqrtPolicy2(), digits)
}

// sqrt_gcf_api2.go v1
