// ulft_stream.go v9
package cf

import (
	"fmt"
	"math/big"
)

// ULFTStream transforms a source continued fraction x into the continued fraction
// of T(x), streaming digits safely using Range bounds + SafeDigit.
//
// Error handling:
//   - Next() returns (0,false) once the stream is exhausted OR on first error.
//   - Call Err() to see whether termination was clean or error-induced.
type ULFTStream struct {
	t   ULFT
	src ContinuedFraction
	b   *Bounder

	srcDone bool
	done    bool
	err     error

	// Optional cycle detection (useful while bringing ULFT golden online).
	detectCycles bool
	seen         map[ulftStateKey]int
	maxRepeats   int
	history      *RingBuf // recent FingerprintULFT() values (human-readable)

	// Progress guards (anti-stall).
	//
	// Strict semantics:
	//   -1 => unlimited
	//    0 => no refines allowed
	//   >0 => max refines allowed
	maxRefinesPerDigit int
	maxTotalRefines    int
	refinesThisDigit   int
	refinesTotal       int

	// Set to true whenever a digit is successfully returned
	emittedAny bool
}

type ULFTStreamOptions struct {
	DetectCycles bool
	MaxRepeats   int // if <=0 and DetectCycles, defaults to 2

	// Strict semantics:
	//   -1 => unlimited
	//    0 => no refines allowed
	//   >0 => max refines allowed
	MaxRefinesPerDigit int
	MaxTotalRefines    int
}

func NewULFTStream(t ULFT, src ContinuedFraction, opts ULFTStreamOptions) *ULFTStream {
	// Defaulting rule for guards:
	// If BOTH guard fields are left at zero, treat that as "unset" => unlimited (-1).
	// This keeps ULFTStreamOptions{DetectCycles:true} from accidentally forbidding refines.
	if opts.MaxRefinesPerDigit == 0 && opts.MaxTotalRefines == 0 {
		opts.MaxRefinesPerDigit = -1
		opts.MaxTotalRefines = -1
	}

	max := opts.MaxRepeats
	if opts.DetectCycles && max <= 0 {
		max = 2
	}

	var hist *RingBuf
	if opts.DetectCycles {
		// Keep enough context to see the loop structure.
		// Heuristic: ~8 fingerprints per repeat, clamped.
		n := max * 8
		if n < 16 {
			n = 16
		}
		if n > 256 {
			n = 256
		}
		hist = NewRingBuf(n)
	}

	return &ULFTStream{
		t:                  t,
		src:                src,
		b:                  NewBounder(),
		detectCycles:       opts.DetectCycles,
		seen:               map[ulftStateKey]int{},
		maxRepeats:         max,
		history:            hist,
		maxRefinesPerDigit: opts.MaxRefinesPerDigit,
		maxTotalRefines:    opts.MaxTotalRefines,
	}
}

func (s *ULFTStream) Err() error { return s.err }

// annotateErrULFT appends a best-effort fingerprint context to err.
// It is intentionally non-fatal: if fingerprinting fails, we return err unchanged.
func annotateErrULFT(err error, t ULFT, r Range) error {
	fp, ferr := FingerprintULFT(t, r)
	if ferr != nil {
		return err
	}
	return fmt.Errorf("%w | %s", err, fp)
}

func (s *ULFTStream) ensureInitialValue() bool {
	if s.b.HasValue() || s.srcDone {
		return true
	}

	a, ok := s.src.Next()
	if !ok {
		s.setErr(fmt.Errorf("ULFTStream: empty source CF"))
		return false
	}
	if err := s.b.Ingest(a); err != nil {
		s.setErr(err)
		return false
	}
	return true
}

func (s *ULFTStream) currentRange() (Range, bool) {
	if s.srcDone {
		s.b.Finish()
	}

	xRange, ok, err := s.b.Range()
	if err != nil {
		s.setErr(err)
		return Range{}, false
	}
	if !ok {
		s.setErr(fmt.Errorf("ULFTStream: internal: no range despite HasValue"))
		return Range{}, false
	}
	return xRange, true
}

