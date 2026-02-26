// ulft_stream_progress_test.go v2
package cf

import "testing"

func TestULFTStream_ProgressGuardsTrip(t *testing.T) {
	// Identity transform on a rational requires refinement before the first digit is safe.
	// With MaxTotalRefines=0 (strict: no refines allowed), it must fail immediately.
	id := NewULFT(1, 0, 0, 1)
	x := mustRat(355, 113) // [3;7,16]

	s := NewULFTStream(id, NewRationalCF(x), ULFTStreamOptions{
		DetectCycles:       false,
		MaxTotalRefines:    0, // strict: forbid refining
		MaxRefinesPerDigit: -1,
	})

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected ok=false due to progress guard")
	}
	if s.Err() == nil {
		t.Fatalf("expected error due to progress guard, got nil")
	}
}

// ulft_stream_progress_test.go v2
