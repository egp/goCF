// mvp_sources.go v12
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

func MVPApproxSnapshotFromApproxFunc(
	approxFn MVPFourOverPiApproxFunc,
	prefixTerms int,
) (GCFApprox, error) {
	if approxFn == nil {
		return GCFApprox{}, fmt.Errorf("MVPApproxSnapshotFromApproxFunc: nil approxFn")
	}
	if prefixTerms <= 0 {
		return GCFApprox{}, fmt.Errorf(
			"MVPApproxSnapshotFromApproxFunc: prefixTerms must be > 0, got %d",
			prefixTerms,
		)
	}

	x, err := approxFn(prefixTerms)
	if err != nil {
		return GCFApprox{}, err
	}

	r := NewRange(x, x, true, true)
	return GCFApprox{
		Convergent:  x,
		Range:       &r,
		PrefixTerms: prefixTerms,
	}, nil
}

func MVPRadicandDefaultFourOverPiSnapshot(prefixTerms int) (GCFApprox, error) {
	return MVPApproxSnapshotFromApproxFunc(MVPDefaultFourOverPiApproxFunc(), prefixTerms)
}

func MVPRadicandDefaultEApproxSnapshot(prefixTerms int) (GCFApprox, error) {
	return MVPDefaultEApproxSnapshot(prefixTerms)
}

func MVPRadicandScaledSquareOfFourOverPiApprox(fourOverPi GCFApprox) (GCFApprox, error) {
	sq, err := fourOverPi.Convergent.Mul(fourOverPi.Convergent)
	if err != nil {
		return GCFApprox{}, err
	}

	scale, err := MVPRadicandScaleFactorSnapshot()
	if err != nil {
		return GCFApprox{}, err
	}

	scaled, err := scale.Convergent.Mul(sq)
	if err != nil {
		return GCFApprox{}, err
	}

	r := NewRange(scaled, scaled, true, true)
	return GCFApprox{
		Convergent:  scaled,
		Range:       &r,
		PrefixTerms: fourOverPi.PrefixTerms,
	}, nil
}

func MVPRadicandAssembleFromSnapshots(
	fourOverPi GCFApprox,
	eApprox GCFApprox,
) (GCFApprox, error) {
	scaledSq, err := MVPRadicandScaledSquareOfFourOverPiApprox(fourOverPi)
	if err != nil {
		return GCFApprox{}, err
	}

	sum, err := scaledSq.Convergent.Add(eApprox.Convergent)
	if err != nil {
		return GCFApprox{}, err
	}

	r := NewRange(sum, sum, true, true)
	prefixTerms := fourOverPi.PrefixTerms
	if eApprox.PrefixTerms < prefixTerms || prefixTerms == 0 {
		prefixTerms = eApprox.PrefixTerms
	}

	return GCFApprox{
		Convergent:  sum,
		Range:       &r,
		PrefixTerms: prefixTerms,
	}, nil
}

func MVPRadicandAssembleSnapshotWithFourOverPiApprox(
	fourOverPiFn MVPFourOverPiApproxFunc,
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
) (GCFApprox, error) {
	fourOverPi, err := MVPApproxSnapshotFromApproxFunc(
		fourOverPiFn,
		fourOverPiPrefixTerms,
	)
	if err != nil {
		return GCFApprox{}, err
	}

	eApprox, err := MVPRadicandDefaultEApproxSnapshot(ePrefixTerms)
	if err != nil {
		return GCFApprox{}, err
	}

	return MVPRadicandAssembleFromSnapshots(fourOverPi, eApprox)
}

func MVPRadicandAssembleSnapshot(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
) (GCFApprox, error) {
	return MVPRadicandAssembleSnapshotWithFourOverPiApprox(
		MVPDefaultFourOverPiApproxFunc(),
		fourOverPiPrefixTerms,
		ePrefixTerms,
	)
}

func MVPRadicandAssembleConvergentWithFourOverPiApprox(
	fourOverPiFn MVPFourOverPiApproxFunc,
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
) (Rational, error) {
	radicand, err := MVPRadicandAssembleSnapshotWithFourOverPiApprox(
		fourOverPiFn,
		fourOverPiPrefixTerms,
		ePrefixTerms,
	)
	if err != nil {
		return Rational{}, err
	}

	return radicand.Convergent, nil
}

func MVPRadicandAssembleConvergent(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
) (Rational, error) {
	return MVPRadicandAssembleConvergentWithFourOverPiApprox(
		MVPDefaultFourOverPiApproxFunc(),
		fourOverPiPrefixTerms,
		ePrefixTerms,
	)
}

// Legacy expression-specific wrappers retained temporarily while tests migrate.

func MVPThreeOverPiSquaredPlusEApproxWithFourOverPiApprox(
	fourOverPiFn MVPFourOverPiApproxFunc,
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
) (Rational, error) {
	return MVPRadicandAssembleConvergentWithFourOverPiApprox(
		fourOverPiFn,
		fourOverPiPrefixTerms,
		ePrefixTerms,
	)
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

// Legacy expression-specific wrapper retained temporarily while tests migrate.
func MVPThreeOverPiSquaredPlusEApprox(fourOverPiPrefixTerms, ePrefixTerms int) (Rational, error) {
	return MVPRadicandAssembleConvergent(
		fourOverPiPrefixTerms,
		ePrefixTerms,
	)
}

func MVPExactScalarGCFSource(n int64) (GCFSource, int, error) {
	if n < 0 {
		return nil, 0, fmt.Errorf("MVPExactScalarGCFSource: negative n %d", n)
	}
	return NewSliceGCF([2]int64{n, 1}), 1, nil
}

func MVPExactScalarSnapshotFromSource(n int64) (GCFApprox, error) {
	src, prefixTerms, err := MVPExactScalarGCFSource(n)
	if err != nil {
		return GCFApprox{}, err
	}
	return GCFApproxFromPrefix(src, prefixTerms)
}

func MVPExactScalarSnapshot(n int64) (GCFApprox, error) {
	return MVPExactScalarSnapshotFromSource(n)
}

func MVPRadicandScaleFactorSnapshot() (GCFApprox, error) {
	three, err := MVPExactScalarSnapshotFromSource(3)
	if err != nil {
		return GCFApprox{}, err
	}
	sixteen, err := MVPExactScalarSnapshotFromSource(16)
	if err != nil {
		return GCFApprox{}, err
	}

	x, err := three.Convergent.Div(sixteen.Convergent)
	if err != nil {
		return GCFApprox{}, err
	}

	r := NewRange(x, x, true, true)
	return GCFApprox{
		Convergent:  x,
		Range:       &r,
		PrefixTerms: 1,
	}, nil
}

// mvp_sources.go v12
