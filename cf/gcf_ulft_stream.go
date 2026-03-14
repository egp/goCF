// gcf_ulft_stream.go v5
package cf

import "fmt"

// GCFULFTStream is the new GCF-native unary/LFT path.
//
// Current supported milestone:
//   - finite GCF prefix ingestion
//   - exact tail evidence supplied by caller
//   - output regular CF digits of the exact transformed value
//
// Infinite uncertified GCF emission is intentionally out of scope for now.
type GCFULFTStream struct {
	t       ULFT
	src     GCFSource
	tailSrc GCFTailSource
	opts    GCFULFTStreamOptions
	err     error
	done    bool
	started bool
	exactCF ContinuedFraction
}

// GCFULFTStreamOptions is intentionally small.
// Current meaning:
//   - MaxIngestTerms < 0 : unlimited
//   - MaxIngestTerms = 0 : no source ingestion allowed
//   - MaxIngestTerms > 0 : maximum number of source terms allowed before exhaustion
type GCFULFTStreamOptions struct {
	MaxIngestTerms int
}

// NewGCFULFTStream constructs a GCF-native ULFT stream using explicit tail
// evidence.
//
// Current implementation supports only exact tail evidence.
// Other evidence modes are intentionally deferred.
func NewGCFULFTStream(
	t ULFT,
	src GCFSource,
	tailSrc GCFTailSource,
	opts GCFULFTStreamOptions,
) *GCFULFTStream {
	return &GCFULFTStream{
		t:       t,
		src:     src,
		tailSrc: tailSrc,
		opts:    opts,
	}
}

// NewGCFULFTStreamWithTail is a convenience wrapper for the current exact-tail path.
func NewGCFULFTStreamWithTail(
	t ULFT,
	src GCFSource,
	tail Rational,
	opts GCFULFTStreamOptions,
) *GCFULFTStream {
	return NewGCFULFTStream(t, src, NewExactTailSource(tail), opts)
}

func (s *GCFULFTStream) Err() error { return s.err }

func (s *GCFULFTStream) initializeExactCF() bool {
	if s.started {
		return s.err == nil
	}
	s.started = true

	tail, ok := s.tailSrc.ExactTail()
	if !ok {
		s.err = fmt.Errorf("GCFULFTStream: tail evidence not implemented")
		s.done = true
		return false
	}

	y, _, err := ApplyComposedGCFULFTToTailExact(s.t, s.src, tail, s.opts.MaxIngestTerms)
	if err != nil {
		s.err = fmt.Errorf("GCFULFTStream: %w", err)
		s.done = true
		return false
	}

	s.exactCF = NewRationalCF(y)
	return true
}

func (s *GCFULFTStream) Next() (int64, bool) {
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

// gcf_ulft_stream.go v5
