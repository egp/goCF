// tanh_special.go v2
package cf

// TanhBoundsSqrt5 returns a conservative inside range for tanh(sqrt(5)).
//
// Current v2 rule:
//   - sqrt(5) > 2
//   - tanh is increasing on positive inputs
//   - therefore tanh(sqrt(5)) > tanh(2)
//
// For x > 0, tanh(x) = (e^(2x)-1)/(e^(2x)+1). At x=2:
//
//	tanh(2) = (e^4 - 1)/(e^4 + 1)
//
// Using the Taylor lower bound:
//
//	e^4 > 1 + 4 + 4^2/2! + 4^3/3! + 4^4/4! + 4^5/5! + 4^6/6! > 43
//
// hence:
//
//	tanh(2) > (43 - 1)/(43 + 1) = 21/22
//
// Also tanh(sqrt(5)) < 1 for positive input, so we use [21/22, 1].
func TanhBoundsSqrt5() Range {
	return NewRange(mustRat(21, 22), mustRat(1, 1), true, true)
}

// tanh_special.go v2
