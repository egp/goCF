// blft_stream_guard_test.go v1
package cf

import "testing"

func TestBLFTStreamNext_EmptyXSourceIsError(t *testing.T) {
	bl := NewBLFT(0, 0, 0, 1, 0, 0, 0, 1) // constant 1

	s := NewBLFTStream(
		bl,
		NewSliceCF(),
		NewSliceCF(1),
		BLFTStreamOptions{},
	)

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected no digit from empty X source")
	}
	if err := s.Err(); err == nil {
		t.Fatalf("expected non-nil error")
	}
}

func TestBLFTStreamNext_EmptyYSourceIsError(t *testing.T) {
	bl := NewBLFT(0, 0, 0, 1, 0, 0, 0, 1) // constant 1

	s := NewBLFTStream(
		bl,
		NewSliceCF(1),
		NewSliceCF(),
		BLFTStreamOptions{},
	)

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected no digit from empty Y source")
	}
	if err := s.Err(); err == nil {
		t.Fatalf("expected non-nil error")
	}
}

func TestBLFTStreamNext_MaxRefinesPerDigitExceeded(t *testing.T) {
	// z = x + y. With one ingested digit from each sqrt(2)-style source, both ranges
	// are [1,2], so z is in [2,4] and no digit is yet safe. Refinement is required
	// immediately, so MaxRefinesPerDigit=0 must fail on the first call.
	bl := NewBLFT(0, 1, 1, 0, 0, 0, 0, 1)

	s := NewBLFTStream(
		bl,
		Sqrt2CF(),
		Sqrt2CF(),
		BLFTStreamOptions{
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

func TestBLFTStreamNext_MaxTotalRefinesExceeded(t *testing.T) {
	// Same setup as above, but total refines are globally forbidden.
	bl := NewBLFT(0, 1, 1, 0, 0, 0, 0, 1)

	s := NewBLFTStream(
		bl,
		Sqrt2CF(),
		Sqrt2CF(),
		BLFTStreamOptions{
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

func TestBLFTStreamNext_ExactConstantTerminatesCleanly(t *testing.T) {
	// Constant z = 1 for all x,y, so the stream should emit [1] and terminate cleanly.
	bl := NewBLFT(0, 0, 0, 1, 0, 0, 0, 1)

	s := NewBLFTStream(
		bl,
		Sqrt2CF(),
		Sqrt2CF(),
		BLFTStreamOptions{
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

// blft_stream_guard_test.go v1
