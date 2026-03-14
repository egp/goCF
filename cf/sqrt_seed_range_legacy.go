// sqrt_seed_range_legacy.go v3
package cf

// DefaultSqrtSeed is the legacy compatibility wrapper around SqrtSeedDefault.
func DefaultSqrtSeed(x Rational) (Rational, error) {
	return SqrtSeedDefault(x)
}

// DefaultSqrtSeedFromRange is the legacy compatibility wrapper around SqrtSeedFromRange.
func DefaultSqrtSeedFromRange(r Range) (Rational, error) {
	return SqrtSeedFromRange(r)
}

// DefaultSqrtSeedFromCFPrefix is the legacy compatibility wrapper around SqrtSeedFromCFPrefix.
func DefaultSqrtSeedFromCFPrefix(src ContinuedFraction, prefixTerms int) (Rational, error) {
	return SqrtSeedFromCFPrefix(src, prefixTerms)
}

// sqrt_seed_range_legacy.go v3
