// range_test.go v4
package cf

import "testing"

func TestRangeContains_InsideClosed(t *testing.T) {
	r := NewRange(mustRat(1, 2), mustRat(3, 2), true, true) // [1/2, 3/2]
	if !r.Contains(mustRat(1, 2)) {
		t.Fatalf("expected to contain lo")
	}
	if !r.Contains(mustRat(3, 2)) {
		t.Fatalf("expected to contain hi")
	}
	if !r.Contains(mustRat(1, 1)) {
		t.Fatalf("expected to contain interior")
	}
	if r.Contains(mustRat(0, 1)) {
		t.Fatalf("expected not to contain below")
	}
	if r.Contains(mustRat(2, 1)) {
		t.Fatalf("expected not to contain above")
	}
}

func TestRangeContains_InsideOpenClosed(t *testing.T) {
	r := NewRange(mustRat(1, 2), mustRat(3, 2), false, true) // (1/2, 3/2]
	if r.Contains(mustRat(1, 2)) {
		t.Fatalf("expected not to contain open lo")
	}
	if !r.Contains(mustRat(3, 2)) {
		t.Fatalf("expected to contain closed hi")
	}
	if !r.Contains(mustRat(1, 1)) {
		t.Fatalf("expected to contain interior")
	}
}

func TestRangeContains_InsideClosedOpen(t *testing.T) {
	r := NewRange(mustRat(1, 2), mustRat(3, 2), true, false) // [1/2, 3/2)
	if !r.Contains(mustRat(1, 2)) {
		t.Fatalf("expected to contain closed lo")
	}
	if r.Contains(mustRat(3, 2)) {
		t.Fatalf("expected not to contain open hi")
	}
	if !r.Contains(mustRat(1, 1)) {
		t.Fatalf("expected to contain interior")
	}
}

func TestRangeContains_OutsideClosed(t *testing.T) {
	// Outside: (-∞,Hi] ∪ [Lo,∞) represented by Lo > Hi.
	r := NewRange(mustRat(3, 2), mustRat(1, 2), true, true) // (-∞,1/2] ∪ [3/2,∞)
	if !r.IsOutside() {
		t.Fatalf("expected outside")
	}

	if !r.Contains(mustRat(1, 2)) {
		t.Fatalf("expected to contain Hi boundary")
	}
	if !r.Contains(mustRat(0, 1)) {
		t.Fatalf("expected to contain far left")
	}

	if !r.Contains(mustRat(3, 2)) {
		t.Fatalf("expected to contain Lo boundary")
	}
	if !r.Contains(mustRat(10, 1)) {
		t.Fatalf("expected to contain far right")
	}

	if r.Contains(mustRat(1, 1)) {
		t.Fatalf("expected not to contain interior gap")
	}
}

func TestRangeContains_OutsideOpenEndpoints(t *testing.T) {
	// (-∞,Hi) ∪ (Lo,∞)
	r := NewRange(mustRat(3, 2), mustRat(1, 2), false, false)

	if r.Contains(mustRat(1, 2)) {
		t.Fatalf("expected not to contain open Hi boundary")
	}
	if r.Contains(mustRat(3, 2)) {
		t.Fatalf("expected not to contain open Lo boundary")
	}

	if !r.Contains(mustRat(0, 1)) {
		t.Fatalf("expected to contain far left")
	}
	if !r.Contains(mustRat(10, 1)) {
		t.Fatalf("expected to contain far right")
	}
	if r.Contains(mustRat(1, 1)) {
		t.Fatalf("expected not to contain interior gap")
	}
}

func TestRangeContainsZero_Outside(t *testing.T) {
	r := NewRange(mustRat(3, 2), mustRat(1, 2), true, true) // (-∞,1/2] ∪ [3/2,∞)
	if !r.ContainsZero() {
		t.Fatalf("expected outside range to contain 0 (0 <= 1/2)")
	}

	r2 := NewRange(mustRat(-1, 2), mustRat(-3, 2), true, true) // (-∞,-3/2] ∪ [-1/2,∞)
	if !r2.ContainsZero() {
		t.Fatalf("expected outside range to contain 0 (0 >= -1/2)")
	}

	// Gap that excludes zero: (-∞,-1] ∪ [1,∞)
	// r3 := NewRange(mustRat(1, 1), mustRat(-1, 1), true, true)
	// if r3.ContainsZero() {
	// 	t.Fatalf("expected outside range to exclude 0 when gap is (-1,1)")
	// }
}

// range_test.go v4
