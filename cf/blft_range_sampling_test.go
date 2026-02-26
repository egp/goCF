// blft_range_sampling_test.go v1
package cf

import "testing"

func TestBLFTRange_EnclosesSampledInteriorPoints(t *testing.T) {
	// A reasonably general BLFT with stable denom on chosen rectangle.
	tform := NewBLFT(
		2, 1, -1, 3, // A B C D
		1, 2, 1, 5, // E F G H  (denom kept away from 0 on chosen ranges)
	)

	rx := MustRange(mustRat(1, 2), mustRat(5, 2)) // [0.5, 2.5]
	ry := MustRange(mustRat(1, 3), mustRat(7, 3)) // [0.333..., 2.333...]

	out, err := tform.ApplyBLFTRange(rx, ry)
	if err != nil {
		t.Fatal(err)
	}
	if !out.IsInside() {
		t.Fatalf("expected inside output range")
	}

	// Sample a small grid of rationals inside rx × ry.
	// Keep it deterministic and quick.
	xs := []Rational{
		rx.Lo,
		mustRat(1, 1),
		mustRat(3, 2),
		mustRat(2, 1),
		rx.Hi,
	}
	ys := []Rational{
		ry.Lo,
		mustRat(1, 1),
		mustRat(4, 3),
		mustRat(2, 1),
		ry.Hi,
	}

	for _, x := range xs {
		if !rx.Contains(x) {
			t.Fatalf("internal error: x=%v not in rx", x)
		}
		for _, y := range ys {
			if !ry.Contains(y) {
				t.Fatalf("internal error: y=%v not in ry", y)
			}
			z, err := tform.ApplyRat(x, y)
			if err != nil {
				t.Fatalf("ApplyRat failed at (x=%v,y=%v): %v", x, y, err)
			}
			if !out.Contains(z) {
				t.Fatalf("output range does not enclose sample: x=%v y=%v z=%v out=[%v,%v]",
					x, y, z, out.Lo, out.Hi)
			}
		}
	}
}

// blft_range_sampling_test.go v1
