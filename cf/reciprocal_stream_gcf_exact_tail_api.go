// reciprocal_stream_gcf_exact_tail_api.go v1
package cf

import "fmt"

// ReciprocalGCFExactTailStream returns an inspectable unary reciprocal stream
// over a generalized continued-fraction source with exact tail evidence.
func ReciprocalGCFExactTailStream(src GCFSource, tailSrc GCFTailSource, maxIngestTerms int) (ReciprocalApproxStream, error) {
	if maxIngestTerms == 0 {
		return nil, fmt.Errorf("ReciprocalGCFExactTailStream: maxIngestTerms must be nonzero")
	}
	return NewReciprocalGCFExactTailStream2(src, tailSrc, maxIngestTerms), nil
}

func ReciprocalGCFExactTailStreamWithTail(src GCFSource, tail Rational, maxIngestTerms int) (ReciprocalApproxStream, error) {
	if maxIngestTerms == 0 {
		return nil, fmt.Errorf("ReciprocalGCFExactTailStreamWithTail: maxIngestTerms must be nonzero")
	}
	return NewReciprocalGCFExactTailStreamWithTail2(src, tail, maxIngestTerms), nil
}

// reciprocal_stream_gcf_exact_tail_api.go v1
