// sqrt_policy.go v2
package cf

import "fmt"

type SqrtPolicy struct {
	MaxSteps int
	Tol      Rational
}

func DefaultSqrtPolicy() SqrtPolicy {
	return SqrtPolicy{
		MaxSteps: 5,
		Tol:      mustRat(1, 1_000_000_000_000),
	}
}

func (p SqrtPolicy) Validate() error {
	if p.MaxSteps < 0 {
		return fmt.Errorf("SqrtPolicy: negative MaxSteps %d", p.MaxSteps)
	}
	if p.Tol.Cmp(intRat(0)) < 0 {
		return fmt.Errorf("SqrtPolicy: negative Tol %v", p.Tol)
	}
	return nil
}

// sqrt_policy.go v2
