package cf

import "testing"

func TestRationalCFNext_Zero(t *testing.T) {
	cf := NewRationalCF(mustRat(0, 1))

	d, ok := cf.Next()
	if !ok {
		t.Fatalf("expected first digit")
	}
	if d != 0 {
		t.Fatalf("got %d want 0", d)
	}

	_, ok = cf.Next()
	if ok {
		t.Fatalf("expected termination after [0]")
	}
}

func TestRationalCFNext_ExactInteger(t *testing.T) {
	cf := NewRationalCF(mustRat(5, 1))

	d, ok := cf.Next()
	if !ok {
		t.Fatalf("expected first digit")
	}
	if d != 5 {
		t.Fatalf("got %d want 5", d)
	}

	_, ok = cf.Next()
	if ok {
		t.Fatalf("expected termination after [5]")
	}
}

func TestRationalCFNext_PositiveProperFraction(t *testing.T) {
	// 3/2 = [1;2]
	cf := NewRationalCF(mustRat(3, 2))

	want := []int64{1, 2}
	for i, wd := range want {
		d, ok := cf.Next()
		if !ok {
			t.Fatalf("expected digit %d", i)
		}
		if d != wd {
			t.Fatalf("digit %d: got %d want %d", i, d, wd)
		}
	}

	_, ok := cf.Next()
	if ok {
		t.Fatalf("expected termination after %v", want)
	}
}

func TestRationalCFNext_NegativeProperFraction(t *testing.T) {
	// -3/2 = [-2;2]
	cf := NewRationalCF(mustRat(-3, 2))

	want := []int64{-2, 2}
	for i, wd := range want {
		d, ok := cf.Next()
		if !ok {
			t.Fatalf("expected digit %d", i)
		}
		if d != wd {
			t.Fatalf("digit %d: got %d want %d", i, d, wd)
		}
	}

	_, ok := cf.Next()
	if ok {
		t.Fatalf("expected termination after %v", want)
	}
}

func TestRationalCFNext_MultiDigitFiniteExpansion(t *testing.T) {
	// 17/12 = [1;2,2,2]
	cf := NewRationalCF(mustRat(17, 12))

	want := []int64{1, 2, 2, 2}
	for i, wd := range want {
		d, ok := cf.Next()
		if !ok {
			t.Fatalf("expected digit %d", i)
		}
		if d != wd {
			t.Fatalf("digit %d: got %d want %d", i, d, wd)
		}
	}

	_, ok := cf.Next()
	if ok {
		t.Fatalf("expected termination after %v", want)
	}
}

func TestRationalCFNext_TerminatesAndStaysDone(t *testing.T) {
	cf := NewRationalCF(mustRat(1, 2)) // [0;2]

	d, ok := cf.Next()
	if !ok || d != 0 {
		t.Fatalf("got (%d,%v) want (0,true)", d, ok)
	}
	d, ok = cf.Next()
	if !ok || d != 2 {
		t.Fatalf("got (%d,%v) want (2,true)", d, ok)
	}

	_, ok = cf.Next()
	if ok {
		t.Fatalf("expected first termination")
	}
	_, ok = cf.Next()
	if ok {
		t.Fatalf("expected second termination to remain false")
	}
}
