// range.go v16
package cf

import (
	"fmt"
)

type Range struct {
	Lo    Rational
	Hi    Rational
	IncLo bool
	IncHi bool
}

func NewRange(lo, hi Rational, incLo, incHi bool) Range {
	return Range{Lo: lo, Hi: hi, IncLo: incLo, IncHi: incHi}
}

func (r Range) IsInside() bool  { return r.Lo.Cmp(r.Hi) <= 0 }
func (r Range) IsOutside() bool { return r.Lo.Cmp(r.Hi) > 0 }

func (r Range) String() string {
	kind := "inside"
	if r.IsOutside() {
		kind = "outside"
	}
	return fmt.Sprintf("[%v,%v]{incLo=%t,incHi=%t,%s}", r.Lo, r.Hi, r.IncLo, r.IncHi, kind)
}

// Contains implements:
//   - Inside: standard interval membership with open/closed endpoints.
//   - Outside (Lo>Hi): union-of-rays semantics: (-∞,Hi] ∪ [Lo,∞), honoring endpoint inclusions.
func (r Range) Contains(x Rational) bool {
	cLo := x.Cmp(r.Lo)
	cHi := x.Cmp(r.Hi)

	if r.IsInside() {
		if cLo < 0 {
			return false
		}
		if cLo == 0 && !r.IncLo {
			return false
		}
		if cHi > 0 {
			return false
		}
		if cHi == 0 && !r.IncHi {
			return false
		}
		return true
	}

	// Outside: x <= Hi OR x >= Lo, with endpoint inclusion control.
	if cHi < 0 {
		return true
	}
	if cHi == 0 && r.IncHi {
		return true
	}
	if cLo > 0 {
		return true
	}
	if cLo == 0 && r.IncLo {
		return true
	}
	return false
}

func (r Range) ContainsZero() bool { return r.Contains(mustRat(0, 1)) }

// RefineMetric is a Gosper-style “uncertainty” metric used ONLY as a refinement heuristic.
// It supports ordering:
//
//	inside narrow < inside wide < outside wide < outside narrow
type RefineMetric struct {
	Outside bool
	// Inside: span = Hi-Lo (>=0)
	// Outside: gap  = Lo-Hi (>0), where (Hi,Lo) is the excluded gap.
	Magnitude Rational
}

func (m RefineMetric) String() string {
	k := "inside"
	if m.Outside {
		k = "outside"
	}
	return fmt.Sprintf("RefineMetric{%s,%v}", k, m.Magnitude)
}

func (m RefineMetric) Cmp(o RefineMetric) int {
	if m.Outside != o.Outside {
		// inside < outside
		if m.Outside {
			return 1
		}
		return -1
	}
	if !m.Outside {
		// inside: smaller span => narrower
		return m.Magnitude.Cmp(o.Magnitude)
	}
	// outside: larger excluded gap => narrower (reverse compare)
	return -m.Magnitude.Cmp(o.Magnitude)
}

func (r Range) RefineMetric() (RefineMetric, error) {
	if r.IsInside() {
		span, err := r.Hi.Sub(r.Lo)
		if err != nil {
			return RefineMetric{}, err
		}
		return RefineMetric{Outside: false, Magnitude: span}, nil
	}

	gap, err := r.Lo.Sub(r.Hi)
	if err != nil {
		return RefineMetric{}, err
	}
	return RefineMetric{Outside: true, Magnitude: gap}, nil
}

// FloorBounds returns a conservative pair (flo, fhi) such that for all x in r,
// floor(x) ∈ [flo, fhi]. Used by digit-safety checks.
//
// Conservative uniform implementation using only endpoint floors.
func (r Range) FloorBounds() (int64, int64, error) {
	fLo, err := floorRat(r.Lo)
	if err != nil {
		return 0, 0, err
	}
	fHi, err := floorRat(r.Hi)
	if err != nil {
		return 0, 0, err
	}
	if fLo <= fHi {
		return fLo, fHi, nil
	}
	return fHi, fLo, nil
}

// ApplyULFT maps an inside input range through ULFT and returns a conservative inside output range.
// This is used by SafeDigit(t,r) and must be interval-safe.
//
// Current strictness (by design):
//   - Requires r to be inside. (Bounder produces inside ranges; ULFT streaming expects that.)
//   - Rejects if denom range may include 0 over r (pole / discontinuity).
func (r Range) ApplyULFT(t ULFT) (Range, error) {
	if !r.IsInside() {
		return Range{}, fmt.Errorf("ApplyULFT requires inside range; got %v", r)
	}

	// Denom(x) = Cx + D. Over inside range, denom extrema occur at endpoints.
	denLo, err := evalLinearOnRat(t.C, t.D, r.Lo)
	if err != nil {
		return Range{}, err
	}
	denHi, err := evalLinearOnRat(t.C, t.D, r.Hi)
	if err != nil {
		return Range{}, err
	}

	denRange := NewRange(minRat(denLo, denHi), maxRat(denLo, denHi), true, true)
	if denRange.ContainsZero() {
		return Range{}, fmt.Errorf("ULFT denominator crosses 0 on range %v", r)
	}

	// ULFT is monotone on an interval that avoids poles, so endpoint images bound the image.
	zLo, err := t.ApplyRat(r.Lo)
	if err != nil {
		return Range{}, err
	}
	zHi, err := t.ApplyRat(r.Hi)
	if err != nil {
		return Range{}, err
	}

	outLo := minRat(zLo, zHi)
	outHi := maxRat(zLo, zHi)
	return NewRange(outLo, outHi, true, true), nil
}

// ---- helpers ----

// floorRat computes floor(p/q) for q>0 using Euclidean division.
// Works for negatives correctly.
func floorRat(x Rational) (int64, error) {
	if x.Q <= 0 {
		return 0, fmt.Errorf("floorRat: invalid denominator %d", x.Q)
	}
	p := x.P
	q := x.Q

	// Go truncates toward zero. Adjust for negative remainder.
	quo := p / q
	rem := p % q
	if rem != 0 && p < 0 {
		quo -= 1
	}
	return quo, nil
}

// evalLinearOnRat computes (a*x + b) exactly on x rational, with overflow detection.
func evalLinearOnRat(a, b int64, x Rational) (Rational, error) {
	// (a * (p/q) + b) = (a*p + b*q)/q
	ap, ok := mul64(a, x.P)
	if !ok {
		return Rational{}, ErrOverflow
	}
	bq, ok := mul64(b, x.Q)
	if !ok {
		return Rational{}, ErrOverflow
	}
	num, ok := add64(ap, bq)
	if !ok {
		return Rational{}, ErrOverflow
	}
	return NewRational(num, x.Q)
}

func minRat(a, b Rational) Rational {
	if a.Cmp(b) <= 0 {
		return a
	}
	return b
}

func maxRat(a, b Rational) Rational {
	if a.Cmp(b) >= 0 {
		return a
	}
	return b
}

// range.go v16
