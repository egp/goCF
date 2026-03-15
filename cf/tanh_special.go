// tanh_special.go v1
package cf

// TanhBoundsSqrt5 returns a conservative inside range for tanh(sqrt(5)).
//
// Current MVP rule:
//   - sqrt(5) > 0
//   - for positive input, tanh(x) is strictly between 0 and 1
//   - we use the conservative closed enclosure [0,1]
//
// This is intentionally loose and should be tightened later.
func TanhBoundsSqrt5() Range {
	return NewRange(mustRat(0, 1), mustRat(1, 1), true, true)
}

// tanh_special.go v1
