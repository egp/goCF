// trig_gcf_exact_tail_api.go v2
package cf

import "fmt"

// SinBoundsDegreesFromGCFWithTail2 evaluates a bounded GCF prefix together with
// an explicit exact tail, producing an exact rational degree angle, then applies
// the current degree-based sin bound kernel.
func SinBoundsDegreesFromGCFWithTail2(
	src GCFSource,
	tail Rational,
	maxIngestTerms int,
) (Range, error) {
	angleValue, _, err := EvalGCFWithTailExact(src, tail, maxIngestTerms)
	if err != nil {
		return Range{}, err
	}
	return SinBoundsDegrees(Degrees(angleValue))
}

// SinApproxDegreesFromGCFWithTail2 is the point-result wrapper around
// SinBoundsDegreesFromGCFWithTail2.
func SinApproxDegreesFromGCFWithTail2(
	src GCFSource,
	tail Rational,
	maxIngestTerms int,
) (Rational, error) {
	r, err := SinBoundsDegreesFromGCFWithTail2(src, tail, maxIngestTerms)
	if err != nil {
		return Rational{}, err
	}
	if r.Lo.Cmp(r.Hi) != 0 {
		return Rational{}, fmt.Errorf("SinApproxDegreesFromGCFWithTail2: bounded non-point result")
	}
	return r.Lo, nil
}

// trig_gcf_exact_tail_api.go v2
