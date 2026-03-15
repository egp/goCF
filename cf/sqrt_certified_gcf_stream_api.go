// sqrt_certified_gcf_stream_api.go v1
package cf

import "fmt"

func SqrtCertifiedGCFStream(src GCFSource, maxPrefixTerms int) (SqrtApproxStream, error) {
	if maxPrefixTerms <= 0 {
		return nil, fmt.Errorf("SqrtCertifiedGCFStream: maxPrefixTerms must be > 0, got %d", maxPrefixTerms)
	}
	return NewSqrtCertifiedGCFPrefixStream(src, maxPrefixTerms)
}

// sqrt_certified_gcf_stream_api.go v1
