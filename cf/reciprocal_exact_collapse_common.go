// reciprocal_exact_collapse_common.go v2
package cf

type reciprocalExactCollapseCore struct {
	state   unaryStreamState
	exactCF ContinuedFraction
	approx  *Rational
}

func (c *reciprocalExactCollapseCore) Err() error { return c.state.Err() }

func (c *reciprocalExactCollapseCore) init(eval func() (Rational, error)) bool {
	if c.state.started {
		return c.state.err == nil
	}
	c.state.started = true

	x, err := eval()
	if err != nil {
		c.state.Fail(err)
		return false
	}

	c.approx = &x
	c.exactCF = NewRationalCF(x)
	return true
}

func (c *reciprocalExactCollapseCore) Next(eval func() (Rational, error)) (int64, bool) {
	if c.state.done {
		return 0, false
	}
	if c.state.err != nil {
		c.state.Exhaust()
		return 0, false
	}
	if !c.init(eval) {
		return 0, false
	}

	d, ok := c.exactCF.Next()
	if !ok {
		c.state.Exhaust()
		return 0, false
	}
	return d, true
}

// reciprocal_exact_collapse_common.go v2
