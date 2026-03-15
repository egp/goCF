// sqrt_certified_first_digit_stream.go v1
package cf

import "fmt"

// SqrtCertifiedFirstDigitCFStream is the first genuinely certified-progressive
// sqrt stream milestone.
//
// Current milestone:
//   - refine input CF prefix until the first sqrt digit is certified from
//     SqrtRangeConservative
//   - emit only the first certified digit for non-exact inputs
//   - if the input is exact and has exact rational sqrt, emit the full exact CF
//
// Future work:
//   - certified remainder/update after first digit
//   - repeated certified-digit loop
type SqrtCertifiedFirstDigitCFStream struct {
	err     error
	done    bool
	started bool

	src            ContinuedFraction
	maxPrefixTerms int

	b        *Bounder
	srcDone  bool
	ingested int

	status     SqrtStreamStatus
	approx     *Rational
	exactCF    ContinuedFraction
	firstDigit *int64
	emitted    bool
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

func (s *SqrtCertifiedFirstDigitCFStream) init() bool {
	if s.started {
		return s.err == nil
	}
	s.started = true

	for {
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

		if s.srcDone {
			s.b.Finish()
		}

		xr, ok, err := s.b.Range()
		if err != nil {
			s.err = fmt.Errorf("SqrtCertifiedFirstDigitCFStream: %w", err)
			s.done = true
			s.status = SqrtStreamStatusFailed
			return false
		}
		if !ok {
			s.err = fmt.Errorf("SqrtCertifiedFirstDigitCFStream: no input range")
			s.done = true
			s.status = SqrtStreamStatusFailed
			return false
		}

		yr, err := SqrtRangeConservative(xr)
		if err != nil {
			s.err = fmt.Errorf("SqrtCertifiedFirstDigitCFStream: %w", err)
			s.done = true
			s.status = SqrtStreamStatusFailed
			return false
		}

		lo, hi, err := yr.FloorBounds()
		if err != nil {
			s.err = fmt.Errorf("SqrtCertifiedFirstDigitCFStream: %w", err)
			s.done = true
			s.status = SqrtStreamStatusFailed
			return false
		}

		if lo == hi {
			s.status = SqrtStreamStatusCertifiedProgressive
			d := lo
			s.firstDigit = &d

			if xr.Lo.Cmp(xr.Hi) == 0 {
				if root, ok, err := RationalSqrtExact(xr.Lo); err == nil && ok {
					s.approx = &root
					s.exactCF = NewRationalCF(root)
				}
			}
			return true
		}

		if s.ingested >= s.maxPrefixTerms {
			s.err = fmt.Errorf("SqrtCertifiedFirstDigitCFStream: could not certify first digit within %d terms", s.maxPrefixTerms)
			s.done = true
			s.status = SqrtStreamStatusFailed
			return false
		}

		if s.srcDone {
			s.err = fmt.Errorf("SqrtCertifiedFirstDigitCFStream: exhausted source before certifying first digit")
			s.done = true
			s.status = SqrtStreamStatusFailed
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
	if !s.init() {
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

	if s.firstDigit != nil && !s.emitted {
		s.emitted = true
		s.done = true
		return *s.firstDigit, true
	}

	s.done = true
	return 0, false
}

// sqrt_certified_first_digit_stream.go v1
