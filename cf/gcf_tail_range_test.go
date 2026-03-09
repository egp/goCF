// gcf_tail_range_test.go v2
package cf

import "testing"

type stubTailRangeGCF struct {
	terms [][2]int64
	i     int
	r     Range
}

func (s *stubTailRangeGCF) NextPQ() (int64, int64, bool) {
	if s.i >= len(s.terms) {
		return 0, 0, false
	}
	t := s.terms[s.i]
	s.i++
	return t[0], t[1], true
}

func (s *stubTailRangeGCF) TailRange() Range { return s.r }

func TestGCFBounder_SetTailRange_RejectsOutside(t *testing.T) {
	b := NewGCFBounder()
	err := b.SetTailRange(NewRange(mustRat(2, 1), mustRat(1, 1), true, true))
	if err == nil {
		t.Fatalf("expected error for outside tail range")
	}
}

func TestGCFBounder_SetTailRange_RejectsNonPositive(t *testing.T) {
	b := NewGCFBounder()
	err := b.SetTailRange(NewRange(mustRat(0, 1), mustRat(2, 1), true, true))
	if err == nil {
		t.Fatalf("expected error for non-positive tail range")
	}
}

func TestGCFBounder_ExplicitTailRangeProducesExpectedRange(t *testing.T) {
	// Prefix terms (1,1), (2,1) define:
	// x = 1 + 1/(2 + 1/tail)
	// with explicit tail in [1,2]
	// tail=1 -> 4/3
	// tail=2 -> 7/5
	b := NewGCFBounder()
	if err := b.IngestPQ(1, 1); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}
	if err := b.IngestPQ(2, 1); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}
	if err := b.SetTailRange(NewRange(mustRat(1, 1), mustRat(2, 1), true, true)); err != nil {
		t.Fatalf("SetTailRange failed: %v", err)
	}

	r, ok, err := b.Range()
	if err != nil {
		t.Fatalf("Range failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	wantLo := mustRat(4, 3)
	wantHi := mustRat(7, 5)
	if r.Lo.Cmp(wantLo) != 0 || r.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v] want [%v,%v]", r.Lo, r.Hi, wantLo, wantHi)
	}
}

func TestGCFBounder_ExplicitTailRangeBeatsLowerBoundInPrecision(t *testing.T) {
	// Explicit tail range [1,2] should be tighter than lower-bound-only [1,+∞).
	bRange := NewGCFBounder()
	if err := bRange.IngestPQ(1, 1); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}
	if err := bRange.IngestPQ(2, 1); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}
	if err := bRange.SetTailRange(NewRange(mustRat(1, 1), mustRat(2, 1), true, true)); err != nil {
		t.Fatalf("SetTailRange failed: %v", err)
	}

	rRange, ok, err := bRange.Range()
	if err != nil {
		t.Fatalf("Range failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	bLower := NewGCFBounder()
	if err := bLower.IngestPQ(1, 1); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}
	if err := bLower.IngestPQ(2, 1); err != nil {
		t.Fatalf("IngestPQ failed: %v", err)
	}
	if err := bLower.SetTailLowerBound(mustRat(1, 1)); err != nil {
		t.Fatalf("SetTailLowerBound failed: %v", err)
	}

	rLower, ok, err := bLower.Range()
	if err != nil {
		t.Fatalf("Range failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	spanRange, err := rRange.Hi.Sub(rRange.Lo)
	if err != nil {
		t.Fatalf("spanRange failed: %v", err)
	}
	spanLower, err := rLower.Hi.Sub(rLower.Lo)
	if err != nil {
		t.Fatalf("spanLower failed: %v", err)
	}

	if spanRange.Cmp(spanLower) >= 0 {
		t.Fatalf("expected explicit tail range to be tighter: explicit=%v lower=%v", spanRange, spanLower)
	}
}

func TestIngestGCFPrefix_DoesNotAutoApplyTailRangeMetadata(t *testing.T) {
	src := &stubTailRangeGCF{
		terms: [][2]int64{
			{1, 1},
			{2, 1},
		},
		r: NewRange(mustRat(1, 1), mustRat(2, 1), true, true),
	}

	b, err := IngestGCFPrefix(src, 2)
	if err != nil {
		t.Fatalf("IngestGCFPrefix failed: %v", err)
	}

	r, ok, err := b.Range()
	if err != nil {
		t.Fatalf("Range failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	// Since tail-range metadata is no longer auto-applied, unfinished range falls
	// back to the current convergent placeholder.
	want := mustRat(3, 2)
	if r.Lo.Cmp(want) != 0 || r.Hi.Cmp(want) != 0 {
		t.Fatalf("got [%v,%v] want placeholder [%v,%v]", r.Lo, r.Hi, want, want)
	}
}

// gcf_tail_range_test.go v2
