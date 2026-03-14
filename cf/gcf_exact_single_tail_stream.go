// gcf_exact_single_tail_stream.go v1
package cf

import "fmt"

type exactSingleTailStream struct {
	tailSrc GCFTailSource
	state   exactCFStreamState
}

func (s *exactSingleTailStream) Err() error { return s.state.Err() }

func (s *exactSingleTailStream) nextFromExactTail(
	missingTailMsg string,
	eval func(tail Rational) (Rational, error),
) (int64, bool) {
	return s.state.nextFromExactCF(func() (Rational, error) {
		tail, ok := s.tailSrc.ExactTail()
		if !ok {
			return Rational{}, fmt.Errorf("%s", missingTailMsg)
		}
		return eval(tail)
	})
}

// gcf_exact_single_tail_stream.go v1
