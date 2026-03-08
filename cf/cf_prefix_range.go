// cf_prefix_range.go v1
package cf

import "fmt"

// RangeApproxFromCFPrefix ingests up to prefixTerms terms from src and returns
// the enclosure Range implied by that finite prefix.
//
// Behavior:
//   - prefixTerms < 0 => error
//   - prefixTerms == 0 => error
//   - if src terminates early, the bounder is finished and the returned range
//     collapses to the exact rational represented by the finite source
//   - if src does not terminate within prefixTerms, the returned range is the
//     current enclosure produced by Bounder.Range()
func RangeApproxFromCFPrefix(src ContinuedFraction, prefixTerms int) (Range, error) {
	if prefixTerms < 0 {
		return Range{}, fmt.Errorf("RangeApproxFromCFPrefix: negative prefixTerms %d", prefixTerms)
	}
	if prefixTerms == 0 {
		return Range{}, fmt.Errorf("RangeApproxFromCFPrefix: prefixTerms must be > 0")
	}

	b := NewBounder()

	for i := 0; i < prefixTerms; i++ {
		a, ok := src.Next()
		if !ok {
			if !b.HasValue() {
				return Range{}, fmt.Errorf("RangeApproxFromCFPrefix: empty source")
			}
			b.Finish()
			rng, ok, err := b.Range()
			if err != nil {
				return Range{}, err
			}
			if !ok {
				return Range{}, fmt.Errorf("RangeApproxFromCFPrefix: internal: no range after finish")
			}
			return rng, nil
		}
		if err := b.Ingest(a); err != nil {
			return Range{}, err
		}
	}

	if !b.HasValue() {
		return Range{}, fmt.Errorf("RangeApproxFromCFPrefix: empty source")
	}

	rng, ok, err := b.Range()
	if err != nil {
		return Range{}, err
	}
	if !ok {
		return Range{}, fmt.Errorf("RangeApproxFromCFPrefix: internal: no range after ingest")
	}
	return rng, nil
}

// cf_prefix_range.go v1
