// sqrt_terms.go v1
package cf

import "fmt"

// collectTerms reads up to n CF terms from src.
// It stops early if the source terminates.
func collectTerms(src ContinuedFraction, n int) []int64 {
	if n <= 0 {
		return []int64{}
	}
	out := make([]int64, 0, n)
	for i := 0; i < n; i++ {
		a, ok := src.Next()
		if !ok {
			break
		}
		out = append(out, a)
	}
	return out
}

// SqrtApproxTerms computes a bounded rational Newton approximation to sqrt(x),
// converts that exact rational approximation to a finite CF, and returns up to
// digits terms of that CF.
//
// This is a convenience bridge from exact-rational sqrt approximation to CF
// output terms. It is not yet a true streaming sqrt operator.
func SqrtApproxTerms(x, seed Rational, steps, digits int) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("SqrtApproxTerms: negative digits %d", digits)
	}

	approx, err := SqrtApproxRational(x, seed, steps)
	if err != nil {
		return nil, err
	}

	return collectTerms(NewRationalCF(approx), digits), nil
}

// sqrt_terms.go v1
