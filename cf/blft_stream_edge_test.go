// blft_stream_edge_test.go v1
package cf

import (
	"errors"
	"strings"
	"testing"
)

func TestAnnotateErrBLFT_IncludesOriginalError(t *testing.T) {
	err0 := errors.New("boom")
	tform := NewBLFT(1, 0, 0, 0, 0, 0, 0, 1)
	rx := NewRange(mustRat(1, 1), mustRat(2, 1), true, true)
	ry := NewRange(mustRat(1, 1), mustRat(2, 1), true, true)

	err := annotateErrBLFT(err0, tform, rx, ry)
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected original error text, got %q", err.Error())
	}
}

func TestNewBLFTStream_DefaultsRefineGuardsAndCycleHistory(t *testing.T) {
	tform := NewBLFT(0, 1, 0, 0, 0, 0, 0, 1) // z=x

	s := NewBLFTStream(
		tform,
		NewSliceCF(1),
		NewSliceCF(1),
		BLFTStreamOptions{
			DetectCycles: true,
		},
	)

	if s.maxRefinesPerDigit != -1 {
		t.Fatalf("maxRefinesPerDigit=%d want -1", s.maxRefinesPerDigit)
	}
	if s.maxTotalRefines != -1 {
		t.Fatalf("maxTotalRefines=%d want -1", s.maxTotalRefines)
	}
	if !s.detectCycles {
		t.Fatalf("expected detectCycles=true")
	}
	if s.maxRepeats != 2 {
		t.Fatalf("maxRepeats=%d want 2", s.maxRepeats)
	}
	if s.history == nil {
		t.Fatalf("expected non-nil history")
	}
	if s.history.Cap() < 16 {
		t.Fatalf("history cap=%d want at least 16", s.history.Cap())
	}
}

func TestBLFTStreamNext_EmptyXSourceSetsErr(t *testing.T) {
	tform := NewBLFT(0, 1, 0, 0, 0, 0, 0, 1) // z=x

	s := NewBLFTStream(
		tform,
		NewSliceCF(), // empty X
		NewSliceCF(1),
		BLFTStreamOptions{},
	)

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected termination")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(s.Err().Error(), "empty X source CF") {
		t.Fatalf("unexpected error: %v", s.Err())
	}
	if !s.done {
		t.Fatalf("expected stream done=true")
	}
}

func TestBLFTStreamNext_EmptyYSourceSetsErr(t *testing.T) {
	tform := NewBLFT(0, 1, 0, 0, 0, 0, 0, 1) // z=x

	s := NewBLFTStream(
		tform,
		NewSliceCF(1),
		NewSliceCF(), // empty Y
		BLFTStreamOptions{},
	)

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected termination")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(s.Err().Error(), "empty Y source CF") {
		t.Fatalf("unexpected error: %v", s.Err())
	}
	if !s.done {
		t.Fatalf("expected stream done=true")
	}
}

func TestBLFTStreamEmitDigitBLFT_KnownTransform(t *testing.T) {
	// z = x+y = N/D with:
	// N = 0*xy + 1*x + 1*y + 0
	// D = 0*xy + 0*x + 0*y + 1
	s := NewBLFTStream(
		NewBLFT(0, 1, 1, 0, 0, 0, 0, 1),
		NewSliceCF(1),
		NewSliceCF(1),
		BLFTStreamOptions{},
	)

	got, err := s.emitDigitBLFT(1)
	if err != nil {
		t.Fatalf("emitDigitBLFT failed: %v", err)
	}

	// z' = 1/(z-1)
	// new numerator = old denominator = (0,0,0,1)
	// new denominator = old numerator - 1*old denominator = (0,1,1,-1)
	want := NewBLFT(0, 0, 0, 1, 0, 1, 1, -1)

	if got != want {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestBLFTStreamEmitDigitBLFT_Overflow(t *testing.T) {
	s := NewBLFTStream(
		NewBLFT(0, 0, 0, 0, 0, 0, 0, 0),
		NewSliceCF(1),
		NewSliceCF(1),
		BLFTStreamOptions{},
	)
	// make d * H overflow
	s.t.H = 1 << 62

	_, err := s.emitDigitBLFT(4)
	if err == nil {
		t.Fatalf("expected overflow error")
	}
	if !errors.Is(err, ErrOverflow) {
		t.Fatalf("got %v, want ErrOverflow", err)
	}
}

func TestBLFTStreamNext_HitsRefineLimitAndAnnotates(t *testing.T) {
	// x*y over sqrt(2),sqrt(2) needs refinement; with zero refine budget it should fail fast.
	tform := NewBLFT(1, 0, 0, 0, 0, 0, 0, 1)

	s := NewBLFTStream(
		tform,
		Sqrt2CF(),
		Sqrt2CF(),
		BLFTStreamOptions{
			MaxRefinesPerDigit: 0,
			MaxTotalRefines:    0,
		},
	)
	// Disable the "both zero means unlimited" default by setting one explicitly after construction.
	s.maxRefinesPerDigit = 0
	s.maxTotalRefines = 0

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected termination")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}
	msg := s.Err().Error()
	if !strings.Contains(msg, "exceeded MaxRefinesPerDigit=0") {
		t.Fatalf("unexpected error: %v", s.Err())
	}
	if !strings.Contains(msg, "BLFT[") {
		t.Fatalf("expected annotated fingerprint, got %v", s.Err())
	}
}

func TestBLFTStreamNext_CycleDetectionPath(t *testing.T) {
	// Force the cycle-detection path to execute; exact outcome may be a cycle error or another
	// bounded termination, but history/count path should be exercised.
	tform := NewBLFT(1, 0, 0, 0, 0, 0, 0, 1)

	s := NewBLFTStream(
		tform,
		Sqrt2CF(),
		Sqrt2CF(),
		BLFTStreamOptions{
			DetectCycles:       true,
			MaxRepeats:         1,
			HistorySize:        4,
			MaxRefinesPerDigit: 3,
			MaxTotalRefines:    3,
		},
	)

	_, _ = s.Next()

	if s.history == nil {
		t.Fatalf("expected non-nil history")
	}
	if s.history.Len() == 0 {
		t.Fatalf("expected history to receive at least one fingerprint")
	}
}

// blft_stream_edge_test.go v1
