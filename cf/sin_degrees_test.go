// sin_degrees_test.go v4
package cf

import (
	"strings"
	"testing"
)

func TestSinApproxDegrees_ExactTable(t *testing.T) {
	cases := []struct {
		name  string
		angle Angle
		want  Rational
	}{
		{"0", Degrees(mustRat(0, 1)), mustRat(0, 1)},
		{"30", Degrees(mustRat(30, 1)), mustRat(1, 2)},
		{"90", Degrees(mustRat(90, 1)), mustRat(1, 1)},
		{"180", Degrees(mustRat(180, 1)), mustRat(0, 1)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := SinApproxDegrees(tc.angle)
			if err != nil {
				t.Fatalf("SinApproxDegrees failed: %v", err)
			}
			if got.Cmp(tc.want) != 0 {
				t.Fatalf("got %v want %v", got, tc.want)
			}
		})
	}
}

func TestSinBoundsDegrees_69IsFurtherTightenedConservativeInsideRange(t *testing.T) {
	got, err := SinBoundsDegrees(Degrees(mustRat(69, 1)))
	if err != nil {
		t.Fatalf("SinBoundsDegrees failed: %v", err)
	}

	want := NewRange(mustRat(14, 15), mustRat(131, 140), true, true)
	if got.Lo.Cmp(want.Lo) != 0 || got.Hi.Cmp(want.Hi) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
	if !got.IsInside() {
		t.Fatalf("expected inside range, got %v", got)
	}
}

func TestSinApproxDegrees_69IsCurrentlyBoundedNotPoint(t *testing.T) {
	_, err := SinApproxDegrees(Degrees(mustRat(69, 1)))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "bounded non-point result") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSinApproxDegrees_RejectsRadians(t *testing.T) {
	_, err := SinApproxDegrees(Radians(mustRat(1, 1)))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "degrees") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSinBoundsDegrees_RejectsRadians(t *testing.T) {
	_, err := SinBoundsDegrees(Radians(mustRat(1, 1)))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "degrees") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// sin_degrees_test.go v4
