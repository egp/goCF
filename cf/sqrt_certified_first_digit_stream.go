// sqrt_certified_first_digit_stream.go v5
package cf

import "fmt"

// SqrtCertifiedFirstDigitCFStream is the first genuinely certified-progressive
// sqrt stream milestone.
//
// Current milestone:
//   - refine input CF prefix until sqrt(x) is enclosed conservatively
//   - maintain a persistent certified remainder-state emitter for the current
//     sqrt output range
//   - when current certification runs out, ingest more source terms, rebuild the
//     emitter from the refined sqrt range, and continue provided the certified
//     prefix remains stable
//   - if the input is exact and has exact rational sqrt, emit the full exact CF
//
// Future work:
//   - tighter linkage to transform/diagonal state
//   - avoid rebuilding emitter from scratch after each input refinement
type SqrtCertifiedFirstDigitCFStream struct {
	err     error
	done    bool
	started bool

	src            ContinuedFraction
	maxPrefixTerms int

	b        *Bounder
	srcDone  bool
	ingested int

	status  SqrtStreamStatus
	approx  *Rational
	exactCF ContinuedFraction

	emitted []int64
	emitPos int

	emitter *CertifiedCFRangeEmitter
}

func NewSqrtCertifiedFirstDigitCFStream(src ContinuedFraction, maxPrefixTerms int) (SqrtApproxStream, error) {
	if maxPrefixTerms <= 0 {
		return nil, fmt.Errorf("SqrtCertifiedFirstDigitCFStream: maxPrefixTerms must be > 0, got %d", maxPrefixTerms)
	}
	return &SqrtCertifiedFirstDigitCFStream{
		src:            src,
		maxPrefixTerms: maxPrefixTerms,
		b:              NewBounder(),
		status:         SqrtStreamStatusUnstarted,
	}, nil
}

func (s *SqrtCertifiedFirstDigitCFStream) Err() error { return s.err }

func (s *SqrtCertifiedFirstDigitCFStream) Snapshot() SqrtApproxStreamSnapshot {
	var approxCopy *Rational
	if s.approx != nil {
		v := *s.approx
		approxCopy = &v
	}
	return SqrtApproxStreamSnapshot{
		Status:      s.status,
		Started:     s.started,
		PrefixTerms: s.ingested,
		Approx:      approxCopy,
	}
}

func (s *SqrtCertifiedFirstDigitCFStream) ingestOne() error {
	a, ok := s.src.Next()
	if !ok {
		s.srcDone = true
		s.b.Finish()
		return nil
	}
	if err := s.b.Ingest(a); err != nil {
		return err
	}
	s.ingested++
	return nil
}

func (s *SqrtCertifiedFirstDigitCFStream) currentInputRange() (Range, error) {
	if !s.b.HasValue() {
		return Range{}, fmt.Errorf("no input value")
	}
	if s.srcDone {
		s.b.Finish()
	}

	xr, ok, err := s.b.Range()
	if err != nil {
		return Range{}, err
	}
	if !ok {
		return Range{}, fmt.Errorf("no input range")
	}
	return xr, nil
}

func (s *SqrtCertifiedFirstDigitCFStream) rebuildEmitterFromCurrentRange() (bool, error) {
	xr, err := s.currentInputRange()
	if err != nil {
		return false, err
	}

	yr, err := SqrtRangeConservative(xr)
	if err != nil {
		return false, err
	}

	lo, hi, err := yr.FloorBounds()
	if err != nil {
		return false, err
	}
	if lo != hi {
		return false, nil
	}

	s.status = SqrtStreamStatusCertifiedProgressive

	if xr.Lo.Cmp(xr.Hi) == 0 {
		if root, ok, err := RationalSqrtExact(xr.Lo); err == nil && ok {
			s.approx = &root
			if s.exactCF == nil {
				s.exactCF = NewRationalCF(root)
			}
			return true, nil
		}
	}

	// Preserve already-emitted prefix by replaying it on a fresh emitter.
	e, err := NewCertifiedCFRangeEmitter(yr)
	if err != nil {
		return false, err
	}
	for i, want := range s.emitted {
		got, ok := e.Next()
		if !ok {
			return false, fmt.Errorf("certified prefix shrank at position %d", i)
		}
		if got != want {
			return false, fmt.Errorf("certified prefix changed at %d: got %d want %d", i, got, want)
		}
	}
	s.emitter = e
	return true, nil
}

func (s *SqrtCertifiedFirstDigitCFStream) ensureReady() bool {
	if s.started {
		return s.err == nil
	}
	s.started = true

	if !s.b.HasValue() {
		if err := s.ingestOne(); err != nil {
			s.err = fmt.Errorf("SqrtCertifiedFirstDigitCFStream: %w", err)
			s.done = true
			s.status = SqrtStreamStatusFailed
			return false
		}
		if !s.b.HasValue() {
			s.err = fmt.Errorf("SqrtCertifiedFirstDigitCFStream: empty source")
			s.done = true
			s.status = SqrtStreamStatusFailed
			return false
		}
	}
	return true
}

func (s *SqrtCertifiedFirstDigitCFStream) ensureAvailableDigit() bool {
	if !s.ensureReady() {
		return false
	}
	if s.exactCF != nil {
		return true
	}
	if s.emitPos < len(s.emitted) {
		return true
	}

	for {
		if s.emitter != nil {
			if d, ok := s.emitter.Next(); ok {
				s.emitted = append(s.emitted, d)
				return true
			}
			if err := s.emitter.Err(); err != nil {
				s.err = fmt.Errorf("SqrtCertifiedFirstDigitCFStream: %w", err)
				s.done = true
				s.status = SqrtStreamStatusFailed
				return false
			}
			s.emitter = nil
		}

		available, err := s.rebuildEmitterFromCurrentRange()
		if err != nil {
			s.err = fmt.Errorf("SqrtCertifiedFirstDigitCFStream: %w", err)
			s.done = true
			s.status = SqrtStreamStatusFailed
			return false
		}
		if available {
			if s.exactCF != nil {
				return true
			}
			continue
		}

		if s.ingested >= s.maxPrefixTerms {
			s.done = true
			return false
		}
		if s.srcDone {
			s.done = true
			return false
		}

		if err := s.ingestOne(); err != nil {
			s.err = fmt.Errorf("SqrtCertifiedFirstDigitCFStream: %w", err)
			s.done = true
			s.status = SqrtStreamStatusFailed
			return false
		}
	}
}

func (s *SqrtCertifiedFirstDigitCFStream) Next() (int64, bool) {
	if s.done {
		return 0, false
	}
	if s.err != nil {
		s.done = true
		return 0, false
	}
	if !s.ensureAvailableDigit() {
		return 0, false
	}

	if s.exactCF != nil {
		d, ok := s.exactCF.Next()
		if !ok {
			s.done = true
			return 0, false
		}
		return d, true
	}

	if s.emitPos < len(s.emitted) {
		d := s.emitted[s.emitPos]
		s.emitPos++
		return d, true
	}

	s.done = true
	return 0, false
}

// sqrt_certified_first_digit_stream.go v5
