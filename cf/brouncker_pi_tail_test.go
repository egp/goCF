// brouncker_pi_tail_test.go v3
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

func TestBrouncker4OverPiTailRangeAfterPrefix_Prefix2(t *testing.T) {
	got, ok, err := Brouncker4OverPiTailRangeAfterPrefix(2)
	if err != nil {
		t.Fatalf("Brouncker4OverPiTailRangeAfterPrefix failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	wantLo := mustRat(2, 1)
	wantHi := mustRat(29, 2)
	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v] want [%v,%v]", got.Lo, got.Hi, wantLo, wantHi)
	}
}

func TestBrouncker4OverPiTailRangeAfterPrefix_Prefix4NotYetSpecialized(t *testing.T) {
	_, ok, err := Brouncker4OverPiTailRangeAfterPrefix(4)
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

func TestBrouncker4OverPiApproxFromPrefix_Prefix2UsesSpecializedTailRange(t *testing.T) {
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

	// Prefix 2 means x = 1 + 1/(2 + 9/tail), with tail in [2, 29/2].
	// Then 2 + 9/tail in [2 + 18/29, 13/2] = [76/29, 13/2],
	// so x in [1 + 2/13, 1 + 29/76] = [15/13, 105/76].
	wantLo := mustRat(15, 13)
	wantHi := mustRat(105, 76)
	if got.Range.Lo.Cmp(wantLo) != 0 || got.Range.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got range [%v,%v] want [%v,%v]", got.Range.Lo, got.Range.Hi, wantLo, wantHi)
	}
}

func TestBrouncker4OverPiApproxFromPrefix_RejectsZeroPrefixTerms(t *testing.T) {
	_, err := Brouncker4OverPiApproxFromPrefix(0)
	if err == nil {
		t.Fatalf("expected error for zero prefixTerms")
	}
}

func TestBrouncker4OverPiTailRangeAfterPrefix_Prefix3(t *testing.T) {
	got, ok, err := Brouncker4OverPiTailRangeAfterPrefix(3)
	if err != nil {
		t.Fatalf("Brouncker4OverPiTailRangeAfterPrefix failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	wantLo := mustRat(2, 1)
	wantHi := mustRat(53, 2)
	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v] want [%v,%v]", got.Lo, got.Hi, wantLo, wantHi)
	}
}

func TestBrouncker4OverPiApproxFromPrefix_Prefix3UsesSpecializedTailRange(t *testing.T) {
	got, err := Brouncker4OverPiApproxFromPrefix(3)
	if err != nil {
		t.Fatalf("Brouncker4OverPiApproxFromPrefix failed: %v", err)
	}

	wantConv := mustRat(15, 13)
	if got.Convergent.Cmp(wantConv) != 0 {
		t.Fatalf("got convergent %v want %v", got.Convergent, wantConv)
	}
	if got.Range == nil {
		t.Fatalf("expected non-nil range")
	}

	// Prefix 3 means x = 1 + 1/(2 + 9/(2 + 25/tail)), with tail in [2, 53/2].
	// This yields x in [105/76, 987/710].
	wantLo := mustRat(315, 263)
	wantHi := mustRat(105, 76)
	if got.Range.Lo.Cmp(wantLo) != 0 || got.Range.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got range [%v,%v] want [%v,%v]", got.Range.Lo, got.Range.Hi, wantLo, wantHi)
	}
}

// brouncker_pi_tail_test.go v3
