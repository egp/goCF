// pi_gcf_compare_test.go v2
package cf

import "testing"

func TestPiGCF_BrounckerAndLambert_BoundedPrefixConvergents(t *testing.T) {
	bTests := []struct {
		prefix int
		want   Rational
	}{
		{2, mustRat(3, 2)},
		{3, mustRat(15, 13)},
		{4, mustRat(105, 76)},
	}

	for _, tc := range bTests {
		got, err := GCFSourceConvergent(NewBrouncker4OverPiGCFSource(), tc.prefix)
		if err != nil {
			t.Fatalf("Brouncker prefix %d: %v", tc.prefix, err)
		}
		if got.Cmp(tc.want) != 0 {
			t.Fatalf("Brouncker prefix %d: got %v want %v", tc.prefix, got, tc.want)
		}
	}

	lTests := []struct {
		prefix int
		want   Rational
	}{
		{2, mustRat(1, 1)},
		{3, mustRat(3, 4)},
		{4, mustRat(19, 24)},
	}

	for _, tc := range lTests {
		got, err := GCFSourceConvergent(NewLambertPiOver4GCFSource(), tc.prefix)
		if err != nil {
			t.Fatalf("Lambert prefix %d: %v", tc.prefix, err)
		}
		if got.Cmp(tc.want) != 0 {
			t.Fatalf("Lambert prefix %d: got %v want %v", tc.prefix, got, tc.want)
		}
	}
}

func TestPiGCF_BrounckerAndLambert_AsRegularCFTerms(t *testing.T) {
	type tc struct {
		name   string
		src    GCFSource
		prefix int
		want   []int64
	}

	tests := []tc{
		{
			name:   "Brouncker prefix 4",
			src:    NewBrouncker4OverPiGCFSource(),
			prefix: 4,
			want:   []int64{1, 2, 1, 1, 1, 1, 1, 3}, // 105/76
		},
		{
			name:   "Lambert prefix 4",
			src:    NewLambertPiOver4GCFSource(),
			prefix: 4,
			want:   []int64{0, 1, 3, 1, 4}, // 19/24
		},
	}

	for _, tc := range tests {
		got, err := GCFSourceTerms(tc.src, tc.prefix, 16)
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if len(got) != len(tc.want) {
			t.Fatalf("%s: len(got)=%d want=%d got=%v", tc.name, len(got), len(tc.want), got)
		}
		for i := range tc.want {
			if got[i] != tc.want[i] {
				t.Fatalf("%s: got[%d]=%d want=%d full=%v", tc.name, i, got[i], tc.want[i], got)
			}
		}
	}
}

func TestPiGCF_BrounckerVsLambert_DifferentObjects(t *testing.T) {
	b, err := GCFSourceConvergent(NewBrouncker4OverPiGCFSource(), 4)
	if err != nil {
		t.Fatalf("Brouncker: %v", err)
	}
	l, err := GCFSourceConvergent(NewLambertPiOver4GCFSource(), 4)
	if err != nil {
		t.Fatalf("Lambert: %v", err)
	}

	if b.Cmp(l) == 0 {
		t.Fatalf("expected different bounded convergents, both got %v", b)
	}
}

// pi_gcf_compare_test.go v2
