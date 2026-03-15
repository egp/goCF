// blft_stream.go v13
package cf

import "fmt"

type BLFTStream struct {
	t  BLFT
	xs ContinuedFraction
	ys ContinuedFraction
	xb *Bounder
	yb *Bounder

	xDone bool
	yDone bool

	done bool
	err  error

	emittedAny bool

	alt  bool
	opts BLFTStreamOptions

	tail ContinuedFraction

	finalizeTried bool

	maxRefinesPerDigit int
	maxTotalRefines    int
	refinesThisDigit   int
	refinesTotal       int

	// Optional cycle detection (diagnostic safety guard).
	detectCycles bool
	maxRepeats   int
	history      *RingBuf // recent FingerprintBLFT() values (human-readable)
}

type BLFTStreamOptions struct {
	// If >0, BLFTStream may drain X and Y (up to this many digits each) to exact rationals,
	// compute exact z=ApplyRat(x,y), and then stream z via NewRationalCF(z).
	//
	// This is intended for rational inputs and prevents range-based denom-guard failures
	// that can occur before the bounders shrink to the exact point.
	MaxFinalizeDigits  int
	MaxRefinesPerDigit int
	MaxTotalRefines    int

	// Optional cycle detection (debugging / safety).
	DetectCycles bool
	MaxRepeats   int // if <=0 and DetectCycles, defaults to 2
	HistorySize  int // if <=0, auto-sized
}

func NewBLFTStream(t BLFT, xs, ys ContinuedFraction, opts BLFTStreamOptions) *BLFTStream {
	max := opts.MaxRepeats
	if opts.DetectCycles && max <= 0 {
		max = 2
	}

	var hist *RingBuf
	if opts.DetectCycles {
		n := opts.HistorySize
		if n <= 0 {
			// Heuristic: enough context to see small loops.
			n = max * 8
			if n < 16 {
				n = 16
			}
			if n > 256 {
				n = 256
			}
		}
		hist = NewRingBuf(n)
	}

	// Defaulting rule for refine guards:
	// If BOTH are zero, treat that as "unset" => unlimited (-1).
	if opts.MaxRefinesPerDigit == 0 && opts.MaxTotalRefines == 0 {
		opts.MaxRefinesPerDigit = -1
		opts.MaxTotalRefines = -1
	}

	return &BLFTStream{
		t:                  t,
		xs:                 xs,
		ys:                 ys,
		xb:                 NewBounder(),
		yb:                 NewBounder(),
		opts:               opts,
		maxRefinesPerDigit: opts.MaxRefinesPerDigit,
		maxTotalRefines:    opts.MaxTotalRefines,
		detectCycles:       opts.DetectCycles,
		maxRepeats:         max,
		history:            hist,
	}
}

func (s *BLFTStream) Err() error { return s.err }

func (s *BLFTStream) binaryClass() binaryStreamClass {
	return binaryStreamClass{
		Operator: binaryOperatorUnknown,
		Input:    binaryInputCF,
		Progress: binaryProgressProgressiveCertified,
	}
}

// annotateErrBLFT appends a best-effort fingerprint context to err.
func annotateErrBLFT(err error, t BLFT, rx, ry Range) error {
	fp, ferr := FingerprintBLFT(t, rx, ry)
	if ferr != nil {
		return err
	}
	return fmt.Errorf("%w | %s", err, fp)
}

func (s *BLFTStream) ensureInitialValues() bool {
	if !s.xb.HasValue() && !s.xDone {
		a, ok := s.xs.Next()
		if !ok {
			s.setErr(fmt.Errorf("BLFTStream: empty X source CF"))
			return false
		}
		if err := s.xb.Ingest(a); err != nil {
			s.setErr(err)
			return false
		}
	}

	if !s.yb.HasValue() && !s.yDone {
		a, ok := s.ys.Next()
		if !ok {
			s.setErr(fmt.Errorf("BLFTStream: empty Y source CF"))
			return false
		}
		if err := s.yb.Ingest(a); err != nil {
			s.setErr(err)
			return false
		}
	}

	return true
}

func (s *BLFTStream) currentRanges() (Range, Range, bool) {
	if s.xDone {
		s.xb.Finish()
	}
	if s.yDone {
		s.yb.Finish()
	}

	xr, ok, err := s.xb.Range()
	if err != nil {
		s.setErr(err)
		return Range{}, Range{}, false
	}
	if !ok {
		s.setErr(fmt.Errorf("BLFTStream: internal: no xRange"))
		return Range{}, Range{}, false
	}

	yr, ok, err := s.yb.Range()
	if err != nil {
		s.setErr(err)
		return Range{}, Range{}, false
	}
	if !ok {
		s.setErr(fmt.Errorf("BLFTStream: internal: no yRange"))
		return Range{}, Range{}, false
	}

	return xr, yr, true
}

