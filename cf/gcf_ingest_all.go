// gcf_ingest_all.go v1
package cf

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

// gcf_ingest_all.go v1
