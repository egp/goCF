// lft.go v1
package cf

import "fmt"

// ULFT: (A x + B) / (C x + D)
type ULFT struct {
	A, B, C, D int64
}

func NewULFT(a, b, c, d int64) ULFT {
	return ULFT{A: a, B: b, C: c, D: d}
}

// ApplyRat evaluates the ULFT exactly on a rational x, with overflow detection.
//
// (A*p/q + B) / (C*p/q + D) = (A*p + B*q) / (C*p + D*q)
func (t ULFT) ApplyRat(x Rational) (Rational, error) {
	if err := t.Validate(); err != nil {
		return Rational{}, err
	}

	ap, ok := mul64(t.A, x.P)
	if !ok {
		return Rational{}, ErrOverflow
	}
	bq, ok := mul64(t.B, x.Q)
	if !ok {
		return Rational{}, ErrOverflow
	}
	num, ok := add64(ap, bq)
	if !ok {
		return Rational{}, ErrOverflow
	}

	cp, ok := mul64(t.C, x.P)
	if !ok {
		return Rational{}, ErrOverflow
	}
	dq, ok := mul64(t.D, x.Q)
	if !ok {
		return Rational{}, ErrOverflow
	}
	den, ok := add64(cp, dq)
	if !ok {
		return Rational{}, ErrOverflow
	}

	return NewRational(num, den)
}

func (t ULFT) String() string {
	return fmt.Sprintf("[[%d %d],[%d %d]]", t.A, t.B, t.C, t.D)
}

// BLFT: (A x y + B x + C y + D) / (E x y + F x + G y + H)
type BLFT struct {
	A, B, C, D int64
	E, F, G, H int64
}

func NewBLFT(a, b, c, d, e, f, g, h int64) BLFT {
	return BLFT{A: a, B: b, C: c, D: d, E: e, F: f, G: g, H: h}
}

// ApplyRat evaluates the BLFT exactly on rationals x=p/q and y=u/v, with overflow detection.
//
// N = A*p*u + B*p*v + C*u*q + D*q*v
// D = E*p*u + F*p*v + G*u*q + H*q*v
func (t BLFT) ApplyRat(x, y Rational) (Rational, error) {
	p, q := x.P, x.Q
	u, v := y.P, y.Q

	apu, ok := mul64(t.A, p)
	if !ok {
		return Rational{}, ErrOverflow
	}
	apu, ok = mul64(apu, u)
	if !ok {
		return Rational{}, ErrOverflow
	}

	bpv, ok := mul64(t.B, p)
	if !ok {
		return Rational{}, ErrOverflow
	}
	bpv, ok = mul64(bpv, v)
	if !ok {
		return Rational{}, ErrOverflow
	}

	cuq, ok := mul64(t.C, u)
	if !ok {
		return Rational{}, ErrOverflow
	}
	cuq, ok = mul64(cuq, q)
	if !ok {
		return Rational{}, ErrOverflow
	}

	dqv, ok := mul64(t.D, q)
	if !ok {
		return Rational{}, ErrOverflow
	}
	dqv, ok = mul64(dqv, v)
	if !ok {
		return Rational{}, ErrOverflow
	}

	N, ok := add64(apu, bpv)
	if !ok {
		return Rational{}, ErrOverflow
	}
	N, ok = add64(N, cuq)
	if !ok {
		return Rational{}, ErrOverflow
	}
	N, ok = add64(N, dqv)
	if !ok {
		return Rational{}, ErrOverflow
	}

	epu, ok := mul64(t.E, p)
	if !ok {
		return Rational{}, ErrOverflow
	}
	epu, ok = mul64(epu, u)
	if !ok {
		return Rational{}, ErrOverflow
	}

	fpv, ok := mul64(t.F, p)
	if !ok {
		return Rational{}, ErrOverflow
	}
	fpv, ok = mul64(fpv, v)
	if !ok {
		return Rational{}, ErrOverflow
	}

	guq, ok := mul64(t.G, u)
	if !ok {
		return Rational{}, ErrOverflow
	}
	guq, ok = mul64(guq, q)
	if !ok {
		return Rational{}, ErrOverflow
	}

	hqv, ok := mul64(t.H, q)
	if !ok {
		return Rational{}, ErrOverflow
	}
	hqv, ok = mul64(hqv, v)
	if !ok {
		return Rational{}, ErrOverflow
	}

	Dd, ok := add64(epu, fpv)
	if !ok {
		return Rational{}, ErrOverflow
	}
	Dd, ok = add64(Dd, guq)
	if !ok {
		return Rational{}, ErrOverflow
	}
	Dd, ok = add64(Dd, hqv)
	if !ok {
		return Rational{}, ErrOverflow
	}

	return NewRational(N, Dd)
}

// lft.go v1
