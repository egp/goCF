// sqrt_gcf_api2_test.go v1
package cf

import "testing"

func TestSqrtApproxFromGCFWithTail2_EmptyPrefixUsesTail(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	got, err := SqrtApproxFromGCFWithTail2(NewSliceGCF(), mustRat(9, 16), 8, p)
	if err != nil {
		t.Fatalf("SqrtApproxFromGCFWithTail2 failed: %v", err)
	}

	want := mustRat(3, 4)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSqrtApproxFromGCFWithTail2_FinitePrefixMatchesExactInputThenSqrt(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	// x = 3 + 2/11 = 35/11
	got, err := SqrtApproxFromGCFWithTail2(
		NewSliceGCF([2]int64{3, 2}),
		mustRat(11, 1),
		8,
		p,
	)
	if err != nil {
		t.Fatalf("SqrtApproxFromGCFWithTail2 failed: %v", err)
	}

	// Compare against the direct rational path on x = 35/11
	want, err := SqrtApproxWithPolicy2(mustRat(35, 11), p)
	if err != nil {
		t.Fatalf("SqrtApproxWithPolicy2 failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSqrtApproxCFFromGCFWithTail2_ProducesExpectedLeadingDigit(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	cf, err := SqrtApproxCFFromGCFWithTail2(
		NewSliceGCF([2]int64{3, 2}),
		mustRat(11, 1),
		8,
		p,
	)
	if err != nil {
		t.Fatalf("SqrtApproxCFFromGCFWithTail2 failed: %v", err)
	}

	got := collectTerms(cf, 8)
	if len(got) == 0 {
		t.Fatalf("expected non-empty CF")
	}
	if got[0] != 1 {
		t.Fatalf("got first digit %d want 1 full=%v", got[0], got)
	}
}

func TestSqrtApproxTermsFromGCFWithTailDefault2_RejectsNegativeDigits(t *testing.T) {
	_, err := SqrtApproxTermsFromGCFWithTailDefault2(
		NewSliceGCF([2]int64{3, 2}),
		mustRat(11, 1),
		8,
		-1,
	)
	if err == nil {
		t.Fatalf("expected error for negative digits")
	}
}

func TestSqrtApproxFromGCFWithTail2_HonorsBoundFailure(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	_, err := SqrtApproxFromGCFWithTail2(
		NewUnitPArithmeticQGCFSource(1, 1),
		mustRat(1, 1),
		3,
		p,
	)
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
}

// sqrt_gcf_api2_test.go v1
