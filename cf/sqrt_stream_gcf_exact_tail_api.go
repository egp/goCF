// sqrt_stream_gcf_exact_tail_api.go v1
package cf

import "fmt"

// SqrtGCFExactTailStream returns an inspectable bounded sqrt stream over a
// generalized continued-fraction source with exact tail evidence.
//
// Current milestone:
//   - bounded GCF ingestion
//   - exact tail evidence only
//   - CF output via exact rational collapse
//
// Future work:
//   - richer tail evidence
//   - stronger progressive certification
func SqrtGCFExactTailStream(src GCFSource, tailSrc GCFTailSource, maxIngestTerms int, p SqrtPolicy) (SqrtApproxStream, error) {
	if maxIngestTerms == 0 {
		return nil, fmt.Errorf("SqrtGCFExactTailStream: maxIngestTerms must be nonzero")
	}
	return NewSqrtGCFExactTailStream2(src, tailSrc, maxIngestTerms, sqrtPolicy2FromOld(p)), nil
}

func SqrtGCFExactTailStreamWithTail(src GCFSource, tail Rational, maxIngestTerms int, p SqrtPolicy) (SqrtApproxStream, error) {
	if maxIngestTerms == 0 {
		return nil, fmt.Errorf("SqrtGCFExactTailStreamWithTail: maxIngestTerms must be nonzero")
	}
	return NewSqrtGCFExactTailStreamWithTail2(src, tail, maxIngestTerms, sqrtPolicy2FromOld(p)), nil
}

// sqrt_stream_gcf_exact_tail_api.go v1
