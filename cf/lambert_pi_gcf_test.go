// lambert_pi_gcf_test.go v1
package cf

import "testing"

func TestLambertPiOver4GCFSource_FirstTerms(t *testing.T) {
	s := NewLambertPiOver4GCFSource()

	want := [][2]int64{
		{0, 1},
		{1, 1},
		{3, 4},
		{5, 9},
		{7, 16},
		{9, 25},
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

func TestLambertPiOver4GCFSource_Prefix1(t *testing.T) {
	got, err := GCFSourceConvergent(NewLambertPiOver4GCFSource(), 1)
	if err != nil {
		t.Fatalf("GCFSourceConvergent failed: %v", err)
	}

	want := mustRat(0, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestLambertPiOver4GCFSource_Prefix2(t *testing.T) {
	got, err := GCFSourceConvergent(NewLambertPiOver4GCFSource(), 2)
	if err != nil {
		t.Fatalf("GCFSourceConvergent failed: %v", err)
	}

	// 0 + 1/1 = 1
	want := mustRat(1, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestLambertPiOver4GCFSource_Prefix3(t *testing.T) {
	got, err := GCFSourceConvergent(NewLambertPiOver4GCFSource(), 3)
	if err != nil {
		t.Fatalf("GCFSourceConvergent failed: %v", err)
	}

	// 1 / (1 + 1/3) = 3/4
	want := mustRat(3, 4)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestLambertPiOver4GCFSource_AsRegularCFTerms(t *testing.T) {
	got, err := GCFSourceTerms(NewLambertPiOver4GCFSource(), 3, 8)
	if err != nil {
		t.Fatalf("GCFSourceTerms failed: %v", err)
	}

	// 3/4 = [0; 1, 3]
	want := []int64{0, 1, 3}
	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d full=%v", i, got[i], want[i], got)
		}
	}
}
