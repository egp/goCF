// range_test.go v3
package cf

import "testing"

func TestRange_IsInsideOutside(t *testing.T) {
	in := NewRange(mustRat(1, 3), mustRat(1, 2))
	if !in.IsInside() || in.IsOutside() {
		t.Fatalf("expected inside range")
	}

	out := NewRange(mustRat(1, 2), mustRat(1, 3)) // Lo > Hi => outside
	if out.IsInside() || !out.IsOutside() {
		t.Fatalf("expected outside range")
	}
}

func TestMustRange_PanicsOnOutside(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("expected panic")
		}
	}()
	_ = MustRange(mustRat(1, 2), mustRat(1, 3))
}

func TestRange_ContainsInside(t *testing.T) {
	r := MustRange(mustRat(1, 3), mustRat(1, 2))

	if !r.Contains(mustRat(1, 3)) {
		t.Fatalf("expected to contain Lo")
	}
	if !r.Contains(mustRat(1, 2)) {
		t.Fatalf("expected to contain Hi")
	}
	if !r.Contains(mustRat(2, 5)) {
		t.Fatalf("expected to contain 2/5")
	}
	if r.Contains(mustRat(1, 4)) {
		t.Fatalf("did not expect to contain 1/4")
	}
}

func TestRange_ContainsOutside(t *testing.T) {
	// Outside range: Lo > Hi => (-∞,Hi] ∪ [Lo,∞)
	// Example: Lo=5, Hi=3 => (-∞,3] ∪ [5,∞)
	r := NewRange(mustRat(5, 1), mustRat(3, 1))

	if !r.Contains(mustRat(0, 1)) {
		t.Fatalf("expected to contain 0 (<=Hi)")
	}
	if !r.Contains(mustRat(3, 1)) {
		t.Fatalf("expected to contain Hi endpoint")
	}
	if r.Contains(mustRat(4, 1)) {
		t.Fatalf("did not expect to contain 4 (in the hole)")
	}
	if !r.Contains(mustRat(5, 1)) {
		t.Fatalf("expected to contain Lo endpoint")
	}
	if !r.Contains(mustRat(100, 1)) {
		t.Fatalf("expected to contain large values (>=Lo)")
	}
	if !r.Contains(mustRat(-100, 1)) {
		t.Fatalf("expected to contain small values (<=Hi)")
	}
}

func TestRange_ContainsZero_Inside(t *testing.T) {
	r := MustRange(mustRat(-1, 2), mustRat(1, 3))
	if !r.ContainsZero() {
		t.Fatalf("expected inside range to contain 0")
	}
	r2 := MustRange(mustRat(1, 10), mustRat(3, 10))
	if r2.ContainsZero() {
		t.Fatalf("did not expect positive inside range to contain 0")
	}
}

func TestRange_ContainsZero_Outside(t *testing.T) {
	// Outside: (-∞,Hi] ∪ [Lo,∞)
	// Case 1: Hi >= 0 => contains 0
	r := NewRange(mustRat(5, 1), mustRat(3, 1))
	if !r.ContainsZero() {
		t.Fatalf("expected outside range to contain 0 when Hi>=0")
	}

	// Case 2: Lo <= 0 => contains 0
	r2 := NewRange(mustRat(-1, 1), mustRat(-3, 1)) // (-∞,-3] ∪ [-1,∞) includes 0
	if !r2.ContainsZero() {
		t.Fatalf("expected outside range to contain 0 when Lo<=0")
	}

	// Case 3: Hi < 0 < Lo => does NOT contain 0
	r3 := NewRange(mustRat(2, 1), mustRat(-2, 1)) // (-∞,-2] ∪ [2,∞) excludes 0
	if r3.ContainsZero() {
		t.Fatalf("did not expect outside range with Hi<0<Lo to contain 0")
	}
}

func TestRange_WidthOutsideErrors(t *testing.T) {
	out := NewRange(mustRat(1, 2), mustRat(1, 3))
	if _, err := out.Width(); err == nil {
		t.Fatalf("expected error for width on outside range")
	}
}

func TestRange_FloorBounds(t *testing.T) {
	r := MustRange(mustRat(1, 3), mustRat(5, 2))
	lo, hi, err := r.FloorBounds()
	if err != nil {
		t.Fatal(err)
	}
	if lo != 0 || hi != 2 {
		t.Fatalf("floor bounds got (%d,%d), want (0,2)", lo, hi)
	}

	neg := MustRange(mustRat(-4, 3), mustRat(-1, 2))
	lo, hi, err = neg.FloorBounds()
	if err != nil {
		t.Fatal(err)
	}
	if lo != -2 || hi != -1 {
		t.Fatalf("floor bounds got (%d,%d), want (-2,-1)", lo, hi)
	}

	out := NewRange(mustRat(2, 1), mustRat(-2, 1))
	if _, _, err := out.FloorBounds(); err == nil {
		t.Fatalf("expected error for FloorBounds on outside range")
	}
}

// range_test.go v3
