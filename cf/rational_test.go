package cf

import "testing"

func TestNewRational_NormalizesSignAndGCD(t *testing.T) {
	r, err := NewRational(-2, -4)
	if err != nil {
		t.Fatal(err)
	}
	if r.Cmp(mustRat(1, 2)) != 0 {
		t.Fatalf("got %v, want 1/2", r)
	}

	r, err = NewRational(2, -4)
	if err != nil {
		t.Fatal(err)
	}
	if r.Cmp(mustRat(-1, 2)) != 0 {
		t.Fatalf("got %v, want -1/2", r)
	}
}

func TestRational_Arithmetic(t *testing.T) {
	a, _ := NewRational(1, 3)
	b, _ := NewRational(1, 6)

	sum, _ := a.Add(b)
	if sum.Cmp(mustRat(1, 2)) != 0 {
		t.Fatalf("sum got %v, want 1/2", sum)
	}

	prod, _ := a.Mul(b)
	if prod.Cmp(mustRat(1, 18)) != 0 {
		t.Fatalf("prod got %v, want 1/18", prod)
	}
}

func TestRangeFloorBounds_PositiveInterval(t *testing.T) {
	r := NewRange(mustRat(3, 2), mustRat(11, 4), true, true) // [1.5, 2.75]

	lo, hi, err := r.FloorBounds()
	if err != nil {
		t.Fatalf("FloorBounds failed: %v", err)
	}
	if lo != 1 || hi != 2 {
		t.Fatalf("got (%d,%d) want (1,2)", lo, hi)
	}
}

func TestRangeFloorBounds_NegativeInterval(t *testing.T) {
	r := NewRange(mustRat(-7, 3), mustRat(-1, 2), true, true) // [-2.333..., -0.5]

	lo, hi, err := r.FloorBounds()
	if err != nil {
		t.Fatalf("FloorBounds failed: %v", err)
	}
	if lo != -3 || hi != -1 {
		t.Fatalf("got (%d,%d) want (-3,-1)", lo, hi)
	}
}

func TestRangeFloorBounds_ExactIntegerPoint(t *testing.T) {
	r := NewRange(mustRat(2, 1), mustRat(2, 1), true, true)

	lo, hi, err := r.FloorBounds()
	if err != nil {
		t.Fatalf("FloorBounds failed: %v", err)
	}
	if lo != 2 || hi != 2 {
		t.Fatalf("got (%d,%d) want (2,2)", lo, hi)
	}
}

func TestRangeFloorBounds_MixedSignInterval(t *testing.T) {
	r := NewRange(mustRat(-1, 2), mustRat(3, 2), true, true) // [-0.5, 1.5]

	lo, hi, err := r.FloorBounds()
	if err != nil {
		t.Fatalf("FloorBounds failed: %v", err)
	}
	if lo != -1 || hi != 1 {
		t.Fatalf("got (%d,%d) want (-1,1)", lo, hi)
	}
}

func TestRangeString_Inside(t *testing.T) {
	r := NewRange(mustRat(1, 2), mustRat(3, 2), true, false)

	got := r.String()
	want := "[1/2,3/2]{incLo=true,incHi=false,inside}"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestRangeString_Outside(t *testing.T) {
	r := NewRange(mustRat(3, 2), mustRat(1, 2), false, true)

	got := r.String()
	want := "[3/2,1/2]{incLo=false,incHi=true,outside}"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestRefineMetricCmp_InsideVsInside(t *testing.T) {
	m1 := RefineMetric{Outside: false, Magnitude: mustRat(1, 2)}
	m2 := RefineMetric{Outside: false, Magnitude: mustRat(3, 2)}

	if got := m1.Cmp(m2); got >= 0 {
		t.Fatalf("expected narrower inside metric to compare smaller, got %d", got)
	}
	if got := m2.Cmp(m1); got <= 0 {
		t.Fatalf("expected wider inside metric to compare larger, got %d", got)
	}
}

func TestRefineMetricCmp_OutsideVsOutside(t *testing.T) {
	// Outside: larger excluded gap => narrower => compares smaller.
	m1 := RefineMetric{Outside: true, Magnitude: mustRat(3, 1)}
	m2 := RefineMetric{Outside: true, Magnitude: mustRat(1, 1)}

	if got := m1.Cmp(m2); got >= 0 {
		t.Fatalf("expected larger outside gap to compare smaller, got %d", got)
	}
	if got := m2.Cmp(m1); got <= 0 {
		t.Fatalf("expected smaller outside gap to compare larger, got %d", got)
	}
}

func TestRefineMetricCmp_InsideVsOutside(t *testing.T) {
	inside := RefineMetric{Outside: false, Magnitude: mustRat(10, 1)}
	outside := RefineMetric{Outside: true, Magnitude: mustRat(1, 1)}

	if got := inside.Cmp(outside); got >= 0 {
		t.Fatalf("expected inside to compare smaller than outside, got %d", got)
	}
	if got := outside.Cmp(inside); got <= 0 {
		t.Fatalf("expected outside to compare larger than inside, got %d", got)
	}
}
