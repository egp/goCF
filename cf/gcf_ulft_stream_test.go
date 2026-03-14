// gcf_ulft_stream_test.go v2
package cf

import (
	"strings"
	"testing"
)

// composeGCFIntoULFT is a test helper: starting from base transform t,
// ingest every (p,q) from src so the result is the transformed ULFT in the
// remaining tail variable.
func composeGCFIntoULFT(t ULFT, src GCFSource) (ULFT, error) {
	cur := t
	for {
		p, q, ok := src.NextPQ()
		if !ok {
			return cur, nil
		}
		var err error
		cur, err = cur.IngestGCF(p, q)
		if err != nil {
			return ULFT{}, err
		}
	}
}

func TestGCFULFTStreamWithTail_Identity_FinitePrefixMatchesExactRational(t *testing.T) {
	// x = 1 + 2/(3 + 4/5) = 29/19
	src := NewSliceGCF(
		[2]int64{1, 2},
		[2]int64{3, 4},
	)
	tail := mustRat(5, 1)

	id := NewULFT(mustBig(1), mustBig(0), mustBig(0), mustBig(1))

	// Exact reference:
	composed, err := composeGCFIntoULFT(id, NewSliceGCF(
		[2]int64{1, 2},
		[2]int64{3, 4},
	))
	if err != nil {
		t.Fatalf("composeGCFIntoULFT failed: %v", err)
	}
	wantRat, err := composed.ApplyRat(tail)
	if err != nil {
		t.Fatalf("ApplyRat failed: %v", err)
	}
	want := collectAll(NewRationalCF(wantRat))

	s := NewGCFULFTStreamWithTail(
		id,
		src,
		tail,
		GCFULFTStreamOptions{
			MaxIngestTerms: 8,
		},
	)

	got := collectAll(s)
	if err := s.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}

	if !equalSlice(got, want) {
		t.Fatalf("got %v want %v wantRat=%v", got, want, wantRat)
	}
}

func TestGCFULFTStreamWithTail_GeneralULFT_FinitePrefixMatchesExactRational(t *testing.T) {
	// x = 2 + 3/(4 + 5/6)
	src := NewSliceGCF(
		[2]int64{2, 3},
		[2]int64{4, 5},
	)
	tail := mustRat(6, 1)

	u := NewULFT(mustBig(2), mustBig(1), mustBig(3), mustBig(4)) // (2x+1)/(3x+4)

	// Exact reference:
	composed, err := composeGCFIntoULFT(u, NewSliceGCF(
		[2]int64{2, 3},
		[2]int64{4, 5},
	))
	if err != nil {
		t.Fatalf("composeGCFIntoULFT failed: %v", err)
	}
	wantRat, err := composed.ApplyRat(tail)
	if err != nil {
		t.Fatalf("ApplyRat failed: %v", err)
	}
	want := collectAll(NewRationalCF(wantRat))

	s := NewGCFULFTStreamWithTail(
		u,
		src,
		tail,
		GCFULFTStreamOptions{
			MaxIngestTerms: 8,
		},
	)

	got := collectAll(s)
	if err := s.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}

	if !equalSlice(got, want) {
		t.Fatalf("got %v want %v wantRat=%v", got, want, wantRat)
	}
}

func TestGCFULFTStreamWithTail_ExhaustedStreamStaysExhausted(t *testing.T) {
	src := NewSliceGCF(
		[2]int64{1, 1},
	)
	tail := mustRat(2, 1)
	id := NewULFT(mustBig(1), mustBig(0), mustBig(0), mustBig(1))

	s := NewGCFULFTStreamWithTail(
		id,
		src,
		tail,
		GCFULFTStreamOptions{
			MaxIngestTerms: 8,
		},
	)

	// Drain completely.
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

func TestGCFULFTStreamWithTail_InfiniteSourceRequiresBoundedIngestion(t *testing.T) {
	id := NewULFT(mustBig(1), mustBig(0), mustBig(0), mustBig(1))

	s := NewGCFULFTStreamWithTail(
		id,
		NewUnitPArithmeticQGCFSource(1, 1), // infinite source
		mustRat(1, 1),
		GCFULFTStreamOptions{
			MaxIngestTerms: 3,
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

func TestGCFULFTStreamWithTail_ZeroMaxIngestTermsFailsImmediately(t *testing.T) {
	id := NewULFT(mustBig(1), mustBig(0), mustBig(0), mustBig(1))

	s := NewGCFULFTStreamWithTail(
		id,
		NewSliceGCF([2]int64{1, 1}),
		mustRat(2, 1),
		GCFULFTStreamOptions{
			MaxIngestTerms: 0,
		},
	)

	_, ok := s.Next()
	if ok {
		t.Fatalf("expected failure, not a digit")
	}
	if s.Err() == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(s.Err().Error(), "exceeded MaxIngestTerms=0") {
		t.Fatalf("unexpected error: %v", s.Err())
	}
}

// gcf_ulft_stream_test.go v2
