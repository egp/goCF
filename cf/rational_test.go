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
