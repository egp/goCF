// sqrt_policy.go v3
package cf

import "fmt"

type SqrtPolicy struct {
	MaxSteps int
	Tol      Rational
	Seed     *Rational
}

func DefaultSqrtPolicy() SqrtPolicy {
	return SqrtPolicy{
		MaxSteps: 5,
		Tol:      mustRat(1, 1_000_000_000_000),
		Seed:     nil,
	}
}

func (p SqrtPolicy) Validate() error {
	if p.MaxSteps < 0 {
		return fmt.Errorf("SqrtPolicy: negative MaxSteps %d", p.MaxSteps)
	}
	if p.Tol.Cmp(intRat(0)) < 0 {
		return fmt.Errorf("SqrtPolicy: negative Tol %v", p.Tol)
	}
	if p.Seed != nil {
		if p.Seed.r.Sign() == 0 {
			return fmt.Errorf("SqrtPolicy: zero Seed")
		}
		if p.Seed.r.Sign() < 0 {
			return fmt.Errorf("SqrtPolicy: negative Seed %v", *p.Seed)
		}
	}
	return nil
}

// sqrt_policy.go v3
