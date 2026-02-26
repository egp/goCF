// rational.go v1
package cf

import (
	"fmt"
)

type Rational struct {
	P int64 // numerator
	Q int64 // denominator, always > 0
}

func NewRational(p, q int64) (Rational, error) {
	if q == 0 {
		return Rational{}, fmt.Errorf("zero denominator")
	}
	if q < 0 {
		p = -p
		q = -q
	}
	if p == 0 {
		return Rational{P: 0, Q: 1}, nil
	}
	g := gcd(abs(p), q)
	return Rational{P: p / g, Q: q / g}, nil
}

func (r Rational) String() string {
	return fmt.Sprintf("%d/%d", r.P, r.Q)
}

func (r Rational) Add(o Rational) (Rational, error) {
	// p/q + u/v = (pv + uq) / (qv)
	p := r.P*o.Q + o.P*r.Q
	q := r.Q * o.Q
	return NewRational(p, q)
}

func (r Rational) Sub(o Rational) (Rational, error) {
	p := r.P*o.Q - o.P*r.Q
	q := r.Q * o.Q
	return NewRational(p, q)
}

func (r Rational) Mul(o Rational) (Rational, error) {
	p := r.P * o.P
	q := r.Q * o.Q
	return NewRational(p, q)
}

func (r Rational) Div(o Rational) (Rational, error) {
	if o.P == 0 {
		return Rational{}, fmt.Errorf("division by zero")
	}
	p := r.P * o.Q
	q := r.Q * o.P
	return NewRational(p, q)
}

func (r Rational) Cmp(o Rational) int {
	// compare r ? o via cross-multiply (beware overflow in real life; OK for now)
	left := r.P * o.Q
	right := o.P * r.Q
	switch {
	case left < right:
		return -1
	case left > right:
		return 1
	default:
		return 0
	}
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

// rational.go v1
