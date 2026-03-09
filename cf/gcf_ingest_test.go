// gcf_ingest_test.go v1
package cf

import (
	"math/big"
	"testing"
)

func TestULFTIngestGCF_Identity(t *testing.T) {
	// T(x)=x, ingest x = 3 + 2/x'
	t0 := NewULFT(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1))

	got, err := t0.IngestGCF(3, 2)
	if err != nil {
		t.Fatalf("IngestGCF failed: %v", err)
	}

	// Expected: (1*3+0)x' + 1*2  over  (0*3+1)x' + 0*2
	//         = (3x' + 2) / x'
	want := NewULFT(big.NewInt(3), big.NewInt(2), big.NewInt(1), big.NewInt(0))

	if got.A.Cmp(want.A) != 0 || got.B.Cmp(want.B) != 0 || got.C.Cmp(want.C) != 0 || got.D.Cmp(want.D) != 0 {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestULFTIngestGCF_ComposesCorrectlyByEvaluation(t *testing.T) {
	// T(x) = (2x+5)/(7x+11)
	t0 := NewULFT(big.NewInt(2), big.NewInt(5), big.NewInt(7), big.NewInt(11))

	got, err := t0.IngestGCF(3, 2)
	if err != nil {
		t.Fatalf("IngestGCF failed: %v", err)
	}

	// Check equality at x'=5, so x = 3 + 2/5 = 17/5.
	xp := mustRat(5, 1)
	x := mustRat(17, 5)

	y0, err := t0.ApplyRat(x)
	if err != nil {
		t.Fatalf("ApplyRat original failed: %v", err)
	}
	y1, err := got.ApplyRat(xp)
	if err != nil {
		t.Fatalf("ApplyRat rewritten failed: %v", err)
	}

	if y0.Cmp(y1) != 0 {
		t.Fatalf("composition mismatch: original=%v rewritten=%v", y0, y1)
	}
}

func TestULFTIngestGCF_RejectsNonPositiveQ(t *testing.T) {
	t0 := NewULFT(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1))

	if _, err := t0.IngestGCF(3, 0); err == nil {
		t.Fatalf("expected error for q=0")
	}
	if _, err := t0.IngestGCF(3, -1); err == nil {
		t.Fatalf("expected error for q<0")
	}
}

func TestDiagBLFTIngestGCF_IdentitySquare(t *testing.T) {
	// T(x)=x^2
	t0 := NewDiagBLFT(
		big.NewInt(1), big.NewInt(0), big.NewInt(0),
		big.NewInt(0), big.NewInt(0), big.NewInt(1),
	)

	got, err := t0.IngestGCF(3, 2)
	if err != nil {
		t.Fatalf("IngestGCF failed: %v", err)
	}

	// x = (3x'+2)/x'
	// x^2 = (9x'^2 + 12x' + 4)/x'^2
	want := NewDiagBLFT(
		big.NewInt(9), big.NewInt(12), big.NewInt(4),
		big.NewInt(1), big.NewInt(0), big.NewInt(0),
	)

	if got.A.Cmp(want.A) != 0 || got.B.Cmp(want.B) != 0 || got.C.Cmp(want.C) != 0 ||
		got.D.Cmp(want.D) != 0 || got.E.Cmp(want.E) != 0 || got.F.Cmp(want.F) != 0 {
		t.Fatalf("got=%v want=%v", got, want)
	}
}

func TestDiagBLFTIngestGCF_ComposesCorrectlyByEvaluation(t *testing.T) {
	// T(x) = (2x^2 + 3x + 5)/(7x^2 + 11x + 13)
	t0 := NewDiagBLFT(
		big.NewInt(2), big.NewInt(3), big.NewInt(5),
		big.NewInt(7), big.NewInt(11), big.NewInt(13),
	)

	got, err := t0.IngestGCF(3, 2)
	if err != nil {
		t.Fatalf("IngestGCF failed: %v", err)
	}

	xp := mustRat(5, 1)
	x := mustRat(17, 5)

	y0, err := t0.ApplyRat(x)
	if err != nil {
		t.Fatalf("ApplyRat original failed: %v", err)
	}
	y1, err := got.ApplyRat(xp)
	if err != nil {
		t.Fatalf("ApplyRat rewritten failed: %v", err)
	}

	if y0.Cmp(y1) != 0 {
		t.Fatalf("composition mismatch: original=%v rewritten=%v", y0, y1)
	}
}

func TestDiagBLFTIngestGCF_RejectsNonPositiveQ(t *testing.T) {
	t0 := NewDiagBLFT(
		big.NewInt(1), big.NewInt(0), big.NewInt(0),
		big.NewInt(0), big.NewInt(0), big.NewInt(1),
	)

	if _, err := t0.IngestGCF(3, 0); err == nil {
		t.Fatalf("expected error for q=0")
	}
	if _, err := t0.IngestGCF(3, -1); err == nil {
		t.Fatalf("expected error for q<0")
	}
}

// gcf_ingest_test.go v1
