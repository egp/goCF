// ulft_stream_infinite_test.go v3
package cf

import "testing"

func TestULFTStream_Infinite_IdentityOnSqrt2(t *testing.T) {
	id := NewULFT(bi(1), bi(0), bi(0), bi(1))
	s := NewULFTStream(id, Sqrt2CF(), ULFTStreamOptions{
		DetectCycles: true,
		MaxRepeats:   3,
	})

	got, err := takeN(s, 10)
	if err != nil {
		t.Fatalf("takeN: %v", err)
	}
	if s.Err() != nil {
		t.Fatalf("stream error: %v", s.Err())
	}

	want := []int64{1, 2, 2, 2, 2, 2, 2, 2, 2, 2}
	if !equalSlice(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestULFTStream_Infinite_Plus1OnSqrt2(t *testing.T) {
	// T(x) = x + 1
	plus1 := NewULFT(bi(1), bi(1), bi(0), bi(1))

	// sqrt(2) = [1; (2)] so (sqrt(2)+1) = [2; (2)]
	s := NewULFTStream(plus1, Sqrt2CF(), ULFTStreamOptions{
		DetectCycles: true,
		MaxRepeats:   3,
	})

	got, err := takeN(s, 10)
	if err != nil {
		t.Fatalf("takeN: %v", err)
	}
	if s.Err() != nil {
		t.Fatalf("stream error: %v", s.Err())
	}

	want := []int64{2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
	if !equalSlice(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestULFTStream_Infinite_ReciprocalOnSqrt2(t *testing.T) {
	// T(x) = 1/x
	recip := NewULFT(bi(0), bi(1), bi(1), bi(0))

	// 1/sqrt(2) = [0; 1, (2)]  (i.e., 0;1,2,2,2,...)
	s := NewULFTStream(recip, Sqrt2CF(), ULFTStreamOptions{
		DetectCycles: true,
		MaxRepeats:   3,
	})

	got, err := takeN(s, 10)
	if err != nil {
		t.Fatalf("takeN: %v", err)
	}
	if s.Err() != nil {
		t.Fatalf("stream error: %v", s.Err())
	}

	want := []int64{0, 1, 2, 2, 2, 2, 2, 2, 2, 2}
	if !equalSlice(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// ulft_stream_infinite_test.go v3
