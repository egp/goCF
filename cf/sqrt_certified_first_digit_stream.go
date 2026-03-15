// sqrt_certified_first_digit_stream.go v2
package cf

import "fmt"

// SqrtCertifiedFirstDigitCFStream is the first genuinely certified-progressive
// sqrt stream milestone.
//
// Current milestone:
//   - refine input CF prefix until the first sqrt digit is certified from
//     SqrtRangeConservative
//   - for non-exact inputs, certify and emit up to two digits using remainder
//     range propagation
//   - if the input is exact and has exact rational sqrt, emit the full exact CF
//
// Future work:
//   - general repeated certified-digit loop
//   - tighter linkage to transform/diagonal state
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

func (s *SqrtCertifiedFirstDigitCFStream) certifyDigitsFromRange(yr Range) error {
	a0lo, a0hi, err := yr.FloorBounds()
	if err != nil {
		return err
	}
	if a0lo != a0hi {
		return fmt.Errorf("uncertified first digit")
	}
	s.emitted = append(s.emitted, a0lo)

	r1, err := CertifiedRemainderRange(yr, a0lo)
	if err != nil {
		return nil // first digit still useful; second digit not yet available
	}

	a1lo, a1hi, err := r1.FloorBounds()
	if err != nil {
		return nil
	}
	if a1lo == a1hi {
		s.emitted = append(s.emitted, a1lo)
	}
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

			if xr.Lo.Cmp(xr.Hi) == 0 {
				if root, ok, err := RationalSqrtExact(xr.Lo); err == nil && ok {
					s.approx = &root
					s.exactCF = NewRationalCF(root)
					return true
				}
			}

			if err := s.certifyDigitsFromRange(yr); err != nil {
				s.err = fmt.Errorf("SqrtCertifiedFirstDigitCFStream: %w", err)
				s.done = true
				s.status = SqrtStreamStatusFailed
				return false
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

	if s.emitPos < len(s.emitted) {
		d := s.emitted[s.emitPos]
		s.emitPos++
		if s.emitPos >= len(s.emitted) {
			s.done = true
		}
		return d, true
	}

	s.done = true
	return 0, false
}

// sqrt_certified_first_digit_stream.go v2
