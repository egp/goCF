# MasterPlan.md

# goCF Master Plan

## Mission
Build a mathematically correct, testable continued-fraction arithmetic library with strong GCF support, robust unary/operator plumbing, and eventually real streaming operators that emit certified digits rather than only collapsing bounded approximations to rationals.

## Current status
The project has moved from broad sqrt refactor work into MVP target construction and exception isolation for a concrete expression:

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
- numerator path works through an explicit finite radicand bridge, a GCFApprox snapshot, and unary sqrt ingestion
- denominator path works through explicit finite exact-tail 69° helpers plus degree-aware sin and tanh(sqrt(5)) bounds
- full MVP target currently returns a positive inside range

What is not complete:
- numerator radicand 3/pi^2 + e still collapses to a bounded rational before finite bridge adaptation
- denominator still relies on an explicit finite exact-tail exception for 69°
- full target is still a bounded non-point result, not a certified point
- GCF-side certified-progressive sqrt is still not the main production engine for the MVP target path
- final public naming / wrapper retirement / docs are deferred

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

### MVP target exception isolation
The current MVP target path has been narrowed and clarified:
- canonical radicand seam for the numerator now routes through GCF-facing source/snapshot helpers
- explicit finite-bridge names now own the current numerator radicand bridge implementation
- legacy numerator bridge names are compatibility wrappers
- explicit finite exact-tail names now own the current 69° denominator exception
- legacy 69° names are compatibility wrappers
- target-level bridge stability tests are green
- current and sharper target budgets overlap
- sharper numerator budgets do not widen the target range

## Current milestone
Finish the MVP target path so that every unary/operator entry is GCF-ingesting, all temporary exceptions are explicit and isolated, and the next architectural replacement can target one remaining temporary seam locally.

## Critical path to completion

1. Keep all unary/operator entry points GCF-ingesting
2. Preserve explicit naming and isolation of temporary MVP exceptions
3. Replace either the numerator finite radicand bridge or the denominator 69° finite exact-tail exception with a more source-driven construction
4. Improve target inspection/output so the MVP target can be reported as:
   - regular continued fraction terms
   - approximate rational
   - approximate float
5. Advance from bounded non-point target range toward a tighter/certified point result where mathematically justified
6. Resume broader operator completion and eventual streaming/certified-progressive operator unification

## Immediate next technical focus
Choose and execute the next real architectural replacement step:
- likely replace the numerator finite radicand bridge with a more source-driven radicand construction
- alternatively replace the denominator 69° exact-tail exception if that path proves more tractable

Also prepare a target-level inspection surface so the current MVP result can be rendered as:
- RCF terms
- approximate rational
- approximate float

## Known risks / unresolved design questions
- How quickly can 3/pi^2 + e move from bounded-rational collapse to a more source-driven construction?
- Is the denominator 69° finite exact-tail exception best retired before or after the numerator bridge?
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