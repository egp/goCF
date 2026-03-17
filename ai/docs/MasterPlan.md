# MasterPlan.md

# goCF Master Plan

## Mission
Build a mathematically correct, testable continued-fraction arithmetic library with strong GCF support, robust unary/operator plumbing, and eventually real streaming operators that emit certified digits rather than only collapsing bounded approximations to rationals.

## Current status
The project has moved from broad sqrt refactor work into MVP target construction and cleanup for a concrete expression:

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
- canonical/internal bounded sqrt substrate
- public bounded sqrt streams for CF prefix, GCF prefix, and GCF exact-tail input
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

What is not complete:
- production still retains legacy finite-bridge compatibility helpers in the radicand / numerator area
- production still retains duplicated or transitional helper layers around the radicand / numerator API
- full target is still a bounded non-point result, not a certified point
- GCF-side certified-progressive sqrt is still not the main production engine for the MVP target path
- final public naming / wrapper retirement / docs are deferred
- the target formula is still represented with MVP-specific scaffolding rather than as a thin validation client of the real operator architecture

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
Public/user-facing bounded sqrt constructors now exist:
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

## Current milestone
Retire remaining MVP-specific scaffolding in the numerator/radicand area so the target formula becomes a thin validation client of the real operator architecture.

## Critical path to completion

1. Retire remaining finite-bridge and legacy compatibility helpers from production in the radicand / numerator path
2. Collapse duplicated radicand / numerator helper layers onto one explicit snapshot-assembly path
3. Reduce MVP-specific scaffolding until the target formula is only a validation client of the real operator architecture
4. Advance from bounded non-point target range toward a tighter/certified point result where mathematically justified
5. Resume broader operator completion and eventual streaming/certified-progressive operator unification
6. Continue toward general arithmetic on infinite GCF streams with operators consuming stream-shaped inputs and emitting RCF

## Immediate next technical focus
Take the next big bite on the critical path:
- retire the remaining finite-bridge compatibility helpers in production
- preserve green tests
- keep the live numerator/radicand path on the unified snapshot assembly pipeline
- convert old bridge helpers to thin wrappers or remove them if no longer needed

## Known risks / unresolved design questions
- How aggressively should legacy finite-bridge helpers be removed versus kept temporarily as thin wrappers during the cleanup phase?
- What is the cleanest final boundary between durable production operator APIs and milestone-specific target-formula helpers?
- How far should MVP go toward a point result versus accepting a mathematically justified bounded range?
- When operator stabilizes, should public names keep compatibility wrappers or collapse onto canonical names?

## Deferred work / future ideas
- final naming/API cleanup after exception replacement
- retirement of remaining compatibility wrappers
- broader documentation rewrite
- decimal digit emission beyond MVP reporting
- broader unary-operator family beyond reciprocal/sqrt
- eventually real streaming sqrt/operator machinery as the main engine rather than bounded-collapse staging

## Practical guidance
- Do not restart broad refactoring
- Prefer mathematically justified fixes over API cosmetics
- Use the fast full suite (`go test ./cf`) frequently
- For new work, target explicit temporary seams, not broad redesign
- Keep Brouncker canonical unless explicitly changed
- Keep steps small and test-driven
- When changing code, prefer whole-function replacements or clearly identified append locations

# EOF MasterPlan.md