// ulft_validate_test.go v1
package cf

import "testing"

func TestULFTValidate(t *testing.T) {
	tform := NewULFT(
		mustBig(1),
		mustBig(2),
		mustBig(3),
		mustBig(4),
	)

	if err := tform.Validate(); err != nil {
		t.Fatalf("unexpected validation failure: %v", err)
	}
}

func TestULFTDeterminant(t *testing.T) {
	tform := NewULFT(
		mustBig(1),
		mustBig(0),
		mustBig(0),
		mustBig(1),
	)

	d, err := tform.Determinant()
	if err != nil {
		t.Fatalf("determinant failed: %v", err)
	}

	if d != 1 {
		t.Fatalf("got %d want 1", d)
	}
}
