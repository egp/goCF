// fingerprint.go v1
package cf

import (
	"fmt"
	"strings"
)

// FingerprintULFT returns a human-readable, stable fingerprint for debugging.
// It canonicalizes the ULFT (divide out gcd, normalize sign) to reduce noise.
func FingerprintULFT(t ULFT, r Range) (string, error) {
	if err := t.Validate(); err != nil {
		return "", err
	}
	tc := canonULFT(t)

	return fmt.Sprintf(
		"ULFT[%s] x=%s",
		fmtULFT(tc),
		fmtRange(r),
	), nil
}

// FingerprintBLFT returns a human-readable, stable fingerprint for debugging.
// It canonicalizes the BLFT (divide out gcd, normalize sign) to reduce noise.
func FingerprintBLFT(t BLFT, rx, ry Range) (string, error) {
	// No Validate() method today; keep this lightweight and trust callers.
	tc := canonBLFT(t)

	return fmt.Sprintf(
		"BLFT[%s] x=%s y=%s",
		fmtBLFT(tc),
		fmtRange(rx),
		fmtRange(ry),
	), nil
}

func fmtULFT(t ULFT) string {
	return fmt.Sprintf("a=%d b=%d c=%d d=%d", t.A, t.B, t.C, t.D)
}

func fmtBLFT(t BLFT) string {
	// (Axy + Bx + Cy + D) / (Exy + Fx + Gy + H)
	return fmt.Sprintf(
		"A=%d B=%d C=%d D=%d E=%d F=%d G=%d H=%d",
		t.A, t.B, t.C, t.D, t.E, t.F, t.G, t.H,
	)
}

func fmtRange(r Range) string {
	kind := "inside"
	if r.IsOutside() {
		kind = "outside"
	}

	// Use a compact, explicit representation (including endpoint flags).
	var b strings.Builder
	b.WriteString("[")
	b.WriteString(r.Lo.String())
	b.WriteString(",")
	b.WriteString(r.Hi.String())
	b.WriteString("]")
	b.WriteString(fmt.Sprintf("{incLo=%t,incHi=%t,%s}", r.IncLo, r.IncHi, kind))
	return b.String()
}

func canonBLFT(t BLFT) BLFT {
	// Divide out gcd of all coefficients (if >1) and normalize sign:
	// make the first non-zero coefficient positive.
	g := fpGcd8(
		fpAbs(t.A), fpAbs(t.B), fpAbs(t.C), fpAbs(t.D),
		fpAbs(t.E), fpAbs(t.F), fpAbs(t.G), fpAbs(t.H),
	)
	if g > 1 {
		t.A /= g
		t.B /= g
		t.C /= g
		t.D /= g
		t.E /= g
		t.F /= g
		t.G /= g
		t.H /= g
	}

	sign := int64(0)
	for _, v := range []int64{t.A, t.B, t.C, t.D, t.E, t.F, t.G, t.H} {
		if v != 0 {
			sign = v
			break
		}
	}
	if sign < 0 {
		t.A = -t.A
		t.B = -t.B
		t.C = -t.C
		t.D = -t.D
		t.E = -t.E
		t.F = -t.F
		t.G = -t.G
		t.H = -t.H
	}
	return t
}

func fpAbs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func fpGcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func fpGcd8(a, b, c, d, e, f, g, h int64) int64 {
	x := fpGcd(a, b)
	x = fpGcd(x, c)
	x = fpGcd(x, d)
	x = fpGcd(x, e)
	x = fpGcd(x, f)
	x = fpGcd(x, g)
	x = fpGcd(x, h)
	return x
}

// fingerprint.go v1
