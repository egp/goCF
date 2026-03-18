// cf/sqrt_unary_digit_test.go v2
package cf

import "testing"

func TestSqrtUnaryNextDigitForced_PointTwoToTwo_ForcesDigitTwo(t *testing.T) {
	r := NewRange(mustRat(2, 1), mustRat(2, 1), true, true)

	d, ok, err := sqrtUnaryNextDigitIfForced(r)
	if err != nil {
		t.Fatalf("sqrtUnaryNextDigitIfForced failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected forced digit")
	}
	if d != 2 {
		t.Fatalf("got %d want 2", d)
	}
}

func TestSqrtUnaryNextDigitForced_RangeFourThirdsToThreeHalves_ForcesDigitOne(t *testing.T) {
	r := NewRange(mustRat(4, 3), mustRat(3, 2), true, true)

	d, ok, err := sqrtUnaryNextDigitIfForced(r)
	if err != nil {
		t.Fatalf("sqrtUnaryNextDigitIfForced failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected forced digit")
	}
	if d != 1 {
		t.Fatalf("got %d want 1", d)
	}
}

func TestSqrtUnaryNextDigitForced_RangeNineTenthsToElevenTenths_DoesNotForceDigit(t *testing.T) {
	r := NewRange(mustRat(9, 10), mustRat(11, 10), true, true)

	d, ok, err := sqrtUnaryNextDigitIfForced(r)
	if err != nil {
		t.Fatalf("sqrtUnaryNextDigitIfForced failed: %v", err)
	}
	if ok {
		t.Fatalf("unexpected forced digit %d", d)
	}
}

func TestSqrtUnaryOperator_NextDigitIfForced_InitialStateNotForced(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	d, ok, err := op.nextDigitIfForced()
	if err != nil {
		t.Fatalf("nextDigitIfForced failed: %v", err)
	}
	if ok {
		t.Fatalf("unexpected forced digit %d", d)
	}
}

func TestSqrtUnaryOperator_NextDigitIfForced_AfterOneIngest_ForcesDigitOne(t *testing.T) {
	op, err := newSqrtUnaryOperator(NewECFGSource(), mustRat(1, 1), defaultSqrtUnaryRefinementPolicy())
	if err != nil {
		t.Fatalf("newSqrtUnaryOperator failed: %v", err)
	}

	if err := op.ingestOneAndRefine(); err != nil {
		t.Fatalf("ingestOneAndRefine failed: %v", err)
	}

	d, ok, err := op.nextDigitIfForced()
	if err != nil {
		t.Fatalf("nextDigitIfForced failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected forced digit")
	}
	if d != 1 {
		t.Fatalf("got %d want 1", d)
	}
}

// cf/sqrt_unary_digit_test.go v2
