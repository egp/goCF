// sqrt_api_equiv_test.go v1
package cf

import "testing"

func TestSqrtApproxWithPolicy_MatchesExplicitSeedAndPolicy(t *testing.T) {
	seed := mustRat(1, 1)
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
		Seed:     &seed,
	}

	gotPolicy, err := SqrtApproxWithPolicy(mustRat(2, 1), p)
	if err != nil {
		t.Fatalf("SqrtApproxWithPolicy failed: %v", err)
	}

	gotExplicit, err := SqrtApproxWithSeedAndPolicy(mustRat(2, 1), seed, p)
	if err != nil {
		t.Fatalf("SqrtApproxWithSeedAndPolicy failed: %v", err)
	}

	if gotPolicy.Cmp(gotExplicit) != 0 {
		t.Fatalf("policy=%v explicit=%v", gotPolicy, gotExplicit)
	}
}

func TestSqrtApproxCFWithPolicy_MatchesExplicitSeedAndPolicy(t *testing.T) {
	seed := mustRat(1, 1)
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
		Seed:     &seed,
	}

	cfPolicy, err := SqrtApproxCFWithPolicy(mustRat(2, 1), p)
	if err != nil {
		t.Fatalf("SqrtApproxCFWithPolicy failed: %v", err)
	}
	cfExplicit, err := SqrtApproxCFWithSeedAndPolicy(mustRat(2, 1), seed, p)
	if err != nil {
		t.Fatalf("SqrtApproxCFWithSeedAndPolicy failed: %v", err)
	}

	gotPolicy := collectTerms(cfPolicy, 16)
	gotExplicit := collectTerms(cfExplicit, 16)

	if len(gotPolicy) != len(gotExplicit) {
		t.Fatalf("len(policy)=%d len(explicit)=%d policy=%v explicit=%v",
			len(gotPolicy), len(gotExplicit), gotPolicy, gotExplicit)
	}
	for i := range gotPolicy {
		if gotPolicy[i] != gotExplicit[i] {
			t.Fatalf("policy[%d]=%d explicit[%d]=%d policy=%v explicit=%v",
				i, gotPolicy[i], i, gotExplicit[i], gotPolicy, gotExplicit)
		}
	}
}

func TestSqrtApproxTermsWithPolicy_MatchesExplicitSeedAndPolicy(t *testing.T) {
	seed := mustRat(1, 1)
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
		Seed:     &seed,
	}

	gotPolicy, err := SqrtApproxTermsWithPolicy(mustRat(2, 1), p, 16)
	if err != nil {
		t.Fatalf("SqrtApproxTermsWithPolicy failed: %v", err)
	}
	gotExplicit, err := SqrtApproxTermsWithSeedAndPolicy(mustRat(2, 1), seed, p, 16)
	if err != nil {
		t.Fatalf("SqrtApproxTermsWithSeedAndPolicy failed: %v", err)
	}

	if len(gotPolicy) != len(gotExplicit) {
		t.Fatalf("len(policy)=%d len(explicit)=%d policy=%v explicit=%v",
			len(gotPolicy), len(gotExplicit), gotPolicy, gotExplicit)
	}
	for i := range gotPolicy {
		if gotPolicy[i] != gotExplicit[i] {
			t.Fatalf("policy[%d]=%d explicit[%d]=%d policy=%v explicit=%v",
				i, gotPolicy[i], i, gotExplicit[i], gotPolicy, gotExplicit)
		}
	}
}

// sqrt_api_equiv_test.go v1
