// mvp_target_output_test.go v2
package cf

import (
	"fmt"
	"testing"
)

func mvpTestRangeMidpoint(r Range) (Rational, error) {
	sum, err := r.Lo.Add(r.Hi)
	if err != nil {
		return Rational{}, err
	}
	return sum.Div(mustRat(2, 1))
}

func mvpTestTargetMidpointApproxDefault() (Rational, error) {
	r, err := mvpTestTargetBoundsDefault()
	if err != nil {
		return Rational{}, err
	}
	return mvpTestRangeMidpoint(r)
}

func mvpTestTargetMidpointApproxTermsDefault(digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("mvpTestTargetMidpointApproxTermsDefault: negative digits %d", digits)
	}

	approx, err := mvpTestTargetMidpointApproxDefault()
	if err != nil {
		return nil, err
	}
	return collectTerms(NewRationalCF(approx), digits), nil
}

func TestMVPTargetMidpointApproxDefault_IsInsideTargetBounds(t *testing.T) {
	got, err := mvpTestTargetMidpointApproxDefault()
	if err != nil {
		t.Fatalf("mvpTestTargetMidpointApproxDefault failed: %v", err)
	}

	bounds, err := mvpTestTargetBoundsDefault()
	if err != nil {
		t.Fatalf("mvpTestTargetBoundsDefault failed: %v", err)
	}

	if !bounds.Contains(got) {
		t.Fatalf("got %v want inside %v", got, bounds)
	}
}

func TestMVPTargetMidpointApproxDefault_IsPositive(t *testing.T) {
	got, err := mvpTestTargetMidpointApproxDefault()
	if err != nil {
		t.Fatalf("mvpTestTargetMidpointApproxDefault failed: %v", err)
	}

	if got.Cmp(intRat(0)) <= 0 {
		t.Fatalf("got %v want positive", got)
	}
}

func TestMVPTargetMidpointApproxDefault_ExceedsOne(t *testing.T) {
	got, err := mvpTestTargetMidpointApproxDefault()
	if err != nil {
		t.Fatalf("mvpTestTargetMidpointApproxDefault failed: %v", err)
	}

	if got.Cmp(intRat(1)) <= 0 {
		t.Fatalf("got %v want > 1", got)
	}
}

func TestMVPTargetMidpointApprox_MatchesManualRangeMidpoint(t *testing.T) {
	got, err := mvpTestTargetMidpointApproxDefault()
	if err != nil {
		t.Fatalf("mvpTestTargetMidpointApproxDefault failed: %v", err)
	}

	r, err := mvpTestTargetBoundsDefault()
	if err != nil {
		t.Fatalf("mvpTestTargetBoundsDefault failed: %v", err)
	}

	sum, err := r.Lo.Add(r.Hi)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}
	want, err := sum.Div(mustRat(2, 1))
	if err != nil {
		t.Fatalf("Div failed: %v", err)
	}

	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPTargetMidpointApproxTermsDefault_MatchesMidpointRationalCF(t *testing.T) {
	got, err := mvpTestTargetMidpointApproxTermsDefault(12)
	if err != nil {
		t.Fatalf("mvpTestTargetMidpointApproxTermsDefault failed: %v", err)
	}

	approx, err := mvpTestTargetMidpointApproxDefault()
	if err != nil {
		t.Fatalf("mvpTestTargetMidpointApproxDefault failed: %v", err)
	}

	want := collectTerms(NewRationalCF(approx), 12)

	if len(got) != len(want) {
		t.Fatalf("len mismatch: got=%v want=%v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("digit %d: got=%v want=%v", i, got, want)
		}
	}
}

// mvp_target_output_test.go v2
