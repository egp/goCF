// gcf_diag_stream.go v2
package cf

import "fmt"

// GCFDiagStream is the GCF-native exact-tail stream for diagonal transforms.
//
// Current supported milestone:
//   - finite GCF prefix ingestion
//   - exact tail evidence supplied by caller
//   - output regular CF digits of the exact transformed value
//
// Infinite uncertified GCF emission is intentionally out of scope for now.
type GCFDiagStream struct {
	t       DiagBLFT
	src     GCFSource
	tailSrc GCFTailSource
	opts    GCFULFTStreamOptions
	err     error
	done    bool
	started bool
	exactCF ContinuedFraction
}

// NewGCFDiagStream constructs a GCF-native diagonal stream using explicit tail
// evidence.
//
// Current implementation supports only exact tail evidence.
func NewGCFDiagStream(
	t DiagBLFT,
	src GCFSource,
	tailSrc GCFTailSource,
	opts GCFULFTStreamOptions,
) *GCFDiagStream {
	return &GCFDiagStream{
		t:       t,
		src:     src,
		tailSrc: tailSrc,
		opts:    opts,
	}
}

// NewGCFDiagStreamWithTail is a convenience wrapper for the current exact-tail path.
func NewGCFDiagStreamWithTail(
	t DiagBLFT,
	src GCFSource,
	tail Rational,
	opts GCFULFTStreamOptions,
) *GCFDiagStream {
	return NewGCFDiagStream(t, src, NewExactTailSource(tail), opts)
}

func (s *GCFDiagStream) Err() error { return s.err }

func (s *GCFDiagStream) initializeExactCF() bool {
	if s.started {
		return s.err == nil
	}
	s.started = true

	tail, ok := s.tailSrc.ExactTail()
	if !ok {
		s.err = fmt.Errorf("GCFDiagStream: tail evidence not implemented")
		s.done = true
		return false
	}

	y, _, err := ApplyComposedGCFDiagBLFTToTailExact(s.t, s.src, tail, s.opts.MaxIngestTerms)
	if err != nil {
		s.err = fmt.Errorf("GCFDiagStream: %w", err)
		s.done = true
		return false
	}

	s.exactCF = NewRationalCF(y)
	return true
}

func (s *GCFDiagStream) Next() (int64, bool) {
	if s.done {
		return 0, false
	}
	if s.err != nil {
		s.done = true
		return 0, false
	}
	if !s.initializeExactCF() {
		return 0, false
	}

	d, ok := s.exactCF.Next()
	if !ok {
		s.done = true
		return 0, false
	}
	return d, true
}

// gcf_diag_stream.go v2
