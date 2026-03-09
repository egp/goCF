// lambert_pi_tail_test.go v2
package cf

import "testing"

func TestLambertPiOver4TailRangeAfterPrefix_RejectsNegative(t *testing.T) {
	_, ok, err := LambertPiOver4TailRangeAfterPrefix(-1)
	if err == nil {
		t.Fatalf("expected error for negative prefixTerms")
	}
	if ok {
		t.Fatalf("expected ok=false on error")
	}
}

func TestLambertPiOver4TailRangeAfterPrefix_Prefix0(t *testing.T) {
	got, ok, err := LambertPiOver4TailRangeAfterPrefix(0)
	if err != nil {
		t.Fatalf("LambertPiOver4TailRangeAfterPrefix failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	wantLo := mustRat(3, 4)
	wantHi := mustRat(1, 1)
	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v] want [%v,%v]", got.Lo, got.Hi, wantLo, wantHi)
	}
}

func TestLambertPiOver4TailRangeAfterPrefix_Prefix1(t *testing.T) {
	got, ok, err := LambertPiOver4TailRangeAfterPrefix(1)
	if err != nil {
		t.Fatalf("LambertPiOver4TailRangeAfterPrefix failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	wantLo := mustRat(1, 1)
	wantHi := mustRat(4, 3)
	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v] want [%v,%v]", got.Lo, got.Hi, wantLo, wantHi)
	}
}

func TestLambertPiOver4TailRangeAfterPrefix_Prefix2NotYetSpecialized(t *testing.T) {
	_, ok, err := LambertPiOver4TailRangeAfterPrefix(2)
	if err != nil {
		t.Fatalf("LambertPiOver4TailRangeAfterPrefix failed: %v", err)
	}
	if ok {
		t.Fatalf("expected ok=false for unsupported specialized prefix")
	}
}

func TestLambertPiOver4ApproxFromPrefix_Prefix1UsesSpecializedTailRange(t *testing.T) {
	got, err := LambertPiOver4ApproxFromPrefix(1)
	if err != nil {
		t.Fatalf("LambertPiOver4ApproxFromPrefix failed: %v", err)
	}

	wantConv := mustRat(0, 1)
	if got.Convergent.Cmp(wantConv) != 0 {
		t.Fatalf("got convergent %v want %v", got.Convergent, wantConv)
	}
	if got.Range == nil {
		t.Fatalf("expected non-nil range")
	}

	// Prefix 1 means x = 1/tail, with tail in [1, 4/3]
	// so x in [3/4, 1].
	wantLo := mustRat(3, 4)
	wantHi := mustRat(1, 1)
	if got.Range.Lo.Cmp(wantLo) != 0 || got.Range.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got range [%v,%v] want [%v,%v]", got.Range.Lo, got.Range.Hi, wantLo, wantHi)
	}
}

func TestLambertPiOver4ApproxFromPrefix_Prefix2FallsBackToGenericLowerBound(t *testing.T) {
	got, err := LambertPiOver4ApproxFromPrefix(2)
	if err != nil {
		t.Fatalf("LambertPiOver4ApproxFromPrefix failed: %v", err)
	}

	wantConv := mustRat(1, 1)
	if got.Convergent.Cmp(wantConv) != 0 {
		t.Fatalf("got convergent %v want %v", got.Convergent, wantConv)
	}
	if got.Range == nil {
		t.Fatalf("expected non-nil range")
	}

	// Generic lower-bound-only behavior:
	// x = 1/(1 + 1/tail), tail >= 1 => [1/2, 1]
	wantLo := mustRat(1, 2)
	wantHi := mustRat(1, 1)
	if got.Range.Lo.Cmp(wantLo) != 0 || got.Range.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got range [%v,%v] want [%v,%v]", got.Range.Lo, got.Range.Hi, wantLo, wantHi)
	}
	if !got.Range.Contains(mustRat(3, 4)) {
		t.Fatalf("expected range %v to contain a valid unfinished Lambert value like 3/4", *got.Range)
	}
}

func TestLambertPiOver4ApproxFromPrefix_RejectsZeroPrefixTerms(t *testing.T) {
	_, err := LambertPiOver4ApproxFromPrefix(0)
	if err == nil {
		t.Fatalf("expected error for zero prefixTerms")
	}
}

// lambert_pi_tail_test.go v2
