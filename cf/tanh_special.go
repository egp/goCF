// tanh_special.go v3
package cf

// TanhBoundsSqrt5 returns a conservative inside range for tanh(sqrt(5)).
//
// Current v3 rule:
//
// First bracket sqrt(5):
//
//	11/5 < sqrt(5) < 9/4
//
// since:
//
//	(11/5)^2 = 121/25 = 4.84 < 5 < 81/16 = (9/4)^2
//
// tanh is increasing on positive inputs, so:
//
//	tanh(11/5) < tanh(sqrt(5)) < tanh(9/4)
//
// Lower bound:
//   - tanh(x) = (e^(2x)-1)/(e^(2x)+1)
//   - to prove tanh(11/5) > 39/40, it suffices to prove e^(22/5) > 79
//   - the first nine positive Taylor terms for e^(22/5) already exceed 79
//
// Upper bound:
//   - to prove tanh(9/4) < 49/50, it suffices to prove e^(9/2) < 99
//   - summing the exponential series for e^(9/2) through the 9th term and
//     bounding the remaining tail geometrically gives a total < 99
//
// Therefore we use the certified enclosure:
//
//	tanh(sqrt(5)) in [39/40, 49/50]
func TanhBoundsSqrt5() Range {
	return NewRange(mustRat(39, 40), mustRat(49, 50), true, true)
}

// tanh_special.go v3
