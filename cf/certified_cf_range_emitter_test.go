// certified_cf_range_emitter_test.go v1
package cf

import "testing"

func TestCertifiedCFRangeEmitter_Sqrt2LikeRange_EmitsPrefix(t *testing.T) {
	r := NewRange(mustRat(181, 128), mustRat(362, 255), true, true)

	e, err := NewCertifiedCFRangeEmitter(r)
	if err != nil {
		t.Fatalf("NewCertifiedCFRangeEmitter failed: %v", err)
	}

	got := collectTerms(e, 8)
	if err := e.Err(); err != nil {
		t.Fatalf("emitter error: %v", err)
	}

	wantPrefix := []int64{1, 2}
	if len(got) < len(wantPrefix) {
		t.Fatalf("len(got)=%d want at least %d got=%v", len(got), len(wantPrefix), got)
	}
	for i := range wantPrefix {
		if got[i] != wantPrefix[i] {
			t.Fatalf("got[%d]=%d want=%d full=%v", i, got[i], wantPrefix[i], got)
		}
	}
}

func TestCertifiedCFRangeEmitter_ExactInteger_StopsAfterOneDigit(t *testing.T) {
	r := NewRange(mustRat(2, 1), mustRat(2, 1), true, true)

	e, err := NewCertifiedCFRangeEmitter(r)
	if err != nil {
		t.Fatalf("NewCertifiedCFRangeEmitter failed: %v", err)
	}

	a0, ok := e.Next()
	if !ok {
		t.Fatalf("expected first digit")
	}
	if a0 != 2 {
		t.Fatalf("got %d want 2", a0)
	}

	if _, ok := e.Next(); ok {
		t.Fatalf("expected clean exhaustion after exact integer")
	}
	if err := e.Err(); err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
}

func TestCertifiedCFRangeEmitter_UncertifiedRangeEmitsNothing(t *testing.T) {
	r := NewRange(mustRat(1, 1), mustRat(2, 1), true, true)

	e, err := NewCertifiedCFRangeEmitter(r)
	if err != nil {
		t.Fatalf("NewCertifiedCFRangeEmitter failed: %v", err)
	}

	if _, ok := e.Next(); ok {
		t.Fatalf("expected no certified digit")
	}
	if err := e.Err(); err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
}

// certified_cf_range_emitter_test.go v1
