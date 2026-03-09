// test_helpers_test.go v3
package cf

import (
	"fmt"
	"math/big"
)

// takeN reads exactly n terms from a ContinuedFraction.
// Returns an error if the source terminates early.
func takeN(src ContinuedFraction, n int) ([]int64, error) {
	out := make([]int64, 0, n)
	for i := 0; i < n; i++ {
		a, ok := src.Next()
		if !ok {
			return nil, fmt.Errorf("takeN: source terminated early at i=%d (wanted n=%d)", i, n)
		}
		out = append(out, a)
	}
	return out, nil
}

func bi(n int64) *big.Int { return big.NewInt(n) }

func mustBig(n int64) *big.Int {
	return big.NewInt(n)
}

// EOF test_helpers_test.go v3
