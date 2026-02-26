// checked_test.go v1
package cf

import (
	"math"
	"testing"
)

func TestAdd64_NoOverflow(t *testing.T) {
	got, ok := add64(10, 20)
	if !ok || got != 30 {
		t.Fatalf("add64(10,20) got (%d,%v), want (30,true)", got, ok)
	}

	got, ok = add64(-10, -20)
	if !ok || got != -30 {
		t.Fatalf("add64(-10,-20) got (%d,%v), want (-30,true)", got, ok)
	}

	got, ok = add64(-10, 20)
	if !ok || got != 10 {
		t.Fatalf("add64(-10,20) got (%d,%v), want (10,true)", got, ok)
	}
}

func TestAdd64_Overflow(t *testing.T) {
	_, ok := add64(math.MaxInt64, 1)
	if ok {
		t.Fatalf("expected overflow for MaxInt64+1")
	}
	_, ok = add64(math.MinInt64, -1)
	if ok {
		t.Fatalf("expected overflow for MinInt64-1")
	}
}

func TestSub64_NoOverflow(t *testing.T) {
	got, ok := sub64(20, 10)
	if !ok || got != 10 {
		t.Fatalf("sub64(20,10) got (%d,%v), want (10,true)", got, ok)
	}

	got, ok = sub64(-20, -10)
	if !ok || got != -10 {
		t.Fatalf("sub64(-20,-10) got (%d,%v), want (-10,true)", got, ok)
	}

	got, ok = sub64(-10, 20)
	if !ok || got != -30 {
		t.Fatalf("sub64(-10,20) got (%d,%v), want (-30,true)", got, ok)
	}
}

func TestSub64_Overflow(t *testing.T) {
	_, ok := sub64(math.MinInt64, 1)
	if ok {
		t.Fatalf("expected overflow for MinInt64-1 via sub64(MinInt64,1)")
	}
	_, ok = sub64(math.MaxInt64, -1)
	if ok {
		t.Fatalf("expected overflow for MaxInt64+1 via sub64(MaxInt64,-1)")
	}
}

func TestMul64_NoOverflow(t *testing.T) {
	got, ok := mul64(7, 6)
	if !ok || got != 42 {
		t.Fatalf("mul64(7,6) got (%d,%v), want (42,true)", got, ok)
	}

	got, ok = mul64(-7, 6)
	if !ok || got != -42 {
		t.Fatalf("mul64(-7,6) got (%d,%v), want (-42,true)", got, ok)
	}

	got, ok = mul64(-7, -6)
	if !ok || got != 42 {
		t.Fatalf("mul64(-7,-6) got (%d,%v), want (42,true)", got, ok)
	}

	got, ok = mul64(0, math.MaxInt64)
	if !ok || got != 0 {
		t.Fatalf("mul64(0,MaxInt64) got (%d,%v), want (0,true)", got, ok)
	}
}

func TestMul64_Overflow(t *testing.T) {
	_, ok := mul64(math.MaxInt64, 2)
	if ok {
		t.Fatalf("expected overflow for MaxInt64*2")
	}

	_, ok = mul64(math.MinInt64, -1)
	if ok {
		t.Fatalf("expected overflow for MinInt64*-1")
	}

	// Large magnitude overflow with opposite signs.
	_, ok = mul64(-(1 << 62), 4)
	if ok {
		t.Fatalf("expected overflow for -(2^62)*4")
	}
}

func TestMulAdd64_NoOverflow(t *testing.T) {
	// 7*6 + 5 = 47
	got, ok := mulAdd64(7, 6, 5)
	if !ok || got != 47 {
		t.Fatalf("mulAdd64(7,6,5) got (%d,%v), want (47,true)", got, ok)
	}

	// (-7)*6 + 5 = -37
	got, ok = mulAdd64(-7, 6, 5)
	if !ok || got != -37 {
		t.Fatalf("mulAdd64(-7,6,5) got (%d,%v), want (-37,true)", got, ok)
	}
}

func TestMulAdd64_Overflow(t *testing.T) {
	// MaxInt64*2 definitely overflows in multiplication.
	_, ok := mulAdd64(math.MaxInt64, 2, 0)
	if ok {
		t.Fatalf("expected overflow for MaxInt64*2+0")
	}

	// MaxInt64 + 1 overflows in addition stage.
	_, ok = mulAdd64(1, 1, math.MaxInt64) // 1 + MaxInt64 overflows
	if ok {
		t.Fatalf("expected overflow for 1*1+MaxInt64")
	}
}

// checked_test.go v1
