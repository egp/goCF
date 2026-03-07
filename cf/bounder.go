// bounder.go v3
package cf

import (
	"fmt"
	"math/big"
)

// Bounder incrementally ingests continued-fraction terms and maintains a
// shrinking enclosure Range for the represented real number.
type Bounder struct {
	// Convergent recurrence state:
	// hm1/km1 = h_n/k_n (latest)
	// hm2/km2 = h_{n-1}/k_{n-1} (previous)
	hm2, hm1 *big.Int
	km2, km1 *big.Int

	count int  // number of ingested terms
	done  bool // true => exact rational (range collapses)
}

func NewBounder() *Bounder {
	return &Bounder{
		// h[-2]=0, h[-1]=1 ; k[-2]=1, k[-1]=0
		hm2: big.NewInt(0),
		hm1: big.NewInt(1),
		km2: big.NewInt(1),
		km1: big.NewInt(0),
	}
}

// Ingest adds the next continued-fraction term.
// If Finish() was already called, this returns an error.
func (b *Bounder) Ingest(a int64) error {
	if b.done {
		return fmt.Errorf("bounder: cannot ingest after Finish()")
	}

	ai := big.NewInt(a)

	// h = a*hm1 + hm2
	h := new(big.Int).Mul(ai, b.hm1)
	h.Add(h, b.hm2)

	// k = a*km1 + km2
	k := new(big.Int).Mul(ai, b.km1)
	k.Add(k, b.km2)

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
	return newRationalBig(new(big.Int).Set(b.hm1), new(big.Int).Set(b.km1))
}

// Range returns a closed interval [Lo, Hi] that contains the represented value.
//
// If b.HasValue()==false, returns (Range{}, false, nil).
// If b.done==true, returns an exact point range.
func (b *Bounder) Range() (Range, bool, error) {
	if !b.HasValue() {
		return Range{}, false, nil
	}

	conv, err := newRationalBig(new(big.Int).Set(b.hm1), new(big.Int).Set(b.km1))
	if err != nil {
		return Range{}, false, err
	}

	if b.done {
		return Range{Lo: conv, Hi: conv, IncLo: true, IncHi: true}, true, nil
	}

	// Mediant endpoint at r=1:
	// (h_n + h_{n-1}) / (k_n + k_{n-1})
	hn := new(big.Int).Add(b.hm1, b.hm2)
	kn := new(big.Int).Add(b.km1, b.km2)

	med, err := newRationalBig(hn, kn)
	if err != nil {
		return Range{}, false, err
	}

	// Order them by comparison (no parity assumptions).
	if conv.Cmp(med) <= 0 {
		return Range{Lo: conv, Hi: med, IncLo: true, IncHi: true}, true, nil
	}
	return Range{Lo: med, Hi: conv, IncLo: true, IncHi: true}, true, nil
}

// bounder.go v3
