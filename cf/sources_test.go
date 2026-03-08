// sources_test.go v5
package cf

import "testing"

func readTerms(t *testing.T, s ContinuedFraction, n int) []int64 {
	t.Helper()
	out := make([]int64, 0, n)
	for i := 0; i < n; i++ {
		a, ok := s.Next()
		if !ok {
			t.Fatalf("expected infinite CF, got exhausted at i=%d", i)
		}
		out = append(out, a)
	}
	return out
}

func assertPrefix(t *testing.T, got []int64, want []int64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("prefix len mismatch: got=%d want=%d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("prefix mismatch at i=%d: got=%v want=%v", i, got, want)
		}
	}
}

func TestSources_PhiPrefix(t *testing.T) {
	got := readTerms(t, PhiCF(), 8)
	want := []int64{1, 1, 1, 1, 1, 1, 1, 1}
	assertPrefix(t, got, want)
}

func TestSources_Sqrt2Prefix(t *testing.T) {
	got := readTerms(t, Sqrt2CF(), 8)
	want := []int64{1, 2, 2, 2, 2, 2, 2, 2}
	assertPrefix(t, got, want)
}

func TestSources_Sqrt3Prefix(t *testing.T) {
	got := readTerms(t, Sqrt3CF(), 10)
	want := []int64{1, 1, 2, 1, 2, 1, 2, 1, 2, 1}
	assertPrefix(t, got, want)
}

func TestSources_Sqrt5Prefix(t *testing.T) {
	got := readTerms(t, Sqrt5CF(), 8)
	want := []int64{2, 4, 4, 4, 4, 4, 4, 4}
	assertPrefix(t, got, want)
}

func TestSources_Sqrt6Prefix(t *testing.T) {
	got := readTerms(t, Sqrt6CF(), 10)
	want := []int64{2, 2, 4, 2, 4, 2, 4, 2, 4, 2}
	assertPrefix(t, got, want)
}

func TestSources_Sqrt7Prefix(t *testing.T) {
	got := readTerms(t, Sqrt7CF(), 12)
	want := []int64{2, 1, 1, 1, 4, 1, 1, 1, 4, 1, 1, 1}
	assertPrefix(t, got, want)
}

func TestSources_Sqrt2SquaresTo2(t *testing.T) { t.Skip("pending algebraic-source support") }
func TestSources_Sqrt3SquaresTo3(t *testing.T) { t.Skip("pending algebraic-source support") }
func TestSources_Sqrt5SquaresTo5(t *testing.T) { t.Skip("pending algebraic-source support") }
func TestSources_Sqrt6SquaresTo6(t *testing.T) { t.Skip("pending algebraic-source support") }
func TestSources_Sqrt7SquaresTo7(t *testing.T) { t.Skip("pending algebraic-source support") }

// sources_test.go v5
