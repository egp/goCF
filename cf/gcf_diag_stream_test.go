// gcf_diag_stream_test.go v1
package cf

import (
	"strings"
	"testing"
)

func TestNewGCFDiagStream_UsesExactTailSource(t *testing.T) {
	// D(x)=x
	base := NewDiagBLFT(
		mustBig(0), mustBig(1), mustBig(0),
		mustBig(0), mustBig(0), mustBig(1),
	)

	src := NewSliceGCF(
		[2]int64{3, 2},
		[2]int64{5, 7},
	)
	tailSrc := NewExactTailSource(mustRat(11, 1))

	wantRat, ingested, err := ApplyComposedGCFDiagBLFTToTailExact(
		base,
		NewSliceGCF(
			[2]int64{3, 2},
			[2]int64{5, 7},
		),
		mustRat(11, 1),
		8,
	)
	if err != nil {
		t.Fatalf("ApplyComposedGCFDiagBLFTToTailExact failed: %v", err)
	}
	if ingested != 2 {
		t.Fatalf("got ingested=%d want 2", ingested)
	}
	want := collectAll(NewRationalCF(wantRat))

	s := NewGCFDiagStream(
		base,
		src,
		tailSrc,
		GCFULFTStreamOptions{MaxIngestTerms: 8},
	)

	got := collectAll(s)
	if err := s.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}
	if !equalSlice(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestNewGCFDiagStream_MissingTailEvidenceIsError(t *testing.T) {
	base := NewDiagBLFT(
		mustBig(0), mustBig(1), mustBig(0),
		mustBig(0), mustBig(0), mustBig(1),
	)

	s := NewGCFDiagStream(
		base,
		NewSliceGCF([2]int64{3, 2}),
		NoTailSource{},
		GCFULFTStreamOptions{MaxIngestTerms: 8},
	)

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected failure, not a digit")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(s.Err().Error(), "tail evidence not implemented") {
		t.Fatalf("unexpected error: %v", s.Err())
	}
}

func TestGCFDiagStreamWithTail_InfiniteSourceRequiresBoundedIngestion(t *testing.T) {
	base := NewDiagBLFT(
		mustBig(0), mustBig(1), mustBig(0),
		mustBig(0), mustBig(0), mustBig(1),
	)

	s := NewGCFDiagStreamWithTail(
		base,
		NewUnitPArithmeticQGCFSource(1, 1),
		mustRat(1, 1),
		GCFULFTStreamOptions{MaxIngestTerms: 3},
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

func TestGCFDiagStreamWithTail_ExhaustedStreamStaysExhausted(t *testing.T) {
	base := NewDiagBLFT(
		mustBig(0), mustBig(1), mustBig(0),
		mustBig(0), mustBig(0), mustBig(1),
	)

	s := NewGCFDiagStreamWithTail(
		base,
		NewSliceGCF([2]int64{3, 2}),
		mustRat(11, 1),
		GCFULFTStreamOptions{MaxIngestTerms: 8},
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

// gcf_diag_stream_test.go v1
