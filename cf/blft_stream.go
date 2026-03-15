// blft_stream.go v1
package cf

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

func (s *BLFTStream) binaryClass() binaryStreamClass {
	return binaryStreamClass{
		Operator: binaryOperatorUnknown,
		Input:    binaryInputCF,
		Progress: binaryProgressProgressiveCertified,
	}
}

// blft_stream.go v1
