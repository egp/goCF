// mvp_sources.go v2
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

// MVP69DegreeGCFSource returns a GCF prefix source that, together with
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

// MVPThreeOverPiSquaredPlusEApprox returns a bounded-prefix rational approximation
// for:
//
//	3/pi^2 + e
//
// using the current canonical MVP sources:
//
//	(3/16) * (4/pi)^2 + e
func MVPThreeOverPiSquaredPlusEApprox(fourOverPiPrefixTerms, ePrefixTerms int) (Rational, error) {
	if fourOverPiPrefixTerms <= 0 {
		return Rational{}, fmt.Errorf(
			"MVPThreeOverPiSquaredPlusEApprox: fourOverPiPrefixTerms must be > 0, got %d",
			fourOverPiPrefixTerms,
		)
	}
	if ePrefixTerms <= 0 {
		return Rational{}, fmt.Errorf(
			"MVPThreeOverPiSquaredPlusEApprox: ePrefixTerms must be > 0, got %d",
			ePrefixTerms,
		)
	}

	fourOverPi, err := GCFSourceConvergent(MVPReciprocalPiGCFSource(), fourOverPiPrefixTerms)
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

// mvp_sources.go v2
