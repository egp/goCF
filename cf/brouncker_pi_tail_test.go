// brouncker_pi_tail_test.go v2
package cf

import "testing"

func TestBrouncker4OverPiTailRangeAfterPrefix_RejectsNegative(t *testing.T) {
	_, ok, err := Brouncker4OverPiTailRangeAfterPrefix(-1)
	if err == nil {
		t.Fatalf("expected error for negative prefixTerms")
	}
	if ok {
		t.Fatalf("expected ok=false on error")
	}
}

func TestBrouncker4OverPiTailRangeAfterPrefix_Prefix0(t *testing.T) {
	got, ok, err := Brouncker4OverPiTailRangeAfterPrefix(0)
	if err != nil {
		t.Fatalf("Brouncker4OverPiTailRangeAfterPrefix failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	wantLo := mustRat(1, 1)
	wantHi := mustRat(3, 2)
	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v] want [%v,%v]", got.Lo, got.Hi, wantLo, wantHi)
	}
}

func TestBrouncker4OverPiTailRangeAfterPrefix_Prefix1(t *testing.T) {
	got, ok, err := Brouncker4OverPiTailRangeAfterPrefix(1)
	if err != nil {
		t.Fatalf("Brouncker4OverPiTailRangeAfterPrefix failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	wantLo := mustRat(2, 1)
	wantHi := mustRat(5, 2)
	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v] want [%v,%v]", got.Lo, got.Hi, wantLo, wantHi)
	}
}

func TestBrouncker4OverPiTailRangeAfterPrefix_Prefix2NotYetSpecialized(t *testing.T) {
	_, ok, err := Brouncker4OverPiTailRangeAfterPrefix(2)
	if err != nil {
		t.Fatalf("Brouncker4OverPiTailRangeAfterPrefix failed: %v", err)
	}
	if ok {
		t.Fatalf("expected ok=false for unsupported specialized prefix")
	}
}

func TestBrouncker4OverPiApproxFromPrefix_Prefix1UsesSpecializedTailRange(t *testing.T) {
	got, err := Brouncker4OverPiApproxFromPrefix(1)
	if err != nil {
		t.Fatalf("Brouncker4OverPiApproxFromPrefix failed: %v", err)
	}

	wantConv := mustRat(1, 1)
	if got.Convergent.Cmp(wantConv) != 0 {
		t.Fatalf("got convergent %v want %v", got.Convergent, wantConv)
	}
	if got.Range == nil {
		t.Fatalf("expected non-nil range")
	}

	// Prefix 1 means x = 1 + 1/tail, with tail in [2, 5/2]
	// so x in [7/5, 3/2].
	wantLo := mustRat(7, 5)
	wantHi := mustRat(3, 2)
	if got.Range.Lo.Cmp(wantLo) != 0 || got.Range.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got range [%v,%v] want [%v,%v]", got.Range.Lo, got.Range.Hi, wantLo, wantHi)
	}
}

func TestBrouncker4OverPiApproxFromPrefix_Prefix2FallsBackToGenericLowerBound(t *testing.T) {
	got, err := Brouncker4OverPiApproxFromPrefix(2)
	if err != nil {
		t.Fatalf("Brouncker4OverPiApproxFromPrefix failed: %v", err)
	}

	wantConv := mustRat(3, 2)
	if got.Convergent.Cmp(wantConv) != 0 {
		t.Fatalf("got convergent %v want %v", got.Convergent, wantConv)
	}
	if got.Range == nil {
		t.Fatalf("expected non-nil range")
	}

	// Corrected Brouncker source:
	// x = 1 + 1/(2 + 9/tail), tail >= 1 => [12/11, 3/2]
	wantLo := mustRat(12, 11)
	wantHi := mustRat(3, 2)
	if got.Range.Lo.Cmp(wantLo) != 0 || got.Range.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got range [%v,%v] want [%v,%v]", got.Range.Lo, got.Range.Hi, wantLo, wantHi)
	}
	if !got.Range.Contains(mustRat(15, 13)) {
		t.Fatalf("expected range %v to contain a valid unfinished Brouncker value like 15/13", *got.Range)
	}
}

func TestBrouncker4OverPiApproxFromPrefix_RejectsZeroPrefixTerms(t *testing.T) {
	_, err := Brouncker4OverPiApproxFromPrefix(0)
	if err == nil {
		t.Fatalf("expected error for zero prefixTerms")
	}
}

// brouncker_pi_tail_test.go v2
