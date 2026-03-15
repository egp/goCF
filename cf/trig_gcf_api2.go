// trig_gcf_api2.go v3
package cf

import "fmt"

// SinBoundsDegreesFromGCFPrefix2 consumes exactly prefixTerms from src and applies
// the current degree-based sin bound kernel to the resulting exact finite angle.
//
// Exactness rule for this entry point:
//   - after consuming prefixTerms terms, the source must be exhausted
//   - otherwise the angle is not yet exact
//
// Current intended use:
//   - finite exact degree angles such as 69° represented as GCF input
//   - regular CF adapted to GCF via AdaptCFToGCF
func SinBoundsDegreesFromGCFPrefix2(src GCFSource, prefixTerms int) (Range, error) {
	if prefixTerms <= 0 {
		return Range{}, fmt.Errorf(
			"SinBoundsDegreesFromGCFPrefix2: prefixTerms must be > 0, got %d",
			prefixTerms,
		)
	}

	b, err := IngestGCFPrefix(src, prefixTerms)
	if err != nil {
		return Range{}, err
	}
	if !b.HasValue() {
		return Range{}, fmt.Errorf("SinBoundsDegreesFromGCFPrefix2: empty source")
	}

	// For this exact-angle entry point, the source must be exhausted after the
	// requested prefix. Peek one more term to determine that.
	if _, _, ok := src.NextPQ(); ok {
		return Range{}, fmt.Errorf(
			"SinBoundsDegreesFromGCFPrefix2: angle not exact at prefixTerms=%d",
			prefixTerms,
		)
	}

	angleValue, err := b.Convergent()
	if err != nil {
		return Range{}, err
	}
	return SinBoundsDegrees(Degrees(angleValue))
}

// SinApproxDegreesFromGCFPrefix2 is the point-result wrapper around
// SinBoundsDegreesFromGCFPrefix2.
func SinApproxDegreesFromGCFPrefix2(src GCFSource, prefixTerms int) (Rational, error) {
	r, err := SinBoundsDegreesFromGCFPrefix2(src, prefixTerms)
	if err != nil {
		return Rational{}, err
	}
	if r.Lo.Cmp(r.Hi) != 0 {
		return Rational{}, fmt.Errorf("SinApproxDegreesFromGCFPrefix2: bounded non-point result")
	}
	return r.Lo, nil
}

// TanhBoundsSpecialFromGCF2 is the current GCF-ingesting entry point for special
// tanh bounds driven by source metadata.
//
// Current v1 support:
//   - sources that advertise themselves as sqrt(5), including AdaptCFToGCF(Sqrt5CF())
func TanhBoundsSpecialFromGCF2(src GCFSource) (Range, error) {
	qr, ok := src.(interface {
		Radicand() (int64, bool)
	})
	if !ok {
		return Range{}, fmt.Errorf("TanhBoundsSpecialFromGCF2: source has no quadratic-radical metadata")
	}

	n, ok := qr.Radicand()
	if !ok {
		return Range{}, fmt.Errorf("TanhBoundsSpecialFromGCF2: source has no quadratic-radical metadata")
	}
	if n != 5 {
		return Range{}, fmt.Errorf("TanhBoundsSpecialFromGCF2: only sqrt(5) is currently supported, got sqrt(%d)", n)
	}

	return TanhBoundsSqrt5(), nil
}

// trig_gcf_api2.go v3
