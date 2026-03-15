// blft_stream_finalize.go v1
package cf

import "fmt"

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

// blft_stream_finalize.go v1
