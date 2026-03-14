// sqrt_source_api2_test.go v1
package cf

import "testing"

func TestSqrtApproxFromApproxRangeSeed2_HonorsExplicitSeedOverride(t *testing.T) {
	a := CFApprox{
		Convergent:  mustRat(3, 2),
		Range:       NewRange(mustRat(4, 3), mustRat(3, 2), true, true),
		PrefixTerms: 2,
	}

	seed := mustRat(1, 1)
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
		Seed:     &seed,
	}

	got, err := SqrtApproxFromApproxRangeSeed2(a, p)
	if err != nil {
		t.Fatalf("SqrtApproxFromApproxRangeSeed2 failed: %v", err)
	}

	// With explicit seed=1 and residual stopping, iteration stops early at 49/40
	// because |(49/40)^2 - 3/2| = 1/1600 <= 1/1000.
	want := mustRat(49, 40)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSqrtApproxFromApproxRangeSeed2_UsesRangeDerivedSeedWhenMissing(t *testing.T) {
	a := CFApprox{
		Convergent:  mustRat(3, 2),
		Range:       NewRange(mustRat(4, 3), mustRat(3, 2), true, true),
		PrefixTerms: 2,
	}

	p := SqrtPolicy2{
		MaxSteps: 1,
		Tol:      mustRat(1, 1000),
		Seed:     nil,
	}

	got, err := SqrtApproxFromApproxRangeSeed2(a, p)
	if err != nil {
		t.Fatalf("SqrtApproxFromApproxRangeSeed2 failed: %v", err)
	}

	if got.Cmp(intRat(0)) <= 0 {
		t.Fatalf("expected positive approximation, got %v", got)
	}
}

func TestSqrtApproxCFFromApproxRangeSeed2_HasExpectedLeadingDigit(t *testing.T) {
	a := CFApprox{
		Convergent:  mustRat(3, 2),
		Range:       NewRange(mustRat(4, 3), mustRat(3, 2), true, true),
		PrefixTerms: 2,
	}

	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	cf, err := SqrtApproxCFFromApproxRangeSeed2(a, p)
	if err != nil {
		t.Fatalf("SqrtApproxCFFromApproxRangeSeed2 failed: %v", err)
	}

	got := collectTerms(cf, 8)
	if len(got) == 0 {
		t.Fatalf("expected non-empty CF")
	}
	if got[0] != 1 {
		t.Fatalf("got first digit %d want 1 full=%v", got[0], got)
	}
}

func TestSqrtApproxTermsFromApproxRangeSeed2_RejectsNegativeDigits(t *testing.T) {
	a := CFApprox{
		Convergent:  mustRat(3, 2),
		Range:       NewRange(mustRat(4, 3), mustRat(3, 2), true, true),
		PrefixTerms: 2,
	}

	_, err := SqrtApproxTermsFromApproxRangeSeed2(a, DefaultSqrtPolicy2(), -1)
	if err == nil {
		t.Fatalf("expected error for negative digits")
	}
}

// sqrt_source_api2_test.go v1
