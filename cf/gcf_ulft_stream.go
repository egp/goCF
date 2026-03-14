// gcf_ulft_stream.go v3
package cf

import "fmt"

// GCFULFTStream is the new GCF-native unary/LFT path.
//
// Current supported milestone:
//   - finite GCF prefix ingestion
//   - explicit exact tail Rational supplied by caller
//   - output regular CF digits of the exact transformed value
//
// Infinite uncertified GCF emission is intentionally out of scope for now.
type GCFULFTStream struct {
	t       ULFT
	src     GCFSource
	tail    Rational
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

// NewGCFULFTStreamWithTail constructs a GCF-native ULFT stream for the current
// exact finite-prefix milestone.
//
// Semantics for current version:
//   - src must exhaust within MaxIngestTerms (if bounded)
//   - tail is the exact remaining tail variable value after src exhaustion
//   - the stream emits the regular CF digits of T(x), where
//     x is represented by the ingested GCF prefix followed by exact tail
func NewGCFULFTStreamWithTail(
	t ULFT,
	src GCFSource,
	tail Rational,
	opts GCFULFTStreamOptions,
) *GCFULFTStream {
	return &GCFULFTStream{
		t:    t,
		src:  src,
		tail: tail,
		opts: opts,
	}
}

func (s *GCFULFTStream) Err() error { return s.err }

func (s *GCFULFTStream) initializeExactCF() bool {
	if s.started {
		return s.err == nil
	}
	s.started = true

	cur := s.t
	ingested := 0

	for {
		if s.opts.MaxIngestTerms >= 0 && ingested >= s.opts.MaxIngestTerms {
			s.err = fmt.Errorf(
				"GCFULFTStream: exceeded MaxIngestTerms=%d before source exhaustion",
				s.opts.MaxIngestTerms,
			)
			s.done = true
			return false
		}

		p, q, ok := s.src.NextPQ()
		if !ok {
			break
		}

		var err error
		cur, err = cur.IngestGCF(p, q)
		if err != nil {
			s.err = fmt.Errorf("GCFULFTStream: ingest (%d,%d): %w", p, q, err)
			s.done = true
			return false
		}
		ingested++
	}

	y, err := cur.ApplyRat(s.tail)
	if err != nil {
		s.err = fmt.Errorf("GCFULFTStream: apply exact tail: %w", err)
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

// gcf_ulft_stream.go v3
