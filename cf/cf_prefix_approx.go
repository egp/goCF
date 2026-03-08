// cf_prefix_approx.go v1
package cf

import "fmt"

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
	if prefixTerms < 0 {
		return Rational{}, Range{}, fmt.Errorf("ApproxFromCFPrefix: negative prefixTerms %d", prefixTerms)
	}
	if prefixTerms == 0 {
		return Rational{}, Range{}, fmt.Errorf("ApproxFromCFPrefix: prefixTerms must be > 0")
	}

	b := NewBounder()

	for i := 0; i < prefixTerms; i++ {
		a, ok := src.Next()
		if !ok {
			if !b.HasValue() {
				return Rational{}, Range{}, fmt.Errorf("ApproxFromCFPrefix: empty source")
			}
			b.Finish()
			conv, err := b.Convergent()
			if err != nil {
				return Rational{}, Range{}, err
			}
			rng, ok, err := b.Range()
			if err != nil {
				return Rational{}, Range{}, err
			}
			if !ok {
				return Rational{}, Range{}, fmt.Errorf("ApproxFromCFPrefix: internal: no range after finish")
			}
			return conv, rng, nil
		}
		if err := b.Ingest(a); err != nil {
			return Rational{}, Range{}, err
		}
	}

	if !b.HasValue() {
		return Rational{}, Range{}, fmt.Errorf("ApproxFromCFPrefix: empty source")
	}

	conv, err := b.Convergent()
	if err != nil {
		return Rational{}, Range{}, err
	}
	rng, ok, err := b.Range()
	if err != nil {
		return Rational{}, Range{}, err
	}
	if !ok {
		return Rational{}, Range{}, fmt.Errorf("ApproxFromCFPrefix: internal: no range after ingest")
	}
	return conv, rng, nil
}

// cf_prefix_approx.go v1
