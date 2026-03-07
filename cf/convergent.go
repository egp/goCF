// convergent.go v3
package cf

import "math/big"

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
// v3: arbitrary precision (no ErrOverflow from recurrence growth).
func RationalFromTerms(terms []int64) (Rational, error) {
	if len(terms) == 0 {
		return NewRational(0, 1)
	}

	// big.Int recurrence
	hm2 := big.NewInt(0) // h[-2]
	hm1 := big.NewInt(1) // h[-1]
	km2 := big.NewInt(1) // k[-2]
	km1 := big.NewInt(0) // k[-1]

	for _, a := range terms {
		ai := big.NewInt(a)

		// h = a*hm1 + hm2
		h := new(big.Int).Mul(ai, hm1)
		h.Add(h, hm2)

		// k = a*km1 + km2
		k := new(big.Int).Mul(ai, km1)
		k.Add(k, km2)

		hm2, hm1 = hm1, h
		km2, km1 = km1, k
	}

	return newRationalBig(hm1, km1)
}

// convergent.go v3
