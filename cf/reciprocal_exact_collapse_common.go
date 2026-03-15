// reciprocal_exact_collapse_common.go v1
package cf

type reciprocalExactCollapseCore struct {
	err     error
	done    bool
	started bool
	exactCF ContinuedFraction
	approx  *Rational
}

func (c *reciprocalExactCollapseCore) Err() error { return c.err }

func (c *reciprocalExactCollapseCore) init(eval func() (Rational, error)) bool {
	if c.started {
		return c.err == nil
	}
	c.started = true

	x, err := eval()
	if err != nil {
		c.err = err
		c.done = true
		return false
	}

	c.approx = &x
	c.exactCF = NewRationalCF(x)
	return true
}

func (c *reciprocalExactCollapseCore) Next(eval func() (Rational, error)) (int64, bool) {
	if c.done {
		return 0, false
	}
	if c.err != nil {
		c.done = true
		return 0, false
	}
	if !c.init(eval) {
		return 0, false
	}

	d, ok := c.exactCF.Next()
	if !ok {
		c.done = true
		return 0, false
	}
	return d, true
}

// reciprocal_exact_collapse_common.go v1
