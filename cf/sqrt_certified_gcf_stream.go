// sqrt_certified_gcf_stream.go v5
package cf

import "fmt"

// SqrtCertifiedGCFPrefixStream progressively certifies sqrt CF digits from a
// bounded GCF prefix source.
//
// Current milestone:
//   - refine bounded GCF input
//   - conservatively enclose sqrt over the current GCF input range
//   - emit as many certified CF digits as possible using a persistent range
//     emitter
//   - refine input and continue until bounds or source exhaustion stop progress
//
// Future work:
//   - carry transformed remainder state without replay
//   - tighter linkage to diagonal / transform-driven machinery
type SqrtCertifiedGCFPrefixStream struct {
	state unaryStreamState

	src            GCFSource
	maxPrefixTerms int

	b        *GCFBounder
	srcDone  bool
	ingested int

	status  SqrtStreamStatus
	approx  *Rational
	exactCF ContinuedFraction

	emitted []int64
	emitPos int

	emitter    *CertifiedCFRangeEmitter
	lastApprox *GCFApprox
}

func NewSqrtCertifiedGCFPrefixStream(src GCFSource, maxPrefixTerms int) (SqrtApproxStream, error) {
	if maxPrefixTerms <= 0 {
		return nil, fmt.Errorf("SqrtCertifiedGCFPrefixStream: maxPrefixTerms must be > 0, got %d", maxPrefixTerms)
	}

	b := NewGCFBounder()

	// Install any static unfinished-tail metadata the source advertises.
	if s, ok := src.(TailRangeBoundedGCFSource); ok {
		if err := b.SetTailRange(s.TailRange()); err != nil {
			return nil, fmt.Errorf("SqrtCertifiedGCFPrefixStream: %w", err)
		}
	} else if s, ok := src.(PositiveTailLowerBoundedGCFSource); ok {
		if err := b.SetTailLowerBound(s.TailLowerBound()); err != nil {
			return nil, fmt.Errorf("SqrtCertifiedGCFPrefixStream: %w", err)
		}
	}

	return &SqrtCertifiedGCFPrefixStream{
		src:            src,
		maxPrefixTerms: maxPrefixTerms,
		b:              b,
		status:         SqrtStreamStatusUnstarted,
	}, nil
}

func (s *SqrtCertifiedGCFPrefixStream) Err() error { return s.state.Err() }

func (s *SqrtCertifiedGCFPrefixStream) unaryClass() unaryStreamClass {
	return unaryStreamClass{
		Operator: unaryOperatorSqrt,
		Input:    unaryInputGCFPrefix,
		Progress: unaryProgressProgressiveCertified,
	}
}

func (s *SqrtCertifiedGCFPrefixStream) Snapshot() SqrtApproxStreamSnapshot {
	var approxCopy *Rational
	if s.approx != nil {
		v := *s.approx
		approxCopy = &v
	}

	var gcfApproxCopy *GCFApprox
	if s.lastApprox != nil {
		v := *s.lastApprox
		gcfApproxCopy = &v
	}

	return SqrtApproxStreamSnapshot{
		Status:         s.status,
		Started:        s.state.started,
		PrefixTerms:    s.ingested,
		Approx:         approxCopy,
		GCFInputApprox: gcfApproxCopy,
	}
}

func (s *SqrtCertifiedGCFPrefixStream) ingestOne() error {
	p, q, ok := s.src.NextPQ()
	if !ok {
		s.srcDone = true
		s.b.Finish()
		return nil
	}
	if err := s.b.IngestPQ(p, q); err != nil {
		return err
	}
	s.ingested++
	return nil
}

func (s *SqrtCertifiedGCFPrefixStream) currentApprox() (GCFApprox, error) {
	if !s.b.HasValue() {
		return GCFApprox{}, fmt.Errorf("no input value")
	}
	if s.srcDone {
		s.b.Finish()
	}
	return gcfApproxFromBounder(s.b, s.ingested, "no input value")
}

func (s *SqrtCertifiedGCFPrefixStream) rebuildEmitterFromCurrentRange() (bool, error) {
	a, err := s.currentApprox()
	if err != nil {
		return false, err
	}
	s.lastApprox = &a

	if a.Range == nil {
		return false, nil
	}

	yr, err := SqrtRangeConservative(*a.Range)
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

	if a.Range.Lo.Cmp(a.Range.Hi) == 0 {
		if root, ok, err := RationalSqrtExact(a.Range.Lo); err == nil && ok {
			s.approx = &root
			if s.exactCF == nil {
				s.exactCF = NewRationalCF(root)
			}
			return true, nil
		}
	}

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
	if err := e.Err(); err != nil {
		return false, err
	}

	d, ok := e.Next()
	if err := e.Err(); err != nil {
		return false, err
	}
	if !ok {
		s.emitter = nil
		return false, nil
	}

	s.emitted = append(s.emitted, d)
	s.emitter = e
	return true, nil
}

func (s *SqrtCertifiedGCFPrefixStream) ensureReady() bool {
	if !s.state.Begin() {
		return false
	}

	if !s.b.HasValue() {
		if err := s.ingestOne(); err != nil {
			s.state.Fail(fmt.Errorf("SqrtCertifiedGCFPrefixStream: %w", err))
			s.status = SqrtStreamStatusFailed
			return false
		}
		if !s.b.HasValue() {
			s.state.Fail(fmt.Errorf("SqrtCertifiedGCFPrefixStream: empty source"))
			s.status = SqrtStreamStatusFailed
			return false
		}
	}
	return true
}

func (s *SqrtCertifiedGCFPrefixStream) ensureAvailableDigit() bool {
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
				s.state.Fail(fmt.Errorf("SqrtCertifiedGCFPrefixStream: %w", err))
				s.status = SqrtStreamStatusFailed
				return false
			}
			s.emitter = nil
		}

		available, err := s.rebuildEmitterFromCurrentRange()
		if err != nil {
			s.state.Fail(fmt.Errorf("SqrtCertifiedGCFPrefixStream: %w", err))
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
			s.state.Exhaust()
			return false
		}
		if s.srcDone {
			s.state.Exhaust()
			return false
		}

		if err := s.ingestOne(); err != nil {
			s.state.Fail(fmt.Errorf("SqrtCertifiedGCFPrefixStream: %w", err))
			s.status = SqrtStreamStatusFailed
			return false
		}
	}
}

func (s *SqrtCertifiedGCFPrefixStream) Next() (int64, bool) {
	if s.state.done {
		return 0, false
	}
	if s.state.err != nil {
		s.state.Exhaust()
		return 0, false
	}
	if !s.ensureAvailableDigit() {
		return 0, false
	}

	if s.exactCF != nil {
		d, ok := s.exactCF.Next()
		if !ok {
			s.state.Exhaust()
			return 0, false
		}
		return d, true
	}

	if s.emitPos < len(s.emitted) {
		d := s.emitted[s.emitPos]
		s.emitPos++
		return d, true
	}

	s.state.Exhaust()
	return 0, false
}

// sqrt_certified_gcf_stream.go v5
