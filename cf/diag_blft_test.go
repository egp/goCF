// diag_blft_test.go v2
package cf

import (
	"testing"
)

// func bi(n int64) *big.Int { return big.NewInt(n) }

func equalDiagBLFT(a, b DiagBLFT) bool {
	return a.A.Cmp(b.A) == 0 &&
		a.B.Cmp(b.B) == 0 &&
		a.C.Cmp(b.C) == 0 &&
		a.D.Cmp(b.D) == 0 &&
		a.E.Cmp(b.E) == 0 &&
		a.F.Cmp(b.F) == 0
}

func TestDiagFromBLFT_MultiplySpecializesToSquare(t *testing.T) {
	mul := NewBLFT(1, 0, 0, 0, 0, 0, 0, 1)

	got := DiagFromBLFT(mul)
	want := NewDiagBLFT(bi(1), bi(0), bi(0), bi(0), bi(0), bi(1))

	if !equalDiagBLFT(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestDiagBLFT_ApplyRat_SquareExactRational(t *testing.T) {
	tform := NewDiagBLFT(bi(1), bi(0), bi(0), bi(0), bi(0), bi(1))

	x := mustRat(7, 5)
	got, err := tform.ApplyRat(x)
	if err != nil {
		t.Fatalf("ApplyRat failed: %v", err)
	}

	want := mustRat(49, 25)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestDiagBLFT_ApplyRange_SquarePositiveInterval(t *testing.T) {
	tform := NewDiagBLFT(bi(1), bi(0), bi(0), bi(0), bi(0), bi(1))

	r := NewRange(mustRat(3, 2), mustRat(2, 1), true, true)
	got, err := tform.ApplyRange(r)
	if err != nil {
		t.Fatalf("ApplyRange failed: %v", err)
	}

	wantLo := mustRat(9, 4)
	wantHi := mustRat(4, 1)

	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v], want [%v,%v]", got.Lo, got.Hi, wantLo, wantHi)
	}
}

func TestDiagBLFT_ApplyRange_UsesVertex(t *testing.T) {
	// x^2 - 2x on [0,2] has vertex at x=1 with minimum -1.
	tform := NewDiagBLFT(bi(1), bi(-2), bi(0), bi(0), bi(0), bi(1))

	r := NewRange(mustRat(0, 1), mustRat(2, 1), true, true)
	got, err := tform.ApplyRange(r)
	if err != nil {
		t.Fatalf("ApplyRange failed: %v", err)
	}

	wantLo := mustRat(-1, 1)
	wantHi := mustRat(0, 1)

	if got.Lo.Cmp(wantLo) != 0 || got.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v], want [%v,%v]", got.Lo, got.Hi, wantLo, wantHi)
	}
}

func TestDiagBLFT_ApplyRange_RejectsNonConstantDenominator(t *testing.T) {
	tform := NewDiagBLFT(bi(1), bi(0), bi(0), bi(0), bi(1), bi(1)) // x^2 / (x + 1)

	r := NewRange(mustRat(1, 1), mustRat(2, 1), true, true)
	_, err := tform.ApplyRange(r)
	if err == nil {
		t.Fatalf("expected rejection for non-constant denominator")
	}
}

func TestDiagBLFT_EmitDigitDiag_SquareAfterEmitting2(t *testing.T) {
	// z = x^2. After emitting 2, z' = 1/(x^2 - 2).
	tform := NewDiagBLFT(bi(1), bi(0), bi(0), bi(0), bi(0), bi(1))

	got, err := tform.emitDigitDiag(2)
	if err != nil {
		t.Fatalf("emitDigitDiag failed: %v", err)
	}

	want := NewDiagBLFT(bi(0), bi(0), bi(1), bi(1), bi(0), bi(-2))
	if !equalDiagBLFT(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// diag_blft_test.go v2
