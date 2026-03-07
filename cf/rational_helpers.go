// rational_helpers.go v1
package cf

import "fmt"

// intRat constructs an exact integer rational n/1.
// Panics only if your Rational constructor rejects q=1 (it shouldn't).
func intRat(n int64) Rational {
	r, err := NewRational(n, 1)
	if err != nil {
		panic(err)
	}
	return r
}

// Convenience: make a Rational from ints and panic on error (tests only).
func mustRat(p, q int64) Rational {
	r, err := NewRational(p, q)
	if err != nil {
		panic(fmt.Sprintf("bad rational %d/%d: %v", p, q, err))
	}
	return r
}

// rational_helpers.go v1
