// sqrt_seed_from_range_test.go v1
package cf

import "testing"

func TestApproxFromCFPrefix_FiniteSourceExact(t *testing.T) {
	conv, rng, err := ApproxFromCFPrefix(NewSliceCF(3, 7, 16), 10)
	if err != nil {
		t.Fatalf("ApproxFromCFPrefix failed: %v", err)
	}

	want := mustRat(355, 113)
	if conv.Cmp(want) != 0 {
		t.Fatalf("conv got %v, want %v", conv, want)
	}
	if rng.Lo.Cmp(want) != 0 || rng.Hi.Cmp(want) != 0 {
		t.Fatalf("range got [%v,%v], want exact [%v,%v]", rng.Lo, rng.Hi, want, want)
	}
}

func TestApproxFromCFPrefix_InfiniteSourceTwoTerms(t *testing.T) {
	conv, rng, err := ApproxFromCFPrefix(Sqrt2CF(), 2)
	if err != nil {
		t.Fatalf("ApproxFromCFPrefix failed: %v", err)
	}

	wantConv := mustRat(3, 2)
	wantLo := mustRat(4, 3)
	wantHi := mustRat(3, 2)

	if conv.Cmp(wantConv) != 0 {
		t.Fatalf("conv got %v, want %v", conv, wantConv)
	}
	if rng.Lo.Cmp(wantLo) != 0 || rng.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("range got [%v,%v], want [%v,%v]", rng.Lo, rng.Hi, wantLo, wantHi)
	}
}

func TestDefaultSqrtSeedFromRange_Sqrt2PrefixTwoTerms(t *testing.T) {
	rng := NewRange(mustRat(4, 3), mustRat(3, 2), true, true)

	got, err := DefaultSqrtSeedFromRange(rng)
	if err != nil {
		t.Fatalf("DefaultSqrtSeedFromRange failed: %v", err)
	}

	// midpoint = (4/3 + 3/2)/2 = 17/12
	// one Newton step toward sqrt(17/12) from base 17/12 gives:
	// ((17/12) + 1) / 2 = 29/24
	want := mustRat(29, 24)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestDefaultSqrtSeedFromCFPrefix_Sqrt2TwoTerms(t *testing.T) {
	got, err := DefaultSqrtSeedFromCFPrefix(Sqrt2CF(), 2)
	if err != nil {
		t.Fatalf("DefaultSqrtSeedFromCFPrefix failed: %v", err)
	}

	want := mustRat(29, 24)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestNewSqrtApproxCFFromSourceRangeSeed_Sqrt2PrefixTwoTerms(t *testing.T) {
	p := SqrtPolicy{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	cf, err := NewSqrtApproxCFFromSourceRangeSeed(Sqrt2CF(), 2, p)
	if err != nil {
		t.Fatalf("NewSqrtApproxCFFromSourceRangeSeed failed: %v", err)
	}

	got := collectTerms(cf, 8)
	if len(got) == 0 {
		t.Fatalf("expected non-empty CF")
	}
	if got[0] != 1 {
		t.Fatalf("got first digit %d, want 1; full=%v", got[0], got)
	}
}

func TestSqrtApproxTermsFromSourceRangeSeedDefault_RejectsNegativeDigits(t *testing.T) {
	_, err := SqrtApproxTermsFromSourceRangeSeedDefault(Sqrt2CF(), 2, -1)
	if err == nil {
		t.Fatalf("expected error for negative digits")
	}
}

// sqrt_seed_from_range_test.go v1
