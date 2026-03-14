// sqrt_gcf_tail_api2_test.go v1
package cf

import (
	"strings"
	"testing"
)

func TestSqrtApproxFromGCFTailSource2_ExactTailSource_EmptyPrefix(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	got, err := SqrtApproxFromGCFTailSource2(
		NewSliceGCF(),
		NewExactTailSource(mustRat(9, 16)),
		8,
		p,
	)
	if err != nil {
		t.Fatalf("SqrtApproxFromGCFTailSource2 failed: %v", err)
	}

	want := mustRat(3, 4)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSqrtApproxFromGCFTailSource2_MissingTailEvidenceIsError(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	_, err := SqrtApproxFromGCFTailSource2(
		NewSliceGCF([2]int64{3, 2}),
		NoTailSource{},
		8,
		p,
	)
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if !strings.Contains(err.Error(), "tail evidence not implemented") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSqrtApproxCFFromGCFTailSource2_HasExpectedLeadingDigit(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	cf, err := SqrtApproxCFFromGCFTailSource2(
		NewSliceGCF([2]int64{3, 2}),
		NewExactTailSource(mustRat(11, 1)),
		8,
		p,
	)
	if err != nil {
		t.Fatalf("SqrtApproxCFFromGCFTailSource2 failed: %v", err)
	}

	got := collectTerms(cf, 8)
	if len(got) == 0 {
		t.Fatalf("expected non-empty CF")
	}
	if got[0] != 1 {
		t.Fatalf("got first digit %d want 1 full=%v", got[0], got)
	}
}

func TestSqrtApproxTermsFromGCFTailSourceDefault2_RejectsNegativeDigits(t *testing.T) {
	_, err := SqrtApproxTermsFromGCFTailSourceDefault2(
		NewSliceGCF([2]int64{3, 2}),
		NewExactTailSource(mustRat(11, 1)),
		8,
		-1,
	)
	if err == nil {
		t.Fatalf("expected error for negative digits")
	}
}

func TestSqrtApproxFromGCFTailSource2_HonorsBoundFailure(t *testing.T) {
	p := SqrtPolicy2{
		MaxSteps: 3,
		Tol:      mustRat(1, 1000),
	}

	_, err := SqrtApproxFromGCFTailSource2(
		NewUnitPArithmeticQGCFSource(1, 1),
		NewExactTailSource(mustRat(1, 1)),
		3,
		p,
	)
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
}

// sqrt_gcf_tail_api2_test.go v1
