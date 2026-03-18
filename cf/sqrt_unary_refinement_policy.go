// cf/sqrt_unary_refinement_policy.go v2
package cf

import "fmt"

type sqrtUnaryRefinementPolicy struct {
	StepsPerInput int
}

func defaultSqrtUnaryRefinementPolicy() sqrtUnaryRefinementPolicy {
	return sqrtUnaryRefinementPolicy{
		StepsPerInput: 1,
	}
}

func (p sqrtUnaryRefinementPolicy) validate() error {
	if p.StepsPerInput <= 0 {
		return fmt.Errorf("sqrtUnaryRefinementPolicy: StepsPerInput must be > 0, got %d", p.StepsPerInput)
	}
	return nil
}

// cf/sqrt_unary_refinement_policy.go v2
