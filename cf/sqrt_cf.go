// sqrt_cf.go v1
package cf

// NewSqrtApproxCF returns a ContinuedFraction source for the bounded rational
// Newton approximation to sqrt(x) produced after the requested number of steps
// from the given seed.
//
// This is a convenience adapter:
//
//	sqrt target -> exact rational approximation -> RationalCF
//
// It is not yet a true streaming sqrt operator.
func NewSqrtApproxCF(x, seed Rational, steps int) (ContinuedFraction, error) {
	approx, err := SqrtApproxRational(x, seed, steps)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

// sqrt_cf.go v1
