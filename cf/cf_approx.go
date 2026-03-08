// cf_approx.go v1
package cf

import "fmt"

// CFApprox bundles the exact convergent and enclosure implied by a finite CF prefix.
type CFApprox struct {
	Convergent  Rational
	Range       Range
	PrefixTerms int
}

// ApproxFromCFPrefix ingests up to prefixTerms terms from src and returns both:
//
//   - the current convergent as an exact Rational
//   - the current enclosure Range implied by that finite prefix
//
// Behavior:
//   - prefixTerms < 0 => error
//   - prefixTerms == 0 => error
//   - if src terminates early, the bounder is finished and both values are exact
func ApproxFromCFPrefix(src ContinuedFraction, prefixTerms int) (Rational, Range, error) {
	approx, err := CFApproxFromPrefix(src, prefixTerms)
	if err != nil {
		return Rational{}, Range{}, err
	}
	return approx.Convergent, approx.Range, nil
}

// CFApproxFromPrefix ingests up to prefixTerms terms from src and returns a bundled CFApprox.
func CFApproxFromPrefix(src ContinuedFraction, prefixTerms int) (CFApprox, error) {
	if prefixTerms < 0 {
		return CFApprox{}, fmt.Errorf("CFApproxFromPrefix: negative prefixTerms %d", prefixTerms)
	}
	if prefixTerms == 0 {
		return CFApprox{}, fmt.Errorf("CFApproxFromPrefix: prefixTerms must be > 0")
	}

	b := NewBounder()

	for i := 0; i < prefixTerms; i++ {
		a, ok := src.Next()
		if !ok {
			if !b.HasValue() {
				return CFApprox{}, fmt.Errorf("CFApproxFromPrefix: empty source")
			}
			b.Finish()
			conv, err := b.Convergent()
			if err != nil {
				return CFApprox{}, err
			}
			rng, ok, err := b.Range()
			if err != nil {
				return CFApprox{}, err
			}
			if !ok {
				return CFApprox{}, fmt.Errorf("CFApproxFromPrefix: internal: no range after finish")
			}
			return CFApprox{
				Convergent:  conv,
				Range:       rng,
				PrefixTerms: prefixTerms,
			}, nil
		}
		if err := b.Ingest(a); err != nil {
			return CFApprox{}, err
		}
	}

	if !b.HasValue() {
		return CFApprox{}, fmt.Errorf("CFApproxFromPrefix: empty source")
	}

	conv, err := b.Convergent()
	if err != nil {
		return CFApprox{}, err
	}
	rng, ok, err := b.Range()
	if err != nil {
		return CFApprox{}, err
	}
	if !ok {
		return CFApprox{}, fmt.Errorf("CFApproxFromPrefix: internal: no range after ingest")
	}
	return CFApprox{
		Convergent:  conv,
		Range:       rng,
		PrefixTerms: prefixTerms,
	}, nil
}

// cf_approx.go v1