func (s *BLFTStream) maybeTerminateExactPoint(xr, yr Range) (done bool) {
	if !(s.xDone && s.yDone) || xr.Lo.Cmp(xr.Hi) != 0 || yr.Lo.Cmp(yr.Hi) != 0 {
		return false
	}

	den, err := s.t.denomAt(xr.Lo, yr.Lo)
	if err != nil {
		s.setErr(annotateErrBLFT(err, s.t, xr, yr))
		return true
	}
	if den.Cmp(intRat(0)) != 0 {
		return false
	}

	done, terr := exactPointTermination(
		"BLFTStream:",
		s.emittedAny,
		fmt.Sprintf("denominator is zero at exact point x=%v y=%v", xr.Lo, yr.Lo),
	)
	if done {
		s.done = true
		return true
	}

	s.setErr(annotateErrBLFT(terr, s.t, xr, yr))
	return true
}

func (s *BLFTStream) checkCycle(xr, yr Range) bool {
	if !(s.detectCycles && s.history != nil) {
		return true
	}

	fp, ferr := FingerprintBLFT(s.t, xr, yr)
	if ferr != nil {
		s.setErr(ferr)
		return false
	}
	s.history.Add(fp)
	if s.history.Count(fp) > s.maxRepeats {
		s.setErr(fmt.Errorf(
			"BLFTStream: cycle detected (repeats>%d): %s\nrecent:\n%s",
			s.maxRepeats, fp, s.history.Dump(),
		))
		return false
	}
	return true
}

func (s *BLFTStream) maybeFinalizeToTail(xr, yr Range) (int64, bool, bool) {
	if s.opts.MaxFinalizeDigits <= 0 {
		return 0, false, false
	}

	if switched, ferr := s.tryFinalizeToTail(); ferr != nil {
		s.setErr(annotateErrBLFT(ferr, s.t, xr, yr))
		return 0, false, true
	} else if switched {
		a, ok := s.tail.Next()
		if !ok {
			s.done = true
			return 0, false, true
		}
		return a, true, true
	}

	return 0, false, false
}

func (s *BLFTStream) chooseRefinement(xr, yr Range) (bool, bool, bool) {
	if s.xDone && s.yDone {
		s.setErr(annotateErrBLFT(
			fmt.Errorf("BLFTStream: cannot refine further (both sources finished) and digit not safe"),
			s.t, xr, yr,
		))
		return false, false, false
	}

	refineX := false
	refineY := false

	if s.xDone {
		refineY = true
	} else if s.yDone {
		refineX = true
	} else {
		mx, err := xr.RefineMetric()
		if err != nil {
			s.setErr(annotateErrBLFT(err, s.t, xr, yr))
			return false, false, false
		}
		my, err := yr.RefineMetric()
		if err != nil {
			s.setErr(annotateErrBLFT(err, s.t, xr, yr))
			return false, false, false
		}
		c := mx.Cmp(my)
		if c > 0 {
			refineX = true
		} else if c < 0 {
			refineY = true
		} else {
			if s.alt {
				refineX = true
			} else {
				refineY = true
			}
			s.alt = !s.alt
		}
	}

	return refineX, refineY, true
}

func (s *BLFTStream) consumeRefine(xr, yr Range) bool {
	if err := consumeRefineBudget(
		"BLFTStream:",
		&s.refinesThisDigit,
		&s.refinesTotal,
		s.maxRefinesPerDigit,
		s.maxTotalRefines,
	); err != nil {
		s.setErr(annotateErrBLFT(err, s.t, xr, yr))
		return false
	}
	return true
}

func (s *BLFTStream) refineChosenSource(refineX, refineY bool) bool {
	if refineX {
		a, ok := s.xs.Next()
		if ok {
			if err := s.xb.Ingest(a); err != nil {
				s.setErr(err)
				return false
			}
			return true
		}
		s.xDone = true
		return true
	}
	if refineY {
		a, ok := s.ys.Next()
		if ok {
			if err := s.yb.Ingest(a); err != nil {
				s.setErr(err)
				return false
			}
			return true
		}
		s.yDone = true
		return true
	}
	return true
}

