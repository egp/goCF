// diag_stream.go v2
package cf

import "fmt"

type DiagBLFTStream struct {
	t   DiagBLFT
	src ContinuedFraction
	b   *Bounder

	srcDone bool
	done    bool
	err     error

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

func (s *DiagBLFTStream) Next() (int64, bool) {
	if s.done {
		return 0, false
	}
	if s.err != nil {
		s.done = true
		return 0, false
	}

	s.refinesThisDigit = 0

	for {
		if !s.b.HasValue() && !s.srcDone {
			a, ok := s.src.Next()
			if !ok {
				s.setErr(fmt.Errorf("DiagBLFTStream: empty source CF"))
				return 0, false
			}
			if err := s.b.Ingest(a); err != nil {
				s.setErr(err)
				return 0, false
			}
		}

		if s.srcDone {
			s.b.Finish()
		}

		xr, ok, err := s.b.Range()
		if err != nil {
			s.setErr(err)
			return 0, false
		}
		if !ok {
			s.setErr(fmt.Errorf("DiagBLFTStream: internal: no xRange"))
			return 0, false
		}

		img, err := s.t.ApplyRange(xr)
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

			// Exact integer termination short-circuit.
			if img.Lo.Cmp(img.Hi) == 0 && img.Lo.Cmp(intRat(d)) == 0 {
				s.done = true
				return d, true
			}

			tp, err := s.t.emitDigitDiag(d)
			if err != nil {
				s.setErr(err)
				return 0, false
			}
			s.t = tp
			return d, true
		}

		if s.srcDone {
			s.setErr(fmt.Errorf("DiagBLFTStream: cannot refine further (source finished) and digit not safe"))
			return 0, false
		}

		s.refinesThisDigit++
		s.refinesTotal++

		if s.maxRefinesPerDigit >= 0 && s.refinesThisDigit > s.maxRefinesPerDigit {
			s.setErr(fmt.Errorf("DiagBLFTStream: exceeded MaxRefinesPerDigit=%d", s.maxRefinesPerDigit))
			return 0, false
		}
		if s.maxTotalRefines >= 0 && s.refinesTotal > s.maxTotalRefines {
			s.setErr(fmt.Errorf("DiagBLFTStream: exceeded MaxTotalRefines=%d", s.maxTotalRefines))
			return 0, false
		}

		a, ok := s.src.Next()
		if ok {
			if err := s.b.Ingest(a); err != nil {
				s.setErr(err)
				return 0, false
			}
			continue
		}

		s.srcDone = true
	}
}

func (s *DiagBLFTStream) setErr(err error) {
	if s.err == nil {
		s.err = err
	}
	s.done = true
}

// diag_stream.go v2
