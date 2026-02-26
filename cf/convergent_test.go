// convergent_test.go v1
package cf

import "testing"

func TestRationalFromTerms_Basic(t *testing.T) {
	// [1;2] = 3/2
	r, err := RationalFromTerms([]int64{1, 2})
	if err != nil {
		t.Fatal(err)
	}
	want := mustRat(3, 2)
	if r.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", r, want)
	}

	// [3;7,16] = 355/113
	r, err = RationalFromTerms([]int64{3, 7, 16})
	if err != nil {
		t.Fatal(err)
	}
	want = mustRat(355, 113)
	if r.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", r, want)
	}

	// [-2;1,2] = -4/3  (floor convention)
	r, err = RationalFromTerms([]int64{-2, 1, 2})
	if err != nil {
		t.Fatal(err)
	}
	want = mustRat(-4, 3)
	if r.Cmp(want) != 0 {
		t.Fatalf("got %v, want %v", r, want)
	}
}

func TestRational_CF_RoundTrip(t *testing.T) {
	testCases := []Rational{
		mustRat(3, 2),
		mustRat(355, 113),
		mustRat(-4, 3),
		mustRat(-1, 2),
		mustRat(7, 5),
	}

	for _, original := range testCases {
		terms := collectCF(NewRationalCF(original))
		back, err := RationalFromTerms(terms)
		if err != nil {
			t.Fatal(err)
		}
		if back.Cmp(original) != 0 {
			t.Fatalf("roundtrip failed: start=%v terms=%v back=%v",
				original, terms, back)
		}
	}
}

// convergent_test.go v1
