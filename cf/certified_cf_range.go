// certified_cf_range.go v1
package cf

import "fmt"

// CertifyCFDigitsFromRange repeatedly certifies continued-fraction digits from a
// conservative inside range.
//
// It emits digits while the current range certifies a unique floor. After each
// emitted digit d, it maps the range through the remainder transform:
//
//	z' = 1 / (z - d)
//
// and continues.
//
// This is a generic certification helper for operator-style streams.
func CertifyCFDigitsFromRange(r Range, maxDigits int) ([]int64, Range, error) {
	if !r.IsInside() {
		return nil, Range{}, fmt.Errorf("CertifyCFDigitsFromRange: requires inside range; got %v", r)
	}
	if maxDigits < 0 {
		return nil, Range{}, fmt.Errorf("CertifyCFDigitsFromRange: negative maxDigits %d", maxDigits)
	}
	if maxDigits == 0 {
		return []int64{}, r, nil
	}

	cur := r
	out := make([]int64, 0, maxDigits)

	for len(out) < maxDigits {
		lo, hi, err := cur.FloorBounds()
		if err != nil {
			return nil, Range{}, err
		}
		if lo != hi {
			return out, cur, nil
		}

		out = append(out, lo)

		next, err := CertifiedRemainderRange(cur, lo)
		if err != nil {
			// This is normal at the end of an exact rational or when the current
			// certified range does not support another remainder step. Return the
			// digits certified so far and the current range.
			return out, cur, nil
		}
		cur = next
	}

	return out, cur, nil
}

// certified_cf_range.go v1
