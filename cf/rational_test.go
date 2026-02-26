package cf

import "testing"

func TestNewRational_NormalizesSignAndGCD(t *testing.T) {
	r, err := NewRational(-2, -4)
	if err != nil {
		t.Fatal(err)
	}
	if r.P != 1 || r.Q != 2 {
		t.Fatalf("got %v, want 1/2", r)
	}

	r, err = NewRational(2, -4)
	if err != nil {
		t.Fatal(err)
	}
	if r.P != -1 || r.Q != 2 {
		t.Fatalf("got %v, want -1/2", r)
	}
}

func TestRational_Arithmetic(t *testing.T) {
	a, _ := NewRational(1, 3)
	b, _ := NewRational(1, 6)

	sum, _ := a.Add(b)
	if sum.P != 1 || sum.Q != 2 {
		t.Fatalf("sum got %v, want 1/2", sum)
	}

	prod, _ := a.Mul(b)
	if prod.P != 1 || prod.Q != 18 {
		t.Fatalf("prod got %v, want 1/18", prod)
	}
}
