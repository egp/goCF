// sqrt_stream_gcf_api.go v2
package cf

import "fmt"

// SqrtGCFStream returns a bounded sqrt stream over a generalized
// continued-fraction source.
//
// Current milestone:
//   - bounded prefix ingestion
//   - range-seeded bounded approximation
//   - CF output via exact rational collapse
//
// Future work:
//   - stronger progressive certification
//   - exact-tail and richer tail-evidence public variants
func SqrtGCFStream(src GCFSource, prefixTerms int, p SqrtPolicy) (SqrtApproxStream, error) {
	if prefixTerms <= 0 {
		return nil, fmt.Errorf("SqrtGCFStream: prefixTerms must be > 0, got %d", prefixTerms)
	}
	return NewSqrtGCFPrefixStream2(src, prefixTerms, sqrtPolicy2FromOld(p)), nil
}

// sqrt_stream_gcf_api.go v2
