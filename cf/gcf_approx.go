// gcf_approx.go v2
package cf

import "fmt"

// GCFApprox bundles the exact convergent implied by a finite GCF prefix.
type GCFApprox struct {
	Convergent  Rational
	PrefixTerms int
}

// GCFApproxFromPrefix ingests up to prefixTerms terms from src and returns a bundled GCFApprox.
//
// Behavior:
//   - prefixTerms < 0 => error
//   - prefixTerms == 0 => error
//   - if src terminates early, the returned convergent is exact for the whole finite source
//   - otherwise the returned convergent is exact for the finite prefix only
func GCFApproxFromPrefix(src GCFSource, prefixTerms int) (GCFApprox, error) {
	if prefixTerms < 0 {
		return GCFApprox{}, fmt.Errorf("GCFApproxFromPrefix: negative prefixTerms %d", prefixTerms)
	}
	if prefixTerms == 0 {
		return GCFApprox{}, fmt.Errorf("GCFApproxFromPrefix: prefixTerms must be > 0")
	}

	b, err := IngestGCFPrefix(src, prefixTerms)
	if err != nil {
		return GCFApprox{}, err
	}
	if !b.HasValue() {
		return GCFApprox{}, fmt.Errorf("GCFApproxFromPrefix: empty source")
	}

	conv, err := b.Convergent()
	if err != nil {
		return GCFApprox{}, err
	}

	return GCFApprox{
		Convergent:  conv,
		PrefixTerms: prefixTerms,
	}, nil
}

// GCFApproxCF returns a regular continued-fraction source for the exact rational
// convergent carried by the GCFApprox snapshot.
func GCFApproxCF(a GCFApprox) ContinuedFraction {
	return NewRationalCF(a.Convergent)
}

// GCFApproxTerms returns up to digits regular CF terms for the exact rational
// convergent carried by the GCFApprox snapshot.
func GCFApproxTerms(a GCFApprox, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("GCFApproxTerms: negative digits %d", digits)
	}
	return collectTerms(GCFApproxCF(a), digits), nil
}

// gcf_approx.go v2
