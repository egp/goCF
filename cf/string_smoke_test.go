// string_smoke_test.go v1
package cf

import (
	"math/big"
	"strings"
	"testing"
)

func TestDiagBLFTString_Smoke(t *testing.T) {
	d := NewDiagBLFT(
		big.NewInt(1), big.NewInt(2), big.NewInt(3),
		big.NewInt(4), big.NewInt(5), big.NewInt(6),
	)

	got := d.String()
	if got == "" {
		t.Fatalf("expected non-empty string")
	}
	for _, want := range []string{"1", "2", "3", "4", "5", "6"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q to contain %q", got, want)
		}
	}
}

func TestULFTString_Smoke(t *testing.T) {
	u := NewULFT(
		big.NewInt(7), big.NewInt(8),
		big.NewInt(9), big.NewInt(10),
	)

	got := u.String()
	if got == "" {
		t.Fatalf("expected non-empty string")
	}
	for _, want := range []string{"7", "8", "9", "10"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected %q to contain %q", got, want)
		}
	}
}

func TestRangeString_InsideAndOutside(t *testing.T) {
	inside := NewRange(mustRat(1, 2), mustRat(3, 4), true, false)
	gotInside := inside.String()
	if gotInside == "" {
		t.Fatalf("expected non-empty inside string")
	}
	for _, want := range []string{"1/2", "3/4", "inside"} {
		if !strings.Contains(gotInside, want) {
			t.Fatalf("expected %q to contain %q", gotInside, want)
		}
	}

	outside := NewRange(mustRat(3, 4), mustRat(1, 2), false, true)
	gotOutside := outside.String()
	if gotOutside == "" {
		t.Fatalf("expected non-empty outside string")
	}
	for _, want := range []string{"3/4", "1/2", "outside"} {
		if !strings.Contains(gotOutside, want) {
			t.Fatalf("expected %q to contain %q", gotOutside, want)
		}
	}
}

// string_smoke_test.go v1
