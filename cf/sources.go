// sources.go v3
package cf

// QuadraticRadicalSource is an optional interface for sources that know they
// represent the positive square root of an integer radicand.
//
// Current narrow meaning:
//   - Radicand() returns n,true  => source is sqrt(n)
//   - (_,false)                 => no algebraic metadata available
type QuadraticRadicalSource interface {
	ContinuedFraction
	Radicand() (int64, bool)
}

// PeriodicCF is a simple continued fraction source with a finite prefix
// followed by an infinitely repeating period.
//
// Example: sqrt(2) = [1; (2)] is prefix=[1], period=[2].
// Example: phi    = [1; (1)] is prefix=[1], period=[1].
type PeriodicCF struct {
	prefix   []int64
	period   []int64
	i        int
	radicand int64 // 0 => no algebraic metadata
}

func NewPeriodicCF(prefix []int64, period []int64) *PeriodicCF {
	// Caller responsibility: period must be non-empty for an infinite source.
	return &PeriodicCF{prefix: prefix, period: period}
}

func newPeriodicCFWithRadicand(prefix []int64, period []int64, radicand int64) *PeriodicCF {
	return &PeriodicCF{prefix: prefix, period: period, radicand: radicand}
}

func (p *PeriodicCF) Next() (int64, bool) {
	// prefix first
	if p.i < len(p.prefix) {
		v := p.prefix[p.i]
		p.i++
		return v, true
	}
	// then repeating period forever
	if len(p.period) == 0 {
		return 0, false
	}
	j := (p.i - len(p.prefix)) % len(p.period)
	v := p.period[j]
	p.i++
	return v, true
}

func (p *PeriodicCF) Radicand() (int64, bool) {
	if p.radicand > 0 {
		return p.radicand, true
	}
	return 0, false
}

// PhiCF returns the infinite continued fraction for the golden ratio φ:
//
//	φ = [1; (1)]
func PhiCF() ContinuedFraction {
	return NewPeriodicCF([]int64{1}, []int64{1})
}

// Sqrt2CF returns the infinite continued fraction for sqrt(2):
//
//	sqrt(2) = [1; (2)]
func Sqrt2CF() ContinuedFraction {
	return newPeriodicCFWithRadicand([]int64{1}, []int64{2}, 2)
}

// Sqrt3CF returns the infinite continued fraction for sqrt(3):
//
//	sqrt(3) = [1; (1,2)]
func Sqrt3CF() ContinuedFraction {
	return newPeriodicCFWithRadicand([]int64{1}, []int64{1, 2}, 3)
}

// Sqrt5CF returns the infinite continued fraction for sqrt(5):
//
//	sqrt(5) = [2; (4)]
func Sqrt5CF() ContinuedFraction {
	return newPeriodicCFWithRadicand([]int64{2}, []int64{4}, 5)
}

// Sqrt6CF returns the infinite continued fraction for sqrt(6):
//
//	sqrt(6) = [2; (2,4)]
func Sqrt6CF() ContinuedFraction {
	return newPeriodicCFWithRadicand([]int64{2}, []int64{2, 4}, 6)
}

// Sqrt7CF returns the infinite continued fraction for sqrt(7):
//
//	sqrt(7) = [2; (1,1,1,4)]
func Sqrt7CF() ContinuedFraction {
	return newPeriodicCFWithRadicand([]int64{2}, []int64{1, 1, 1, 4}, 7)
}

// sources.go v3
