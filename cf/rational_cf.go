// rational_cf.go v3
package cf

import "math/big"

// RationalCF streams the continued-fraction terms of a Rational.
type RationalCF struct {
	r    Rational
	done bool
}

func NewRationalCF(r Rational) *RationalCF {
	return &RationalCF{r: r}
}

func (c *RationalCF) Next() (int64, bool) {
	if c.done {
		return 0, false
	}

	num, den := c.r.ratNumDen()
	if den.Sign() == 0 {
		// Should be impossible if constructed from Rational, but keep it safe.
		c.done = true
		return 0, false
	}

	a, rem := floorDivModBig(num, den)
	if !a.IsInt64() {
		c.done = true
		return 0, false
	}
	digit := a.Int64()

	if rem.Sign() == 0 {
		c.done = true
		return digit, true
	}

	// Next step: x = den/rem with rem > 0
	next, err := newRationalBig(den, rem)
	if err != nil {
		c.done = true
		return 0, false
	}
	c.r = next
	return digit, true
}

// floorDivModBig returns (a, r) such that:
//
//	p = a*q + r
//
// with 0 <= r < |q| and a = floor(p/q) for q > 0.
// We assume q > 0 in this library (Rational invariant), but keep it general-ish.
func floorDivModBig(p, q *big.Int) (a, r *big.Int) {
	if q.Sign() == 0 {
		panic("floorDivModBig: q=0")
	}

	pp := new(big.Int).Set(p)
	qq := new(big.Int).Set(q)
	if qq.Sign() < 0 {
		pp.Neg(pp)
		qq.Neg(qq)
	}

	a = new(big.Int)
	r = new(big.Int)
	a.QuoRem(pp, qq, r)

	// big.Int QuoRem truncates toward zero. We need floor for negatives.
	if r.Sign() < 0 {
		r.Add(r, qq)
		a.Sub(a, big.NewInt(1))
	}
	return a, r
}

// rational_cf.go v3
