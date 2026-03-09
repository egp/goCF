# goCF API Spec Draft v1

## Scope

goCF provides exact rational arithmetic, continued-fraction sources and streams,
LFT/BLFT-based transform machinery, bounded sqrt approximation, and early-stage
generalized continued fraction (GCF) ingestion.

This draft focuses on public interfaces and their current semantics.

## Core numeric type

### Rational
Exact arbitrary-precision rational.

Key operations:
- NewRational(p, q int64) (Rational, error)
- String() string
- Add/Sub/Mul/Div
- Cmp

Notes:
- denominator is normalized positive
- arithmetic is exact
- Div rejects division by zero

## Regular continued fractions

### ContinuedFraction
Interface:
- Next() (int64, bool)

Meaning:
- returns next regular CF term
- bool=false means exhausted

### Basic sources
- NewSliceCF(...int64) *SliceCF
- NewPeriodicCF(prefix []int64, period []int64) *PeriodicCF
- PhiCF() ContinuedFraction
- Sqrt2CF(), Sqrt3CF(), Sqrt5CF(), Sqrt6CF(), Sqrt7CF()

Optional metadata:
- QuadraticRadicalSource
  - Radicand() (int64, bool)

## Exact rational to CF

### RationalCF
- NewRationalCF(r Rational) *RationalCF

Meaning:
- expands an exact rational into its finite regular CF

## Range and enclosure support

### Range
Represents an interval-like object:
- NewRange(lo, hi Rational, incLo, incHi bool) Range
- IsInside / IsOutside
- Contains / ContainsZero
- RefineMetric
- FloorBounds
- ApplyULFT

Notes:
- “inside” means Lo <= Hi
- “outside” means union-of-rays semantics
- some transform APIs currently require inside ranges only

### Bounder
Incrementally ingests regular CF terms and maintains:
- current convergent
- current enclosure range

Public methods:
- NewBounder() *Bounder
- Ingest(a int64) error
- Finish()
- HasValue() bool
- Convergent() (Rational, error)
- Range() (Range, bool, error)

Notes:
- after Finish, Range collapses to an exact point
- before Finish, Range is a conservative inside enclosure for regular CF prefixes

## Unary and binary transforms

### ULFT
Represents:
- (A*x + B) / (C*x + D)

Key methods:
- NewULFT(a, b, c, d *big.Int) ULFT
- ApplyRat(x Rational) (Rational, error)
- String() string
- IngestGCF(p, q int64) (ULFT, error)

### BLFT
Represents:
- (A*x*y + B*x + C*y + D) / (E*x*y + F*x + G*y + H)

Key methods:
- NewBLFT(...)
- ApplyRat(x, y Rational) (Rational, error)

### DiagBLFT
Diagonal specialization:
- (A*x^2 + B*x + C) / (D*x^2 + E*x + F)

Key methods:
- NewDiagBLFT(...)
- DiagFromBLFT(t BLFT) DiagBLFT
- ApplyRat(x Rational) (Rational, error)
- ApplyRange(r Range) (Range, error)
- String() string
- IngestGCF(p, q int64) (DiagBLFT, error)

## Streaming transform engines

### ULFTStream
Transforms a regular CF source through a ULFT.

Public surface:
- NewULFTStream(t ULFT, src ContinuedFraction, opts ULFTStreamOptions) *ULFTStream
- Next() (int64, bool)
- Err() error

Notes:
- uses range-based safe-digit logic
- supports optional cycle detection
- supports refine guards
- returns false on clean exhaustion or first error; inspect Err()

### BLFTStream
Transforms two regular CF sources through a BLFT.

Public surface:
- NewBLFTStream(t BLFT, xs, ys ContinuedFraction, opts BLFTStreamOptions) *BLFTStream
- Next() (int64, bool)
- Err() error

Notes:
- supports optional rational finalization
- supports optional cycle detection
- supports refine guards
- exact integer short-circuit is implemented for some cases

### DiagBLFTStream
Transforms one regular CF source through a diagonal BLFT.

Public surface:
- NewDiagBLFTStream(t DiagBLFT, src ContinuedFraction, opts DiagBLFTStreamOptions) *DiagBLFTStream
- Next() (int64, bool)
- Err() error

Notes:
- includes a narrow algebraic shortcut for sqrt(n)-like sources under x^2 + k forms
- current interval support is deliberately limited

## Sqrt support

Production layout:
- sqrt_newton.go
- sqrt_api.go
- sqrt_cf.go
- sqrt_range.go

### Exact and Newton-based primitives
- RationalSqrtExact(x Rational) (Rational, bool, error)
- NewtonSqrtStep(x, y Rational) (Rational, error)
- NewtonSqrtIterates(x, seed Rational, steps int) ([]Rational, error)
- SqrtResidual / SqrtResidualAbs

### Bounded rational sqrt approximation
- SqrtApproxRational(...)
- SqrtApproxRationalUntilExact(...)
- SqrtApproxRationalUntilExactDefault(...)
- SqrtApproxRationalUntilResidual(...)
- SqrtApproxRationalUntilResidualDefault(...)

### Policy
- type SqrtPolicy
- DefaultSqrtPolicy() SqrtPolicy
- (SqrtPolicy) Validate() error

