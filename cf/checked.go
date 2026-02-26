// checked.go v1
package cf

import (
	"errors"
	"math"
	"math/bits"
)

var ErrOverflow = errors.New("int64 overflow")

func add64(a, b int64) (int64, bool) {
	ua, ub := uint64(a), uint64(b)
	sum, _ := bits.Add64(ua, ub, 0)
	r := int64(sum)
	// overflow if a and b have same sign and result has different sign
	if (a^b) >= 0 && (a^r) < 0 {
		return 0, false
	}
	return r, true
}

func sub64(a, b int64) (int64, bool) {
	ua, ub := uint64(a), uint64(b)
	diff, _ := bits.Sub64(ua, ub, 0)
	r := int64(diff)
	// overflow if a and b have different sign and result has different sign than a
	if (a^b) < 0 && (a^r) < 0 {
		return 0, false
	}
	return r, true
}

func mul64(a, b int64) (int64, bool) {
	// Fast-path zeros.
	if a == 0 || b == 0 {
		return 0, true
	}
	// Special-case MinInt64 * -1 which overflows.
	if (a == math.MinInt64 && b == -1) || (b == math.MinInt64 && a == -1) {
		return 0, false
	}

	neg := (a < 0) != (b < 0)
	ua := absU64(a)
	ub := absU64(b)

	hi, lo := bits.Mul64(ua, ub)
	if hi != 0 {
		return 0, false
	}

	// lo now fits in uint64; must fit in signed int64 magnitude.
	if !neg {
		if lo > uint64(math.MaxInt64) {
			return 0, false
		}
		return int64(lo), true
	}

	// negative result must fit down to MinInt64 (magnitude MaxInt64+1)
	if lo > uint64(math.MaxInt64)+1 {
		return 0, false
	}
	if lo == uint64(math.MaxInt64)+1 {
		return math.MinInt64, true
	}
	return -int64(lo), true
}

// mulAdd64 computes a*b + c with overflow detection.
func mulAdd64(a, b, c int64) (int64, bool) {
	p, ok := mul64(a, b)
	if !ok {
		return 0, false
	}
	return add64(p, c)
}

// absU64 returns |x| as uint64, handling MinInt64 correctly.
func absU64(x int64) uint64 {
	if x >= 0 {
		return uint64(x)
	}
	// For MinInt64, -x overflows int64; compute via two's complement.
	return uint64(^x) + 1
}

// checked.go v1
