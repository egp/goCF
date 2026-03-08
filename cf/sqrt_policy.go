// sqrt_policy.go v1
package cf

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

// sqrt_policy.go v1
