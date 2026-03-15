// reciprocal_stream_gcf_prefix_api.go v1
package cf

import "fmt"

// ReciprocalGCFPrefixStream returns an inspectable unary reciprocal stream over
// a bounded generalized continued-fraction prefix.
func ReciprocalGCFPrefixStream(src GCFSource, prefixTerms int) (ReciprocalApproxStream, error) {
	if prefixTerms <= 0 {
		return nil, fmt.Errorf("ReciprocalGCFPrefixStream: prefixTerms must be > 0, got %d", prefixTerms)
	}
	return NewReciprocalGCFPrefixStream2(src, prefixTerms), nil
}

// reciprocal_stream_gcf_prefix_api.go v1
