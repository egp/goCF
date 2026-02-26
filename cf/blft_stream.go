// blft_stream.go v5
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
}

type BLFTStreamOptions struct {
	// If >0, BLFTStream may drain X and Y (up to this many digits each) to exact rationals,
	// compute exact z=ApplyRat(x,y), and then stream z via NewRationalCF(z).
	//
	// This is intended for rational inputs and prevents range-based denom-guard failures
	// that can occur before the bounders shrink to the exact point.
	MaxFinalizeDigits int
}

func NewBLFTStream(t BLFT, xs, ys ContinuedFraction, opts BLFTStreamOptions) *BLFTStream {
	return &BLFTStream{
		t:    t,
		xs:   xs,
		ys:   ys,
		xb:   NewBounder(),
		yb:   NewBounder(),
		opts: opts,
	}
}

func (s *BLFTStream) Err() error { return s.err }

func (s *BLFTStream) Next() (int64, bool) {
	if s.done {
		return 0, false
	}
	if s.err != nil {
		s.done = true
		return 0, false
	}

	// If we have an exact tail CF, delegate.
	if s.tail != nil {
		a, ok := s.tail.Next()
		if !ok {
			s.done = true
			return 0, false
		}
		return a, true
	}

	// v5: early rational finalization attempt (bounded, opt-in).
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

		img, err := s.t.ApplyBLFTRange(xr, yr)
		if err != nil {
			// If denom guard trips and finalization is enabled, try finalizing now.
			if s.opts.MaxFinalizeDigits > 0 {
				if switched, ferr := s.tryFinalizeToTail(); ferr != nil {
					s.setErr(ferr)
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
			s.setErr(err)
			return 0, false
		}

		lo, hi, err := img.FloorBounds()
		if err != nil {
			s.setErr(err)
			return 0, false
		}

		if lo == hi {
			d := lo
			tp, err := s.emitDigitBLFT(d)
			if err != nil {
				s.setErr(err)
				return 0, false
			}
			s.t = tp
			return d, true
		}

		if s.xDone && s.yDone {
			s.setErr(fmt.Errorf("BLFTStream: cannot refine further (both sources finished) and digit not safe"))
			return 0, false
		}

		refineX := false
		refineY := false

		if s.xDone {
			refineY = true
		} else if s.yDone {
			refineX = true
		} else {
			wx, err := xr.Width()
			if err != nil {
				s.setErr(err)
				return 0, false
			}
			wy, err := yr.Width()
			if err != nil {
				s.setErr(err)
				return 0, false
			}
			c := wx.Cmp(wy)
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

// blft_stream.go v5
