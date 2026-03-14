// sqrt_gcf_tail_api2.go v1
package cf

import "fmt"

// SqrtApproxFromGCFTailSource2 evaluates a bounded finite GCF prefix using
// explicit tail evidence, then computes a bounded sqrt approximation under the
// supplied policy.
//
// Current implementation supports only exact tail evidence.
func SqrtApproxFromGCFTailSource2(src GCFSource, tailSrc GCFTailSource, maxIngestTerms int, p SqrtPolicy2) (Rational, error) {
	tail, ok := tailSrc.ExactTail()
	if !ok {
		return Rational{}, fmt.Errorf("SqrtApproxFromGCFTailSource2: tail evidence not implemented")
	}
	return SqrtApproxFromGCFWithTail2(src, tail, maxIngestTerms, p)
}

// SqrtApproxCFFromGCFTailSource2 evaluates a bounded finite GCF prefix using
// explicit tail evidence, then returns a ContinuedFraction source for the
// bounded sqrt approximation under the supplied policy.
func SqrtApproxCFFromGCFTailSource2(src GCFSource, tailSrc GCFTailSource, maxIngestTerms int, p SqrtPolicy2) (ContinuedFraction, error) {
	approx, err := SqrtApproxFromGCFTailSource2(src, tailSrc, maxIngestTerms, p)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

// SqrtApproxTermsFromGCFTailSource2 evaluates a bounded finite GCF prefix using
// explicit tail evidence, then returns up to digits CF terms for the bounded
// sqrt approximation under the supplied policy.
func SqrtApproxTermsFromGCFTailSource2(src GCFSource, tailSrc GCFTailSource, maxIngestTerms int, p SqrtPolicy2, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromGCFTailSource2: negative digits %d", digits)
	}
	cf, err := SqrtApproxCFFromGCFTailSource2(src, tailSrc, maxIngestTerms, p)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// SqrtApproxFromGCFTailSourceDefault2 is the default-policy wrapper around
// SqrtApproxFromGCFTailSource2.
func SqrtApproxFromGCFTailSourceDefault2(src GCFSource, tailSrc GCFTailSource, maxIngestTerms int) (Rational, error) {
	return SqrtApproxFromGCFTailSource2(src, tailSrc, maxIngestTerms, DefaultSqrtPolicy2())
}

// SqrtApproxCFFromGCFTailSourceDefault2 is the default-policy wrapper around
// SqrtApproxCFFromGCFTailSource2.
func SqrtApproxCFFromGCFTailSourceDefault2(src GCFSource, tailSrc GCFTailSource, maxIngestTerms int) (ContinuedFraction, error) {
	return SqrtApproxCFFromGCFTailSource2(src, tailSrc, maxIngestTerms, DefaultSqrtPolicy2())
}

// SqrtApproxTermsFromGCFTailSourceDefault2 is the default-policy wrapper around
// SqrtApproxTermsFromGCFTailSource2.
func SqrtApproxTermsFromGCFTailSourceDefault2(src GCFSource, tailSrc GCFTailSource, maxIngestTerms int, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTermsFromGCFTailSourceDefault2: negative digits %d", digits)
	}
	return SqrtApproxTermsFromGCFTailSource2(src, tailSrc, maxIngestTerms, DefaultSqrtPolicy2(), digits)
}

// sqrt_gcf_tail_api2.go v1
