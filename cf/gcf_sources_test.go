// gcf_sources_test.go v1
package cf

import "testing"

func TestNewSliceGCF_Empty(t *testing.T) {
	s := NewSliceGCF()

	_, _, ok := s.NextPQ()
	if ok {
		t.Fatalf("expected empty source to terminate")
	}
}

func TestNewSliceGCF_NextPQ(t *testing.T) {
	s := NewSliceGCF(
		[2]int64{3, 2},
		[2]int64{5, 7},
	)

	p, q, ok := s.NextPQ()
	if !ok {
		t.Fatalf("expected first term")
	}
	if p != 3 || q != 2 {
		t.Fatalf("got (%d,%d), want (3,2)", p, q)
	}

	p, q, ok = s.NextPQ()
	if !ok {
		t.Fatalf("expected second term")
	}
	if p != 5 || q != 7 {
		t.Fatalf("got (%d,%d), want (5,7)", p, q)
	}

	_, _, ok = s.NextPQ()
	if ok {
		t.Fatalf("expected termination after two terms")
	}
}

func TestNewSliceGCF_CopiesInput(t *testing.T) {
	terms := [][2]int64{
		{3, 2},
		{5, 7},
	}
	s := NewSliceGCF(terms...)

	terms[0] = [2]int64{99, 99}

	p, q, ok := s.NextPQ()
	if !ok {
		t.Fatalf("expected first term")
	}
	if p != 3 || q != 2 {
		t.Fatalf("got (%d,%d), want copied original (3,2)", p, q)
	}
}

// gcf_sources_test.go v1
