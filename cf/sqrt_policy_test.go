// sqrt_policy_test.go v3
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
	if p.Seed != nil {
		t.Fatalf("expected nil Seed by default")
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

	seed := mustRat(1, 1)
	okSeed := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
		Seed:     &seed,
	}
	if err := okSeed.Validate(); err != nil {
		t.Fatalf("expected valid seeded policy, got %v", err)
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

	zero := mustRat(0, 1)
	badZeroSeed := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
		Seed:     &zero,
	}
	if err := badZeroSeed.Validate(); err == nil {
		t.Fatalf("expected error for zero Seed")
	}

	neg := mustRat(-1, 1)
	badNegSeed := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
		Seed:     &neg,
	}
	if err := badNegSeed.Validate(); err == nil {
		t.Fatalf("expected error for negative Seed")
	}
}

// sqrt_policy_test.go v3
