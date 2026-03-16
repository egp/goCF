# UserGuide.md

# goCF User Guide

## 1. Overview

goCF is a Go package for exact arithmetic and transform-based experimentation
with continued fractions.

It currently supports four closely related areas:

- exact arbitrary-precision rational arithmetic
- ordinary continued fractions (CF)
- bounded square-root approximation over exact rationals and CF-derived inputs
- generalized continued fractions (GCF), currently focused on ingestion and
  bounded prefix analysis

The package is aimed at mathematically careful experimentation. Some parts are
already exact and mature. Other parts are intentionally narrow, bounded, or
heuristic and should be used with that in mind.

## 2. Core ideas

The package revolves around a few central concepts:

- `Rational`: exact arbitrary-precision rational values
- `ContinuedFraction`: a stream of regular CF terms
- `Range`: an interval-like enclosure used for safe digit logic
- `Bounder`: ingests regular CF terms and maintains convergents and enclosures
- `ULFT`, `BLFT`, `DiagBLFT`: linear / bilinear transform machinery
- `GCFSource`: a stream of generalized continued-fraction `(p,q)` terms

Two different notions of “streaming” appear in this package:

- source streaming: values are generated term-by-term
- transform streaming: digits are emitted only when they are safe to emit

## 3. Recommended starting points

If you are new to the package, start here:

- exact rationals:
  - `NewRational`
- finite and periodic regular CF sources:
  - `NewSliceCF`
  - `NewPeriodicCF`
  - `PhiCF`, `Sqrt2CF`, `Sqrt3CF`, `Sqrt5CF`, `Sqrt6CF`, `Sqrt7CF`
- exact rational to regular CF:
  - `NewRationalCF`
- prefix convergents and enclosures:
  - `Bounder`
  - `CFApproxFromPrefix`
- bounded sqrt:
  - `SqrtApprox`
  - `SqrtApproxCF`
  - `SqrtApproxTermsAuto`
- GCF experimentation:
  - `NewSliceGCF`
  - `GCFApproxFromPrefix`
  - `GCFSourceConvergent`
  - `GCFSourceTerms`

## 4. Exact rationals

### What `Rational` is

`Rational` is an exact arbitrary-precision rational number.

Important properties:

- denominator is normalized positive
- arithmetic is exact
- values are reduced by `big.Rat`

### Main API

- `NewRational(p, q int64) (Rational, error)`
- `Add`, `Sub`, `Mul`, `Div`
- `Cmp`
- `String`

### Notes

Use `Rational` whenever correctness matters more than floating-point speed.
This package is built around exact symbolic arithmetic, not binary floating
approximations.

## 5. Regular continued fractions

### `ContinuedFraction`

A `ContinuedFraction` is any source that provides:

- `Next() (int64, bool)`

It returns the next regular CF term, or `false` when the finite source ends.

### Built-in sources

Finite source:

- `NewSliceCF(...)`

Periodic / infinite source:

- `NewPeriodicCF(prefix, period)`

Named periodic sources:

- `PhiCF()`
- `Sqrt2CF()`
- `Sqrt3CF()`
- `Sqrt5CF()`
- `Sqrt6CF()`
- `Sqrt7CF()`

Some sources also implement `QuadraticRadicalSource`, which currently provides
optional metadata for positive square roots of integer radicands.

## 6. Rational to regular CF

`NewRationalCF(r)` expands an exact rational into its finite regular CF.

This is useful for:

- debugging exact intermediate values
- comparing approximations as CF term lists
- inspecting GCF prefix convergents through the regular CF pipeline

## 7. Prefix ingestion and enclosures for regular CF

### `Bounder`

`Bounder` ingests regular CF terms and provides:

- current convergent
- current enclosure range
- exact point range after `Finish()`

Main methods:

- `NewBounder()`
- `Ingest(a)`
- `Finish()`
- `HasValue()`
- `Convergent()`
- `Range()`

### `CFApprox`

`CFApproxFromPrefix(src, prefixTerms)` returns a bundled prefix snapshot
containing:

- exact convergent
- current `Range`
- prefix length used

This is the preferred bounded regular-CF snapshot API.

## 8. Range model

### `Range`

`Range` represents either:

- an inside interval, where `Lo <= Hi`
- an outside union-of-rays object, where `Lo > Hi`

Main API includes:

- `NewRange`
- `IsInside`
- `IsOutside`
- `Contains`
- `ContainsZero`
- `RefineMetric`
- `FloorBounds`
- `ApplyULFT`

### Notes

