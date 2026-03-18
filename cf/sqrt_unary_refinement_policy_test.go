// cf/sqrt_unary_refinement_policy_test.go v2
package cf

import (
	"strings"
	"testing"
)

func TestDefaultSqrtUnaryRefinementPolicy_HasOneStepPerInput(t *testing.T) {
	p := defaultSqrtUnaryRefinementPolicy()
	if p.StepsPerInput != 1 {
		t.Fatalf("StepsPerInput: got %d want 1", p.StepsPerInput)
	}
}

func TestSqrtUnaryRefinementPolicy_Validate_RejectsNonpositiveSteps(t *testing.T) {
	cases := []struct {
		name  string
		steps int
	}{
		{"zero", 0},
		{"minus_one", -1},
		{"minus_three", -3},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := sqrtUnaryRefinementPolicy{StepsPerInput: tc.steps}
			err := p.validate()
			if err == nil {
				t.Fatalf("expected error")
			}
			if !strings.Contains(err.Error(), "StepsPerInput must be > 0") {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestSqrtUnaryRefinementPolicy_Validate_AcceptsPositiveSteps(t *testing.T) {
	cases := []int{1, 2, 3}

	for _, steps := range cases {
		p := sqrtUnaryRefinementPolicy{StepsPerInput: steps}
		if err := p.validate(); err != nil {
			t.Fatalf("steps=%d validate failed: %v", steps, err)
		}
	}
}

// cf/sqrt_unary_refinement_policy_test.go v2
