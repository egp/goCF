// gcf_stream_test_helpers_test.go v2
package cf

import "testing"

type finitePrefixGCFSource struct {
	src   GCFSource
	limit int
	n     int
}

func newFinitePrefixGCFSource(src GCFSource, limit int) *finitePrefixGCFSource {
	return &finitePrefixGCFSource{
		src:   src,
		limit: limit,
	}
}

func (s *finitePrefixGCFSource) NextPQ() (int64, int64, bool) {
	if s.n >= s.limit {
		return 0, 0, false
	}
	p, q, ok := s.src.NextPQ()
	if !ok {
		return 0, 0, false
	}
	s.n++
	return p, q, true
}

func collectFinitePrefixTerms(src GCFSource, n int) [][2]int64 {
	var out [][2]int64
	for i := 0; i < n; i++ {
		p, q, ok := src.NextPQ()
		if !ok {
			break
		}
		out = append(out, [2]int64{p, q})
	}
	return out
}

func exactDigitsFromFinitePrefix(
	t *testing.T,
	srcFactory func() GCFSource,
	prefixLen int,
	maxDigits int,
) []int64 {
	t.Helper()

	terms := collectFinitePrefixTerms(srcFactory(), prefixLen)
	rat, err := EvaluateFiniteGCF(NewSliceGCF(terms...))
	if err != nil {
		t.Fatalf("EvaluateFiniteGCF failed for prefixLen=%d: %v", prefixLen, err)
	}
	return collectTerms(NewRationalCF(rat), maxDigits)
}

// gcf_stream_test_helpers_test.go v2
