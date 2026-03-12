// sources.go v15
package cf

// GCFSource streams generalized continued-fraction terms (p,q), using the convention:
//
//	x = p + q/x'
//
// with q > 0 for each emitted term.
//
// NextPQ returns (p, q, true) or (_, _, false) if exhausted.
type GCFSource interface {
	NextPQ() (int64, int64, bool)
}

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

// SliceGCF is a trivial finite generalized continued-fraction source, useful for tests.
//
// Terms are interpreted using x = p + q/x' with q > 0 expected by downstream ingestion logic.
type SliceGCF struct {
	terms [][2]int64
	i     int
}

func NewSliceGCF(terms ...[2]int64) *SliceGCF {
	cp := append([][2]int64(nil), terms...)
	return &SliceGCF{terms: cp}
}

func (s *SliceGCF) NextPQ() (int64, int64, bool) {
	if s.i >= len(s.terms) {
		return 0, 0, false
	}
	t := s.terms[s.i]
	s.i++
	return t[0], t[1], true
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

// CFGCFAdapter adapts an ordinary continued fraction source into a generalized
// continued-fraction source by mapping each regular term a to generalized term
// (p,q) = (a,1).
type CFGCFAdapter struct {
	src ContinuedFraction
}

// AdaptCFToGCF wraps a ContinuedFraction as a GCFSource using (a,1) terms.
func AdaptCFToGCF(src ContinuedFraction) GCFSource {
	return &CFGCFAdapter{src: src}
}

func (a *CFGCFAdapter) NextPQ() (int64, int64, bool) {
	v, ok := a.src.Next()
	if !ok {
		return 0, 0, false
	}
	return v, 1, true
}

// PeriodicGCF is a generalized continued-fraction source with a finite prefix
// followed by an infinitely repeating period of (p,q) terms.
type PeriodicGCF struct {
	prefix [][2]int64
	period [][2]int64
	i      int
}

func NewPeriodicGCF(prefix [][2]int64, period [][2]int64) *PeriodicGCF {
	// Caller responsibility: period must be non-empty for an infinite source.
	pfx := append([][2]int64(nil), prefix...)
	per := append([][2]int64(nil), period...)
	return &PeriodicGCF{prefix: pfx, period: per}
}

func (p *PeriodicGCF) NextPQ() (int64, int64, bool) {
	if p.i < len(p.prefix) {
		t := p.prefix[p.i]
		p.i++
		return t[0], t[1], true
	}
	if len(p.period) == 0 {
		return 0, 0, false
	}
	j := (p.i - len(p.prefix)) % len(p.period)
	t := p.period[j]
	p.i++
	return t[0], t[1], true
}

// FuncGCFSource is a generalized continued-fraction source backed by a generator function.
type FuncGCFSource struct {
	fn func(i int) (p, q int64, ok bool)
	i  int
}

func NewFuncGCFSource(fn func(i int) (p, q int64, ok bool)) *FuncGCFSource {
	return &FuncGCFSource{fn: fn}
}

func (s *FuncGCFSource) NextPQ() (int64, int64, bool) {
	p, q, ok := s.fn(s.i)
	if !ok {
		return 0, 0, false
	}
	s.i++
	return p, q, true
}

// ECFGSource is an algorithmic generalized continued-fraction source for e,
// emitted via the regular continued-fraction pattern mapped into GCF terms
// (p,q) = (a,1).
//
// Regular CF for e:
//
//	[2; 1,2,1, 1,4,1, 1,6,1, ...]
type ECFGSource struct {
	i int
}

func NewECFGSource() *ECFGSource {
	return &ECFGSource{}
}

func (s *ECFGSource) NextPQ() (int64, int64, bool) {
	// a0 = 2
	if s.i == 0 {
		s.i++
		return 2, 1, true
	}

	// For n>=1:
	// positions 2,5,8,... (1-based after a0) are 2,4,6,...
	// in zero-based global indexing, that's i % 3 == 2.
	var a int64
	if s.i%3 == 2 {
		a = 2 * int64((s.i+1)/3)
	} else {
		a = 1
	}

	s.i++
	return a, 1, true
}

// UnitPArithmeticQGCFSource is an infinite generalized continued-fraction source
// emitting:
//
//	(1,startQ), (1,startQ+step), (1,startQ+2*step), ...
//
// with step > 0 and startQ > 0.
//
// This is a simple nontrivial algorithmic GCF source used to exercise genuine
// generalized ingestion where q is not always 1.
type UnitPArithmeticQGCFSource struct {
	nextQ int64
	step  int64
}

func NewUnitPArithmeticQGCFSource(startQ, step int64) *UnitPArithmeticQGCFSource {
	// Caller responsibility for now: startQ > 0 and step > 0.
	return &UnitPArithmeticQGCFSource{
		nextQ: startQ,
		step:  step,
	}
}

func (s *UnitPArithmeticQGCFSource) NextPQ() (int64, int64, bool) {
	q := s.nextQ
	s.nextQ += s.step
	return 1, q, true
}

// Brouncker4OverPiGCFSource is Brouncker's generalized continued fraction for 4/pi:
//
//	4/pi = 1 + 1/(2 + 9/(2 + 25/(2 + 49/(2 + ...))))
//
// In (p,q) terms under x = p + q/x', this emits:
//
//	(1,1), (2,9), (2,25), (2,49), ...
type Brouncker4OverPiGCFSource struct {
	i int
}

func NewBrouncker4OverPiGCFSource() *Brouncker4OverPiGCFSource {
	return &Brouncker4OverPiGCFSource{}
}

func (s *Brouncker4OverPiGCFSource) NextPQ() (int64, int64, bool) {
	if s.i == 0 {
		s.i++
		return 1, 1, true
	}

	// After the leading 1 + 1/(...), Brouncker's 4/pi GCF is:
	//
	//   1 + 1 / (2 + 3^2 / (2 + 5^2 / (2 + 7^2 / ... )))
	//
	// So the emitted (p,q) terms are:
	//   (1,1), (2,9), (2,25), (2,49), ...
	//
	// with odd values 3,5,7,... after the initial leading term.
	odd := int64(2*s.i + 1) // 3,5,7,... for i=1,2,3,...
	q := odd * odd
	s.i++
	return 2, q, true
}

// LambertPiOver4GCFSource is Lambert's generalized continued fraction for pi/4:
//
//	pi/4 = 1 / (1 + 1/(3 + 4/(5 + 9/(7 + 16/(9 + ...)))))
//
// In (p,q) terms under x = p + q/x', this emits:
//
//	(0,1), (1,1), (3,4), (5,9), (7,16), (9,25), ...
type LambertPiOver4GCFSource struct {
	i int
}

func NewLambertPiOver4GCFSource() *LambertPiOver4GCFSource {
	return &LambertPiOver4GCFSource{}
}

func (s *LambertPiOver4GCFSource) NextPQ() (int64, int64, bool) {
	if s.i == 0 {
		s.i++
		return 0, 1, true
	}

	n := int64(s.i) // 1,2,3,...
	p := 2*n - 1
	q := n * n
	s.i++
	return p, q, true
}

// PositiveTailLowerBoundedGCFSource is an optional interface for GCF sources
// whose unfinished tails are known to satisfy tail >= L for some positive L.
//
// This metadata is used to derive conservative unfinished-prefix enclosures.
type PositiveTailLowerBoundedGCFSource interface {
	GCFSource
	TailLowerBound() Rational
}

func (a *CFGCFAdapter) TailLowerBound() Rational {
	// Regular CF tails are >= 1 when interpreted after a prefix in the usual
	// positive-term setting. This adapter is intended for that standard usage.
	return mustRat(1, 1)
}

func (p *PeriodicGCF) TailLowerBound() Rational {
	return mustRat(1, 1)
}

func (s *ECFGSource) TailLowerBound() Rational {
	return mustRat(1, 1)
}

func (s *Brouncker4OverPiGCFSource) TailLowerBound() Rational {
	return mustRat(1, 1)
}

func (s *LambertPiOver4GCFSource) TailLowerBound() Rational {
	return mustRat(1, 1)
}

func (s *UnitPArithmeticQGCFSource) TailLowerBound() Rational {
	return mustRat(1, 1)
}

type TailRangeBoundedGCFSource interface {
	GCFSource
	TailRange() Range
}

// sources.go v15
