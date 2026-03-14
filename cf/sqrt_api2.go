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
	return sqrtApproxCanonical(x)
}

// SqrtApproxCF2 returns a ContinuedFraction source for the bounded default
// sqrt approximation produced by SqrtApprox2.
func SqrtApproxCF2(x Rational) (ContinuedFraction, error) {
	return sqrtApproxCFCanonical(x)
}

// SqrtApproxTerms2 returns up to digits CF terms for the bounded default
// sqrt approximation produced by SqrtApprox2.
func SqrtApproxTerms2(x Rational, digits int) ([]int64, error) {
	return sqrtApproxTermsCanonical(x, digits)
}

// SqrtApproxWithPolicy2 computes a bounded rational approximation to sqrt(x)
// using the supplied policy.
func SqrtApproxWithPolicy2(x Rational, p SqrtPolicy2) (Rational, error) {
	return sqrtApproxWithPolicyCanonical(x, p)
}

// sqrt_api2.go v1
