// gcf_blft_stream_test.go v2
package cf

import (
	"strings"
	"testing"
)

func TestNewGCFBLFTStream_UsesExactTailSources(t *testing.T) {
	// B(x,y)=x+y  => (0xy + 1x + 1y + 0)/(0xy + 0x + 0y + 1)
	base := NewBLFT(0, 1, 1, 0, 0, 0, 0, 1)

	xSrc := NewSliceGCF(
		[2]int64{3, 2},
	)
	xTail := mustRat(11, 1)

	ySrc := NewSliceGCF(
		[2]int64{2, 3},
		[2]int64{4, 5},
	)
	yTail := mustRat(7, 1)

	wantRat, xIngested, yIngested, err := ApplyComposedGCFXYBLFTToTailsExact(
		base,
		NewSliceGCF([2]int64{3, 2}), xTail, 8,
		NewSliceGCF(
			[2]int64{2, 3},
			[2]int64{4, 5},
		), yTail, 8,
	)
	if err != nil {
		t.Fatalf("ApplyComposedGCFXYBLFTToTailsExact failed: %v", err)
	}
	if xIngested != 1 {
		t.Fatalf("got xIngested=%d want 1", xIngested)
	}
	if yIngested != 2 {
		t.Fatalf("got yIngested=%d want 2", yIngested)
	}
	want := collectAll(NewRationalCF(wantRat))

	s := NewGCFBLFTStream(
		base,
		xSrc,
		NewExactTailSource(xTail),
		ySrc,
		NewExactTailSource(yTail),
		GCFBLFTStreamOptions{
			MaxXIngestTerms: 8,
			MaxYIngestTerms: 8,
		},
	)

	got := collectAll(s)
	if err := s.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}
	if !equalSlice(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestNewGCFBLFTStream_MissingXTailEvidenceIsError(t *testing.T) {
	base := NewBLFT(0, 1, 1, 0, 0, 0, 0, 1)

	s := NewGCFBLFTStream(
		base,
		NewSliceGCF([2]int64{3, 2}),
		NoTailSource{},
		NewSliceGCF([2]int64{2, 3}),
		NewExactTailSource(mustRat(5, 1)),
		GCFBLFTStreamOptions{
			MaxXIngestTerms: 8,
			MaxYIngestTerms: 8,
		},
	)

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected failure, not a digit")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(s.Err().Error(), "x tail evidence not implemented") {
		t.Fatalf("unexpected error: %v", s.Err())
	}
}

func TestNewGCFBLFTStream_MissingYTailEvidenceIsError(t *testing.T) {
	base := NewBLFT(0, 1, 1, 0, 0, 0, 0, 1)

	s := NewGCFBLFTStream(
		base,
		NewSliceGCF([2]int64{3, 2}),
		NewExactTailSource(mustRat(5, 1)),
		NewSliceGCF([2]int64{2, 3}),
		NoTailSource{},
		GCFBLFTStreamOptions{
			MaxXIngestTerms: 8,
			MaxYIngestTerms: 8,
		},
	)

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected failure, not a digit")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(s.Err().Error(), "y tail evidence not implemented") {
		t.Fatalf("unexpected error: %v", s.Err())
	}
}

func TestGCFBLFTStreamWithTails_InfiniteSourceRequiresBoundedIngestion(t *testing.T) {
	base := NewBLFT(0, 1, 1, 0, 0, 0, 0, 1)

	s := NewGCFBLFTStreamWithTails(
		base,
		NewUnitPArithmeticQGCFSource(1, 1), mustRat(1, 1),
		NewSliceGCF([2]int64{2, 3}), mustRat(5, 1),
		GCFBLFTStreamOptions{
			MaxXIngestTerms: 3,
			MaxYIngestTerms: 8,
		},
	)

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected failure, not a digit")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(s.Err().Error(), "exceeded MaxIngestTerms=3") {
		t.Fatalf("unexpected error: %v", s.Err())
	}
}

func TestGCFBLFTStreamWithTails_ExhaustedStreamStaysExhausted(t *testing.T) {
	base := NewBLFT(0, 1, 1, 0, 0, 0, 0, 1)

	s := NewGCFBLFTStreamWithTails(
		base,
		NewSliceGCF([2]int64{3, 2}), mustRat(11, 1),
		NewSliceGCF([2]int64{2, 3}), mustRat(7, 1),
		GCFBLFTStreamOptions{
			MaxXIngestTerms: 8,
			MaxYIngestTerms: 8,
		},
	)

	for {
		_, ok := s.Next()
		if !ok {
			break
		}
	}

	if _, ok := s.Next(); ok {
		t.Fatalf("expected exhausted stream to stay exhausted")
	}
	if _, ok := s.Next(); ok {
		t.Fatalf("expected exhausted stream to stay exhausted on repeated calls")
	}
}

func TestGCFBLFTStreamWithTails_SeparateXBoundIsHonored(t *testing.T) {
	base := NewBLFT(0, 1, 1, 0, 0, 0, 0, 1)

	s := NewGCFBLFTStreamWithTails(
		base,
		NewUnitPArithmeticQGCFSource(1, 1), mustRat(1, 1),
		NewSliceGCF([2]int64{2, 3}), mustRat(5, 1),
		GCFBLFTStreamOptions{
			MaxXIngestTerms: 3,
			MaxYIngestTerms: 8,
		},
	)

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected failure, not a digit")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(s.Err().Error(), "exceeded MaxIngestTerms=3") {
		t.Fatalf("unexpected error: %v", s.Err())
	}
}

func TestGCFBLFTStreamWithTails_SeparateYBoundIsHonored(t *testing.T) {
	base := NewBLFT(0, 1, 1, 0, 0, 0, 0, 1)

	s := NewGCFBLFTStreamWithTails(
		base,
		NewSliceGCF([2]int64{2, 3}), mustRat(5, 1),
		NewUnitPArithmeticQGCFSource(1, 1), mustRat(1, 1),
		GCFBLFTStreamOptions{
			MaxXIngestTerms: 8,
			MaxYIngestTerms: 3,
		},
	)

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected failure, not a digit")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(s.Err().Error(), "exceeded MaxIngestTerms=3") {
		t.Fatalf("unexpected error: %v", s.Err())
	}
}

// gcf_blft_stream_test.go v2
