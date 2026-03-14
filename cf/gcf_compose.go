// gcf_compose.go v2
package cf

import (
	"fmt"
	"math/big"
)

// ComposeGCFIntoULFTBounded ingests up to maxIngestTerms from src into base.
//
// Contract:
//   - If src exhausts before the bound is hit, returns (tFinal, ingested, true, nil).
//   - If the bound is hit before src exhausts, returns (zero, ingested, false, err).
//   - If any ingested term is invalid, returns (zero, ingestedSoFar, false, err).
//
// Bound semantics:
//   - maxIngestTerms < 0 : unlimited
//   - maxIngestTerms = 0 : no ingestion allowed
//   - maxIngestTerms > 0 : maximum permitted ingested source terms
func ComposeGCFIntoULFTBounded(
	base ULFT,
	src GCFSource,
	maxIngestTerms int,
) (ULFT, int, bool, error) {
	cur := base

	ingested, err := ingestGCFBounded(
		"ComposeGCFIntoULFTBounded",
		src,
		maxIngestTerms,
		func(p, q int64) error {
			var ierr error
			cur, ierr = cur.IngestGCF(p, q)
			return ierr
		},
	)
	if err != nil {
		return ULFT{}, ingested, false, err
	}
	return cur, ingested, true, nil
}

// ApplyComposedGCFULFTToTailExact boundedly ingests src into base, requires
// source exhaustion within the bound, then applies the resulting ULFT exactly
// to the supplied tail rational.
func ApplyComposedGCFULFTToTailExact(
	base ULFT,
	src GCFSource,
	tail Rational,
	maxIngestTerms int,
) (Rational, int, error) {
	composed, ingested, exhausted, err := ComposeGCFIntoULFTBounded(base, src, maxIngestTerms)
	if err != nil {
		return Rational{}, ingested, err
	}
	if !exhausted {
		return Rational{}, ingested, fmt.Errorf(
			"ApplyComposedGCFULFTToTailExact: internal: bounded compose returned !exhausted without error",
		)
	}

	y, err := composed.ApplyRat(tail)
	if err != nil {
		return Rational{}, ingested, err
	}
	return y, ingested, nil
}

// EvalGCFWithTailExact evaluates a bounded finite GCF prefix with an explicit
// exact tail rational.
//
// It computes x represented by:
//
//	p0 + q0/(p1 + q1/(... tail ...))
//
// by first composing the prefix into a ULFT and then applying that ULFT to tail.
func EvalGCFWithTailExact(
	src GCFSource,
	tail Rational,
	maxIngestTerms int,
) (Rational, int, error) {
	id := NewULFT(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(1))

	composed, ingested, exhausted, err := ComposeGCFIntoULFTBounded(id, src, maxIngestTerms)
	if err != nil {
		return Rational{}, ingested, err
	}
	if !exhausted {
		return Rational{}, ingested, fmt.Errorf(
			"EvalGCFWithTailExact: internal: bounded compose returned !exhausted without error",
		)
	}

	x, err := composed.ApplyRat(tail)
	if err != nil {
		return Rational{}, ingested, err
	}
	return x, ingested, nil
}

// gcf_compose.go v2