Current transform streaming mostly expects inside ranges. Outside ranges exist
for interval semantics and heuristics, but some operations intentionally reject
them.

## 9. Transform types

### `ULFT`

Unary linear fractional transform:

- `(A*x + B) / (C*x + D)`

Used for:

- exact transform evaluation on rationals
- range mapping
- streaming CF digit emission
- generalized CF ingestion rewrites

### `BLFT`

Binary linear fractional transform:

- `(A*x*y + B*x + C*y + D) / (E*x*y + F*x + G*y + H)`

Used for:

- binary arithmetic over CF sources
- product / quotient / other bilinear operations

### `DiagBLFT`

Diagonal specialization of a BLFT on `T(x,x)`:

- `(A*x^2 + B*x + C) / (D*x^2 + E*x + F)`

Used for:

- square / quadratic forms
- diagonal specialization experiments
- narrow algebraic shortcuts for some sqrt sources

## 10. Streaming transform engines

### `ULFTStream`

Transforms a regular CF source through a ULFT and emits regular CF digits when
safe.

### `BLFTStream`

Transforms two regular CF sources through a BLFT.

### `DiagBLFTStream`

Transforms one regular CF source through a diagonal BLFT.

### Important warning about all stream engines

`Next()` returns `(0, false)` for both:

- clean exhaustion
- error termination

Always inspect `Err()` after a stream stops if you need to know why it ended.

### Stream options

The stream types expose options for:

- cycle detection
- refine guards
- bounded rational finalization in some cases

These are important when experimenting with incomplete algorithms or debugging
stalls.

## 11. Square root support

The sqrt subsystem is currently bounded and layered. It is useful, but it is
not yet a full true streaming exact-real square-root operator.

### Exact / low-level sqrt functions

- `RationalSqrtExact`
- `NewtonSqrtStep`
- `NewtonSqrtIterates`
- `SqrtResidual`
- `SqrtResidualAbs`

### Bounded rational sqrt approximation

- `SqrtApproxRational`
- `SqrtApproxRationalUntilExact`
- `SqrtApproxRationalUntilResidual`

Default-seed wrappers also exist.

### Policy-driven top-level sqrt APIs

Recommended starting points:

- `SqrtApprox`
- `SqrtApproxCF`
- `SqrtApproxTermsAuto`

Advanced variants:

- `SqrtApproxWithPolicy`
- `SqrtApproxCFWithPolicy`
- `SqrtApproxTermsWithPolicy`

Compatibility wrappers with explicit seeds also exist.

### Source-based sqrt helpers

The package also supports bounded sqrt approximations derived from regular-CF
prefixes, using:

- convergent-targeted approximation
- range-seeded approximation
- midpoint-targeted approximation

These APIs are powerful, but more experimental than the basic rational sqrt
entry points.

### Sqrt range helpers

- `SqrtRangeExact`
- `SqrtRangeHeuristic`

Important warning:

- `SqrtRangeExact` is exact only in narrow supported cases
- `SqrtRangeHeuristic` is not proof-safe and should not be treated as a formal
  conservative enclosure

## 12. Generalized continued fractions (GCF)

### Convention used in this package

The package uses:

- `x = p + q/x'`

with `q > 0` required by ingestion logic.

### `GCFSource`

A generalized continued-fraction source provides:

- `NextPQ() (int64, int64, bool)`

### Built-in GCF sources

Finite:

- `NewSliceGCF(...)`

Periodic:

- `NewPeriodicGCF(prefix, period)`

Function-backed:

- `NewFuncGCFSource(fn)`

Regular-CF adapter:

- `AdaptCFToGCF(src)`

Named / algorithmic sources currently include:

- `NewECFGSource()`
- `NewBrouncker4OverPiGCFSource()`
- `NewLambertPiOver4GCFSource()`
- `NewUnitPArithmeticQGCFSource(startQ, step)`
## 13. GCF ingestion, evaluation, and bounded prefix semantics

### GCF convention

Generalized continued fractions in this package use:

- `x = p + q/x'`

with `q > 0` required by ingestion logic.

### Forward composition

- `ComposeGCFToULFT(src)`

Composes a finite GCF prefix into a ULFT without buffering the whole stream.

This is useful for exact algebraic reasoning and for future enclosure logic.

### Exact finite evaluation

- `EvaluateFiniteGCF(src)`

Evaluates a finite GCF exactly under the package’s finite-tail convention.

### `GCFBounder`

`GCFBounder` is the current bounded-ingestion object for generalized continued fractions.

Main methods:

