// sources_sqrt2_mul_test.go v2
package cf

import (
	"testing"
	"time"
)

func nextBLFTWithTimeout(t *testing.T, s *BLFTStream, d time.Duration) (int64, bool) {
	t.Helper()

	type result struct {
		a  int64
		ok bool
	}
	ch := make(chan result, 1)

	go func() {
		a, ok := s.Next()
		ch <- result{a: a, ok: ok}
	}()

	select {
	case r := <-ch:
		return r.a, r.ok
	case <-time.After(d):
		t.Fatalf("BLFTStream.Next() timed out after %v; likely stuck refining without progress; err=%v", d, s.Err())
		return 0, false
	}
}

func TestBLFTStream_Sqrt2TimesSqrt2_Equals2Exactly(t *testing.T) {
	// Multiply via BLFT: z = x*y
	// (1*x*y + 0*x + 0*y + 0) / (0*x*y + 0*x + 0*y + 1)
	mul := NewBLFT(1, 0, 0, 0, 0, 0, 0, 1)

	s := NewBLFTStream(
		mul,
		Sqrt2CF(),
		Sqrt2CF(),
		BLFTStreamOptions{
			DetectCycles:       false,
			MaxRefinesPerDigit: 1000,
			MaxTotalRefines:    10000,
		},
	)

	a0, ok := nextBLFTWithTimeout(t, s, 2*time.Second)
	if !ok {
		t.Fatalf("expected first digit, got termination; err=%v", s.Err())
	}
	if a0 != 2 {
		t.Fatalf("expected first digit 2, got %d (err=%v)", a0, s.Err())
	}

	// Should terminate immediately for exact 2 = [2].
	a1, ok := nextBLFTWithTimeout(t, s, 2*time.Second)
	if ok {
		t.Fatalf("expected termination after [2], but got another digit %d", a1)
	}
	if err := s.Err(); err != nil {
		t.Fatalf("expected clean termination, got err=%v", err)
	}
}

// sources_sqrt2_mul_test.go v2
