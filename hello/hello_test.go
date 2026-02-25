package hello

import "testing"

func TestGreet(t *testing.T) {
	got := Greet("World")
	want := "Hello, World!"

	if got != want {
		t.Fatalf("Greet() = %q, want %q", got, want)
	}
}
