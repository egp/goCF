package cf

// ContinuedFraction streams terms a0, a1, a2, ...
// Next returns (ai, true) or (_, false) if exhausted (rationals terminate).
type ContinuedFraction interface {
	Next() (int64, bool)
}

// SliceCF is a trivial CF source (useful for tests).
type SliceCF struct {
	terms []int64
	i     int
}

func NewSliceCF(terms ...int64) *SliceCF {
	return &SliceCF{terms: terms}
}

func (s *SliceCF) Next() (int64, bool) {
	if s.i >= len(s.terms) {
		return 0, false
	}
	v := s.terms[s.i]
	s.i++
	return v, true
}
