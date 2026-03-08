// sources_sqrt2_mul_test.go v5
package cf

import "testing"

func TestDiagBLFTStream_Sqrt2Square_Equals2Exactly(t *testing.T) {
	sq := NewDiagBLFT(
		bi(1), bi(0), bi(0),
		bi(0), bi(0), bi(1),
	)

	s := NewDiagBLFTStream(sq, Sqrt2CF(), DiagBLFTStreamOptions{})

	a0, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, got termination; err=%v", s.Err())
	}
	if a0 != 2 {
		t.Fatalf("expected first digit 2, got %d (err=%v)", a0, s.Err())
	}

	a1, ok := s.Next()
	if ok {
		t.Fatalf("expected termination after [2], got extra digit %d", a1)
	}
	if err := s.Err(); err != nil {
		t.Fatalf("expected clean termination, got err=%v", err)
	}
}

// sources_sqrt2_mul_test.go v5
