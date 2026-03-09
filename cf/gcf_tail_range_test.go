// gcf_tail_range_test.go v1
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

func TestIngestGCFPrefix_PrefersTailRangeMetadata(t *testing.T) {
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

	// x = 1 + 1/(2 + 1/tail), tail in [1,2]
	// tail=1 -> 4/3
	// tail=2 -> 7/5
	wantLo := mustRat(4, 3)
	wantHi := mustRat(7, 5)
	if r.Lo.Cmp(wantLo) != 0 || r.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v] want [%v,%v]", r.Lo, r.Hi, wantLo, wantHi)
	}
}

func TestTailRangeBeatsLowerBoundInPrecision(t *testing.T) {
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

	got, ok, err := b.Range()
	if err != nil {
		t.Fatalf("Range failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected ok=true")
	}

	// Compare with the old ray-based lower-bound enclosure [4/3, 3/2].
	oldHi := mustRat(3, 2)
	if got.Hi.Cmp(oldHi) >= 0 {
		t.Fatalf("expected tighter upper bound than %v, got %v", oldHi, got.Hi)
	}
}
