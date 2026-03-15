// gcf_diag_stream.go v3
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
	t    DiagBLFT
	src  GCFSource
	opts GCFULFTStreamOptions
	exactSingleTailStream
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
		t:    t,
		src:  src,
		opts: opts,
		exactSingleTailStream: exactSingleTailStream{
			tailSrc: tailSrc,
		},
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

func (s *GCFDiagStream) Err() error { return s.exactSingleTailStream.Err() }

func (s *GCFDiagStream) Next() (int64, bool) {
	return s.nextFromExactTail(
		"GCFDiagStream: tail evidence not implemented",
		func(tail Rational) (Rational, error) {
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
		},
	)
}

func (s *GCFDiagStream) binaryClass() binaryStreamClass {
	return binaryStreamClass{
		Operator: binaryOperatorUnknown,
		Input:    binaryInputGCF,
		Progress: binaryProgressExactCollapse,
	}
}

// gcf_diag_stream.go v3
