// mvp_denominator_test.go v1
package cf

import (
	"strings"
	"testing"
)

func TestMVPDenominatorApproxDefault_UsesDegreesByDefault(t *testing.T) {
	_, err := MVPDenominatorApproxDefault()
	if err == nil {
		t.Fatalf("expected stub error")
	}
	if !strings.Contains(err.Error(), "not implemented") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMVPDenominatorApprox_RejectsRadiansForMVP(t *testing.T) {
	_, err := MVPDenominatorApprox(
		DefaultSqrtPolicy2(),
		Radians(mustRat(69, 1)),
	)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "degrees") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMVPDenominatorApprox_AcceptsDegreesShape(t *testing.T) {
	_, err := MVPDenominatorApprox(
		DefaultSqrtPolicy2(),
		Degrees(mustRat(69, 1)),
	)
	if err == nil {
		t.Fatalf("expected stub error")
	}
	if !strings.Contains(err.Error(), "not implemented") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Full target formula intentionally remains in test code for now.
// This test fixes only the denominator shape:
//
//	tanh(sqrt(5)) - sin(69°)
//
// mvp_denominator_test.go v1
