// bounder_test.go v1
package cf

import "testing"

func TestBounder_NoTerms(t *testing.T) {
	b := NewBounder()
	_, ok, err := b.Range()
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatalf("expected ok=false before any terms ingested")
	}
}

func TestBounder_SingleTermRange(t *testing.T) {
	// Prefix [3; ...] implies x in [3,4] (since tail r in [1,∞): 3 + 1/r)
	b := NewBounder()
	if err := b.Ingest(3); err != nil {
		t.Fatal(err)
	}

	rng, ok, err := b.Range()
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatalf("expected ok=true after ingest")
	}

	if rng.Lo.Cmp(mustRat(3, 1)) != 0 || rng.Hi.Cmp(mustRat(4, 1)) != 0 {
		t.Fatalf("got [%v,%v], want [3/1,4/1]", rng.Lo, rng.Hi)
	}
}

func TestBounder_RangeShrinksAsWeIngest(t *testing.T) {
	// sqrt(2) = [1;2,2,2,...]
	b := NewBounder()

	// [1; ...] => [1,2]
	_ = b.Ingest(1)
	r1, ok, err := b.Range()
	if err != nil || !ok {
		t.Fatalf("range1 err=%v ok=%v", err, ok)
	}
	w1, _ := r1.RefineMetric()

	// [1;2; ...] => between 4/3 and 3/2
	_ = b.Ingest(2)
	r2, ok, err := b.Range()
	if err != nil || !ok {
		t.Fatalf("range2 err=%v ok=%v", err, ok)
	}
	w2, _ := r2.RefineMetric()

	// [1;2,2; ...]
	_ = b.Ingest(2)
	r3, ok, err := b.Range()
	if err != nil || !ok {
		t.Fatalf("range3 err=%v ok=%v", err, ok)
	}
	w3, _ := r3.RefineMetric()

	// Expect strict shrink: w3 < w2 < w1
	if w2.Cmp(w1) >= 0 {
		t.Fatalf("expected w2 < w1, got w1=%v w2=%v", w1, w2)
	}
	if w3.Cmp(w2) >= 0 {
		t.Fatalf("expected w3 < w2, got w2=%v w3=%v", w2, w3)
	}
}

func TestBounder_KnownPrefixMatchesConvergentAndMediant(t *testing.T) {
	// Prefix [1;2; ...] => conv=3/2, mediant=(3+1)/(2+1)=4/3, so range [4/3,3/2]
	b := NewBounder()
	_ = b.Ingest(1)
	_ = b.Ingest(2)

	rng, ok, err := b.Range()
	if err != nil || !ok {
		t.Fatalf("err=%v ok=%v", err, ok)
	}

	wantLo := mustRat(4, 3)
	wantHi := mustRat(3, 2)

	if rng.Lo.Cmp(wantLo) != 0 || rng.Hi.Cmp(wantHi) != 0 {
		t.Fatalf("got [%v,%v], want [%v,%v]", rng.Lo, rng.Hi, wantLo, wantHi)
	}
}

func TestBounder_FinishCollapsesToExactRational(t *testing.T) {
	// 355/113 = [3;7,16]
	b := NewBounder()
	_ = b.Ingest(3)
	_ = b.Ingest(7)
	_ = b.Ingest(16)
	b.Finish()

	rng, ok, err := b.Range()
	if err != nil || !ok {
		t.Fatalf("err=%v ok=%v", err, ok)
	}

	want := mustRat(355, 113)
	if rng.Lo.Cmp(want) != 0 || rng.Hi.Cmp(want) != 0 {
		t.Fatalf("got [%v,%v], want exact [%v,%v]", rng.Lo, rng.Hi, want, want)
	}
}

// bounder_test.go v1
