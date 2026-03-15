// sqrt_stream_api.go v4
package cf

import "fmt"

// SqrtStream returns a bounded sqrt stream over a continued-fraction source.
//
// Current milestone:
//   - bounded prefix ingestion
//   - range-seeded bounded approximation
//   - CF output via exact rational collapse
//
// Future work:
//   - stronger progressive certification
//   - less collapse, more direct operator behavior
func SqrtStream(src ContinuedFraction, prefixTerms int, p SqrtPolicy) (SqrtApproxStream, error) {
	if prefixTerms <= 0 {
		return nil, fmt.Errorf("SqrtStream: prefixTerms must be > 0, got %d", prefixTerms)
	}
	return NewSqrtCFPrefixStream2(src, prefixTerms, sqrtPolicy2FromOld(p)), nil
}

// sqrt_stream_api.go v4
