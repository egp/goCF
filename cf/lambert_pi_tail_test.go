// lambert_pi_tail_test.go v4
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

func TestLambertPiOver4TailRangeAfterPrefix_Prefix2(t *testing.T) {
	got, ok, err := LambertPiOver4TailRangeAfterPrefix(2)
	if err != nil {
		t.Fatalf("LambertPiOver4TailRangeAfterPrefix failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	wantLo := mustRat(3, 1)
	wantHi := mustRat(5, 1)
	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v] want [%v,%v]", got.Lo, got.Hi, wantLo, wantHi)
	}
}

func TestLambertPiOver4TailRangeAfterPrefix_Prefix3(t *testing.T) {
	got, ok, err := LambertPiOver4TailRangeAfterPrefix(3)
	if err != nil {
		t.Fatalf("LambertPiOver4TailRangeAfterPrefix failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	wantLo := mustRat(5, 1)
	wantHi := mustRat(34, 5)
	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v] want [%v,%v]", got.Lo, got.Hi, wantLo, wantHi)
	}
}

func TestLambertPiOver4TailRangeAfterPrefix_Prefix4NotYetSpecialized(t *testing.T) {
	_, ok, err := LambertPiOver4TailRangeAfterPrefix(4)
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

	wantLo := mustRat(3, 4)
	wantHi := mustRat(1, 1)
	if got.Range.Lo.Cmp(wantLo) != 0 || got.Range.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got range [%v,%v] want [%v,%v]", got.Range.Lo, got.Range.Hi, wantLo, wantHi)
	}
}

func TestLambertPiOver4ApproxFromPrefix_Prefix2UsesSpecializedTailRange(t *testing.T) {
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

	// After two terms, x = 1 / (1 + 1/tail) with tail in [3, 5],
	// so x = tail/(tail+1) in [3/4, 5/6].
	wantLo := mustRat(3, 4)
	wantHi := mustRat(5, 6)
	if got.Range.Lo.Cmp(wantLo) != 0 || got.Range.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got range [%v,%v] want [%v,%v]", got.Range.Lo, got.Range.Hi, wantLo, wantHi)
	}
}

func TestLambertPiOver4ApproxFromPrefix_Prefix3UsesSpecializedTailRange(t *testing.T) {
	got, err := LambertPiOver4ApproxFromPrefix(3)
	if err != nil {
		t.Fatalf("LambertPiOver4ApproxFromPrefix failed: %v", err)
	}

	wantConv := mustRat(3, 4)
	if got.Convergent.Cmp(wantConv) != 0 {
		t.Fatalf("got convergent %v want %v", got.Convergent, wantConv)
	}
	if got.Range == nil {
		t.Fatalf("expected non-nil range")
	}

	// After three terms:
	//
	//	x = 1 / (1 + 1/(3 + 4/tail))
	//
	// with tail in [5, 34/5].
	//
	// Let y = 3 + 4/tail, so y in [61/17, 19/5].
	// Then x = y/(y+1), giving x in [61/78, 19/24].
	wantLo := mustRat(61, 78)
	wantHi := mustRat(19, 24)
	if got.Range.Lo.Cmp(wantLo) != 0 || got.Range.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got range [%v,%v] want [%v,%v]", got.Range.Lo, got.Range.Hi, wantLo, wantHi)
	}
}

func TestLambertPiOver4ApproxFromPrefix_RejectsZeroPrefixTerms(t *testing.T) {
	_, err := LambertPiOver4ApproxFromPrefix(0)
	if err == nil {
		t.Fatalf("expected error for zero prefixTerms")
	}
}

// lambert_pi_tail_test.go v4
