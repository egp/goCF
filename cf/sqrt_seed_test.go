// sqrt_seed_test.go v1
package cf

import "testing"

func TestDefaultSqrtSeed_NegativeRejected(t *testing.T) {
	_, err := DefaultSqrtSeed(mustRat(-1, 1))
	if err == nil {
		t.Fatalf("expected error for negative input")
	}
}

func TestDefaultSqrtSeed_ZeroGivesOne(t *testing.T) {
	got, err := DefaultSqrtSeed(mustRat(0, 1))
	if err != nil {
		t.Fatalf("DefaultSqrtSeed failed: %v", err)
	}
	if got.Cmp(intRat(1)) != 0 {
		t.Fatalf("got %v, want 1", got)
	}
}

func TestDefaultSqrtSeed_FractionBelowOneGivesOne(t *testing.T) {
	got, err := DefaultSqrtSeed(mustRat(3, 4))
	if err != nil {
		t.Fatalf("DefaultSqrtSeed failed: %v", err)
	}
	if got.Cmp(intRat(1)) != 0 {
		t.Fatalf("got %v, want 1", got)
	}
}

func TestDefaultSqrtSeed_AboveOneReturnsX(t *testing.T) {
	x := mustRat(5, 2)
	got, err := DefaultSqrtSeed(x)
	if err != nil {
		t.Fatalf("DefaultSqrtSeed failed: %v", err)
	}
	if got.Cmp(x) != 0 {
		t.Fatalf("got %v, want %v", got, x)
	}
}

func TestNewSqrtApproxCFDefault_Sqrt2_ThreeSteps(t *testing.T) {
	cf, err := NewSqrtApproxCFDefault(mustRat(2, 1), 3)
	if err != nil {
		t.Fatalf("NewSqrtApproxCFDefault failed: %v", err)
	}

	got := collectTerms(cf, 16)
	want := []int64{1, 2, 2, 2, 2, 2, 2, 2} // 577/408

	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d full=%v", i, got[i], want[i], got)
		}
	}
}

// sqrt_seed_test.go v1
