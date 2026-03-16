// mvp_target_output_test.go v1
package cf

import "testing"

func TestMVPTargetMidpointApproxDefault_IsInsideTargetBounds(t *testing.T) {
	got, err := MVPTargetMidpointApproxDefault()
	if err != nil {
		t.Fatalf("MVPTargetMidpointApproxDefault failed: %v", err)
	}

	bounds, err := MVPTargetBoundsDefault()
	if err != nil {
		t.Fatalf("MVPTargetBoundsDefault failed: %v", err)
	}

	if !bounds.Contains(got) {
		t.Fatalf("got %v want inside %v", got, bounds)
	}
}

func TestMVPTargetMidpointApproxDefault_IsPositive(t *testing.T) {
	got, err := MVPTargetMidpointApproxDefault()
	if err != nil {
		t.Fatalf("MVPTargetMidpointApproxDefault failed: %v", err)
	}

	if got.Cmp(intRat(0)) <= 0 {
		t.Fatalf("got %v want positive", got)
	}
}

func TestMVPTargetMidpointApproxDefault_ExceedsOne(t *testing.T) {
	got, err := MVPTargetMidpointApproxDefault()
	if err != nil {
		t.Fatalf("MVPTargetMidpointApproxDefault failed: %v", err)
	}

	if got.Cmp(intRat(1)) <= 0 {
		t.Fatalf("got %v want > 1", got)
	}
}

func TestMVPTargetMidpointApprox_MatchesManualRangeMidpoint(t *testing.T) {
	got, err := MVPTargetMidpointApproxDefault()
	if err != nil {
		t.Fatalf("MVPTargetMidpointApproxDefault failed: %v", err)
	}

	r, err := MVPTargetBoundsDefault()
	if err != nil {
		t.Fatalf("MVPTargetBoundsDefault failed: %v", err)
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
	got, err := MVPTargetMidpointApproxTermsDefault(12)
	if err != nil {
		t.Fatalf("MVPTargetMidpointApproxTermsDefault failed: %v", err)
	}

	approx, err := MVPTargetMidpointApproxDefault()
	if err != nil {
		t.Fatalf("MVPTargetMidpointApproxDefault failed: %v", err)
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

// mvp_target_output_test.go v1
