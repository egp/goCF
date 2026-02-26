package cf

import "testing"

func collectCF(cf ContinuedFraction) []int64 {
	var out []int64
	for {
		v, ok := cf.Next()
		if !ok {
			return out
		}
		out = append(out, v)
	}
}

func TestRationalCF_PositiveExamples(t *testing.T) {
	// 3/2 = [1;2]
	got := collectCF(NewRationalCF(mustRat(3, 2)))
	want := []int64{1, 2}
	if !equalSlice(got, want) {
		t.Fatalf("3/2 got %v, want %v", got, want)
	}

	// 355/113 = [3;7,16] (classic pi convergent)
	got = collectCF(NewRationalCF(mustRat(355, 113)))
	want = []int64{3, 7, 16}
	if !equalSlice(got, want) {
		t.Fatalf("355/113 got %v, want %v", got, want)
	}
}

func TestRationalCF_NegativeExamples_FloorConvention(t *testing.T) {
	// Using floor convention:
	// -4/3 = [-2;1,2] because -4/3 = -2 + 2/3, and 3/2 = [1;2]
	got := collectCF(NewRationalCF(mustRat(-4, 3)))
	want := []int64{-2, 1, 2}
	if !equalSlice(got, want) {
		t.Fatalf("-4/3 got %v, want %v", got, want)
	}

	// -1/2 = [-1;2] because -0.5 = -1 + 1/2
	got = collectCF(NewRationalCF(mustRat(-1, 2)))
	want = []int64{-1, 2}
	if !equalSlice(got, want) {
		t.Fatalf("-1/2 got %v, want %v", got, want)
	}
}

func equalSlice(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
