# goCF Master Plan

## Mission
Build a mathematically correct, testable continued-fraction arithmetic library with strong GCF support, robust unary/operator plumbing, and eventually real streaming operators that emit certified RCF terms rather than only collapsing bounded approximations to rationals.

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
- a third sqrt rebuild produced green exact-arithmetic support/scaffolding:
  - exact Newton step over `Rational`
  - mutable unary sqrt state
  - GCF prefix state
  - operator refinement policy
  - residual snapshots
  - point-input and ranged-input sqrt enclosures
  - big-int floor bounds for certification
  - forced-term certification helper
  - first forced-term transition
  - demand-driven first-term forcing helper
  - first operator-level emit transition
- emitted outputs for the new sqrt work are RCF terms as `*big.Int`, not decimal digits
- the current sqrt try3 scaffolding is useful support code, but it is not yet Gosper-centered enough to be the final operator core

What is not complete:
- production still retains legacy finite-bridge compatibility helpers in the radicand / numerator area
- production still retains duplicated or transitional helper layers around the radicand / numerator API
- full target is still a bounded non-point result, not a certified point
- GCF-side certified-progressive sqrt is still not the main production engine for the MVP target path
- final public naming / wrapper retirement / docs are deferred
- the target formula is still represented with MVP-specific scaffolding rather than as a thin validation client of the real operator architecture
- legacy sqrt code still exists in parallel with newer sqrt work
- the current sqrt try3 operator core is still Newton/enclosure-centered rather than Gosper transform-centered
- terminology cleanup is still needed in the newer sqrt work: “digit” should become “term”
- a true Gosper-centered transform-first sqrt core has not yet been implemented
- certified nontrivial lazy sqrt term emission for inputs like `sqrt(2)` is not implemented yet

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

### sqrt try3 support scaffold
Started over on sqrt with a new minimal path rather than extending the legacy sqrt jungle, then evolved it into a green support scaffold:
- exact Newton step over `Rational`
- mutable unary sqrt state
- GCF prefix state
- unary sqrt operator scaffold with explicit refinement policy
- residual snapshots
- point-input sqrt enclosure from iterate/reciprocal
- conservative ranged-input sqrt enclosure
- big-int floor bounds for certification
- forced first-term certification helper
- forced first-term transition
- demand-driven first-term forcing helper
- first operator-level emit transition
- operator snapshots carrying current approximation/enclosure/certification state

This work is currently treated as support/scaffolding, not as the final Gosper-centered sqrt core.

## Current milestone
Stop extending the current sqrt try3 core as the main operator, start a parallel Gosper-centered transform-first sqrt core, and then retire legacy sqrt and remaining MVP-specific scaffolding.

## Critical path to completion

1. Freeze the current sqrt try3 core as helper/support code rather than the final operator core
2. Start a parallel transform-first sqrt core based directly on Gosper’s square-root method
3. Use Newton only as the small per-term fixed-point solver inside the transform method, not as the primary operator state
4. First target: reproduce Gosper’s rational example `sqrt(17/10) = [1;(3,3,2)]` with the transform-first core
5. After the transform-first rational case works, connect that core to continued-fraction radicands and GCF ingestion
6. Make the transform-first sqrt path primary and delete or isolate legacy sqrt APIs/files/tests
7. Retire remaining finite-bridge and legacy compatibility helpers from production in the radicand / numerator path
8. Collapse duplicated radicand / numerator helper layers onto one explicit snapshot-assembly path
9. Reduce MVP-specific scaffolding until the target formula is only a validation client of the real operator architecture
10. Advance from bounded non-point target range toward a tighter/certified point result where mathematically justified
11. Resume broader operator completion and eventual streaming/certified-progressive operator unification
12. Continue toward general arithmetic on infinite GCF streams with operators consuming stream-shaped inputs and emitting RCF

## Immediate next technical focus
Take the next big bite on the critical path in sqrt:
- do not extend the current Newton/enclosure-centered try3 core
- obtain the current `range.go` and relevant transform files
- start a parallel transform-first fixed-point sqrt core
- write red tests for Gosper’s rational example `sqrt(17/10) = [1;(3,3,2)]`
- keep Newton only as the tiny per-term solver inside the transform method
- keep try3 green and available as support/scaffolding while the new core is built

## Known risks / unresolved design questions
- What is the smallest correct transform-first state representation for Gosper sqrt using the existing ULFT / DiagBLFT machinery?
- For the transform-first core, what exact left/right state should be compared to decide whether the next RCF term is forced?
- How should Newton averaging be embedded as the small per-term fixed-point solver without becoming the primary operator state?
- Which existing try3 helpers are worth reusing inside the transform core, and which should be frozen and ignored?
- How aggressively should legacy finite-bridge helpers be removed versus kept temporarily as thin wrappers during the cleanup phase?
- What is the cleanest final boundary between durable production operator APIs and milestone-specific target-formula helpers?
- How far should MVP go toward a point result versus accepting a mathematically justified bounded range?
- When operator stabilizes, should public names keep compatibility wrappers or collapse onto canonical names?
- Should any public budget/tuning API exist for sqrt, or should all such controls remain internal?

## Deferred work / future ideas
- terminology cleanup in the newer sqrt work: rename “digit” to correct RCF “term”
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