func (s *ULFTStream) maybeTerminateExactPoint(xRange Range) (done bool) {
	if !s.srcDone || xRange.Lo.Cmp(xRange.Hi) != 0 {
		return false
	}

	den, err := evalLinearOnRat(s.t.C, s.t.D, xRange.Lo)
	if err != nil {
		s.setErr(annotateErrULFT(err, s.t, xRange))
		return true
	}
	if den.Cmp(intRat(0)) != 0 {
		return false
	}

	done, terr := exactPointTermination(
		"ULFTStream:",
		s.emittedAny,
		fmt.Sprintf("denominator is zero at exact point x=%v", xRange.Lo),
	)
	if done {
		s.done = true
		return true
	}

	s.setErr(annotateErrULFT(terr, s.t, xRange))
	return true
}

func (s *ULFTStream) checkCycle(xRange Range) bool {
	if !s.detectCycles {
		return true
	}

	fp, ferr := FingerprintULFT(s.t, xRange)
	if ferr != nil {
		s.setErr(ferr)
		return false
	}

	if s.history != nil {
		s.history.Add(fp)
		if s.history.Count(fp) > s.maxRepeats {
			s.setErr(fmt.Errorf(
				"ULFTStream: cycle detected (repeats>%d): %s\nrecent:\n%s",
				s.maxRepeats, fp, s.history.Dump(),
			))
			return false
		}
		return true
	}

	key, kerr := ulftFingerprint(s.t, xRange)
	if kerr != nil {
		s.setErr(kerr)
		return false
	}
	s.seen[key]++
	if s.seen[key] > s.maxRepeats {
		s.setErr(fmt.Errorf("ULFTStream: cycle detected (repeats>%d): %v", s.maxRepeats, key))
		return false
	}
	return true
}

func (s *ULFTStream) refineForCurrentDigit(xRange Range) bool {
	if s.srcDone {
		s.setErr(annotateErrULFT(
			fmt.Errorf("ULFTStream: cannot refine further (source finished) and digit not safe"),
			s.t, xRange,
		))
		return false
	}

	if err := consumeRefineBudget(
		"ULFTStream:",
		&s.refinesThisDigit,
		&s.refinesTotal,
		s.maxRefinesPerDigit,
		s.maxTotalRefines,
	); err != nil {
		s.setErr(annotateErrULFT(err, s.t, xRange))
		return false
	}

	a, okSrc := s.src.Next()
	if okSrc {
		if err := s.b.Ingest(a); err != nil {
			s.setErr(err)
			return false
		}
		return true
	}

	s.srcDone = true
	return true
}

func (s *ULFTStream) refineAfterSafeDigitError(xRange Range) bool {
	if s.srcDone {
		return false
	}

	if err := consumeRefineBudget(
		"ULFTStream:",
		&s.refinesThisDigit,
		&s.refinesTotal,
		s.maxRefinesPerDigit,
		s.maxTotalRefines,
	); err != nil {
		s.setErr(annotateErrULFT(err, s.t, xRange))
		return false
	}

	a, okSrc := s.src.Next()
	if okSrc {
		if err := s.b.Ingest(a); err != nil {
			s.setErr(err)
			return false
		}
		return true
	}

	s.srcDone = true
	return true
}

