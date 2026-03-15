// mvp_sources.go v7
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

// MVPFourOverPiApproxWithSource returns a bounded rational approximation of 4/pi
// using the supplied source factory and bounded prefix.
func MVPFourOverPiApproxWithSource(
	srcFn func() GCFSource,
	fourOverPiPrefixTerms int,
) (Rational, error) {
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
// using a supplied 4/pi source:
//
//	(3/16) * (4/pi)^2 + e
func MVPThreeOverPiSquaredPlusEApproxWithFourOverPiSource(
	fourOverPiSrcFn func() GCFSource,
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
) (Rational, error) {
	if fourOverPiPrefixTerms <= 0 {
		return Rational{}, fmt.Errorf(
			"MVPThreeOverPiSquaredPlusEApproxWithFourOverPiSource: fourOverPiPrefixTerms must be > 0, got %d",
			fourOverPiPrefixTerms,
		)
	}
	if ePrefixTerms <= 0 {
		return Rational{}, fmt.Errorf(
			"MVPThreeOverPiSquaredPlusEApproxWithFourOverPiSource: ePrefixTerms must be > 0, got %d",
			ePrefixTerms,
		)
	}

	fourOverPi, err := MVPFourOverPiApproxWithSource(fourOverPiSrcFn, fourOverPiPrefixTerms)
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

// MVPThreeOverPiSquaredPlusEApprox returns a bounded-prefix rational approximation
// for:
//
//	3/pi^2 + e
//
// using the current canonical MVP sources:
//
//	(3/16) * (4/pi)^2 + e
func MVPThreeOverPiSquaredPlusEApprox(fourOverPiPrefixTerms, ePrefixTerms int) (Rational, error) {
	return MVPThreeOverPiSquaredPlusEApproxWithFourOverPiSource(
		MVPReciprocalPiGCFSource,
		fourOverPiPrefixTerms,
		ePrefixTerms,
	)
}

// mvp_sources.go v7
