// gcf_tail.go v1
package cf

// GCFTailSource is the minimal tail-evidence interface for the new GCF-native
// unary stream path.
//
// Current supported evidence:
//   - an exact tail rational
//
// Future slices may add certified range evidence or richer metadata, but that
// is intentionally deferred.
type GCFTailSource interface {
	ExactTail() (Rational, bool)
}

// ExactTailSource provides exact tail evidence.
type ExactTailSource struct {
	tail Rational
}

func NewExactTailSource(tail Rational) ExactTailSource {
	return ExactTailSource{tail: tail}
}

func (s ExactTailSource) ExactTail() (Rational, bool) {
	return s.tail, true
}

// NoTailSource is a test/helper implementation representing absence of usable
// tail evidence.
type NoTailSource struct{}

func (NoTailSource) ExactTail() (Rational, bool) {
	return Rational{}, false
}

// gcf_tail.go v1
