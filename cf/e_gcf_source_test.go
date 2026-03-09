// e_gcf_source_test.go v1
package cf

import "testing"

func TestECFGSource_FirstTerms(t *testing.T) {
	s := NewECFGSource()

	want := [][2]int64{
		{2, 1},
		{1, 1},
		{2, 1},
		{1, 1},
		{1, 1},
		{4, 1},
		{1, 1},
		{1, 1},
		{6, 1},
		{1, 1},
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

func TestECFGSource_ApproxPrefix3(t *testing.T) {
	// Prefix terms: (2,1), (1,1), (2,1)
	// => 2 + 1/(1 + 1/2) = 2 + 2/3 = 8/3
	got, err := GCFApproxFromPrefix(NewECFGSource(), 3)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}

	want := mustRat(8, 3)
	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

func TestECFGSource_ApproxPrefix6(t *testing.T) {
	// Prefix regular CF [2;1,2,1,1,4] = 87/32
	got, err := GCFApproxFromPrefix(NewECFGSource(), 6)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}

	want := mustRat(87, 32)
	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

// e_gcf_source_test.go v1
