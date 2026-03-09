// brouncker_pi_gcf_test.go v1
package cf

import "testing"

func TestBrouncker4OverPiGCFSource_FirstTerms(t *testing.T) {
	s := NewBrouncker4OverPiGCFSource()

	want := [][2]int64{
		{1, 1},
		{2, 1},
		{2, 9},
		{2, 25},
		{2, 49},
		{2, 81},
	}

	for i, w := range want {
		p, q, ok := s.NextPQ()
		if !ok {
			t.Fatalf("expected term %d", i)
		}
		if p != w[0] || q != w[1] {
			t.Fatalf("term %d: got (%d,%d), want (%d,%d)", i, p, q, w[0], w[1])
		}
	}
}

func TestBrouncker4OverPiGCFSource_Prefix2(t *testing.T) {
	s := NewBrouncker4OverPiGCFSource()

	got, err := GCFApproxFromPrefix(s, 2)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}

	// 1 + 1/2 = 3/2
	want := mustRat(3, 2)
	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

func TestBrouncker4OverPiGCFSource_Prefix3(t *testing.T) {
	s := NewBrouncker4OverPiGCFSource()

	got, err := GCFApproxFromPrefix(s, 3)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}

	// 1 + 1/(2 + 1/2) ??? No.
	// Prefix terms are (1,1), (2,1), (2,9):
	// finite convention => 1 + 1/(2 + 1/2) = 7/5
	want := mustRat(7, 5)
	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

func TestBrouncker4OverPiGCFSource_Prefix4(t *testing.T) {
	s := NewBrouncker4OverPiGCFSource()

	got, err := GCFApproxFromPrefix(s, 4)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}

	// Terms: (1,1), (2,1), (2,9), (2,25)
	// finite convention:
	// v3 = 2
	// v2 = 2 + 9/2 = 13/2
	// v1 = 2 + 1/(13/2) = 28/13
	// v0 = 1 + 1/(28/13) = 41/28
	want := mustRat(41, 28)
	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

func TestBrouncker4OverPiGCFSource_IngestPrefix(t *testing.T) {
	b, err := IngestGCFPrefix(NewBrouncker4OverPiGCFSource(), 4)
	if err != nil {
		t.Fatalf("IngestGCFPrefix failed: %v", err)
	}

	got, err := b.Convergent()
	if err != nil {
		t.Fatalf("Convergent failed: %v", err)
	}

	want := mustRat(41, 28)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

// brouncker_pi_gcf_test.go v1
