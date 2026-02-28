// blft_stream_cycle_test.go v1
package cf

import "testing"

func TestBLFTStream_CycleDetection_Smoke(t *testing.T) {
	// z = x + y
	// numerator: (x+y) => A=0 B=1 C=1 D=0
	// denom: 1 => E=0 F=0 G=0 H=1
	tform := NewBLFT(0, 1, 1, 0, 0, 0, 0, 1)

	x := mustRat(1, 2) // 0.5
	y := mustRat(1, 3) // 0.333...
	// z = 5/6 = [0; 1, 5]
	want := NewRationalCF(mustRat(5, 6))

	opts := BLFTStreamOptions{
		MaxFinalizeDigits: 32,
		DetectCycles:      true,
		MaxRepeats:        2,
		HistorySize:       32,
	}

	got := NewBLFTStream(tform, NewRationalCF(x), NewRationalCF(y), opts)

	// Compare a short prefix.
	for i := 0; i < 8; i++ {
		wa, wok := want.Next()
		ga, gok := got.Next()

		if wok != gok {
			t.Fatalf("termination mismatch at i=%d: wantOk=%v gotOk=%v", i, wok, gok)
		}
		if !wok {
			break
		}
		if wa != ga {
			t.Fatalf("digit mismatch at i=%d: want=%d got=%d", i, wa, ga)
		}
	}

	if err := got.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

// blft_stream_cycle_test.go v1
