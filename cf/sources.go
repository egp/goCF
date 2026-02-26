// sources.go v1
package cf

// PeriodicCF is a simple continued fraction source with a finite prefix
// followed by an infinitely repeating period.
//
// Example: sqrt(2) = [1; (2)] is prefix=[1], period=[2].
// Example: phi     = [1; (1)] is prefix=[1], period=[1].
type PeriodicCF struct {
	prefix []int64
	period []int64
	i      int
}

func NewPeriodicCF(prefix []int64, period []int64) *PeriodicCF {
	// Caller responsibility: period must be non-empty for an infinite source.
	return &PeriodicCF{prefix: prefix, period: period}
}

func (p *PeriodicCF) Next() (int64, bool) {
	if p.i < len(p.prefix) {
		v := p.prefix[p.i]
		p.i++
		return v, true
	}
	if len(p.period) == 0 {
		// Degenerate finite CF source.
		return 0, false
	}
	j := (p.i - len(p.prefix)) % len(p.period)
	v := p.period[j]
	p.i++
	return v, true
}

// Sqrt2CF returns the infinite continued fraction for sqrt(2):
//
//	sqrt(2) = [1; (2)]
func Sqrt2CF() ContinuedFraction {
	return NewPeriodicCF([]int64{1}, []int64{2})
}

// PhiCF returns the infinite continued fraction for the golden ratio φ:
//
//	φ = [1; (1)]
func PhiCF() ContinuedFraction {
	return NewPeriodicCF([]int64{1}, []int64{1})
}

// sources.go v1
