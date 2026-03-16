# MasterPlan.md

# goCF Master Plan

## Mission
Build a mathematically correct, testable continued-fraction arithmetic library with strong GCF support and eventually a real streaming `sqrt` operator that emits certified digits rather than only collapsing bounded approximations to rationals.

## Current status
The project has largely completed the sqrt migration/refactor phase. The remaining work is now primarily operator completion, not cleanup.

What is working:
- exact-tail transform streams for ULFT / DiagBLFT / BLFT
- canonical/internal bounded sqrt substrate
- public bounded sqrt streams for CF prefix, GCF prefix, and GCF exact-tail input
- stream introspection and status semantics
- proof-safe conservative sqrt enclosure engine
- unary reciprocal operator for GCF exact-tail and bounded GCF-prefix input
- CF-side certified-progressive sqrt stream
- reusable certified CF range certifier/emitter

What is not complete:
- GCF-side certified-progressive sqrt is drafted but currently hangs on adapted `sqrt(2)`
- progressive streams still rebuild/replay rather than carrying a principled persistent transformed remainder state
- diagonal / transform-driven sqrt is not yet the main engine
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

## Current milestone
Finish general sqrt as a real operator by stabilizing GCF-side progressive certification and then replacing replay/rebuild logic with a principled persistent transformed remainder state.

## Critical path to completion

1. Stabilize GCF unfinished-tail enclosure / progress guarantees
2. Finish multi-digit certified-progressive GCF sqrt
3. Replace replay-on-refinement with persistent transformed remainder state
4. Drive remainder/update and emission through diagonal / transform machinery
5. Unify progressive sqrt architecture across CF and GCF input
6. Finalize unary reciprocal as a polished first-class operator family
7. Final naming / wrapper retirement / duplicate-test retirement / docs

## Immediate next technical focus
Diagnose and fix the hang in the GCF progressive sqrt path:
- `TestSqrtCertifiedGCFPrefixStream_Sqrt2_RefinesAndContinues`

Most likely root area:
- unfinished-tail enclosure quality / progress guarantees on the GCF side
- range refinement may be rebuilding without producing additional certifiable progress

The next fix should target GCF unfinished-tail enclosure/progress behavior, not more surface API patching.

## Known risks / unresolved design questions
- Is current GCF unfinished-tail metadata/range information strong enough for progressive certification, especially for adapted CF sources?
- Should GCF progressive sqrt require stronger monotone-progress checks to avoid rebuild loops?
- How much of the eventual remainder/update engine should live in generic range certifier code vs dedicated diagonal/transform state?
- When the operator stabilizes, should public names keep legacy wrappers or collapse onto the canonical/internal naming scheme?

## Deferred work / future ideas
- final naming/API cleanup
- retirement of remaining compatibility wrappers
- broader documentation rewrite
- decimal digit emission (low priority)
- broader unary-operator family beyond reciprocal/sqrt

## Practical guidance
- Do not restart broad refactoring
- Prefer mathematically justified fixes over API cosmetics
- Use the fast full suite (`go test ./cf`) frequently
- For new work, target the operator path, not wrapper proliferation

# EOF MasterPlan.md