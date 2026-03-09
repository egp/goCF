// gcf_approx.go v5
package cf

import "fmt"

// GCFApprox bundles the exact convergent implied by a finite GCF prefix.
// If Range is non-nil, it is a conservative enclosure for that prefix under the
// current unfinished-tail semantics. If Range is nil, no conservative enclosure
// is currently available.
type GCFApprox struct {
	Convergent  Rational
	Range       *Range
	PrefixTerms int
}

// GCFApproxFromPrefix ingests up to prefixTerms terms from src and returns a bundled GCFApprox.
//
// Behavior:
//   - prefixTerms < 0 => error
//   - prefixTerms == 0 => error
//   - if src terminates early, the returned convergent is exact for the whole finite source
//   - otherwise the returned convergent is exact for the finite prefix only
//
// If a conservative enclosure is available from GCFBounder.Range(), it is stored in Range.
// Otherwise Range is nil.
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

	var rp *Range
	if r, ok, err := b.Range(); err != nil {
		return GCFApprox{}, err
	} else if ok {
		rr := r
		rp = &rr
	}

	return GCFApprox{
		Convergent:  conv,
		Range:       rp,
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

// GCFSourceTerms ingests up to prefixTerms terms from src, forms a GCFApprox,
// and returns up to digits regular CF terms for the exact rational convergent.
func GCFSourceTerms(src GCFSource, prefixTerms int, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("GCFSourceTerms: negative digits %d", digits)
	}

	a, err := GCFApproxFromPrefix(src, prefixTerms)
	if err != nil {
		return nil, err
	}
	return GCFApproxTerms(a, digits)
}

// GCFSourceConvergent ingests up to prefixTerms terms from src and returns the
// exact rational convergent of that bounded GCF prefix.
func GCFSourceConvergent(src GCFSource, prefixTerms int) (Rational, error) {
	a, err := GCFApproxFromPrefix(src, prefixTerms)
	if err != nil {
		return Rational{}, err
	}
	return a.Convergent, nil
}

// HasRange reports whether the GCFApprox carries a conservative enclosure.
func (a GCFApprox) HasRange() bool {
	return a.Range != nil
}

// ExactRange reports whether the GCFApprox range exists and is an exact point range.
func (a GCFApprox) ExactRange() bool {
	return a.Range != nil && a.Range.Lo.Cmp(a.Range.Hi) == 0
}

// RangeContainsConvergent reports whether the stored enclosure exists and contains
// the exact convergent.
func (a GCFApprox) RangeContainsConvergent() bool {
	return a.Range != nil && a.Range.Contains(a.Convergent)
}

// GCFInspect bundles a bounded GCF prefix snapshot together with regular CF terms
// of its exact rational convergent.
type GCFInspect struct {
	Approx GCFApprox
	Terms  []int64
}

// InspectGCFSource ingests up to prefixTerms terms from src, forms a GCFApprox,
// and returns that snapshot together with up to digits regular CF terms of the
// exact rational convergent.
func InspectGCFSource(src GCFSource, prefixTerms int, digits int) (GCFInspect, error) {
	if digits < 0 {
		return GCFInspect{}, fmt.Errorf("InspectGCFSource: negative digits %d", digits)
	}

	a, err := GCFApproxFromPrefix(src, prefixTerms)
	if err != nil {
		return GCFInspect{}, err
	}

	terms, err := GCFApproxTerms(a, digits)
	if err != nil {
		return GCFInspect{}, err
	}

	return GCFInspect{
		Approx: a,
		Terms:  terms,
	}, nil
}

// gcf_approx.go v5
