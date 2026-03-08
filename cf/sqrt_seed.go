// sqrt_seed.go v1
package cf

import "fmt"

// DefaultSqrtSeed returns a simple exact positive seed for Newton sqrt.
//
// Policy:
//   - x < 0  => error
//   - x = 0  => 1
//   - 0 < x <= 1 => 1
//   - x > 1  => x
//
// This is intentionally simple and exact. It is a seed policy, not a proof.
func DefaultSqrtSeed(x Rational) (Rational, error) {
	if x.Cmp(intRat(0)) < 0 {
		return Rational{}, fmt.Errorf("sqrt of negative rational: %v", x)
	}
	if x.Cmp(intRat(1)) <= 0 {
		return intRat(1), nil
	}
	return x, nil
}

// NewSqrtApproxCFDefault is a convenience wrapper that uses DefaultSqrtSeed.
func NewSqrtApproxCFDefault(x Rational, steps int) (ContinuedFraction, error) {
	seed, err := DefaultSqrtSeed(x)
	if err != nil {
		return nil, err
	}
	return NewSqrtApproxCF(x, seed, steps)
}

// sqrt_seed.go v1
