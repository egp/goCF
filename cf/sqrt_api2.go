// sqrt_api2.go v1
package cf

import "fmt"

// SqrtPolicy2 is the parallel policy type for the restructured sqrt path.
type SqrtPolicy2 struct {
	MaxSteps int
	Tol      Rational
	Seed     *Rational
}

func DefaultSqrtPolicy2() SqrtPolicy2 {
	return SqrtPolicy2{
		MaxSteps: 5,
		Tol:      mustRat(1, 1_000_000_000_000),
		Seed:     nil,
	}
}

func (p SqrtPolicy2) Validate() error {
	if p.MaxSteps < 0 {
		return fmt.Errorf("SqrtPolicy2: negative MaxSteps %d", p.MaxSteps)
	}
	if p.Tol.Cmp(intRat(0)) < 0 {
		return fmt.Errorf("SqrtPolicy2: negative Tol %v", p.Tol)
	}
	if p.Seed != nil {
		if p.Seed.r.Sign() == 0 {
			return fmt.Errorf("SqrtPolicy2: zero Seed")
		}
		if p.Seed.r.Sign() < 0 {
			return fmt.Errorf("SqrtPolicy2: negative Seed %v", *p.Seed)
		}
	}
	return nil
}

// SqrtApprox2 uses a simple default policy to compute a bounded rational
// approximation to sqrt(x).
func SqrtApprox2(x Rational) (Rational, error) {
	p := DefaultSqrtPolicy2()
	approx, _, err := SqrtCoreApproxRationalUntilResidualDefault2(x, p.MaxSteps, p.Tol)
	return approx, err
}

// SqrtApproxCF2 returns a ContinuedFraction source for the bounded default
// sqrt approximation produced by SqrtApprox2.
func SqrtApproxCF2(x Rational) (ContinuedFraction, error) {
	approx, err := SqrtApprox2(x)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(approx), nil
}

// SqrtApproxTerms2 returns up to digits CF terms for the bounded default
// sqrt approximation produced by SqrtApprox2.
func SqrtApproxTerms2(x Rational, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTerms2: negative digits %d", digits)
	}
	cf, err := SqrtApproxCF2(x)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

// SqrtApproxWithPolicy2 computes a bounded rational approximation to sqrt(x)
// using the supplied policy.
func SqrtApproxWithPolicy2(x Rational, p SqrtPolicy2) (Rational, error) {
	if err := p.Validate(); err != nil {
		return Rational{}, err
	}

	if p.Seed != nil {
		approx, _, err := SqrtCoreApproxRationalUntilResidual(x, *p.Seed, p.MaxSteps, p.Tol)
		return approx, err
	}

	approx, _, err := SqrtCoreApproxRationalUntilResidualDefault2(x, p.MaxSteps, p.Tol)
	return approx, err
}

// SqrtCoreApproxRationalUntilResidualDefault2 performs bounded Newton iteration for
// sqrt(x), using SqrtSeedDefault(x), and stops early once |y^2 - x| <= tol.
func SqrtCoreApproxRationalUntilResidualDefault2(x Rational, maxSteps int, tol Rational) (Rational, bool, error) {
	seed, err := SqrtSeedDefault(x)
	if err != nil {
		return Rational{}, false, err
	}
	return SqrtCoreApproxRationalUntilResidual(x, seed, maxSteps, tol)
}

// sqrt_api2.go v1
