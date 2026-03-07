package cf

import "testing"

func TestULFT_ApplyRat(t *testing.T) {
	// x -> (2x + 1) / (3x + 4)
	tform := NewULFT(bi(2), bi(1), bi(3), bi(4))
	x, _ := NewRational(1, 2)

	got, err := tform.ApplyRat(x)
	if err != nil {
		t.Fatal(err)
	}
	// (2*1/2+1)/(3*1/2+4) = (1+1)/(1.5+4)=2/5.5=4/11
	want, _ := NewRational(4, 11)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestBLFT_ApplyRat(t *testing.T) {
	// (x+y) as BLFT:
	// (0*xy + 1*x + 1*y + 0) / (0*xy + 0*x + 0*y + 1)
	tform := NewBLFT(0, 1, 1, 0, 0, 0, 0, 1)

	x, _ := NewRational(1, 3)
	y, _ := NewRational(1, 6)

	got, err := tform.ApplyRat(x, y)
	if err != nil {
		t.Fatal(err)
	}
	want, _ := NewRational(1, 2)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}
