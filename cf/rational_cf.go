package cf

import "fmt"

// RationalCF streams the continued-fraction terms of a Rational.
type RationalCF struct {
	p, q int64
	done bool
}

func NewRationalCF(r Rational) *RationalCF {
	return &RationalCF{p: r.P, q: r.Q}
}

func (c *RationalCF) Next() (int64, bool) {
	if c.done {
		return 0, false
	}
	if c.q == 0 {
		// Should be impossible if constructed from Rational, but keep it safe.
		c.done = true
		return 0, false
	}

	a, rem := floorDivMod(c.p, c.q)
	if rem == 0 {
		c.done = true
		return a, true
	}

	// Next step: x = q/rem with sign handled by floorDivMod remainder being >= 0
	c.p, c.q = c.q, rem
	return a, true
}

// floorDivMod returns (a, r) such that:
//
//	p = a*q + r
//
// with 0 <= r < |q| and a = floor(p/q) for q > 0.
// We assume q > 0 in this library (Rational invariant), but keep it general-ish.
func floorDivMod(p, q int64) (a, r int64) {
	if q == 0 {
		panic("floorDivMod: q=0")
	}
	if q < 0 {
		p = -p
		q = -q
	}

	// Go / truncates toward zero. We need floor for negatives.
	a = p / q
	r = p % q
	if r < 0 {
		r += q
		a -= 1
	}
	return a, r
}

// Convenience: make a Rational from ints and panic on error (tests only).
func mustRat(p, q int64) Rational {
	r, err := NewRational(p, q)
	if err != nil {
		panic(fmt.Sprintf("bad rational %d/%d: %v", p, q, err))
	}
	return r
}
