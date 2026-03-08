// sqrt_policy_test.go v1
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

// sqrt_policy_test.go v1
