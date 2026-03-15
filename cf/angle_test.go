// angle_test.go v1
package cf

import "testing"

func TestAngle_Degrees(t *testing.T) {
	a := Degrees(mustRat(69, 1))

	if !a.IsDegrees() {
		t.Fatalf("expected degrees")
	}
	if a.Value().Cmp(mustRat(69, 1)) != 0 {
		t.Fatalf("got %v want 69", a.Value())
	}
}

func TestAngle_Radians(t *testing.T) {
	a := Radians(mustRat(1, 2))

	if a.IsDegrees() {
		t.Fatalf("expected radians")
	}
	if a.Value().Cmp(mustRat(1, 2)) != 0 {
		t.Fatalf("got %v want 1/2", a.Value())
	}
}

func TestAngle_Validate(t *testing.T) {
	if err := Degrees(mustRat(69, 1)).Validate(); err != nil {
		t.Fatalf("degrees validate failed: %v", err)
	}
	if err := Radians(mustRat(1, 1)).Validate(); err != nil {
		t.Fatalf("radians validate failed: %v", err)
	}
}

// angle_test.go v1
