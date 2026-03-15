// unary_stream_state_common.go v1
package cf

type unaryStreamState struct {
	err     error
	done    bool
	started bool
}

func (s *unaryStreamState) Err() error { return s.err }

func (s *unaryStreamState) Begin() bool {
	if s.started {
		return s.err == nil
	}
	s.started = true
	return true
}

func (s *unaryStreamState) Fail(err error) {
	s.err = err
	s.done = true
}

func (s *unaryStreamState) Exhaust() {
	s.done = true
}

// unary_stream_state_common.go v1
