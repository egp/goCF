// mvp_sources.go v10
package cf

import "fmt"

// MVPFourOverPiFamily identifies a bounded 4/pi source family for MVP work.
type MVPFourOverPiFamily string

const (
	MVPFourOverPiFamilyBrouncker MVPFourOverPiFamily = "brouncker"
	MVPFourOverPiFamilyLambert   MVPFourOverPiFamily = "lambert"
)

// MVPDefaultFourOverPiFamily is the canonical 4/pi source family for the MVP path.
const MVPDefaultFourOverPiFamily = MVPFourOverPiFamilyBrouncker

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

// MVP69DegreeGCFSource returns the canonical exact finite GCF source for 69°.
//
// Construction:
//   - 68 + 1/1 = 69
//   - encoded directly as a finite two-term GCF
func MVP69DegreeGCFSource() GCFSource {
	return NewSliceGCF(
		[2]int64{68, 1},
		[2]int64{1, 1},
	)
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
// by first approximating pi/4 from Lambert's GCF prefix machinery, then taking
// the reciprocal of the convergent exactly.
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

// MVPFourOverPiApproxFuncForFamily returns the 4/pi approximation function for
// the requested source family.
func MVPFourOverPiApproxFuncForFamily(family MVPFourOverPiFamily) (MVPFourOverPiApproxFunc, error) {
	switch family {
	case MVPFourOverPiFamilyBrouncker:
		return MVPFourOverPiApproxBrouncker, nil
	case MVPFourOverPiFamilyLambert:
		return MVPFourOverPiApproxLambert, nil
	default:
		return nil, fmt.Errorf("MVPFourOverPiApproxFuncForFamily: unsupported family %q", family)
	}
}

// MVPDefaultFourOverPiApproxFunc returns the canonical 4/pi approximation
// function for the MVP path.
func MVPDefaultFourOverPiApproxFunc() MVPFourOverPiApproxFunc {
	return MVPFourOverPiApproxBrouncker
}

type MVPGCFSourceFunc func() GCFSource

func MVPDefaultFourOverPiSourceFunc() MVPGCFSourceFunc {
	return MVPReciprocalPiGCFSource
}

func MVPDefaultESourceFunc() MVPGCFSourceFunc {
	return MVPEGCFSource
}

func MVPNilSafeGCFSourceFromFunc(srcFn MVPGCFSourceFunc) (GCFSource, error) {
	if srcFn == nil {
		return nil, fmt.Errorf("MVPNilSafeGCFSourceFromFunc: nil srcFn")
	}
	src := srcFn()
	if src == nil {
		return nil, fmt.Errorf("MVPNilSafeGCFSourceFromFunc: nil source")
	}
	return src, nil
}

func MVPApproxSnapshotFromSourceFunc(
	srcFn MVPGCFSourceFunc,
	prefixTerms int,
) (GCFApprox, error) {
	if prefixTerms <= 0 {
		return GCFApprox{}, fmt.Errorf(
			"MVPApproxSnapshotFromSourceFunc: prefixTerms must be > 0, got %d",
			prefixTerms,
		)
	}

	src, err := MVPNilSafeGCFSourceFromFunc(srcFn)
	if err != nil {
		return GCFApprox{}, err
	}

	return GCFApproxFromPrefix(src, prefixTerms)
}

func MVPDefaultFourOverPiApproxSnapshot(prefixTerms int) (GCFApprox, error) {
	return MVPApproxSnapshotFromSourceFunc(
		MVPDefaultFourOverPiSourceFunc(),
		prefixTerms,
	)
}

func MVPDefaultEApproxSnapshot(prefixTerms int) (GCFApprox, error) {
	return MVPApproxSnapshotFromSourceFunc(
		MVPDefaultESourceFunc(),
		prefixTerms,
	)
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
	a, err := MVPApproxSnapshotFromSourceFunc(srcFn, fourOverPiPrefixTerms)
	if err != nil {
		return Rational{}, err
	}
	return a.Convergent, nil
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
// using the current canonical MVP source family.
func MVPThreeOverPiSquaredPlusEApprox(fourOverPiPrefixTerms, ePrefixTerms int) (Rational, error) {
	return MVPThreeOverPiSquaredPlusEApproxWithFourOverPiApprox(
		MVPDefaultFourOverPiApproxFunc(),
		fourOverPiPrefixTerms,
		ePrefixTerms,
	)
}

// mvp_sources.go v10