- `NewGCFBounder()`
- `IngestPQ(p, q)`
- `SetTailLowerBound(lower)`
- `Finish()`
- `HasValue()`
- `Convergent()`
- `Range()`

### Exact convergents

`Convergent()` returns the exact rational value of the finite prefix seen so far,
using the finite-tail convention that the last ingested term contributes just
its `p` value.

### Finished vs unfinished prefixes

After `Finish()`:

- `Range()` collapses to an exact point range

Before `Finish()`:

- if no tail lower bound is known, `Range()` currently falls back to a point
  placeholder at the current convergent
- if a positive lower bound for the unfinished tail is known, `Range()` uses a
  conservative ULFT ray-image enclosure

### Positive tail lower bound metadata

Some GCF sources implement:

- `PositiveTailLowerBoundedGCFSource`

with:

- `TailLowerBound() Rational`

This means the unfinished tail is known to satisfy:

- `tail >= L`
- with `L > 0`

For such sources, bounded prefix ingestion can produce a conservative enclosure
for the unfinished prefix by mapping the ray `[L, +∞)` through the composed
prefix ULFT.

This is the first real conservative unfinished-prefix enclosure mechanism for
GCF in the package.

### Current limitations

This enclosure model is still intentionally narrow:

- it depends on a positive lower bound for the unfinished tail
- it currently uses ray-image logic, not a fully general tail enclosure model
- some GCF sources may need richer metadata later than just `TailLowerBound()`

So today, unfinished GCF range semantics are:

- conservative for supported positive-tail sources
- still limited / placeholder-like for unsupported source classes

## 14. Bounded GCF prefix APIs

### Ingestion helpers

- `IngestAllGCF(src)`
- `IngestGCFPrefix(src, prefixTerms)`

If the source also provides positive tail lower-bound metadata, prefix ingestion
will automatically configure the returned `GCFBounder` to use conservative
unfinished-prefix ray enclosures.

### Snapshot API

- `GCFApprox`
- `GCFApproxFromPrefix(src, prefixTerms)`

`GCFApprox` currently stores:

- exact convergent
- prefix length used

Unlike `CFApprox`, it does not yet carry a full enclosure object.

### Inspection helpers

- `GCFApproxCF`
- `GCFApproxTerms`
- `GCFSourceConvergent`
- `GCFSourceTerms`

These are currently the easiest way to inspect bounded GCF prefixes:
- as exact rationals
- or as regular CF term lists derived from those exact rational convergents

## 15. Named GCF sources

Current named or algorithmic GCF sources include:

- `NewECFGSource()`
- `NewBrouncker4OverPiGCFSource()`
- `NewLambertPiOver4GCFSource()`
- `NewUnitPArithmeticQGCFSource(startQ, step)`

### Brouncker and Lambert

The current GCF enclosure design has been exercised on:

- Brouncker for `4/pi`
- Lambert for `pi/4`

using bounded prefixes and conservative unfinished-tail enclosures derived from
positive tail lower-bound metadata.

These sources are now good reference cases for testing and examples.

## 16. Debugging tips

Helpful tools and concepts include:

- `FingerprintULFT`
- `FingerprintBLFT`
- `Err()` on all stream engines
- bounded refine limits
- cycle detection options
- converting exact rational values back to regular CF via `NewRationalCF`

For GCF debugging, inspect:

- bounded prefix convergents as exact rationals
- regular CF term expansions of those convergents

## 17. Current warnings and limitations

- GCF support is currently ingest-first; emit support is not implemented
- unfinished GCF enclosure semantics are not yet conservative
- bounded sqrt APIs are not the same as true exact-real sqrt streaming
- heuristic sqrt-range helpers are not proof-safe
- some transform interval mappings are intentionally narrow
- some stream engines include debugging guards that are primarily for development
  safety rather than polished end-user semantics

## 18. Suggested worked examples to add later

- finite rational to regular CF
- regular CF prefix convergents with `Bounder`
- `ULFTStream` identity and reciprocal examples
- bounded sqrt of a rational using default policy
- bounded sqrt from a regular-CF prefix
- Brouncker bounded prefixes for `4/pi`
- Lambert bounded prefixes for `pi/4`
- comparing regular CF vs adapted GCF ingestion
- using a function-backed GCF source

## 19. Future documentation work

Planned follow-up documentation should include:

- a shorter Quick Start
- a transform cookbook
- a GCF conventions section with derivations
- a “recommended APIs vs advanced APIs” table
- warnings and stability notes for experimental features
- examples for future published GCF sources, including pi-related ones

# End of UserGuide.md