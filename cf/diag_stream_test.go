// diag_stream_test.go v2
package cf

import (
	"math/big"
	"testing"
)

func TestPeriodicCF_Radicand_Metadata(t *testing.T) {
	cf, ok := Sqrt2CF().(QuadraticRadicalSource)
	if !ok {
		t.Fatalf("Sqrt2CF should implement QuadraticRadicalSource")
	}
	n, ok := cf.Radicand()
	if !ok || n != 2 {
		t.Fatalf("got (%d,%v), want (2,true)", n, ok)
	}

	phi, ok := PhiCF().(QuadraticRadicalSource)
	if !ok {
		t.Fatalf("PhiCF should still be a PeriodicCF and implement the interface")
	}
	_, has := phi.Radicand()
	if has {
		t.Fatalf("PhiCF should not advertise a radicand")
	}
}

func TestDiagBLFTStream_ExactSquareShortcut_Sqrt2(t *testing.T) {
	sq := NewDiagBLFT(
		big.NewInt(1), big.NewInt(0), big.NewInt(0),
		big.NewInt(0), big.NewInt(0), big.NewInt(1),
	)

	s := NewDiagBLFTStream(sq, Sqrt2CF(), DiagBLFTStreamOptions{})

	a0, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, got termination; err=%v", s.Err())
	}
	if a0 != 2 {
		t.Fatalf("got %d, want 2", a0)
	}

	if _, ok := s.Next(); ok {
		t.Fatalf("expected clean termination after [2]")
	}
	if err := s.Err(); err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
}

func TestDiagBLFTStream_ExactSquareShortcut_Sqrt7(t *testing.T) {
	sq := NewDiagBLFT(
		big.NewInt(1), big.NewInt(0), big.NewInt(0),
		big.NewInt(0), big.NewInt(0), big.NewInt(1),
	)

	s := NewDiagBLFTStream(sq, Sqrt7CF(), DiagBLFTStreamOptions{})

	a0, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, got termination; err=%v", s.Err())
	}
	if a0 != 7 {
		t.Fatalf("got %d, want 7", a0)
	}

	if _, ok := s.Next(); ok {
		t.Fatalf("expected clean termination after [7]")
	}
	if err := s.Err(); err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
}

func TestDiagBLFTStream_ExactSquarePlusConstShortcut_Sqrt2Plus1(t *testing.T) {
	tform := NewDiagBLFT(
		big.NewInt(1), big.NewInt(0), big.NewInt(1),
		big.NewInt(0), big.NewInt(0), big.NewInt(1),
	)

	s := NewDiagBLFTStream(tform, Sqrt2CF(), DiagBLFTStreamOptions{})

	a0, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, got termination; err=%v", s.Err())
	}
	if a0 != 3 {
		t.Fatalf("got %d, want 3", a0)
	}

	if _, ok := s.Next(); ok {
		t.Fatalf("expected clean termination after [3]")
	}
	if err := s.Err(); err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
}

func TestDiagBLFTStream_ExactSquarePlusConstShortcut_Sqrt5Minus1(t *testing.T) {
	tform := NewDiagBLFT(
		big.NewInt(1), big.NewInt(0), big.NewInt(-1),
		big.NewInt(0), big.NewInt(0), big.NewInt(1),
	)

	s := NewDiagBLFTStream(tform, Sqrt5CF(), DiagBLFTStreamOptions{})

	a0, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, got termination; err=%v", s.Err())
	}
	if a0 != 4 {
		t.Fatalf("got %d, want 4", a0)
	}

	if _, ok := s.Next(); ok {
		t.Fatalf("expected clean termination after [4]")
	}
	if err := s.Err(); err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
}

func TestDiagBLFTStream_NoShortcut_ForNonSquareTransform(t *testing.T) {
	// 2*x^2 + 1, not exactly x^2 + k, so no algebraic shortcut should fire.
	tform := NewDiagBLFT(
		big.NewInt(2), big.NewInt(0), big.NewInt(1),
		big.NewInt(0), big.NewInt(0), big.NewInt(1),
	)

	s := NewDiagBLFTStream(tform, Sqrt2CF(), DiagBLFTStreamOptions{
		MaxRefinesPerDigit: 2,
		MaxTotalRefines:    10,
	})

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected bounded refine failure, got a digit")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil err")
	}
}

// diag_stream_test.go v2
