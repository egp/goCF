# goCF Master Plan

## Mission
Build a mathematically correct, testable continued-fraction arithmetic library with strong GCF support, robust unary/operator plumbing, and eventually real streaming operators that emit certified digits rather than only collapsing bounded approximations to rationals.

## Current status
The project has moved past broad sqrt refactor work and past most MVP target-path cleanup into rebuilding sqrt on a simpler foundation.

The concrete validation expression remains:

$$
\frac{\sqrt{\frac{3}{\pi^2} + e}}{\tanh(\sqrt{5}) - \sin(69^\circ)}
$$

    sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))

Current mathematical float estimates for that target are:

- target ≈ 39.66207809377289
- numerator sqrt(3/pi^2 + e) ≈ 1.738460634983162
- denominator tanh(sqrt(5)) - sin(69°) ≈ 0.04383180908657702

What is working:
- exact-tail transform streams for ULFT / DiagBLFT / BLFT
- canonical/internal bounded sqrt substrate from earlier work still exists in legacy code
- stream introspection and status semantics
- proof-safe conservative sqrt enclosure engine
- unary reciprocal operator for GCF exact-tail and bounded GCF-prefix input
- CF-side certified-progressive sqrt stream
- reusable certified CF range certifier/emitter
- canonical MVP 4/pi family is Brouncker
- Lambert retained as an alternate/parity approximation path at the 4/pi layer
- e is ingested as GCF
- sqrt(5) is ingested as GCF
- 69° is represented as an exact finite GCF source
- denominator path no longer depends on an exact-tail trig entry for 69°
- numerator live path no longer depends on the fake finite radicand bridge source
- numerator radicand is assembled through explicit snapshot helpers
- exact scalar constants used in numerator scale-factor construction are routed through finite GCF ingestion
- radicand approximation helpers route through a unified snapshot assembly path
- full MVP target currently returns a positive inside range
- tests are currently green
- new sqrt work is underway around a single canonical public entry point:
  `SqrtGCF(src GCFSource) (ContinuedFraction, error)`
- `SqrtGCF` is lazy at construction time
- exact finite perfect-square inputs round-trip correctly through the new `SqrtGCF`
- exact finite rational perfect-square inputs such as `1/4` and `9/16` round-trip correctly through the new `SqrtGCF`
- exact finite non-square inputs currently return a bounded Newton rational approximation as a regular CF
- `QuadraticRadicalSource` metadata is preserved through `AdaptCFToGCF`
- the new sqrt path has a first true-lazy metadata fast path for perfect-square radicands like `sqrt(4)`
- the new sqrt stream keeps buffered input progress across calls and uses bounded per-call ingestion

What is not complete:
- production still retains legacy finite-bridge compatibility helpers in the radicand / numerator area
- production still retains duplicated or transitional helper layers around the radicand / numerator API
- full target is still a bounded non-point result, not a certified point
- GCF-side certified-progressive sqrt is still not the main production engine for the MVP target path
- final public naming / wrapper retirement / docs are deferred
- the target formula is still represented with MVP-specific scaffolding rather than as a thin validation client of the real operator architecture
- legacy sqrt code still exists in parallel with the new `SqrtGCF` path
- the new sqrt stream still needs one key semantic fix: unresolved long-input state must remain live and retryable rather than being treated as a terminal error
- certified nontrivial lazy sqrt emission for inputs like `sqrt(2)` is not implemented yet

## Completed work

### Exact-tail transform substrate
Completed bounded exact-tail streams and support for:
- `GCFULFTStream`
- `GCFDiagStream`
- `GCFBLFTStream`

Including:
- exact-tail source plumbing
- bounded ingestion
- shared helper extraction to reduce duplication
- stronger exact-rational-image tests

### sqrt canonicalization
Established the newer canonical sqrt substrate:
- `sqrt_core_exact.go`
- `sqrt_seed_range.go`
- `sqrt_api2.go`
- `sqrt_source_api2.go`
- `sqrt_source_prefix_api2.go`
- `sqrt_gcf_api2.go`
- `sqrt_gcf_range_seed_api2.go`
- `sqrt_gcf_tail_api2.go`
- `sqrt_midpoint_api2.go`
- `sqrt_canonical_api.go`
- `sqrt_canonical_source_api.go`

Migrated much of the old public surface onto this newer substrate.

### Legacy simplification
Substantial retirement of duplicate or obsolete sqrt code/tests:
- old exact/Newton/range bodies reduced to wrappers or retired
- dead helpers removed
- obsolete tests retired
- remaining legacy tests trimmed toward public-surface regression coverage

### Public bounded sqrt surfaces
Public/user-facing bounded sqrt constructors now exist in legacy code:
- `SqrtStream(...)`
- `SqrtGCFStream(...)`
- `SqrtGCFExactTailStream(...)`

All expose inspectable stream state.

### Stream introspection / semantics
Snapshots and status semantics added:
- `unstarted`
- `exact_input`
- `bounded_collapse`
- `certified_progressive`
- `failed`

### Proof-safe sqrt enclosure engine
Implemented first proof-safe conservative sqrt enclosure support:
- `SqrtLowerBoundRational`
- `SqrtUpperBoundRational`
- `SqrtRangeConservative`

Current implementation:
- exact-square fast paths
- negative rejection
- proof-safe scaled integer bracketing for non-square rationals
- conservative enclosure over nonnegative inside ranges

