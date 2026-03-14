// sqrt_seed_test.go v1
package cf

import "testing"

func TestNewSqrtApproxCFDefault_Sqrt2_ThreeSteps(t *testing.T) {
	cf, err := NewSqrtApproxCFDefault(mustRat(2, 1), 3)
	if err != nil {
		t.Fatalf("NewSqrtApproxCFDefault failed: %v", err)
	}

	got := collectTerms(cf, 16)
	want := []int64{1, 2, 2, 2, 2, 2, 2, 2} // 577/408

	if len(got) != len(want) {
		t.Fatalf("len(got)=%d want=%d got=%v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d]=%d want=%d full=%v", i, got[i], want[i], got)
		}
	}
}

// sqrt_seed_test.go v1
