// sin_degrees_test.go v1
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

func TestSinApproxDegrees_RejectsRadians(t *testing.T) {
	_, err := SinApproxDegrees(Radians(mustRat(1, 1)))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "degrees") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSinApproxDegrees_UnsupportedDegreeIsStub(t *testing.T) {
	_, err := SinApproxDegrees(Degrees(mustRat(69, 1)))
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "not implemented") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// sin_degrees_test.go v1
