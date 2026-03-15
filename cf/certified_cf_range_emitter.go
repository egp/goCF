// certified_cf_range_emitter.go v1
package cf

import "fmt"

// CertifiedCFRangeEmitter incrementally emits continued-fraction digits from a
// conservative inside range by maintaining the current certified remainder range.
//
// This is a reusable operator-style component sitting between raw range
// certification and full transform-driven streaming.
type CertifiedCFRangeEmitter struct {
	cur     Range
	started bool
	done    bool
	err     error
}

func NewCertifiedCFRangeEmitter(r Range) (*CertifiedCFRangeEmitter, error) {
	if !r.IsInside() {
		return nil, fmt.Errorf("NewCertifiedCFRangeEmitter: requires inside range; got %v", r)
	}
	return &CertifiedCFRangeEmitter{cur: r}, nil
}

func (e *CertifiedCFRangeEmitter) Err() error { return e.err }

func (e *CertifiedCFRangeEmitter) CurrentRange() (Range, error) {
	if e.err != nil {
		return Range{}, e.err
	}
	return e.cur, nil
}

func (e *CertifiedCFRangeEmitter) Next() (int64, bool) {
	if e.done {
		return 0, false
	}
	if e.err != nil {
		e.done = true
		return 0, false
	}

	e.started = true

	lo, hi, err := e.cur.FloorBounds()
	if err != nil {
		e.err = err
		e.done = true
		return 0, false
	}
	if lo != hi {
		e.done = true
		return 0, false
	}

	d := lo
	next, err := CertifiedRemainderRange(e.cur, d)
	if err != nil {
		// This is a normal stopping condition for exact integers / finite exact
		// rationals where no further reciprocal remainder step is valid.
		e.done = true
		return d, true
	}

	e.cur = next
	return d, true
}

// certified_cf_range_emitter.go v1
