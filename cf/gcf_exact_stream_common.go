// gcf_exact_stream_common.go v1
package cf

type exactCFStreamState struct {
	err     error
	done    bool
	started bool
	exactCF ContinuedFraction
}

func (s *exactCFStreamState) Err() error { return s.err }

func (s *exactCFStreamState) nextFromExactCF(
	init func() (Rational, error),
) (int64, bool) {
	if s.done {
		return 0, false
	}
	if s.err != nil {
		s.done = true
		return 0, false
	}

	if !s.started {
		s.started = true

		r, err := init()
		if err != nil {
			s.err = err
			s.done = true
			return 0, false
		}

		s.exactCF = NewRationalCF(r)
	}

	d, ok := s.exactCF.Next()
	if !ok {
		s.done = true
		return 0, false
	}
	return d, true
}

// gcf_exact_stream_common.go v1
