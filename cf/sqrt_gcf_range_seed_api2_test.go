// sqrt_gcf_range_seed_api2_test.go v1
package cf

import "testing"

func TestSqrtApproxFromGCFApproxRangeSeed2_HonorsExplicitSeedOverride(t *testing.T) {
	a := GCFApprox{
		Convergent:  mustRat(3, 2),
		Range:       ptrRange(NewRange(mustRat(4, 3), mustRat(3, 2), true, true)),
		PrefixTerms: 2,
	}

	seed := mustRat(1, 1)
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
		Seed:     &seed,
	}

	got, err := SqrtApproxFromGCFApproxRangeSeed2(a, p)
	if err != nil {
		t.Fatalf("SqrtApproxFromGCFApproxRangeSeed2 failed: %v", err)
	}

	want := mustRat(49, 40)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSqrtApproxFromGCFApproxRangeSeed2_UsesRangeDerivedSeedWhenMissing(t *testing.T) {
	a := GCFApprox{
		Convergent:  mustRat(3, 2),
		Range:       ptrRange(NewRange(mustRat(4, 3), mustRat(3, 2), true, true)),
		PrefixTerms: 2,
	}

	p := SqrtPolicy2{
		MaxSteps: 1,
		Tol:      mustRat(1, 1000),
	}

	got, err := SqrtApproxFromGCFApproxRangeSeed2(a, p)
	if err != nil {
		t.Fatalf("SqrtApproxFromGCFApproxRangeSeed2 failed: %v", err)
	}
	if got.Cmp(intRat(0)) <= 0 {
		t.Fatalf("expected positive approximation, got %v", got)
	}
}

func TestSqrtApproxFromGCFApproxRangeSeed2_FallsBackWhenRangeMissing(t *testing.T) {
	a := GCFApprox{
		Convergent:  mustRat(3, 2),
		Range:       nil,
		PrefixTerms: 2,
	}

	p := SqrtPolicy2{
		MaxSteps: 1,
		Tol:      mustRat(1, 1000),
	}

	got, err := SqrtApproxFromGCFApproxRangeSeed2(a, p)
	if err != nil {
		t.Fatalf("SqrtApproxFromGCFApproxRangeSeed2 failed: %v", err)
	}
	if got.Cmp(intRat(0)) <= 0 {
		t.Fatalf("expected positive approximation, got %v", got)
	}
}

func TestSqrtApproxCFFromGCFSourceRangeSeed2_HasExpectedLeadingDigit(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	cf, err := SqrtApproxCFFromGCFSourceRangeSeed2(
		AdaptCFToGCF(Sqrt2CF()),
		2,
		p,
	)
	if err != nil {
		t.Fatalf("SqrtApproxCFFromGCFSourceRangeSeed2 failed: %v", err)
	}

	got := collectTerms(cf, 8)
	if len(got) == 0 {
		t.Fatalf("expected non-empty CF")
	}
	if got[0] != 1 {
		t.Fatalf("got first digit %d want 1 full=%v", got[0], got)
	}
}

func TestSqrtApproxTermsFromGCFSourceRangeSeedDefault2_RejectsNegativeDigits(t *testing.T) {
	_, err := SqrtApproxTermsFromGCFSourceRangeSeedDefault2(
		AdaptCFToGCF(Sqrt2CF()),
		2,
		-1,
	)
	if err == nil {
		t.Fatalf("expected error for negative digits")
	}
}

func ptrRange(r Range) *Range {
	return &r
}

// sqrt_gcf_range_seed_api2_test.go v1
