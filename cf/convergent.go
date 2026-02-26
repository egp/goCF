// convergent.go v2
package cf

// RationalFromTerms computes the exact rational value of a finite
// continued fraction [a0; a1, a2, ..., an].
//
// It evaluates via the recurrence:
//
//	h[-2]=0, h[-1]=1
//	k[-2]=1, k[-1]=0
//	h[i] = a[i]*h[i-1] + h[i-2]
//	k[i] = a[i]*k[i-1] + k[i-2]
//
// This version is checked: returns ErrOverflow on int64 overflow.
func RationalFromTerms(terms []int64) (Rational, error) {
	if len(terms) == 0 {
		return NewRational(0, 1)
	}

	// h[-2], h[-1]
	hm2 := int64(0)
	hm1 := int64(1)

	// k[-2], k[-1]
	km2 := int64(1)
	km1 := int64(0)

	for _, a := range terms {
		h, ok := mulAdd64(a, hm1, hm2) // a*hm1 + hm2
		if !ok {
			return Rational{}, ErrOverflow
		}
		k, ok := mulAdd64(a, km1, km2) // a*km1 + km2
		if !ok {
			return Rational{}, ErrOverflow
		}

		hm2, hm1 = hm1, h
		km2, km1 = km1, k
	}

	return NewRational(hm1, km1)
}

// convergent.go v2
