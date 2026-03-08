// sqrt_cfsource.go v1
package cf

import "fmt"

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
	if prefixTerms < 0 {
		return Rational{}, fmt.Errorf("RationalApproxFromCFPrefix: negative prefixTerms %d", prefixTerms)
	}
	if prefixTerms == 0 {
		return Rational{}, fmt.Errorf("RationalApproxFromCFPrefix: prefixTerms must be > 0")
	}

	b := NewBounder()

	for i := 0; i < prefixTerms; i++ {
		a, ok := src.Next()
		if !ok {
			if !b.HasValue() {
				return Rational{}, fmt.Errorf("RationalApproxFromCFPrefix: empty source")
			}
			b.Finish()
			return b.Convergent()
		}
		if err := b.Ingest(a); err != nil {
			return Rational{}, err
		}
	}

	if !b.HasValue() {
		return Rational{}, fmt.Errorf("RationalApproxFromCFPrefix: empty source")
	}
	return b.Convergent()
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

// sqrt_cfsource.go v1
