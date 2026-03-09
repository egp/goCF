// nontrivial_gcf_source_test.go v1
package cf

import "testing"

func TestUnitPArithmeticQGCFSource_FirstTerms(t *testing.T) {
	s := NewUnitPArithmeticQGCFSource(1, 1)

	want := [][2]int64{
		{1, 1},
		{1, 2},
		{1, 3},
		{1, 4},
		{1, 5},
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

func TestUnitPArithmeticQGCFSource_WithGCFApproxPrefix2(t *testing.T) {
	s := NewUnitPArithmeticQGCFSource(1, 1)

	got, err := GCFApproxFromPrefix(s, 2)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}

	// Prefix: (1,1), (1,2)
	// finite convention => 1 + 1/1 = 2
	want := mustRat(2, 1)
	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

func TestUnitPArithmeticQGCFSource_WithGCFApproxPrefix3(t *testing.T) {
	s := NewUnitPArithmeticQGCFSource(1, 1)

	got, err := GCFApproxFromPrefix(s, 3)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}

	// Prefix: (1,1), (1,2), (1,3)
	// finite convention => 1 + 1/(1 + 2/1) = 4/3
	want := mustRat(4, 3)
	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

func TestUnitPArithmeticQGCFSource_WithIngestPrefix(t *testing.T) {
	s := NewUnitPArithmeticQGCFSource(2, 2) // q = 2,4,6,...

	b, err := IngestGCFPrefix(s, 3)
	if err != nil {
		t.Fatalf("IngestGCFPrefix failed: %v", err)
	}

	got, err := b.Convergent()
	if err != nil {
		t.Fatalf("Convergent failed: %v", err)
	}

	// Prefix: (1,2), (1,4), (1,6)
	// finite convention => 1 + 2/(1 + 4/1) = 1 + 2/5 = 7/5
	want := mustRat(7, 5)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

// nontrivial_gcf_source_test.go v1
