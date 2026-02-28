// ringbuf_test.go v1
package cf

import "testing"

func TestRingBuf_BasicLenAndDump(t *testing.T) {
	r := NewRingBuf(3)
	if r.Len() != 0 {
		t.Fatalf("expected empty len=0, got %d", r.Len())
	}

	r.Add("a")
	r.Add("b")
	if r.Len() != 2 {
		t.Fatalf("expected len=2, got %d", r.Len())
	}
	if got := r.Count("a"); got != 1 {
		t.Fatalf("expected count(a)=1, got %d", got)
	}

	r.Add("c")
	if r.Len() != 3 {
		t.Fatalf("expected len=3, got %d", r.Len())
	}

	d := r.Dump()
	if d == "" {
		t.Fatalf("expected non-empty dump")
	}
}

func TestRingBuf_OverwritesOldest(t *testing.T) {
	r := NewRingBuf(2)
	r.Add("a")
	r.Add("b")
	if r.Count("a") != 1 {
		t.Fatalf("expected count(a)=1")
	}

	r.Add("c") // overwrites "a"
	if r.Count("a") != 0 {
		t.Fatalf("expected count(a)=0 after overwrite")
	}
	if r.Count("b") != 1 || r.Count("c") != 1 {
		t.Fatalf("expected b and c present")
	}
}

func TestRingBuf_DumpIsOldestToNewest(t *testing.T) {
	r := NewRingBuf(2)
	r.Add("a")
	r.Add("b")
	r.Add("c") // window should be b,c in that order

	d := r.Dump()
	// We don't assert exact formatting, just ordering.
	if !(containsInOrder(d, "b", "c")) {
		t.Fatalf("expected dump to contain b then c; dump=\n%s", d)
	}
}

func containsInOrder(s, a, b string) bool {
	ia := indexOf(s, a)
	ib := indexOf(s, b)
	return ia >= 0 && ib >= 0 && ia < ib
}

func indexOf(s, sub string) int {
	// minimal search; sufficient for tests
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

// ringbuf_test.go v1
