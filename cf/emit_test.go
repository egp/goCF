// emit_test.go v2
package cf

import "testing"

func TestSafeDigit_TrueWhenFloorsAgree(t *testing.T) {
	// Identity transform: y = x
	id := NewULFT(1, 0, 0, 1)

	// Range [1/3, 1/2] maps to itself; floor bounds are both 0.
	r := MustRange(mustRat(1, 3), mustRat(1, 2))

	d, ok, err := SafeDigit(id, r)
	if err != nil {
		t.Fatal(err)
	}
	if !ok || d != 0 {
		t.Fatalf("got (d=%d ok=%v), want (0,true)", d, ok)
	}
}

func TestSafeDigit_FalseWhenFloorsDiffer(t *testing.T) {
	// Identity transform: y = x
	id := NewULFT(1, 0, 0, 1)

	// Range [1/3, 5/2] has floors 0 and 2 -> not safe.
	r := MustRange(mustRat(1, 3), mustRat(5, 2))

	_, ok, err := SafeDigit(id, r)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatalf("expected ok=false")
	}
}

func TestSafeDigit_ErrorWhenDenominatorCrossesZero(t *testing.T) {
	// y = 1/(x-1)
	tform := NewULFT(0, 1, 1, -1)
	r := MustRange(mustRat(0, 1), mustRat(2, 1))

	_, ok, err := SafeDigit(tform, r)
	if err == nil {
		t.Fatalf("expected error")
	}
	if ok {
		t.Fatalf("expected ok=false when error")
	}
}

func TestEmitDigit_IdentityAfterEmitting0(t *testing.T) {
	// If y=x and we emit a=0, remainder should be z = 1/(x-0) = 1/x
	// New transform should be 1/x which is [[0,1],[1,0]].
	id := NewULFT(1, 0, 0, 1)

	got, err := EmitDigit(id, 0)
	if err != nil {
		t.Fatal(err)
	}
	want := NewULFT(0, 1, 1, 0)

	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestEmitDigit_CompositionMatchesDefinition(t *testing.T) {
	// Check: T'(x) = 1/(T(x) - a) using rationals.
	// Choose T(x)=(2x+1)/(3x+4), a=0, x=1/2
	T := NewULFT(2, 1, 3, 4)
	a := int64(0)
	x := mustRat(1, 2)

	// Direct definition:
	Tx, err := T.ApplyRat(x)
	if err != nil {
		t.Fatal(err)
	}
	// z = 1/(Tx - a)
	TxMinusA, err := Tx.Sub(mustRat(a, 1))
	if err != nil {
		t.Fatal(err)
	}
	zDirect, err := mustRat(1, 1).Div(TxMinusA)
	if err != nil {
		t.Fatal(err)
	}

	// Via EmitDigit:
	Tp, err := EmitDigit(T, a)
	if err != nil {
		t.Fatal(err)
	}
	zVia, err := Tp.ApplyRat(x)
	if err != nil {
		t.Fatal(err)
	}

	if zVia.Cmp(zDirect) != 0 {
		t.Fatalf("zVia=%v, zDirect=%v, T'=%v, T=%v", zVia, zDirect, Tp, T)
	}
}

func TestULFT_DeterminantHelper(t *testing.T) {
	T := NewULFT(2, 1, 3, 4)
	det, err := T.Determinant() // 2*4 - 1*3 = 5
	if err != nil {
		t.Fatal(err)
	}
	if det != 5 {
		t.Fatalf("det got %d, want 5", det)
	}
}

// emit_test.go v2
