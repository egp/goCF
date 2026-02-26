// blft_stream.go v2
package cf

import "fmt"

// BLFTStream streams the continued fraction digits of Z = T(X,Y), where
// T is a BLFT and X,Y are continued fractions.
//
// This v2 implementation targets rational sources (finite CFs) first.
// It uses two Bounders to maintain conservative ranges for X and Y and
// emits digits when floor(img.Lo) == floor(img.Hi).
//
// Refinement policy: refine the wider of xRange/yRange; tie-break alternates.
//
// Termination (rational-safe):
// If both sources are finished AND the current exact Z is integer d,
// emit d and stop WITHOUT updating the transform (to avoid division by zero).
type BLFTStream struct {
	t     BLFT
	xs    ContinuedFraction
	ys    ContinuedFraction
	xb    *Bounder
	yb    *Bounder
	xDone bool
	yDone bool

	done bool
	err  error

	alt bool // tie-break toggle
}

type BLFTStreamOptions struct{}

func NewBLFTStream(t BLFT, xs, ys ContinuedFraction, _ BLFTStreamOptions) *BLFTStream {
	return &BLFTStream{
		t:  t,
		xs: xs,
		ys: ys,
		xb: NewBounder(),
		yb: NewBounder(),
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

	for {
		// Ensure both bounders have at least one term ingested (unless already done).
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

			// Rational termination: if X and Y are exact and Z is exactly integer d, emit and stop.
			if s.xDone && s.yDone && xr.Lo.Cmp(xr.Hi) == 0 && yr.Lo.Cmp(yr.Hi) == 0 {
				z, err := s.t.ApplyRat(xr.Lo, yr.Lo)
				if err != nil {
					s.setErr(err)
					return 0, false
				}
				if z.Q == 1 && z.P == d {
					s.done = true
					return d, true
				}
			}

			tp, err := s.emitDigitBLFT(d)
			if err != nil {
				s.setErr(err)
				return 0, false
			}
			s.t = tp
			return d, true
		}

		// Need refinement; if both done, we cannot refine further.
		if s.xDone && s.yDone {
			s.setErr(fmt.Errorf("BLFTStream: cannot refine further (both sources finished) and digit not safe"))
			return 0, false
		}

		// Choose which to refine.
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

// blft_stream.go v2
