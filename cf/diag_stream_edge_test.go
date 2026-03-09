// diag_stream_edge_test.go v1
package cf

import (
	"math/big"
	"strings"
	"testing"
)

func TestNewDiagBLFTStream_DefaultsRefineGuards(t *testing.T) {
	tform := NewDiagBLFT(
		big.NewInt(1), big.NewInt(0), big.NewInt(0),
		big.NewInt(0), big.NewInt(0), big.NewInt(1),
	)

	s := NewDiagBLFTStream(tform, NewSliceCF(1), DiagBLFTStreamOptions{})

	if s.maxRefinesPerDigit != -1 {
		t.Fatalf("maxRefinesPerDigit=%d want -1", s.maxRefinesPerDigit)
	}
	if s.maxTotalRefines != -1 {
		t.Fatalf("maxTotalRefines=%d want -1", s.maxTotalRefines)
	}
}

func TestDiagBLFTStreamNext_EmptySourceSetsErr(t *testing.T) {
	tform := NewDiagBLFT(
		big.NewInt(1), big.NewInt(0), big.NewInt(0),
		big.NewInt(0), big.NewInt(0), big.NewInt(1),
	)

	s := NewDiagBLFTStream(tform, NewSliceCF(), DiagBLFTStreamOptions{})

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected termination")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(s.Err().Error(), "empty source CF") {
		t.Fatalf("unexpected error: %v", s.Err())
	}
	if !s.done {
		t.Fatalf("expected done=true")
	}
}

func TestDiagBLFTStreamNext_ExactIntegerTermination(t *testing.T) {
	// Constant transform z(x)=2 exactly.
	tform := NewDiagBLFT(
		big.NewInt(0), big.NewInt(0), big.NewInt(2),
		big.NewInt(0), big.NewInt(0), big.NewInt(1),
	)

	s := NewDiagBLFTStream(tform, NewSliceCF(1), DiagBLFTStreamOptions{})

	a0, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if a0 != 2 {
		t.Fatalf("got %d want 2", a0)
	}

	_, ok = s.Next()
	if ok {
		t.Fatalf("expected termination after exact integer")
	}
	if err := s.Err(); err != nil {
		t.Fatalf("expected clean termination, got %v", err)
	}
}

func TestDiagBLFTStreamNext_RefineLimitExceeded(t *testing.T) {
	// x^2 over sqrt(2) does not become safe from plain interval refinement here;
	// with zero refine budget it should fail immediately after first unsafe digit check.
	tform := NewDiagBLFT(
		big.NewInt(1), big.NewInt(0), big.NewInt(0),
		big.NewInt(0), big.NewInt(0), big.NewInt(1),
	)

	s := NewDiagBLFTStream(tform, NewSliceCF(1), DiagBLFTStreamOptions{
		MaxRefinesPerDigit: 0,
		MaxTotalRefines:    1,
	})

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected termination")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(s.Err().Error(), "exceeded MaxRefinesPerDigit=0") {
		t.Fatalf("unexpected error: %v", s.Err())
	}
}

func TestDiagBLFTStreamNext_FiniteSourceCollapsesToExactPoint(t *testing.T) {
	// x^2 over finite source [1] means x=1 exactly, so output is exactly [1].
	tform := NewDiagBLFT(
		big.NewInt(1), big.NewInt(0), big.NewInt(0),
		big.NewInt(0), big.NewInt(0), big.NewInt(1),
	)

	s := NewDiagBLFTStream(tform, NewSliceCF(1), DiagBLFTStreamOptions{
		MaxRefinesPerDigit: 5,
		MaxTotalRefines:    5,
	})

	a0, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if a0 != 1 {
		t.Fatalf("got %d want 1", a0)
	}

	_, ok = s.Next()
	if ok {
		t.Fatalf("expected termination after exact [1]")
	}
	if err := s.Err(); err != nil {
		t.Fatalf("expected clean termination, got %v", err)
	}
}

func TestDiagBLFTStreamNext_ApplyRangeErrorPath(t *testing.T) {
	// Non-constant denominator is not supported by ApplyRange for non-point ranges.
	tform := NewDiagBLFT(
		big.NewInt(1), big.NewInt(0), big.NewInt(0),
		big.NewInt(0), big.NewInt(1), big.NewInt(1), // denominator x + 1
	)

	s := NewDiagBLFTStream(tform, NewSliceCF(1), DiagBLFTStreamOptions{})

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected termination")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(s.Err().Error(), "non-constant denominator not yet supported") {
		t.Fatalf("unexpected error: %v", s.Err())
	}
}

// diag_stream_edge_test.go v1