func (s *BLFTStream) emitSafeDigit(d int64, img, xr, yr Range) (int64, bool) {
	if img.Lo.Cmp(img.Hi) == 0 && img.Lo.Cmp(intRat(d)) == 0 {
		s.done = true
		s.emittedAny = true
		return d, true
	}

	tp, err := s.emitDigitBLFT(d)
	if err != nil {
		s.setErr(annotateErrBLFT(err, s.t, xr, yr))
		return 0, false
	}
	s.t = tp
	s.emittedAny = true
	return d, true
}
func (s *BLFTStream) Next() (int64, bool) {
	if s.done {
		return 0, false
	}
	if s.err != nil {
		s.done = true
		return 0, false
	}

	s.refinesThisDigit = 0

	if s.tail != nil {
		a, ok := s.tail.Next()
		if !ok {
			s.done = true
			return 0, false
		}
		return a, true
	}

	if s.opts.MaxFinalizeDigits > 0 && !s.finalizeTried {
		s.finalizeTried = true
		if switched, err := s.tryFinalizeToTail(); err != nil {
			s.setErr(err)
			return 0, false
		} else if switched {
			a, ok := s.tail.Next()
			if !ok {
				s.done = true
				return 0, false
			}
			return a, true
		}
	}

	for {
		if !s.ensureInitialValues() {
			return 0, false
		}

		xr, yr, ok := s.currentRanges()
		if !ok {
			return 0, false
		}

		if s.maybeTerminateExactPoint(xr, yr) {
			return 0, false
		}

		if !s.checkCycle(xr, yr) {
			return 0, false
		}

		needRefine := false

		img, err := s.t.ApplyBLFTRange(xr, yr)
		if err != nil {
			if a, ok, handled := s.maybeFinalizeToTail(xr, yr); handled {
				return a, ok
			}

			if !(s.xDone && s.yDone) {
				needRefine = true
			} else {
				s.setErr(annotateErrBLFT(err, s.t, xr, yr))
				return 0, false
			}
		}

		if !needRefine {
			lo, hi, err := img.FloorBounds()
			if err != nil {
				s.setErr(annotateErrBLFT(err, s.t, xr, yr))
				return 0, false
			}

			if lo == hi {
				return s.emitSafeDigit(lo, img, xr, yr)
			}
		}

		refineX, refineY, ok := s.chooseRefinement(xr, yr)
		if !ok {
			return 0, false
		}

		if !s.consumeRefine(xr, yr) {
			return 0, false
		}

		if !s.refineChosenSource(refineX, refineY) {
			return 0, false
		}
	}
}

func (s *BLFTStream) setErr(err error) {
	if s.err == nil {
		s.err = err
	}
	s.done = true
}

// tryFinalizeToTail drains up to MaxFinalizeDigits from both inputs. If both terminate
// and collapse to exact points, it switches to an exact rational tail stream.
func (s *BLFTStream) tryFinalizeToTail() (bool, error) {
	limit := s.opts.MaxFinalizeDigits

	if !s.xDone {
		for i := 0; i < limit; i++ {
			a, ok := s.xs.Next()
			if !ok {
				s.xDone = true
				break
			}
			if err := s.xb.Ingest(a); err != nil {
				return false, err
			}
		}
	}
	if !s.yDone {
		for i := 0; i < limit; i++ {
			a, ok := s.ys.Next()
			if !ok {
				s.yDone = true
				break
			}
			if err := s.yb.Ingest(a); err != nil {
				return false, err
			}
		}
	}

	if s.xDone {
		s.xb.Finish()
	}
	if s.yDone {
		s.yb.Finish()
	}

	if !(s.xDone && s.yDone) {
		return false, nil
	}

	xr, ok, err := s.xb.Range()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, fmt.Errorf("BLFTStream: internal: no xRange after finalize")
	}
	yr, ok, err := s.yb.Range()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, fmt.Errorf("BLFTStream: internal: no yRange after finalize")
	}

	if xr.Lo.Cmp(xr.Hi) != 0 || yr.Lo.Cmp(yr.Hi) != 0 {
		return false, nil
	}

	z, err := s.t.ApplyRat(xr.Lo, yr.Lo)
	if err != nil {
		return false, err
	}

	s.tail = NewRationalCF(z)
	return true, nil
}

// emitDigitBLFT updates BLFT coefficients to represent z' = 1/(z - d).
//
// Given z = N/D, z - d = (N - dD)/D, and 1/(z-d) = D/(N - dD).
//
// So new numerator coefficients become old denominator coefficients.
// New denominator coefficients become (old numerator - d*old denominator).
func (s *BLFTStream) emitDigitBLFT(d int64) (BLFT, error) {
	t := s.t

	// New numerator = old denom: (E,F,G,H)
	A2, B2, C2, D2 := t.E, t.F, t.G, t.H

	dE, ok := mul64(d, t.E)
	if !ok {
		return BLFT{}, ErrOverflow
	}
	dF, ok := mul64(d, t.F)
	if !ok {
		return BLFT{}, ErrOverflow
	}
	dG, ok := mul64(d, t.G)
	if !ok {
		return BLFT{}, ErrOverflow
	}
	dH, ok := mul64(d, t.H)
	if !ok {
		return BLFT{}, ErrOverflow
	}

	E2, ok := sub64(t.A, dE)
	if !ok {
		return BLFT{}, ErrOverflow
	}
	F2, ok := sub64(t.B, dF)
	if !ok {
		return BLFT{}, ErrOverflow
	}
	G2, ok := sub64(t.C, dG)
	if !ok {
		return BLFT{}, ErrOverflow
	}
	H2, ok := sub64(t.D, dH)
	if !ok {
		return BLFT{}, ErrOverflow
	}

	return NewBLFT(A2, B2, C2, D2, E2, F2, G2, H2), nil
}

// blft_stream.go v13
