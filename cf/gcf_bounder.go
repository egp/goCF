// gcf_bounder.go v1
package cf

import (
	"fmt"
)

// GCFBounder incrementally ingests generalized continued-fraction terms (p,q)
// under the convention:
//
//	x = p + q/x'
//
// It maintains exact convergents for the finite prefixes seen so far.
//
// Current v1 semantics:
//   - Convergent() returns the exact value of the finite prefix under the
//     finite-tail convention that the last ingested term contributes just p_last.
//   - Range() returns the convergent as a degenerate point range.
//   - After Finish(), that point range is exact for the whole finite source.
//   - Before Finish(), this is not yet a conservative enclosure for an infinite GCF.
type GCFBounder struct {
	terms [][2]int64
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
	b.terms = append(b.terms, [2]int64{p, q})
	return nil
}

// Finish marks the GCF source as exhausted (finite termination).
func (b *GCFBounder) Finish() {
	b.done = true
}

// HasValue reports whether at least one term has been ingested.
func (b *GCFBounder) HasValue() bool {
	return len(b.terms) > 0
}

// Convergent returns the exact rational value of the finite generalized CF
// prefix seen so far, using the finite-tail convention that the last ingested
// term contributes just p_last.
func (b *GCFBounder) Convergent() (Rational, error) {
	if !b.HasValue() {
		return Rational{}, fmt.Errorf("gcfbounder: no terms ingested")
	}

	// Finite evaluation by backward recurrence on buffered exact terms:
	// v = p_last
	// v = p_i + q_i / v
	last := b.terms[len(b.terms)-1]
	v := intRat(last[0])

	for i := len(b.terms) - 2; i >= 0; i-- {
		p := intRat(b.terms[i][0])
		q := intRat(b.terms[i][1])

		qOverV, err := q.Div(v)
		if err != nil {
			return Rational{}, err
		}
		v, err = p.Add(qOverV)
		if err != nil {
			return Rational{}, err
		}
	}

	return v, nil
}

// Range returns the current range information for the ingested GCF prefix.
//
// Current v1 behavior:
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

// gcf_bounder.go v1