### Unary reciprocal operator
Unary reciprocal is now part of the operator path:
- public exact-tail GCF reciprocal stream
- public bounded GCF-prefix reciprocal stream
- reciprocal range operator for proof-safe range work

### Certified-progressive sqrt
CF-side certified-progressive sqrt now exists:
- conservative sqrt range
- certified digit extraction from ranges
- persistent certified CF range emitter
- refinement of input when current certification is exhausted
- continued emission without immediate bounded-rational collapse

Reusable components added:
- `ReciprocalRangeConservative`
- `ShiftRangeByInt`
- `CertifiedRemainderRange`
- `CertifyCFDigitsFromRange`
- `CertifiedCFRangeEmitter`

### MVP target path cleanup
The current MVP target path has been narrowed and clarified:
- e is ingested as GCF via `MVPEGCFSource` / `NewECFGSource`
- sqrt(5) is ingested as GCF via `AdaptCFToGCF(Sqrt5CF())`
- 69° is represented as exact finite GCF input
- denominator no longer uses an exact-tail trig entry
- numerator live path no longer depends on the fake finite radicand bridge source
- radicand construction is split into explicit snapshot assembly steps
- exact scalar snapshots used in the numerator scale factor are routed through finite GCF ingestion
- radicand approximation helpers now route through the unified snapshot assembly path
- target-level bridge stability tests are green
- current and sharper target budgets overlap
- sharper numerator budgets do not widen the target range
- production naming in this area has been pushed toward radicand / radicand-root terminology, with the full target formula kept in tests rather than production

### New sqrt reboot
Started over on sqrt with a new minimal path rather than extending the legacy sqrt jungle:
- introduced a single new public entry point:
  `SqrtGCF(src GCFSource) (ContinuedFraction, error)`
- added a small Newton bootstrap core:
  - `sqrtNewtonStep`
  - `sqrtNewtonApprox`
  - bootstrap state wrapper
- added tests for:
  - exact finite squares
  - exact finite rational squares
  - negative input rejection
  - non-square bootstrap approximation
  - lazy construction
  - lazy metadata shortcut for perfect-square radicands
  - bounded per-call ingestion
  - recorded terminal error behavior
- added first true-lazy shortcut:
  - if the source advertises square-radicand metadata and the radicand is a perfect square, `SqrtGCF` returns the exact root CF without reading source terms

## Current milestone
Make the new `SqrtGCF` path a mathematically coherent live lazy unary operator, then retire legacy sqrt and remaining MVP-specific scaffolding.

## Critical path to completion

1. Finish the new sqrt stream semantics so unresolved long-input state remains live and retryable rather than becoming an immediate terminal error
2. Replace bootstrap-style whole-value thinking with a real unary operator state for the new `SqrtGCF` path
3. Certify the first nontrivial lazy sqrt output for a case such as `sqrt(2)`
4. Make the new sqrt path primary and delete or isolate legacy sqrt APIs/files/tests
5. Retire remaining finite-bridge and legacy compatibility helpers from production in the radicand / numerator path
6. Collapse duplicated radicand / numerator helper layers onto one explicit snapshot-assembly path
7. Reduce MVP-specific scaffolding until the target formula is only a validation client of the real operator architecture
8. Advance from bounded non-point target range toward a tighter/certified point result where mathematically justified
9. Resume broader operator completion and eventual streaming/certified-progressive operator unification
10. Continue toward general arithmetic on infinite GCF streams with operators consuming stream-shaped inputs and emitting RCF

## Immediate next technical focus
Take the next big bite on the critical path in the new `SqrtGCF` path:
- distinguish live unresolved state from terminal error
- preserve bounded per-call ingestion
- preserve progress across repeated `Next()` calls
- make bootstrap-budget exhaustion the first terminal error for long unresolved inputs
- keep exact finite and metadata-fast-path tests green

## Known risks / unresolved design questions
- How aggressively should legacy finite-bridge helpers be removed versus kept temporarily as thin wrappers during the cleanup phase?
- What is the cleanest final boundary between durable production operator APIs and milestone-specific target-formula helpers?
- How far should MVP go toward a point result versus accepting a mathematically justified bounded range?
- When operator stabilizes, should public names keep compatibility wrappers or collapse onto canonical names?
- How much of the existing range machinery should be reused directly for certified lazy sqrt emission versus introducing a cleaner unary-specific layer?
- Should any public budget/tuning API exist for sqrt, or should all such controls remain internal?

## Deferred work / future ideas
- final naming/API cleanup after exception replacement
- retirement of remaining compatibility wrappers
- broader documentation rewrite
- decimal digit emission beyond MVP reporting
- broader unary-operator family beyond reciprocal/sqrt
- eventually real streaming sqrt/operator machinery as the main engine rather than bounded-collapse staging
- cleanup to remove remaining production references where legacy “CF” meant regular CF instead of GCF; production should be GCF-first except for emitted result terms
- later cleanup to remove remaining MVP scaffolding once replacement paths are stable

## Practical guidance
- Do not restart broad refactoring
- Prefer mathematically justified fixes over API cosmetics
- Use the fast full suite (`go test ./cf`) frequently
- For new work, target explicit temporary seams, not broad redesign
- Keep Brouncker canonical unless explicitly changed
- Keep steps small and test-driven
- When changing code, prefer whole-function replacements or clearly identified append locations

# EOF MasterPlan.md