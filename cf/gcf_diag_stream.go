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
	state   exactCFStreamState
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

func (s *GCFDiagStream) Err() error { return s.state.Err() }

func (s *GCFDiagStream) Next() (int64, bool) {
	return s.state.nextFromExactCF(func() (Rational, error) {
		tail, ok := s.tailSrc.ExactTail()
		if !ok {
			return Rational{}, fmt.Errorf("GCFDiagStream: tail evidence not implemented")
		}

		y, _, err := ApplyComposedGCFDiagBLFTToTailExact(
			s.t,
			s.src,
			tail,
			s.opts.MaxIngestTerms,
		)
		if err != nil {
			return Rational{}, fmt.Errorf("GCFDiagStream: %w", err)
		}
		return y, nil
	})
}

// gcf_diag_stream.go v2