Fields:
- MaxSteps int
- Tol Rational
- Seed *Rational

### Top-level sqrt APIs
- SqrtApprox(x Rational) (Rational, error)
- SqrtApproxCF(x Rational) (ContinuedFraction, error)
- SqrtApproxTermsAuto(x Rational, digits int) ([]int64, error)

Policy-driven variants:
- SqrtApproxWithPolicy
- SqrtApproxCFWithPolicy
- SqrtApproxTermsWithPolicy

Compatibility wrappers:
- SqrtApproxWithSeedAndPolicy
- SqrtApproxCFWithSeedAndPolicy
- SqrtApproxTermsWithSeedAndPolicy

### CF/source-facing sqrt helpers
- DefaultSqrtSeed
- DefaultSqrtSeedFromRange
- DefaultSqrtSeedFromCFPrefix
- NewSqrtApproxCF
- NewSqrtApproxCFDefault
- NewSqrtApproxCFUntilResidual
- NewSqrtApproxCFUntilResidualDefault
- SqrtApproxTerms
- SqrtApproxTermsDefault
- SqrtApproxTermsUntilResidual
- SqrtApproxTermsUntilResidualDefault

### CF prefix snapshots
- CFApprox
- CFApproxFromPrefix
- ApproxFromCFPrefix

### Source-based sqrt bridging
- NewSqrtApproxCFFromSource
- NewSqrtApproxCFFromSourceDefault
- SqrtApproxTermsFromSource
- SqrtApproxTermsFromSourceDefault
- NewSqrtApproxCFFromApproxRangeSeed
- NewSqrtApproxCFFromSourceRangeSeed
- NewSqrtApproxCFFromSourceRangeSeedDefault
- SqrtApproxTermsFromApproxRangeSeed
- SqrtApproxTermsFromSourceRangeSeed
- SqrtApproxTermsFromSourceRangeSeedDefault
- NewSqrtApproxCFFromApproxRangeMidpoint
- NewSqrtApproxCFFromSourceRangeMidpoint
- NewSqrtApproxCFFromSourceRangeMidpointDefault
- SqrtApproxTermsFromApproxRangeMidpoint
- SqrtApproxTermsFromSourceRangeMidpoint
- SqrtApproxTermsFromSourceRangeMidpointDefault

### Sqrt range helpers
- SqrtRangeExact
- SqrtRangeExactFromCFApprox
- SqrtRangeHeuristic
- SqrtRangeHeuristicFromCFApprox

Warnings:
- heuristic sqrt-range helpers are not proof-safe conservative enclosures
- bounded sqrt APIs are approximation tools, not true streaming exact-real sqrt yet

## Generalized continued fractions (GCF)

Convention used throughout:
- x = p + q/x'
- q > 0 required by ingestion logic

### GCFSource
Interface:
- NextPQ() (int64, int64, bool)

### GCF sources
- NewSliceGCF(...)
- NewPeriodicGCF(prefix, period)
- NewFuncGCFSource(fn)
- AdaptCFToGCF(src ContinuedFraction) GCFSource
- NewECFGSource()
- NewBrouncker4OverPiGCFSource()
- NewLambertPiOver4GCFSource()
- NewUnitPArithmeticQGCFSource(startQ, step)

### Finite evaluation and composition
- ComposeGCFToULFT(src GCFSource) (ULFT, error)
- EvaluateFiniteGCF(src GCFSource) (Rational, error)

### GCFBounder
- NewGCFBounder() *GCFBounder
- IngestPQ(p, q int64) error
- Finish()
- HasValue() bool
- Convergent() (Rational, error)
- Range() (Range, bool, error)

Warnings:
- current Range is a point-range placeholder, not yet a true conservative enclosure for unfinished infinite GCFs

### GCF ingestion helpers
- IngestAllGCF(src GCFSource) (*GCFBounder, error)
- IngestGCFPrefix(src GCFSource, prefixTerms int) (*GCFBounder, error)

### GCF snapshots and inspection
- GCFApprox
- GCFApproxFromPrefix
- GCFApproxCF
- GCFApproxTerms
- GCFSourceTerms
- GCFSourceConvergent

## Current warnings and limitations

- GCF support is currently ingest-first; emit support is not implemented
- GCFBounder does not yet provide conservative infinite-prefix enclosure semantics
- sqrt support is bounded and partially heuristic in source/range-facing layers
- some diagonal and transform interval mappings are intentionally narrow
- streaming engines terminate with bool=false on both clean exhaustion and error; call Err()

## Recommended usage pattern today

For regular CF work:
- use ContinuedFraction sources + Bounder + ULFT/BLFT streams

For bounded sqrt experimentation:
- use SqrtApproxWithPolicy or CF/source-facing sqrt helpers

For GCF experimentation:
- use GCFSource + GCFApproxFromPrefix + GCFSourceTerms
- compare bounded convergents as rationals and as regular CF term lists

## Documentation work still to do

Future user guide should add:
- worked examples
- transform ingestion diagrams
- warnings around heuristic APIs
- recommended debugging workflow
- comparison examples for regular CF vs GCF
- examples for Brouncker, Lambert, e, and future pi-related generators