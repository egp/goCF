// sources_test.go v7
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

func assertDiagSquaresToInt(t *testing.T, name string, src func() ContinuedFraction, n int64) {
	t.Helper()

	sq := NewDiagBLFT(
		bi(1), bi(0), bi(0),
		bi(0), bi(0), bi(1),
	)

	s := NewDiagBLFTStream(sq, src(), DiagBLFTStreamOptions{})

	a0, ok := s.Next()
	if !ok {
		t.Fatalf("%s^2: expected first digit, stream exhausted; err=%v", name, s.Err())
	}
	if a0 != n {
		t.Fatalf("%s^2: first digit mismatch: got=%d want=%d err=%v", name, a0, n, s.Err())
	}

	if a1, ok := s.Next(); ok {
		t.Fatalf("%s^2: expected termination after first digit; got extra digit %d err=%v", name, a1, s.Err())
	}

	if err := s.Err(); err != nil {
		t.Fatalf("%s^2: stream error: %v", name, err)
	}
}

func TestSources_Sqrt2SquaresTo2(t *testing.T) { assertDiagSquaresToInt(t, "sqrt2", Sqrt2CF, 2) }
func TestSources_Sqrt3SquaresTo3(t *testing.T) { assertDiagSquaresToInt(t, "sqrt3", Sqrt3CF, 3) }
func TestSources_Sqrt5SquaresTo5(t *testing.T) { assertDiagSquaresToInt(t, "sqrt5", Sqrt5CF, 5) }
func TestSources_Sqrt6SquaresTo6(t *testing.T) { assertDiagSquaresToInt(t, "sqrt6", Sqrt6CF, 6) }
func TestSources_Sqrt7SquaresTo7(t *testing.T) { assertDiagSquaresToInt(t, "sqrt7", Sqrt7CF, 7) }

func TestFuncGCFSource_Finite(t *testing.T) {
	s := NewFuncGCFSource(func(i int) (int64, int64, bool) {
		switch i {
		case 0:
			return 3, 2, true
		case 1:
			return 5, 7, true
		default:
			return 0, 0, false
		}
	})

	p, q, ok := s.NextPQ()
	if !ok || p != 3 || q != 2 {
		t.Fatalf("got (%d,%d,%v), want (3,2,true)", p, q, ok)
	}

	p, q, ok = s.NextPQ()
	if !ok || p != 5 || q != 7 {
		t.Fatalf("got (%d,%d,%v), want (5,7,true)", p, q, ok)
	}

	_, _, ok = s.NextPQ()
	if ok {
		t.Fatalf("expected termination")
	}
}

func TestFuncGCFSource_WithGCFApproxFromPrefix(t *testing.T) {
	s := NewFuncGCFSource(func(i int) (int64, int64, bool) {
		switch i {
		case 0:
			return 1, 1, true
		case 1:
			return 2, 1, true
		case 2:
			return 2, 1, true
		default:
			return 0, 0, false
		}
	})

	got, err := GCFApproxFromPrefix(s, 3)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}

	want := mustRat(7, 5)
	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

// sources_test.go v7
