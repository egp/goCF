// gcf_approx.go v1
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

// gcf_approx.go v1