func (s *ULFTStream) emitSafeDigit(d int64, xRange Range) (int64, bool) {
	img, err := xRange.ApplyULFT(s.t)
	if err != nil {
		// Exact-point pole after the final emitted digit is clean exhaustion,
		// not an error.
		if s.srcDone && xRange.Lo.Cmp(xRange.Hi) == 0 {
			den, derr := evalLinearOnRat(s.t.C, s.t.D, xRange.Lo)
			if derr == nil && den.Cmp(intRat(0)) == 0 {
				s.done = true
				return 0, false
			}
		}
		s.setErr(annotateErrULFT(err, s.t, xRange))
		return 0, false
	}

	if img.Lo.Cmp(img.Hi) == 0 && img.Lo.Cmp(intRat(d)) == 0 {
		s.done = true
		s.emittedAny = true
		return d, true
	}

	if s.srcDone && xRange.Lo.Cmp(xRange.Hi) == 0 {
		y, err := s.t.ApplyRat(xRange.Lo)
		if err == nil && y.Cmp(intRat(d)) == 0 {
			s.done = true
			s.emittedAny = true
			return d, true
		}
	}

	tp, err := EmitDigit(s.t, d)
	if err != nil {
		s.setErr(annotateErrULFT(err, s.t, xRange))
		return 0, false
	}
	s.t = tp
	s.emittedAny = true
	return d, true
}

func (s *ULFTStream) Next() (int64, bool) {
	if s.done {
		return 0, false
	}
	if s.err != nil {
		s.done = true
		return 0, false
	}

	s.refinesThisDigit = 0

	for {
		if !s.ensureInitialValue() {
			return 0, false
		}

		xRange, ok := s.currentRange()
		if !ok {
			return 0, false
		}

		if s.maybeTerminateExactPoint(xRange) {
			return 0, false
		}

		if !s.checkCycle(xRange) {
			return 0, false
		}

		d, okDigit, err := SafeDigit(s.t, xRange)
		if err != nil {
			if s.refineAfterSafeDigitError(xRange) {
				continue
			}
			s.setErr(annotateErrULFT(err, s.t, xRange))
			return 0, false
		}

		if okDigit {
			return s.emitSafeDigit(d, xRange)
		}

		if !s.refineForCurrentDigit(xRange) {
			return 0, false
		}
	}
}

func (s *ULFTStream) setErr(err error) {
	if s.err == nil {
		s.err = err
	}
	s.done = true
}

// ---- fingerprinting (cycle detection) ----

type ulftStateKey struct {
	A, B, C, D string
	Lo, Hi     string
}

func ulftFingerprint(t ULFT, r Range) (ulftStateKey, error) {
	if err := t.Validate(); err != nil {
		return ulftStateKey{}, err
	}
	tc := canonULFT(t)

	return ulftStateKey{
		A:  tc.A.String(),
		B:  tc.B.String(),
		C:  tc.C.String(),
		D:  tc.D.String(),
		Lo: r.Lo.String(),
		Hi: r.Hi.String(),
	}, nil
}

func canonULFT(t ULFT) ULFT {
	// Work on copies so we do not mutate caller-owned big.Ints.
	a := new(big.Int).Set(t.A)
	b := new(big.Int).Set(t.B)
	c := new(big.Int).Set(t.C)
	d := new(big.Int).Set(t.D)

	// g = gcd(|a|,|b|,|c|,|d|)
	aa := new(big.Int).Abs(a)
	bb := new(big.Int).Abs(b)
	cc := new(big.Int).Abs(c)
	dd := new(big.Int).Abs(d)

	g := new(big.Int)
	g.GCD(nil, nil, aa, bb)
	g.GCD(nil, nil, g, cc)
	g.GCD(nil, nil, g, dd)

	one := big.NewInt(1)
	if g.Cmp(one) > 0 {
		a.Quo(a, g)
		b.Quo(b, g)
		c.Quo(c, g)
		d.Quo(d, g)
	}

	// Normalize sign: make first non-zero coefficient positive.
	first := a
	if first.Sign() == 0 {
		first = b
		if first.Sign() == 0 {
			first = c
			if first.Sign() == 0 {
				first = d
			}
		}
	}

	if first.Sign() < 0 {
		a.Neg(a)
		b.Neg(b)
		c.Neg(c)
		d.Neg(d)
	}

	return ULFT{A: a, B: b, C: c, D: d}
}

// ulft_stream.go v9
