// binary_stream_contract_test.go v1
package cf

import "testing"

func TestBinaryStreamContract_BLFTStream_IsClassified(t *testing.T) {
	s := NewBLFTStream(
		NewBLFT(0, 1, 1, 0, 0, 0, 0, 1),
		NewRationalCF(mustRat(1, 2)),
		NewRationalCF(mustRat(1, 3)),
		BLFTStreamOptions{MaxFinalizeDigits: 16},
	)

	got, ok := classifyBinaryStream(s)
	if !ok {
		t.Fatalf("expected BLFTStream to be classified")
	}
	want := binaryStreamClass{
		Operator: binaryOperatorUnknown,
		Input:    binaryInputCF,
		Progress: binaryProgressProgressiveCertified,
	}
	if got != want {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestBinaryStreamContract_GCFBLFTStream_IsClassified(t *testing.T) {
	s := NewGCFBLFTStreamWithTails(
		NewBLFT(0, 1, 1, 0, 0, 0, 0, 1),
		NewSliceGCF([2]int64{3, 2}), mustRat(11, 1),
		NewSliceGCF([2]int64{2, 3}), mustRat(7, 1),
		GCFBLFTStreamOptions{MaxXIngestTerms: 8, MaxYIngestTerms: 8},
	)

	got, ok := classifyBinaryStream(s)
	if !ok {
		t.Fatalf("expected GCFBLFTStream to be classified")
	}
	want := binaryStreamClass{
		Operator: binaryOperatorUnknown,
		Input:    binaryInputGCF,
		Progress: binaryProgressExactCollapse,
	}
	if got != want {
		t.Fatalf("got %+v want %+v", got, want)
	}
}

func TestBinaryStreamContract_GCFDiagStream_IsClassified(t *testing.T) {
	s := NewGCFDiagStreamWithTail(
		NewDiagBLFT(
			mustBig(0), mustBig(1), mustBig(0),
			mustBig(0), mustBig(0), mustBig(1),
		),
		NewSliceGCF([2]int64{3, 2}),
		mustRat(11, 1),
		GCFULFTStreamOptions{MaxIngestTerms: 8},
	)

	got, ok := classifyBinaryStream(s)
	if !ok {
		t.Fatalf("expected GCFDiagStream to be classified")
	}
	want := binaryStreamClass{
		Operator: binaryOperatorUnknown,
		Input:    binaryInputGCF,
		Progress: binaryProgressExactCollapse,
	}
	if got != want {
		t.Fatalf("got %+v want %+v", got, want)
	}
}
