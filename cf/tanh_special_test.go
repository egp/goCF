// tanh_special_test.go v3
package cf

import "testing"

func TestTanhBoundsSqrt5_IsInsidePositiveTightenedRange(t *testing.T) {
	got := TanhBoundsSqrt5()
	want := NewRange(mustRat(39, 40), mustRat(49, 50), true, true)

	if got.Lo.Cmp(want.Lo) != 0 || got.Hi.Cmp(want.Hi) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
	if !got.IsInside() {
		t.Fatalf("expected inside range, got %v", got)
	}
}

// tanh_special_test.go v3
