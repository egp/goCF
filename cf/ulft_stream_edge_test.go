// ulft_stream_edge_test.go v1
package cf

import (
	"errors"
	"math/big"
	"strings"
	"testing"
)

func TestAnnotateErrULFT_IncludesOriginalError(t *testing.T) {
	err0 := errors.New("boom")
	tform := NewULFT(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1))
	r := NewRange(mustRat(1, 1), mustRat(2, 1), true, true)

	err := annotateErrULFT(err0, tform, r)
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected original error text, got %q", err.Error())
	}
}

func TestNewULFTStream_DefaultsRefineGuardsAndCycleSetup(t *testing.T) {
	tform := NewULFT(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1))

	s := NewULFTStream(
		tform,
		NewSliceCF(1),
		ULFTStreamOptions{DetectCycles: true},
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
}

func TestULFTStreamNext_EmptySourceSetsErr(t *testing.T) {
	tform := NewULFT(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1))

	s := NewULFTStream(tform, NewSliceCF(), ULFTStreamOptions{})

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

func TestULFTStreamNext_ExactIntegerTermination(t *testing.T) {
	// Identity transform over exact rational source x = [2].
	// Since x is exact and equal to integer 2, ULFTStream should emit 2
	// and terminate cleanly via the exact-point integer short-circuit.
	tform := NewULFT(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1))

	s := NewULFTStream(tform, NewSliceCF(2), ULFTStreamOptions{})

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

func TestULFTStreamNext_RefineLimitExceeded(t *testing.T) {
	// Identity over sqrt(2): first prefix range [1,2] is not digit-safe for x itself.
	tform := NewULFT(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1))

	s := NewULFTStream(
		tform,
		Sqrt2CF(),
		ULFTStreamOptions{
			MaxRefinesPerDigit: 0,
			MaxTotalRefines:    1,
		},
	)

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
	if !strings.Contains(s.Err().Error(), "ULFT[") {
		t.Fatalf("expected annotated fingerprint, got %v", s.Err())
	}
}

func TestULFTStreamNext_CannotRefineFurther(t *testing.T) {
	// T(x)=1/(2x). For exact x=1, output is 1/2 => [0;2], so first digit 0 is safe,
	// second call should terminate cleanly after exact point handling.
	tform := NewULFT(big.NewInt(0), big.NewInt(1), big.NewInt(2), big.NewInt(0))

	s := NewULFTStream(tform, NewSliceCF(1), ULFTStreamOptions{
		MaxRefinesPerDigit: 5,
		MaxTotalRefines:    5,
	})

	a0, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if a0 != 0 {
		t.Fatalf("got %d want 0", a0)
	}
	a1, ok := s.Next()
	if !ok {
		t.Fatalf("expected second digit, err=%v", s.Err())
	}
	if a1 != 2 {
		t.Fatalf("got %d want 2", a1)
	}
	_, ok = s.Next()
	if ok {
		t.Fatalf("expected termination after [0;2]")
	}
	if err := s.Err(); err != nil {
		t.Fatalf("expected clean termination, got %v", err)
	}
}

func TestULFTStreamNext_ApplyRangeErrorPath(t *testing.T) {
	// Denominator crosses zero on initial range [1,2].
	tform := NewULFT(big.NewInt(1), big.NewInt(0), big.NewInt(1), big.NewInt(-1))

	s := NewULFTStream(tform, NewSliceCF(1), ULFTStreamOptions{})

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected termination")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(s.Err().Error(), "denominator crosses 0") {
		t.Fatalf("unexpected error: %v", s.Err())
	}
}

func TestULFTStreamNext_CycleDetectionHistoryPath(t *testing.T) {
	tform := NewULFT(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1))

	s := NewULFTStream(
		tform,
		Sqrt2CF(),
		ULFTStreamOptions{
			DetectCycles:       true,
			MaxRepeats:         1,
			MaxRefinesPerDigit: 2,
			MaxTotalRefines:    2,
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

func TestUlftFingerprint_CanonicalizesSignAndGCD(t *testing.T) {
	r := NewRange(mustRat(1, 1), mustRat(2, 1), true, true)

	t1 := NewULFT(big.NewInt(2), big.NewInt(4), big.NewInt(6), big.NewInt(8))
	t2 := NewULFT(big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(4))

	k1, err := ulftFingerprint(t1, r)
	if err != nil {
		t.Fatalf("ulftFingerprint t1 failed: %v", err)
	}
	k2, err := ulftFingerprint(t2, r)
	if err != nil {
		t.Fatalf("ulftFingerprint t2 failed: %v", err)
	}

	if k1 != k2 {
		t.Fatalf("canonical fingerprints differ: %+v vs %+v", k1, k2)
	}
}

// ulft_stream_edge_test.go v1
