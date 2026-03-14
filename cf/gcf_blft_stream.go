// gcf_blft_stream.go v3
package cf

import "fmt"

// GCFBLFTStream is the GCF-native exact-tail stream for binary transforms.
//
// Current supported milestone:
//   - finite GCF prefix ingestion on x and y
//   - exact tail evidence supplied for x and y
//   - output regular CF digits of the exact transformed value
//
// Infinite uncertified GCF emission is intentionally out of scope for now.
type GCFBLFTStream struct {
	t        BLFT
	xSrc     GCFSource
	xTailSrc GCFTailSource
	ySrc     GCFSource
	yTailSrc GCFTailSource
	opts     GCFBLFTStreamOptions
	state    exactCFStreamState
}

// GCFBLFTStreamOptions controls bounded exact ingestion for the binary stream.
//
// Current meaning:
//
//   - MaxXIngestTerms < 0 : unlimited x-side ingestion
//
//   - MaxXIngestTerms = 0 : no x-side ingestion allowed
//
//   - MaxXIngestTerms > 0 : maximum permitted x-side source terms
//
//   - MaxYIngestTerms < 0 : unlimited y-side ingestion
//
//   - MaxYIngestTerms = 0 : no y-side ingestion allowed
//
//   - MaxYIngestTerms > 0 : maximum permitted y-side source terms
type GCFBLFTStreamOptions struct {
	MaxXIngestTerms int
	MaxYIngestTerms int
}

func NewGCFBLFTStream(
	t BLFT,
	xSrc GCFSource,
	xTailSrc GCFTailSource,
	ySrc GCFSource,
	yTailSrc GCFTailSource,
	opts GCFBLFTStreamOptions,
) *GCFBLFTStream {
	return &GCFBLFTStream{
		t:        t,
		xSrc:     xSrc,
		xTailSrc: xTailSrc,
		ySrc:     ySrc,
		yTailSrc: yTailSrc,
		opts:     opts,
	}
}

func NewGCFBLFTStreamWithTails(
	t BLFT,
	xSrc GCFSource,
	xTail Rational,
	ySrc GCFSource,
	yTail Rational,
	opts GCFBLFTStreamOptions,
) *GCFBLFTStream {
	return NewGCFBLFTStream(
		t,
		xSrc,
		NewExactTailSource(xTail),
		ySrc,
		NewExactTailSource(yTail),
		opts,
	)
}

func (s *GCFBLFTStream) Err() error { return s.state.Err() }

func (s *GCFBLFTStream) Next() (int64, bool) {
	return s.state.nextFromExactCF(func() (Rational, error) {
		xTail, ok := s.xTailSrc.ExactTail()
		if !ok {
			return Rational{}, fmt.Errorf("GCFBLFTStream: x tail evidence not implemented")
		}

		yTail, ok := s.yTailSrc.ExactTail()
		if !ok {
			return Rational{}, fmt.Errorf("GCFBLFTStream: y tail evidence not implemented")
		}

		z, _, _, err := ApplyComposedGCFXYBLFTToTailsExact(
			s.t,
			s.xSrc,
			xTail,
			s.opts.MaxXIngestTerms,
			s.ySrc,
			yTail,
			s.opts.MaxYIngestTerms,
		)
		if err != nil {
			return Rational{}, fmt.Errorf("GCFBLFTStream: %w", err)
		}
		return z, nil
	})
}

// gcf_blft_stream.go v3
