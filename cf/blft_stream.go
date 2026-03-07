// blft_stream.go v11
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

// annotateErrBLFT appends a best-effort fingerprint context to err.
func annotateErrBLFT(err error, t BLFT, rx, ry Range) error {
	fp, ferr := FingerprintBLFT(t, rx, ry)
	if ferr != nil {
		return err
	}
	return fmt.Errorf("%w | %s", err, fp)
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

	// If we have an exact tail CF, delegate.
	if s.tail != nil {
		a, ok := s.tail.Next()
		if !ok {
			s.done = true
			return 0, false
		}
		return a, true
	}

	// Early rational finalization attempt (bounded, opt-in).
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
		// Ensure at least one ingested digit for each side (so Range() is defined).
		if !s.xb.HasValue() && !s.xDone {
			a, ok := s.xs.Next()
			if !ok {
				s.setErr(fmt.Errorf("BLFTStream: empty X source CF"))
				return 0, false
			}
			if err := s.xb.Ingest(a); err != nil {
				s.setErr(err)
				return 0, false
			}
		}
		if !s.yb.HasValue() && !s.yDone {
			a, ok := s.ys.Next()
			if !ok {
				s.setErr(fmt.Errorf("BLFTStream: empty Y source CF"))
				return 0, false
			}
			if err := s.yb.Ingest(a); err != nil {
				s.setErr(err)
				return 0, false
			}
		}

		if s.xDone {
			s.xb.Finish()
		}
		if s.yDone {
			s.yb.Finish()
		}

		xr, ok, err := s.xb.Range()
		if err != nil {
			s.setErr(err)
			return 0, false
		}
		if !ok {
			s.setErr(fmt.Errorf("BLFTStream: internal: no xRange"))
			return 0, false
		}

		yr, ok, err := s.yb.Range()
		if err != nil {
			s.setErr(err)
			return 0, false
		}
		if !ok {
			s.setErr(fmt.Errorf("BLFTStream: internal: no yRange"))
			return 0, false
		}

		// Cycle detection guard (best-effort diagnostic).
		if s.detectCycles && s.history != nil {
			fp, ferr := FingerprintBLFT(s.t, xr, yr)
			if ferr != nil {
				s.setErr(ferr)
				return 0, false
			}
			s.history.Add(fp)
			if s.history.Count(fp) > s.maxRepeats {
				s.setErr(fmt.Errorf(
					"BLFTStream: cycle detected (repeats>%d): %s\nrecent:\n%s",
					s.maxRepeats, fp, s.history.Dump(),
				))
				return 0, false
			}
		}

		img, err := s.t.ApplyBLFTRange(xr, yr)
		if err != nil {
			// If denom guard trips and finalization is enabled, try finalizing now.
			if s.opts.MaxFinalizeDigits > 0 {
				if switched, ferr := s.tryFinalizeToTail(); ferr != nil {
					s.setErr(annotateErrBLFT(ferr, s.t, xr, yr))
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
			s.setErr(annotateErrBLFT(err, s.t, xr, yr))
			return 0, false
		}

		lo, hi, err := img.FloorBounds()
		if err != nil {
			s.setErr(annotateErrBLFT(err, s.t, xr, yr))
			return 0, false
		}

		if lo == hi {
			d := lo

			// v10: integer termination short-circuit
			//
			// If the image interval has collapsed to an exact integer d, then the CF is exactly [d]
			// at this point and we MUST NOT apply emitDigitBLFT (which would compute 1/(d-d)).
			if img.Lo.Cmp(img.Hi) == 0 && img.Lo.Cmp(intRat(d)) == 0 {
				s.done = true
				return d, true
			}

			tp, err := s.emitDigitBLFT(d)
			if err != nil {
				s.setErr(annotateErrBLFT(err, s.t, xr, yr))
				return 0, false
			}
			s.t = tp
			return d, true
		}

		// No safe digit: refine.
		if s.xDone && s.yDone {
			s.setErr(annotateErrBLFT(
				fmt.Errorf("BLFTStream: cannot refine further (both sources finished) and digit not safe"),
				s.t, xr, yr,
			))
			return 0, false
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
				return 0, false
			}
			my, err := yr.RefineMetric()
			if err != nil {
				s.setErr(annotateErrBLFT(err, s.t, xr, yr))
				return 0, false
			}
			c := mx.Cmp(my)
			if c > 0 {
				refineX = true
			} else if c < 0 {
				refineY = true
			} else {
				// Tie-breaker: alternate.
				if s.alt {
					refineX = true
				} else {
					refineY = true
				}
				s.alt = !s.alt
			}
		}

		// Progress guards: each attempted refinement consumes one refine budget.
		s.refinesThisDigit++
		s.refinesTotal++
		if s.maxRefinesPerDigit >= 0 && s.refinesThisDigit > s.maxRefinesPerDigit {
			s.setErr(annotateErrBLFT(
				fmt.Errorf("BLFTStream: exceeded MaxRefinesPerDigit=%d", s.maxRefinesPerDigit),
				s.t, xr, yr,
			))
			return 0, false
		}
		if s.maxTotalRefines >= 0 && s.refinesTotal > s.maxTotalRefines {
			s.setErr(annotateErrBLFT(
				fmt.Errorf("BLFTStream: exceeded MaxTotalRefines=%d", s.maxTotalRefines),
				s.t, xr, yr,
			))
			return 0, false
		}

		if refineX {
			a, ok := s.xs.Next()
			if ok {
				if err := s.xb.Ingest(a); err != nil {
					s.setErr(err)
					return 0, false
				}
				continue
			}
			s.xDone = true
			continue
		}
		if refineY {
			a, ok := s.ys.Next()
			if ok {
				if err := s.yb.Ingest(a); err != nil {
					s.setErr(err)
					return 0, false
				}
				continue
			}
			s.yDone = true
			continue
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

// EOF blft_stream.go v11
