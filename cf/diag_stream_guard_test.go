// diag_stream_guard_test.go v2
package cf

import (
	"math/big"
	"testing"
)

func TestDiagBLFTStreamNext_EmptySourceIsError(t *testing.T) {
	tform := NewDiagBLFT(
		big.NewInt(1), big.NewInt(0), big.NewInt(0),
		big.NewInt(0), big.NewInt(0), big.NewInt(1),
	)

	s := NewDiagBLFTStream(tform, NewSliceCF(), DiagBLFTStreamOptions{})

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected no digit from empty source")
	}
	if err := s.Err(); err == nil {
		t.Fatalf("expected non-nil error")
	}
}

func TestDiagBLFTStreamNext_MaxRefinesPerDigitExceeded(t *testing.T) {
	// Non-constant denominator is not supported on non-point ranges by ApplyRange.
	// For an uncertain source range, the stream must refine immediately.
	// With MaxRefinesPerDigit=0, the first call must fail.
	tform := NewDiagBLFT(
		big.NewInt(1), big.NewInt(0), big.NewInt(0),
		big.NewInt(0), big.NewInt(1), big.NewInt(1), // denominator x + 1
	)

	s := NewDiagBLFTStream(
		tform,
		Sqrt2CF(),
		DiagBLFTStreamOptions{
			MaxRefinesPerDigit: 0,
			MaxTotalRefines:    -1,
		},
	)

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected refinement-guard failure on first digit")
	}
	if err := s.Err(); err == nil {
		t.Fatalf("expected non-nil error")
	}
}

func TestDiagBLFTStreamNext_MaxTotalRefinesExceeded(t *testing.T) {
	// Same setup as above, but total refines are globally forbidden.
	tform := NewDiagBLFT(
		big.NewInt(1), big.NewInt(0), big.NewInt(0),
		big.NewInt(0), big.NewInt(1), big.NewInt(1), // denominator x + 1
	)

	s := NewDiagBLFTStream(
		tform,
		Sqrt2CF(),
		DiagBLFTStreamOptions{
			MaxRefinesPerDigit: -1,
			MaxTotalRefines:    0,
		},
	)

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected total-refine-guard failure on first digit")
	}
	if err := s.Err(); err == nil {
		t.Fatalf("expected non-nil error")
	}
}

func TestDiagBLFTStreamNext_ExactConstantTerminatesCleanly(t *testing.T) {
	// Constant z = 1 for all x, so the stream should emit [1] and terminate cleanly.
	tform := NewDiagBLFT(
		big.NewInt(0), big.NewInt(0), big.NewInt(1),
		big.NewInt(0), big.NewInt(0), big.NewInt(1),
	)

	s := NewDiagBLFTStream(
		tform,
		Sqrt2CF(),
		DiagBLFTStreamOptions{
			MaxRefinesPerDigit: -1,
			MaxTotalRefines:    -1,
		},
	)

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if d != 1 {
		t.Fatalf("got first digit %d want 1", d)
	}

	_, ok = s.Next()
	if ok {
		t.Fatalf("expected termination")
	}
	if err := s.Err(); err != nil {
		t.Fatalf("expected clean termination, got err=%v", err)
	}
}

// diag_stream_guard_test.go v2
