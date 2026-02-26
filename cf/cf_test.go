package cf

import "testing"

func TestSliceCF_Basic(t *testing.T) {
	s := NewSliceCF(1, 2, 3)

	if v, ok := s.Next(); !ok || v != 1 {
		t.Fatalf("first: got (%d,%v), want (1,true)", v, ok)
	}
	if v, ok := s.Next(); !ok || v != 2 {
		t.Fatalf("second: got (%d,%v), want (2,true)", v, ok)
	}
	if v, ok := s.Next(); !ok || v != 3 {
		t.Fatalf("third: got (%d,%v), want (3,true)", v, ok)
	}
	if _, ok := s.Next(); ok {
		t.Fatalf("fourth: got ok=true, want ok=false")
	}
}
