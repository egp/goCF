// diag_stream.go v5
package cf

import (
	"fmt"
	"math/big"
)

type DiagBLFTStream struct {
	t   DiagBLFT
	src ContinuedFraction
	b   *Bounder

	srcDone bool
	done    bool
	err     error

	emittedAny bool

	maxRefinesPerDigit int
	maxTotalRefines    int
	refinesThisDigit   int
	refinesTotal       int
}

type DiagBLFTStreamOptions struct {
	// Strict semantics:
	//   -1 => unlimited
	//    0 => no refines allowed
	//   >0 => max refines allowed
	MaxRefinesPerDigit int
	MaxTotalRefines    int
}

func NewDiagBLFTStream(t DiagBLFT, src ContinuedFraction, opts DiagBLFTStreamOptions) *DiagBLFTStream {
	if opts.MaxRefinesPerDigit == 0 && opts.MaxTotalRefines == 0 {
		opts.MaxRefinesPerDigit = -1
		opts.MaxTotalRefines = -1
	}
	return &DiagBLFTStream{
		t:                  t,
		src:                src,
		b:                  NewBounder(),
		maxRefinesPerDigit: opts.MaxRefinesPerDigit,
		maxTotalRefines:    opts.MaxTotalRefines,
	}
}

func (s *DiagBLFTStream) Err() error { return s.err }

// exactIntFromQuadraticRadical returns z,true when:
//   - src advertises itself as sqrt(n)
//   - t is exactly (x^2 + k) / 1 for some integer k that fits in int64
//
// This is a deliberately narrow algebraic proof hook.
func (s *DiagBLFTStream) exactIntFromQuadraticRadical() (int64, bool) {
	qr, ok := s.src.(QuadraticRadicalSource)
	if !ok {
		return 0, false
	}
	n, ok := qr.Radicand()
	if !ok {
		return 0, false
	}

	// Match: (x^2 + k) / 1
	isSquarePlusConst :=
		s.t.A.Cmp(big.NewInt(1)) == 0 &&
			s.t.B.Sign() == 0 &&
			s.t.D.Sign() == 0 &&
			s.t.E.Sign() == 0 &&
			s.t.F.Cmp(big.NewInt(1)) == 0

	if !isSquarePlusConst {
		return 0, false
	}
	if !s.t.C.IsInt64() {
		return 0, false
	}

	k := s.t.C.Int64()
	z, ok := add64(n, k)
	if !ok {
		return 0, false
	}
	return z, true
}

func (s *DiagBLFTStream) ensureInitialValue() bool {
	if s.b.HasValue() || s.srcDone {
		return true
	}

	a, ok := s.src.Next()
	if !ok {
		s.setErr(fmt.Errorf("DiagBLFTStream: empty source CF"))
		return false
	}
	if err := s.b.Ingest(a); err != nil {
		s.setErr(err)
		return false
	}
	return true
}

func (s *DiagBLFTStream) currentRange() (Range, bool) {
	if s.srcDone {
		s.b.Finish()
	}

	xr, ok, err := s.b.Range()
	if err != nil {
		s.setErr(err)
		return Range{}, false
	}
	if !ok {
		s.setErr(fmt.Errorf("DiagBLFTStream: internal: no xRange"))
		return Range{}, false
	}
	return xr, true
}

func (s *DiagBLFTStream) maybeTerminateExactPoint(xr Range) bool {
	if !s.srcDone || xr.Lo.Cmp(xr.Hi) != 0 {
		return false
	}

	zero, err := diagDenomZeroAt(s.t, xr.Lo)
	if err != nil {
		s.setErr(err)
		return true
	}
	if !zero {
		return false
	}

	done, terr := exactPointTermination(
		"DiagBLFTStream:",
		s.emittedAny,
		fmt.Sprintf("denominator is zero at exact point x=%v", xr.Lo),
	)
	if done {
		s.done = true
		return true
	}

	s.setErr(terr)
	return true
}

func (s *DiagBLFTStream) refineForCurrentDigit() bool {
	if s.srcDone {
		s.setErr(fmt.Errorf("DiagBLFTStream: cannot refine further (source finished) and digit not safe"))
		return false
	}

	if err := consumeRefineBudget(
		"DiagBLFTStream:",
		&s.refinesThisDigit,
		&s.refinesTotal,
		s.maxRefinesPerDigit,
		s.maxTotalRefines,
	); err != nil {
		s.setErr(err)
		return false
	}

	a, ok := s.src.Next()
	if ok {
		if err := s.b.Ingest(a); err != nil {
			s.setErr(err)
			return false
		}
		return true
	}

	s.srcDone = true
	return true
}

func (s *DiagBLFTStream) emitSafeDigit(d int64, img Range) (int64, bool) {
	if img.Lo.Cmp(img.Hi) == 0 && img.Lo.Cmp(intRat(d)) == 0 {
		s.done = true
		s.emittedAny = true
		return d, true
	}

	tp, err := s.t.emitDigitDiag(d)
	if err != nil {
		s.setErr(err)
		return 0, false
	}
	s.t = tp
	s.emittedAny = true
	return d, true
}

func (s *DiagBLFTStream) Next() (int64, bool) {
	if s.done {
		return 0, false
	}
	if s.err != nil {
		s.done = true
		return 0, false
	}

	if n, ok := s.exactIntFromQuadraticRadical(); ok {
		s.done = true
		s.emittedAny = true
		return n, true
	}

	s.refinesThisDigit = 0

	for {
		if !s.ensureInitialValue() {
			return 0, false
		}

		xr, ok := s.currentRange()
		if !ok {
			return 0, false
		}

		if s.maybeTerminateExactPoint(xr) {
			return 0, false
		}

		needRefine := false

		img, err := s.t.ApplyRange(xr)
		if err != nil {
			if !s.srcDone {
				needRefine = true
			} else {
				s.setErr(err)
				return 0, false
			}
		}

		if !needRefine {
			lo, hi, err := img.FloorBounds()
			if err != nil {
				s.setErr(err)
				return 0, false
			}

			if lo == hi {
				return s.emitSafeDigit(lo, img)
			}
		}

		if !s.refineForCurrentDigit() {
			return 0, false
		}
	}
}

func (s *DiagBLFTStream) setErr(err error) {
	if s.err == nil {
		s.err = err
	}
	s.done = true
}

func diagDenomZeroAt(t DiagBLFT, x Rational) (bool, error) {
	var x2 big.Rat
	x2.Mul(&x.r, &x.r)

	var den, term big.Rat
	den.SetInt64(0)

	term.Mul(ratFromBigInt(t.D), &x2)
	den.Add(&den, &term)

	term.Mul(ratFromBigInt(t.E), &x.r)
	den.Add(&den, &term)

	den.Add(&den, ratFromBigInt(t.F))

	return den.Sign() == 0, nil
}

// diag_stream.go v5
