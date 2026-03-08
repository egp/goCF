// sqrt_seed_compare_test.go v1
package cf

import "testing"

func TestRangeSeedAndDefaultSeed_AreComparableOnSqrt2Prefix2(t *testing.T) {
	conv, rng, err := ApproxFromCFPrefix(Sqrt2CF(), 2)
	if err != nil {
		t.Fatalf("ApproxFromCFPrefix failed: %v", err)
	}

	seedDefault, err := DefaultSqrtSeed(conv)
	if err != nil {
		t.Fatalf("DefaultSqrtSeed failed: %v", err)
	}
	seedRange, err := DefaultSqrtSeedFromRange(rng)
	if err != nil {
		t.Fatalf("DefaultSqrtSeedFromRange failed: %v", err)
	}

	if seedDefault.Cmp(intRat(0)) <= 0 {
		t.Fatalf("default seed must be positive, got %v", seedDefault)
	}
	if seedRange.Cmp(intRat(0)) <= 0 {
		t.Fatalf("range seed must be positive, got %v", seedRange)
	}
}

func TestRangeSeedResidualVsDefaultSeed_OnSqrt2Prefix2(t *testing.T) {
	conv, rng, err := ApproxFromCFPrefix(Sqrt2CF(), 2)
	if err != nil {
		t.Fatalf("ApproxFromCFPrefix failed: %v", err)
	}

	seedDefault, err := DefaultSqrtSeed(conv)
	if err != nil {
		t.Fatalf("DefaultSqrtSeed failed: %v", err)
	}
	seedRange, err := DefaultSqrtSeedFromRange(rng)
	if err != nil {
		t.Fatalf("DefaultSqrtSeedFromRange failed: %v", err)
	}

	// Compare both seeds against the same target xApprox = convergent(prefix).
	rd, err := SqrtResidualAbs(conv, seedDefault)
	if err != nil {
		t.Fatalf("SqrtResidualAbs(default) failed: %v", err)
	}
	rr, err := SqrtResidualAbs(conv, seedRange)
	if err != nil {
		t.Fatalf("SqrtResidualAbs(range) failed: %v", err)
	}

	// At this stage we do not assert one is strictly better than the other.
	// We only assert both are exact, nonnegative residuals and comparable.
	if rd.Cmp(intRat(0)) < 0 {
		t.Fatalf("default residual should be nonnegative, got %v", rd)
	}
	if rr.Cmp(intRat(0)) < 0 {
		t.Fatalf("range residual should be nonnegative, got %v", rr)
	}
}

func TestRangeSeedAfterOneNewtonStep_IsComparableToDefaultSeedPath(t *testing.T) {
	conv, rng, err := ApproxFromCFPrefix(Sqrt2CF(), 2)
	if err != nil {
		t.Fatalf("ApproxFromCFPrefix failed: %v", err)
	}

	seedDefault, err := DefaultSqrtSeed(conv)
	if err != nil {
		t.Fatalf("DefaultSqrtSeed failed: %v", err)
	}
	seedRange, err := DefaultSqrtSeedFromRange(rng)
	if err != nil {
		t.Fatalf("DefaultSqrtSeedFromRange failed: %v", err)
	}

	nextDefault, err := NewtonSqrtStep(conv, seedDefault)
	if err != nil {
		t.Fatalf("NewtonSqrtStep(default) failed: %v", err)
	}
	nextRange, err := NewtonSqrtStep(conv, seedRange)
	if err != nil {
		t.Fatalf("NewtonSqrtStep(range) failed: %v", err)
	}

	rd, err := SqrtResidualAbs(conv, nextDefault)
	if err != nil {
		t.Fatalf("SqrtResidualAbs(nextDefault) failed: %v", err)
	}
	rr, err := SqrtResidualAbs(conv, nextRange)
	if err != nil {
		t.Fatalf("SqrtResidualAbs(nextRange) failed: %v", err)
	}

	if rd.Cmp(intRat(0)) < 0 || rr.Cmp(intRat(0)) < 0 {
		t.Fatalf("residuals must be nonnegative, got default=%v range=%v", rd, rr)
	}
}

// sqrt_seed_compare_test.go v1
