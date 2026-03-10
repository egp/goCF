// ulft_stream_guard_test.go v2
package cf

import (
	"math/big"
	"testing"
)

func TestULFTStreamNext_EmptySourceIsError(t *testing.T) {
	tform := NewULFT(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1))

	s := NewULFTStream(tform, NewSliceCF(), ULFTStreamOptions{})

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected no digit from empty source")
	}
	if err := s.Err(); err == nil {
		t.Fatalf("expected non-nil error")
	}
}

func TestULFTStreamNext_MaxRefinesPerDigitExceeded(t *testing.T) {
	// Identity transform on sqrt(2). With the initial source range [1,2],
	// no digit is yet safe, so a refine is required immediately.
	// With MaxRefinesPerDigit=0, the first call must fail.
	tform := NewULFT(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1))

	s := NewULFTStream(
		tform,
		Sqrt2CF(),
		ULFTStreamOptions{
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

func TestULFTStreamNext_MaxTotalRefinesExceeded(t *testing.T) {
	// Same setup as above, but total refines are globally forbidden.
	tform := NewULFT(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1))

	s := NewULFTStream(
		tform,
		Sqrt2CF(),
		ULFTStreamOptions{
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

func TestULFTStreamNext_ExactConstantTerminatesCleanly(t *testing.T) {
	// T(x) = 1 for all x, so the stream should emit [1] and terminate cleanly.
	tform := NewULFT(big.NewInt(0), big.NewInt(1), big.NewInt(0), big.NewInt(1))

	s := NewULFTStream(
		tform,
		Sqrt2CF(),
		ULFTStreamOptions{
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

// ulft_stream_guard_test.go v2
