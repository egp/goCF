// sources_test.go v1
package cf

import "testing"

func takeN(cf ContinuedFraction, n int) ([]int64, bool) {
	out := make([]int64, 0, n)
	for i := 0; i < n; i++ {
		v, ok := cf.Next()
		if !ok {
			return out, false
		}
		out = append(out, v)
	}
	return out, true
}

func TestSqrt2CF_Prefix(t *testing.T) {
	cf := Sqrt2CF()
	got, ok := takeN(cf, 8)
	if !ok {
		t.Fatalf("expected infinite stream")
	}
	want := []int64{1, 2, 2, 2, 2, 2, 2, 2}
	if !equalSlice(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestPhiCF_Prefix(t *testing.T) {
	cf := PhiCF()
	got, ok := takeN(cf, 8)
	if !ok {
		t.Fatalf("expected infinite stream")
	}
	want := []int64{1, 1, 1, 1, 1, 1, 1, 1}
	if !equalSlice(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestPeriodicCF_FiniteIfNoPeriod(t *testing.T) {
	cf := NewPeriodicCF([]int64{3, 7, 16}, nil)
	got, ok := takeN(cf, 3)
	if !ok {
		t.Fatalf("expected to read 3 terms")
	}
	want := []int64{3, 7, 16}
	if !equalSlice(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	_, ok = takeN(cf, 1)
	if ok {
		t.Fatalf("expected stream to be exhausted")
	}
}

// sources_test.go v1
