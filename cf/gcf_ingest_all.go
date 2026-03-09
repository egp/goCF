// gcf_ingest_all.go v3
package cf

import "fmt"

// IngestAllGCF drains a finite generalized continued-fraction source into a new
// GCFBounder and then marks it finished.
//
// This is a convenience helper for exact finite-prefix ingestion.
func IngestAllGCF(src GCFSource) (*GCFBounder, error) {
	b := NewGCFBounder()

	for {
		p, q, ok := src.NextPQ()
		if !ok {
			break
		}
		if err := b.IngestPQ(p, q); err != nil {
			return nil, err
		}
	}

	b.Finish()
	return b, nil
}

// IngestGCFPrefix ingests up to prefixTerms terms from a GCFSource into a new
// GCFBounder.
//
// Behavior:
//   - prefixTerms < 0 => error
//   - prefixTerms == 0 => returns an empty bounder
//   - if the source terminates before prefixTerms terms, the returned bounder is finished
//   - otherwise the returned bounder contains exactly prefixTerms ingested terms
//     and is not yet finished
//
// If src also implements PositiveTailLowerBoundedGCFSource, the returned bounder
// is configured with that lower bound so unfinished prefixes can produce a
// conservative ray-based enclosure.
func IngestGCFPrefix(src GCFSource, prefixTerms int) (*GCFBounder, error) {
	if prefixTerms < 0 {
		return nil, fmt.Errorf("IngestGCFPrefix: negative prefixTerms %d", prefixTerms)
	}

	b := NewGCFBounder()

	// Only apply stable lower-bound metadata automatically.
	//
	// Do NOT auto-apply TailRangeBoundedGCFSource here: for stateful sources,
	// a source-level TailRange() is not stable across consumed prefixes unless
	// the contract is made prefix-aware.
	if bounded, ok := src.(PositiveTailLowerBoundedGCFSource); ok {
		if err := b.SetTailLowerBound(bounded.TailLowerBound()); err != nil {
			return nil, err
		}
	}

	if prefixTerms == 0 {
		return b, nil
	}

	for i := 0; i < prefixTerms; i++ {
		p, q, ok := src.NextPQ()
		if !ok {
			b.Finish()
			return b, nil
		}
		if err := b.IngestPQ(p, q); err != nil {
			return nil, err
		}
	}

	return b, nil
}

// gcf_ingest_all.go v3
