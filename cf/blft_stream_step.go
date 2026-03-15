// blft_stream_step.go v1
package cf

import "fmt"

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

// blft_stream_step.go v1
