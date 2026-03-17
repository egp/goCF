// mvp_denominator_test.go v8
package cf

import (
	"fmt"
	"strings"
	"testing"
)

func mvpTestDenominatorBounds(
	sqrt5Policy SqrtPolicy2,
	angle Angle,
) (Range, error) {
	_ = sqrt5Policy // reserved for later tighter tanh(sqrt(5)) work

	if err := angle.Validate(); err != nil {
		return Range{}, err
	}
	if !angle.IsDegrees() {
		return Range{}, fmt.Errorf("mvpTestDenominatorBounds: angle must be expressed in degrees")
	}

	tanhR, err := TanhBoundsSpecialFromGCF2(AdaptCFToGCF(Sqrt5CF()))
	if err != nil {
		return Range{}, err
	}

	sinR, err := SinBoundsDegreesFromGCFPrefix2(
		MVP69DegreeGCFSource(),
		2,
	)
	if err != nil {
		return Range{}, err
	}

	lo, err := tanhR.Lo.Sub(sinR.Hi)
	if err != nil {
		return Range{}, err
	}
	hi, err := tanhR.Hi.Sub(sinR.Lo)
	if err != nil {
		return Range{}, err
	}

	return NewRange(lo, hi, true, true), nil
}

func mvpTestDenominatorBoundsDefault() (Range, error) {
	return mvpTestDenominatorBounds(
		DefaultSqrtPolicy2(),
		Degrees(mustRat(69, 1)),
	)
}

func mvpTestDenominatorApprox(
	sqrt5Policy SqrtPolicy2,
	angle Angle,
) (Rational, error) {
	r, err := mvpTestDenominatorBounds(sqrt5Policy, angle)
	if err != nil {
		return Rational{}, err
	}
	if r.Lo.Cmp(r.Hi) != 0 {
		return Rational{}, fmt.Errorf("mvpTestDenominatorApprox: bounded non-point result for %v", angle)
	}
	return r.Lo, nil
}

func mvpTestDenominatorApproxDefault() (Rational, error) {
	return mvpTestDenominatorApprox(
		DefaultSqrtPolicy2(),
		Degrees(mustRat(69, 1)),
	)
}

func TestMVPDenominatorBoundsDefault_UsesDegreesByDefault(t *testing.T) {
	got, err := mvpTestDenominatorBoundsDefault()
	if err != nil {
		t.Fatalf("mvpTestDenominatorBoundsDefault failed: %v", err)
	}

	want := NewRange(mustRat(11, 280), mustRat(7, 150), true, true)
	if got.Lo.Cmp(want.Lo) != 0 || got.Hi.Cmp(want.Hi) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
	if got.Contains(intRat(0)) {
		t.Fatalf("got %v want zero excluded", got)
	}
}

func TestMVPDenominatorBounds_RejectsRadiansForMVP(t *testing.T) {
	_, err := mvpTestDenominatorBounds(
		DefaultSqrtPolicy2(),
		Radians(mustRat(69, 1)),
	)
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "degrees") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMVPDenominatorBounds_Accepts69DegreeBoundAndExcludesZero(t *testing.T) {
	got, err := mvpTestDenominatorBounds(
		DefaultSqrtPolicy2(),
		Degrees(mustRat(69, 1)),
	)
	if err != nil {
		t.Fatalf("mvpTestDenominatorBounds failed: %v", err)
	}

	want := NewRange(mustRat(11, 280), mustRat(7, 150), true, true)
	if got.Lo.Cmp(want.Lo) != 0 || got.Hi.Cmp(want.Hi) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
	if got.Contains(intRat(0)) {
		t.Fatalf("got %v want zero excluded", got)
	}
}

func TestMVPDenominatorApprox_CurrentlyReportsBoundedNonPoint(t *testing.T) {
	_, err := mvpTestDenominatorApproxDefault()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "bounded non-point result") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Full target formula intentionally remains in test code for now.
// This test fixes only the denominator shape:
//
//	tanh(sqrt(5)) - sin(69°)
//

func TestMVPDenominatorBounds_UsesExactFiniteGCFPrefixFor69Degrees(t *testing.T) {
	got, err := SinBoundsDegreesFromGCFPrefix2(
		MVP69DegreeGCFSource(),
		2,
	)
	if err != nil {
		t.Fatalf("SinBoundsDegreesFromGCFPrefix2 failed: %v", err)
	}

	want, err := SinBoundsDegrees(Degrees(mustRat(69, 1)))
	if err != nil {
		t.Fatalf("SinBoundsDegrees failed: %v", err)
	}

	if got.Lo.Cmp(want.Lo) != 0 || got.Hi.Cmp(want.Hi) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPDenominatorBounds_NoLongerNeedsExactTailTrigEntryPointFor69Degrees(t *testing.T) {
	got, err := mvpTestDenominatorBoundsDefault()
	if err != nil {
		t.Fatalf("mvpTestDenominatorBoundsDefault failed: %v", err)
	}

	want := NewRange(mustRat(11, 280), mustRat(7, 150), true, true)
	if got.Lo.Cmp(want.Lo) != 0 || got.Hi.Cmp(want.Hi) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVP69DegreeGCFSource_Prefix2EvaluatesExactlyTo69(t *testing.T) {
	got, err := GCFSourceConvergent(MVP69DegreeGCFSource(), 2)
	if err != nil {
		t.Fatalf("GCFSourceConvergent failed: %v", err)
	}

	want := mustRat(69, 1)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestSinBoundsDegreesFromGCFPrefix2_MVP69DegreeGCFSource_IsExactAtPrefix2(t *testing.T) {
	got, err := SinBoundsDegreesFromGCFPrefix2(MVP69DegreeGCFSource(), 2)
	if err != nil {
		t.Fatalf("SinBoundsDegreesFromGCFPrefix2 failed: %v", err)
	}

	want, err := SinBoundsDegrees(Degrees(mustRat(69, 1)))
	if err != nil {
		t.Fatalf("SinBoundsDegrees failed: %v", err)
	}

	if got.Lo.Cmp(want.Lo) != 0 || got.Hi.Cmp(want.Hi) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

// mvp_denominator_test.go v8
