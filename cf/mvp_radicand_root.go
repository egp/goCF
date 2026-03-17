// mvp_radicand_root.go v2
package cf

import "fmt"

const (
	MVPRadicandDefaultFourOverPiPrefixTerms = 6
	MVPRadicandDefaultEPrefixTerms          = 8
	MVPRadicandSnapshotTerms                = 64
)

func MVPRadicandConvergent(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
) (Rational, error) {
	return MVPRadicandAssembleConvergent(
		fourOverPiPrefixTerms,
		ePrefixTerms,
	)
}

func MVPRadicandSnapshot(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	snapshotTerms int,
) (GCFApprox, error) {
	if snapshotTerms <= 0 {
		return GCFApprox{}, fmt.Errorf(
			"MVPRadicandSnapshot: snapshotTerms must be > 0, got %d",
			snapshotTerms,
		)
	}

	return MVPRadicandAssembleSnapshot(
		fourOverPiPrefixTerms,
		ePrefixTerms,
		snapshotTerms,
	)
}

func MVPRadicandRootValueFromSnapshot(
	a GCFApprox,
	sqrtPolicy SqrtPolicy2,
) (Rational, error) {
	return SqrtApproxFromGCFApproxRangeSeed2(a, sqrtPolicy)
}

func MVPRadicandRootValue(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	sqrtPolicy SqrtPolicy2,
) (Rational, error) {
	return MVPRadicandRootValueWithSnapshotTerms(
		fourOverPiPrefixTerms,
		ePrefixTerms,
		sqrtPolicy,
		MVPRadicandSnapshotTerms,
	)
}

func MVPRadicandRootValueWithSnapshotTerms(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	sqrtPolicy SqrtPolicy2,
	snapshotTerms int,
) (Rational, error) {
	a, err := MVPRadicandSnapshot(
		fourOverPiPrefixTerms,
		ePrefixTerms,
		snapshotTerms,
	)
	if err != nil {
		return Rational{}, err
	}

	return MVPRadicandRootValueFromSnapshot(a, sqrtPolicy)
}

func MVPRadicandRootValueDefault(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
) (Rational, error) {
	return MVPRadicandRootValue(
		fourOverPiPrefixTerms,
		ePrefixTerms,
		DefaultSqrtPolicy2(),
	)
}

func MVPRadicandRootValueCurrentDefault() (Rational, error) {
	return MVPRadicandRootValue(
		MVPRadicandDefaultFourOverPiPrefixTerms,
		MVPRadicandDefaultEPrefixTerms,
		DefaultSqrtPolicy2(),
	)
}

func MVPRadicandRootCF(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	sqrtPolicy SqrtPolicy2,
) (ContinuedFraction, error) {
	value, err := MVPRadicandRootValue(fourOverPiPrefixTerms, ePrefixTerms, sqrtPolicy)
	if err != nil {
		return nil, err
	}
	return NewRationalCF(value), nil
}

func MVPRadicandRootCFDefault(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
) (ContinuedFraction, error) {
	return MVPRadicandRootCF(
		fourOverPiPrefixTerms,
		ePrefixTerms,
		DefaultSqrtPolicy2(),
	)
}

func MVPRadicandRootTerms(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	sqrtPolicy SqrtPolicy2,
	digits int,
) ([]int64, error) {
	if digits < 0 {
		return nil, fmt.Errorf("MVPRadicandRootTerms: negative digits %d", digits)
	}

	cf, err := MVPRadicandRootCF(fourOverPiPrefixTerms, ePrefixTerms, sqrtPolicy)
	if err != nil {
		return nil, err
	}
	return collectTerms(cf, digits), nil
}

func MVPRadicandRootTermsDefault(
	fourOverPiPrefixTerms int,
	ePrefixTerms int,
	digits int,
) ([]int64, error) {
	return MVPRadicandRootTerms(
		fourOverPiPrefixTerms,
		ePrefixTerms,
		DefaultSqrtPolicy2(),
		digits,
	)
}

// MVP rooted-radicand shape note:
//
//	The full target formula should remain assembled in tests for now.
//	This production helper intentionally stops at:
//
//	    sqrt(3/pi^2 + e)
//
// mvp_radicand_root.go v2
