// sqrt_policy_test.go v2
package cf

import "testing"

func TestDefaultSqrtPolicy(t *testing.T) {
	p := DefaultSqrtPolicy()

	if p.MaxSteps != 5 {
		t.Fatalf("got MaxSteps=%d, want 5", p.MaxSteps)
	}
	if p.Tol.Cmp(mustRat(1, 1_000_000_000_000)) != 0 {
		t.Fatalf("got Tol=%v, want 1/1000000000000", p.Tol)
	}
}

func TestSqrtPolicyValidate(t *testing.T) {
	ok := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}
	if err := ok.Validate(); err != nil {
		t.Fatalf("expected valid policy, got %v", err)
	}

	badSteps := SqrtPolicy{
		MaxSteps: -1,
		Tol:      mustRat(1, 1000),
	}
	if err := badSteps.Validate(); err == nil {
		t.Fatalf("expected error for negative MaxSteps")
	}

	badTol := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(-1, 1000),
	}
	if err := badTol.Validate(); err == nil {
		t.Fatalf("expected error for negative Tol")
	}
}

// sqrt_policy_test.go v2
