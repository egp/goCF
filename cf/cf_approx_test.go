// cf_approx_test.go v1
package cf

import "testing"

func TestCFApproxFromPrefix_FiniteSourceExact(t *testing.T) {
	got, err := CFApproxFromPrefix(NewSliceCF(3, 7, 16), 10)
	if err != nil {
		t.Fatalf("CFApproxFromPrefix failed: %v", err)
	}

	want := mustRat(355, 113)
	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("conv got %v, want %v", got.Convergent, want)
	}
	if got.Range.Lo.Cmp(want) != 0 || got.Range.Hi.Cmp(want) != 0 {
		t.Fatalf("range got [%v,%v], want exact [%v,%v]", got.Range.Lo, got.Range.Hi, want, want)
	}
	if got.PrefixTerms != 10 {
		t.Fatalf("PrefixTerms got %d, want 10", got.PrefixTerms)
	}
}

func TestCFApproxFromPrefix_InfiniteSourceTwoTerms(t *testing.T) {
	got, err := CFApproxFromPrefix(Sqrt2CF(), 2)
	if err != nil {
		t.Fatalf("CFApproxFromPrefix failed: %v", err)
	}

	wantConv := mustRat(3, 2)
	wantLo := mustRat(4, 3)
	wantHi := mustRat(3, 2)

	if got.Convergent.Cmp(wantConv) != 0 {
		t.Fatalf("conv got %v, want %v", got.Convergent, wantConv)
	}
	if got.Range.Lo.Cmp(wantLo) != 0 || got.Range.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("range got [%v,%v], want [%v,%v]", got.Range.Lo, got.Range.Hi, wantLo, wantHi)
	}
	if got.PrefixTerms != 2 {
		t.Fatalf("PrefixTerms got %d, want 2", got.PrefixTerms)
	}
}

func TestApproxFromCFPrefix_BackCompat(t *testing.T) {
	conv, rng, err := ApproxFromCFPrefix(Sqrt2CF(), 2)
	if err != nil {
		t.Fatalf("ApproxFromCFPrefix failed: %v", err)
	}
	a, err := CFApproxFromPrefix(Sqrt2CF(), 2)
	if err != nil {
		t.Fatalf("CFApproxFromPrefix failed: %v", err)
	}

	if conv.Cmp(a.Convergent) != 0 {
		t.Fatalf("conv mismatch: %v vs %v", conv, a.Convergent)
	}
	if rng.Lo.Cmp(a.Range.Lo) != 0 || rng.Hi.Cmp(a.Range.Hi) != 0 {
		t.Fatalf("range mismatch: [%v,%v] vs [%v,%v]", rng.Lo, rng.Hi, a.Range.Lo, a.Range.Hi)
	}
}

// cf_approx_test.go v1
