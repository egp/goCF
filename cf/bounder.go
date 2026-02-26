// bounder.go v1
package cf

import "fmt"

// Bounder incrementally ingests continued-fraction terms and maintains a
// shrinking enclosure Range for the represented real number.
//
// Convention / assumptions (aligned with our RationalCF):
//   - Simple continued fractions with "tail" >= 1 when not finished.
//   - For canonical expansions of reals: a0 is any integer, and ai >= 1 for i>=1.
//   - For rationals, you should call Finish() when the source is exhausted; then
//     the range collapses to an exact point.
//
// Math:
// Given prefix [a0; a1, ..., an] and unknown tail r in [1, +∞),
// value is (h_n*r + h_{n-1}) / (k_n*r + k_{n-1})
// Endpoints occur at r=1 (mediant) and r→∞ (convergent).
type Bounder struct {
	// Convergent recurrence state:
	// hm1/km1 = h_n/k_n (latest)
	// hm2/km2 = h_{n-1}/k_{n-1} (previous)
	hm2, hm1 int64
	km2, km1 int64

	count int  // number of ingested terms
	done  bool // true => exact rational (range collapses)
}

func NewBounder() *Bounder {
	return &Bounder{
		// h[-2]=0, h[-1]=1 ; k[-2]=1, k[-1]=0
		hm2: 0, hm1: 1,
		km2: 1, km1: 0,
	}
}

// Ingest adds the next continued-fraction term.
// If Finish() was already called, this returns an error.
func (b *Bounder) Ingest(a int64) error {
	if b.done {
		return fmt.Errorf("bounder: cannot ingest after Finish()")
	}
	h := a*b.hm1 + b.hm2
	k := a*b.km1 + b.km2

	b.hm2, b.hm1 = b.hm1, h
	b.km2, b.km1 = b.km1, k
	b.count++
	return nil
}

// Finish marks the CF source as exhausted (rational termination).
// Thereafter Range() returns an exact point range.
func (b *Bounder) Finish() {
	b.done = true
}

// HasValue reports whether at least one term has been ingested.
func (b *Bounder) HasValue() bool { return b.count > 0 }

// Convergent returns the current convergent h_n/k_n as an exact Rational.
// Requires HasValue().
func (b *Bounder) Convergent() (Rational, error) {
	if !b.HasValue() {
		return Rational{}, fmt.Errorf("bounder: no terms ingested")
	}
	return NewRational(b.hm1, b.km1)
}

// Range returns a closed interval [Lo, Hi] that contains the represented value.
//
// If b.HasValue()==false, returns (Range{}, false, nil).
// If b.done==true, returns an exact point range.
func (b *Bounder) Range() (Range, bool, error) {
	if !b.HasValue() {
		return Range{}, false, nil
	}

	conv, err := NewRational(b.hm1, b.km1)
	if err != nil {
		return Range{}, false, err
	}

	if b.done {
		return Range{Lo: conv, Hi: conv}, true, nil
	}

	// Mediant endpoint at r=1:
	// (h_n + h_{n-1}) / (k_n + k_{n-1})
	med, err := NewRational(b.hm1+b.hm2, b.km1+b.km2)
	if err != nil {
		return Range{}, false, err
	}

	// Order them safely by comparison (avoids parity mistakes).
	if conv.Cmp(med) <= 0 {
		return Range{Lo: conv, Hi: med}, true, nil
	}
	return Range{Lo: med, Hi: conv}, true, nil
}

// bounder.go v1
