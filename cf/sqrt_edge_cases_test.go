// sqrt_edge_cases_test.go v1
package cf

import "testing"

func TestSqrtApproxTermsAuto_ZeroDigits(t *testing.T) {
	got, err := SqrtApproxTermsAuto(mustRat(2, 1), 0)
	if err != nil {
		t.Fatalf("SqrtApproxTermsAuto failed: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

func TestSqrtApproxWithPolicy_RejectsZeroSeedInPolicy(t *testing.T) {
	zero := mustRat(0, 1)
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
		Seed:     &zero,
	}

	_, err := SqrtApproxWithPolicy(mustRat(2, 1), p)
	if err == nil {
		t.Fatalf("expected error for zero policy seed")
	}
}

func TestSqrtApproxWithPolicy_RejectsNegativeSeedInPolicy(t *testing.T) {
	neg := mustRat(-1, 1)
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
		Seed:     &neg,
	}

	_, err := SqrtApproxWithPolicy(mustRat(2, 1), p)
	if err == nil {
		t.Fatalf("expected error for negative policy seed")
	}
}

func TestNewSqrtApproxCFUntilResidualDefault_PerfectSquare(t *testing.T) {
	cf, err := NewSqrtApproxCFUntilResidualDefault(mustRat(9, 16), 5, mustRat(0, 1))
	if err != nil {
		t.Fatalf("NewSqrtApproxCFUntilResidualDefault failed: %v", err)
	}

	got := collectTerms(cf, 8)
	want := []int64{0, 1, 3} // 3/4
	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d full=%v", i, got[i], want[i], got)
		}
	}
}

func TestRationalApproxFromCFPrefix_RejectsNegativePrefixTerms(t *testing.T) {
	_, err := RationalApproxFromCFPrefix(Sqrt2CF(), -1)
	if err == nil {
		t.Fatalf("expected error for negative prefixTerms")
	}
}

func TestCFApproxFromPrefix_RejectsNegativePrefixTerms(t *testing.T) {
	_, err := CFApproxFromPrefix(Sqrt2CF(), -1)
	if err == nil {
		t.Fatalf("expected error for negative prefixTerms")
	}
}

func TestDefaultSqrtSeedFromRange_ExactMidpointSquare(t *testing.T) {
	// midpoint = (1 + 3)/2 = 2, so exact sqrt is not rational; use a case with midpoint 1.
	r := NewRange(mustRat(1, 2), mustRat(3, 2), true, true)

	got, err := DefaultSqrtSeedFromRange(r)
	if err != nil {
		t.Fatalf("DefaultSqrtSeedFromRange failed: %v", err)
	}

	// midpoint = 1, exact sqrt = 1
	want := mustRat(1, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestNewSqrtApproxCFFromApproxRangeSeed_HonorsExplicitSeedOverride(t *testing.T) {
	a := CFApprox{
		Convergent:  mustRat(3, 2),
		Range:       NewRange(mustRat(4, 3), mustRat(3, 2), true, true),
		PrefixTerms: 2,
	}

	seed := mustRat(1, 1)
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
		Seed:     &seed,
	}

	cf, err := NewSqrtApproxCFFromApproxRangeSeed(a, p)
	if err != nil {
		t.Fatalf("NewSqrtApproxCFFromApproxRangeSeed failed: %v", err)
	}

	got := collectTerms(cf, 8)
	want := []int64{1, 4, 1, 2} // sqrt(3/2) approx path from explicit seed=1 after 3 steps => 19601/16002
	if len(got) == 0 {
		t.Fatalf("expected non-empty CF")
	}
	if got[0] != want[0] {
		t.Fatalf("got first digit %d, want %d; full=%v", got[0], want[0], got)
	}
}

func TestSqrtRangeExact_RejectsOutsideRange(t *testing.T) {
	r := NewRange(mustRat(2, 1), mustRat(1, 1), true, true) // outside

	_, _, err := SqrtRangeExact(r)
	if err == nil {
		t.Fatalf("expected error for outside range")
	}
}

func TestSqrtRangeHeuristic_PreservesEndpointInclusions(t *testing.T) {
	r := NewRange(mustRat(1, 4), mustRat(9, 16), false, false)

	got, err := SqrtRangeHeuristic(r)
	if err != nil {
		t.Fatalf("SqrtRangeHeuristic failed: %v", err)
	}

	if got.IncLo || got.IncHi {
		t.Fatalf("expected inclusions preserved as false,false; got %v %v", got.IncLo, got.IncHi)
	}
}

func TestSqrtRangeHeuristic_RejectsOutsideRange(t *testing.T) {
	r := NewRange(mustRat(2, 1), mustRat(1, 1), true, true)

	_, err := SqrtRangeHeuristic(r)
	if err == nil {
		t.Fatalf("expected error for outside range")
	}
}

// sqrt_edge_cases_test.go v1
