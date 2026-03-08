// sqrt_until_tol_default_test.go v1
package cf

import "testing"

func TestSqrtApproxRationalUntilResidualDefault_Sqrt2(t *testing.T) {
	got, ok, err := SqrtApproxRationalUntilResidualDefault(
		mustRat(2, 1),
		3,
		mustRat(1, 1000),
	)
	if err != nil {
		t.Fatalf("SqrtApproxRationalUntilResidualDefault failed: %v", err)
	}
	if !ok {
		t.Fatalf("expected satisfied=true")
	}

	want := mustRat(577, 408)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// sqrt_until_tol_default_test.go v1
