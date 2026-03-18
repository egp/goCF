// cf/gcf_prefix_state.go v3
package cf

type gcfPrefixState struct {
	terms         [][2]int64
	currentApprox *GCFApprox
}

func newGcfPrefixState() *gcfPrefixState {
	return &gcfPrefixState{}
}

func (s *gcfPrefixState) ingestOne(p, q int64) error {
	s.terms = append(s.terms, [2]int64{p, q})

	a, err := GCFApproxFromPrefix(NewSliceGCF(s.terms...), len(s.terms))
	if err != nil {
		return err
	}
	s.currentApprox = &a
	return nil
}

func (s *gcfPrefixState) hasApprox() bool {
	return s.currentApprox != nil
}

func (s *gcfPrefixState) approx() GCFApprox {
	if s.currentApprox != nil {
		return *s.currentApprox
	}
	return GCFApprox{}
}

// cf/gcf_prefix_state.go v3
