// gcf_tail_metadata_test.go v1
package cf

import "testing"

func TestPositiveTailLowerBoundedGCFSources_Smoke(t *testing.T) {
	tests := []PositiveTailLowerBoundedGCFSource{
		AdaptCFToGCF(NewSliceCF(1, 2, 3)).(PositiveTailLowerBoundedGCFSource),
		NewPeriodicGCF(nil, [][2]int64{{2, 1}}),
		NewECFGSource(),
		NewBrouncker4OverPiGCFSource(),
		NewLambertPiOver4GCFSource(),
		NewUnitPArithmeticQGCFSource(1, 1),
	}

	for i, src := range tests {
		lb := src.TailLowerBound()
		if lb.Cmp(mustRat(1, 1)) != 0 {
			t.Fatalf("source %d: got lower bound %v want 1", i, lb)
		}
	}
}

// gcf_tail_metadata_test.go v1
