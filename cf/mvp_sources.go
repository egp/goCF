// mvp_sources.go v8
package cf

import "fmt"

// MVPReciprocalPiGCFSource returns the canonical MVP source for reciprocal-pi work.
//
// Current choice:
//   - Brouncker 4/pi
//
// Rationale:
//   - the target expression needs 3/pi^2
//   - 3/pi^2 = (3/16) * (4/pi)^2
//   - this avoids introducing a separate reciprocal-of-pi step for MVP
func MVPReciprocalPiGCFSource() GCFSource {
	return NewBrouncker4OverPiGCFSource()
}

// MVPEGCFSource returns the canonical MVP source for e.
func MVPEGCFSource() GCFSource {
	return NewECFGSource()
}

// MVP69DegreeGCFSource returns a finite GCF prefix source that, together with
// MVP69DegreeTail(), evaluates exactly to 69.
//
// Construction:
//   - 68 + 1/1 = 69
func MVP69DegreeGCFSource() GCFSource {
	return NewSliceGCF([2]int64{68, 1})
}

// MVP69DegreeTail returns the exact tail used with MVP69DegreeGCFSource().
func MVP69DegreeTail() Rational {
	return mustRat(1, 1)
}

// MVPFourOverPiApproxFunc is a subexpression-level hook that returns a bounded
// rational approximation to 4/pi from a chosen source family.
type MVPFourOverPiApproxFunc func(prefixTerms int) (Rational, error)

// MVPFourOverPiApproxBrouncker returns a bounded rational approximation of 4/pi
// from the canonical Brouncker source.
func MVPFourOverPiApproxBrouncker(prefixTerms int) (Rational, error) {
	if prefixTerms <= 0 {
		return Rational{}, fmt.Errorf(
			"MVPFourOverPiApproxBrouncker: prefixTerms must be > 0, got %d",
			prefixTerms,
		)
	}
	return GCFSourceConvergent(NewBrouncker4OverPiGCFSource(), prefixTerms)
}

// MVPFourOverPiApproxLambert returns a bounded rational approximation of 4/pi
// by first approximating pi/4 from Lambert's GCF with exact tail evidence, then
// taking the reciprocal exactly.
//
// Current Lambert tail rule for pi/4:
//   - use exact tail 1
func MVPFourOverPiApproxLambert(prefixTerms int) (Rational, error) {
	if prefixTerms <= 0 {
		return Rational{}, fmt.Errorf(
			"MVPFourOverPiApproxLambert: prefixTerms must be > 0, got %d",
			prefixTerms,
		)
	}

	piOver4Approx, err := LambertPiOver4ApproxFromPrefix(prefixTerms)
	if err != nil {
		return Rational{}, err
	}
	if piOver4Approx.Convergent.Cmp(intRat(0)) == 0 {
		return Rational{}, fmt.Errorf("MVPFourOverPiApproxLambert: reciprocal of zero")
	}

	return intRat(1).Div(piOver4Approx.Convergent)
}

// MVPThreeOverPiSquaredPlusEApproxWithFourOverPiApprox returns a bounded-prefix
// rational approximation for:
//
//	3/pi^2 + e
//
// from a supplied bounded 4/pi approximation function:
//
//	(3/16) * (4/pi)^2 + e
func MVPThreeOverPiSquaredPlusEApproxWithFourOverPiApprox(
	fourOverPiFn MVPFourOverPiApproxFunc,
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
) (Rational, error) {
	if fourOverPiFn == nil {
		return Rational{}, fmt.Errorf("MVPThreeOverPiSquaredPlusEApproxWithFourOverPiApprox: nil fourOverPiFn")
	}
	if fourOverPiPrefixTerms <= 0 {
		return Rational{}, fmt.Errorf(
			"MVPThreeOverPiSquaredPlusEApproxWithFourOverPiApprox: fourOverPiPrefixTerms must be > 0, got %d",
			fourOverPiPrefixTerms,
		)
	}
	if ePrefixTerms <= 0 {
		return Rational{}, fmt.Errorf(
			"MVPThreeOverPiSquaredPlusEApproxWithFourOverPiApprox: ePrefixTerms must be > 0, got %d",
			ePrefixTerms,
		)
	}

	fourOverPi, err := fourOverPiFn(fourOverPiPrefixTerms)
	if err != nil {
		return Rational{}, err
	}
	eApprox, err := GCFSourceConvergent(MVPEGCFSource(), ePrefixTerms)
	if err != nil {
		return Rational{}, err
	}

	fourOverPiSq, err := fourOverPi.Mul(fourOverPi)
	if err != nil {
		return Rational{}, err
	}

	scale := mustRat(3, 16)
	threeOverPiSq, err := scale.Mul(fourOverPiSq)
	if err != nil {
		return Rational{}, err
	}

	sum, err := threeOverPiSq.Add(eApprox)
	if err != nil {
		return Rational{}, err
	}
	return sum, nil
}

// MVPFourOverPiApproxWithSource returns a bounded rational approximation of 4/pi
// using the supplied GCF source factory and bounded prefix.
//
// Transitional note:
//   - this hook fits source families that are natively 4/pi as GCFSource
//   - Lambert parity uses MVPFourOverPiApproxLambert instead
func MVPFourOverPiApproxWithSource(
	srcFn func() GCFSource,
	fourOverPiPrefixTerms int,
) (Rational, error) {
	if srcFn == nil {
		return Rational{}, fmt.Errorf("MVPFourOverPiApproxWithSource: nil srcFn")
	}
	if fourOverPiPrefixTerms <= 0 {
		return Rational{}, fmt.Errorf(
			"MVPFourOverPiApproxWithSource: fourOverPiPrefixTerms must be > 0, got %d",
			fourOverPiPrefixTerms,
		)
	}
	return GCFSourceConvergent(srcFn(), fourOverPiPrefixTerms)
}

// MVPThreeOverPiSquaredPlusEApproxWithFourOverPiSource returns a bounded-prefix
// rational approximation for:
//
//	3/pi^2 + e
//
// using a supplied 4/pi GCF source:
//
//	(3/16) * (4/pi)^2 + e
//
// Transitional note:
//   - this hook fits source families that are natively 4/pi as GCFSource
//   - Lambert parity uses MVPThreeOverPiSquaredPlusEApproxWithFourOverPiApprox instead
func MVPThreeOverPiSquaredPlusEApproxWithFourOverPiSource(
	fourOverPiSrcFn func() GCFSource,
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
) (Rational, error) {
	return MVPThreeOverPiSquaredPlusEApproxWithFourOverPiApprox(
		func(prefixTerms int) (Rational, error) {
			return MVPFourOverPiApproxWithSource(fourOverPiSrcFn, prefixTerms)
		},
		fourOverPiPrefixTerms,
		ePrefixTerms,
	)
}

// MVPThreeOverPiSquaredPlusEApprox returns a bounded-prefix rational approximation
// for:
//
//	3/pi^2 + e
//
// using the current canonical MVP sources:
//
//	(3/16) * (4/pi)^2 + e
func MVPThreeOverPiSquaredPlusEApprox(fourOverPiPrefixTerms, ePrefixTerms int) (Rational, error) {
	return MVPThreeOverPiSquaredPlusEApproxWithFourOverPiApprox(
		MVPFourOverPiApproxBrouncker,
		fourOverPiPrefixTerms,
		ePrefixTerms,
	)
}

// mvp_sources.go v8
