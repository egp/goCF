// gcf_bounder.go v2
package cf

import (
	"fmt"
	"math/big"
)

// GCFBounder incrementally ingests generalized continued-fraction terms (p,q)
// under the convention:
//
//	x = p + q/x'
//
// It maintains exact convergents for the finite prefixes seen so far.
//
// Current v2 semantics:
//   - Convergent() returns the exact value of the finite prefix under the
//     finite-tail convention that the last ingested term contributes just p_last.
//   - Range() returns the convergent as a degenerate point range.
//   - After Finish(), that point range is exact for the whole finite source.
//   - Before Finish(), this is not yet a conservative enclosure for an infinite GCF.
//
// Memory discipline:
//   - constant space in the number of ingested terms
//   - does not buffer the full term stream
type GCFBounder struct {
	// Latest convergent h_n/k_n and previous h_{n-1}/k_{n-1}.
	hm2, hm1 *big.Int
	km2, km1 *big.Int

	prevQ *big.Int // q_n from the latest ingested term
	count int
	done  bool
}

func NewGCFBounder() *GCFBounder {
	return &GCFBounder{}
}

// IngestPQ adds the next generalized term (p,q).
// If Finish() was already called, this returns an error.
// Requires q > 0.
func (b *GCFBounder) IngestPQ(p, q int64) error {
	if b.done {
		return fmt.Errorf("gcfbounder: cannot ingest after Finish()")
	}
	if q <= 0 {
		return fmt.Errorf("gcfbounder: require q>0, got q=%d", q)
	}

	pi := big.NewInt(p)
	qi := big.NewInt(q)

	if b.count == 0 {
		// First term: value is just p0 under the finite-tail convention.
		//
		// Set up recurrence state so that for the next term:
		// h1 = p1*h0 + q0*h[-1]
		// k1 = p1*k0 + q0*k[-1]
		// with h[-1]=1, k[-1]=0.
		b.hm2 = big.NewInt(1)
		b.hm1 = new(big.Int).Set(pi)
		b.km2 = big.NewInt(0)
		b.km1 = big.NewInt(1)
		b.prevQ = new(big.Int).Set(qi)
		b.count = 1
		return nil
	}

	// Standard generalized continuant update:
	// h_n = p_n*h_{n-1} + q_{n-1}*h_{n-2}
	// k_n = p_n*k_{n-1} + q_{n-1}*k_{n-2}
	h := new(big.Int).Mul(pi, b.hm1)
	h.Add(h, new(big.Int).Mul(b.prevQ, b.hm2))

	k := new(big.Int).Mul(pi, b.km1)
	k.Add(k, new(big.Int).Mul(b.prevQ, b.km2))

	b.hm2, b.hm1 = b.hm1, h
	b.km2, b.km1 = b.km1, k
	b.prevQ = new(big.Int).Set(qi)
	b.count++
	return nil
}

// Finish marks the GCF source as exhausted (finite termination).
func (b *GCFBounder) Finish() {
	b.done = true
}

// HasValue reports whether at least one term has been ingested.
func (b *GCFBounder) HasValue() bool {
	return b.count > 0
}

// Convergent returns the exact rational value of the finite generalized CF
// prefix seen so far, using the finite-tail convention that the last ingested
// term contributes just p_last.
func (b *GCFBounder) Convergent() (Rational, error) {
	if !b.HasValue() {
		return Rational{}, fmt.Errorf("gcfbounder: no terms ingested")
	}
	return newRationalBig(new(big.Int).Set(b.hm1), new(big.Int).Set(b.km1))
}

// Range returns the current range information for the ingested GCF prefix.
//
// Current v2 behavior:
//   - if no terms: returns (Range{}, false, nil)
//   - otherwise: returns the convergent as a degenerate point range
//
// This is exact after Finish(), but before Finish() it is only a placeholder
// until true conservative GCF enclosure logic is added later.
func (b *GCFBounder) Range() (Range, bool, error) {
	if !b.HasValue() {
		return Range{}, false, nil
	}

	conv, err := b.Convergent()
	if err != nil {
		return Range{}, false, err
	}

	return NewRange(conv, conv, true, true), true, nil
}

// gcf_bounder.go v2
