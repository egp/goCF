// gcf_stream.go v1
package cf

import "fmt"

type GCFStreamOptions struct{}

type GCFStream struct {
	tail ContinuedFraction
	done bool
	err  error
}

func NewGCFStream(src GCFSource, opts GCFStreamOptions) *GCFStream {
	r, err := EvaluateFiniteGCF(src)
	if err != nil {
		return &GCFStream{err: err, done: true}
	}
	return &GCFStream{
		tail: NewRationalCF(r),
	}
}

func (s *GCFStream) Err() error { return s.err }

func (s *GCFStream) Next() (int64, bool) {
	if s.done {
		return 0, false
	}
	if s.err != nil {
		s.done = true
		return 0, false
	}
	if s.tail == nil {
		s.err = fmt.Errorf("GCFStream: no tail")
		s.done = true
		return 0, false
	}

	d, ok := s.tail.Next()
	if !ok {
		s.done = true
		return 0, false
	}
	return d, true
}

// gcf_stream.go v1
