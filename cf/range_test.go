// range_test.go v6
package cf

import "testing"

func TestRangeRefineMetricOrdering_Gosper(t *testing.T) {
	insideNarrow := NewRange(mustRat(0, 1), mustRat(1, 10), true, true) // span 0.1
	insideWide := NewRange(mustRat(0, 1), mustRat(2, 1), true, true)    // span 2

	// Outside wide (narrower outside): big excluded gap
	outsideWide := NewRange(mustRat(10, 1), mustRat(-10, 1), true, true) // gap 20

	// Outside narrow (wider outside): small excluded gap
	outsideNarrow := NewRange(mustRat(1, 1), mustRat(0, 1), true, true) // gap 1

	mInN, _ := insideNarrow.RefineMetric()
	mInW, _ := insideWide.RefineMetric()
	mOutW, _ := outsideWide.RefineMetric()
	mOutN, _ := outsideNarrow.RefineMetric()

	// inside narrow < inside wide
	if mInN.Cmp(mInW) >= 0 {
		t.Fatalf("expected inside narrow < inside wide: %v vs %v", mInN, mInW)
	}
	// inside wide < outside wide
	if mInW.Cmp(mOutW) >= 0 {
		t.Fatalf("expected inside wide < outside wide: %v vs %v", mInW, mOutW)
	}
	// outside wide < outside narrow (bigger excluded gap => narrower outside)
	if mOutW.Cmp(mOutN) >= 0 {
		t.Fatalf("expected outside wide < outside narrow: %v vs %v", mOutW, mOutN)
	}
}

// range_test.go v6
