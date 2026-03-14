# MasterPlan.md

# goCF Master Plan

## Current strategic goals
- Keep mathematical correctness first.
- Improve streaming certifiability for generalized continued fractions.
- Prefer simple, generic, mathematically justified mechanisms over source-specific special cases.
- Keep production code cohesive and testable.
- Use bounded/fast tests and avoid unbounded hangs.

## Current status summary
- `GCFStream` has been refactored and split into:
  - `gcf_stream.go`
  - `gcf_stream_tail_evidence.go`
- The stream now supports a layered evidence model:
  - `TailEvidence`
  - `CandidateTailEvidence`
  - `RefinedTailEvidence`
  - `PostEmitTailEvidence`
- `GCFStreamOptions` now includes:
  - `MaxRefinementSteps`
  - `Trace`
- Trace-based tests exist for evidence-path observation.
- Tail-evidence-driven certification is materially more capable than earlier versions.
- Unary/LFT work and `sqrt` remain the more central long-term direction for the library.

## Completed recently

### GCFStream structure and policy
- Split `gcf_stream.go` into progression/lifecycle logic and tail-evidence logic.
- Added configurable refinement depth via `GCFStreamOptions`.
- Added trace hook to observe evidence-path selection in tests.
- Added support for multiple successive refinements before further ingestion.
- Added support for candidate tail evidences before refinement or ingestion.

### Lambert improvements
- Added Lambert prefix-2 specialized tail evidence.
- Added Lambert prefix-3 specialized tail evidence.
- Updated Lambert approximation and inspect expectations to match stronger evidence.
- Added Lambert cadence/baseline comparisons against a true lower-bound-only wrapper.
- Demonstrated that Lambert specialization gives visible cadence payoff.

### Brouncker improvements
- Added Brouncker prefix-specific specialized tail evidence.
- Added true lower-bound-only Brouncker wrapper/baseline comparisons.
- Traced a bad Brouncker third-digit path to generic `lower-bound-ray` fallback.
- Disabled weak Brouncker lower-bound-ray fallback by raising `LowerBoundMinPrefix` to an effectively-disabled sentinel.
- Converted several Brouncker infinite-stream tests to bounded finite-prefix versions to avoid hangs after fallback disablement.
- Added exploratory Brouncker candidate/post-emit evidence work.
- Switched Brouncker lookahead helper arithmetic to big-rational arithmetic to avoid `int64` overflow.

### Test organization
- Split oversized `gcf_stream_tail_evidence_test.go` into smaller focused files.
- Continued moving named-source-specific tests into named-source-specific files where appropriate.
- Reworked hanging tests into bounded forms where needed.

### ULFT exact finite-GCF substrate
- Kept the old regular-CF/RCF-based unary path intact rather than mutating it in place.
- Started a new GCF-native ULFT path beside the old path.
- Added `GCFULFTStream` for the current narrow milestone:
  - finite GCF prefix ingestion
  - exact tail evidence only
  - regular-CF digit emission from exact rational result
- Added `GCFTailSource` abstraction for the new GCF-native unary path.
- Added `ExactTailSource` and `NoTailSource`.
- Added bounded-ingestion policy via `MaxIngestTerms` for the new exact finite-prefix path.

### ULFT exact helper ladder
- Added `ComposeGCFIntoULFTBounded`.
- Added `ApplyComposedGCFULFTToTailExact`.
- Added `EvalGCFWithTailExact`.
- Locked down `ULFT.IngestGCF` with direct rewrite-law tests and a bounded property test.

### BLFT exact finite-GCF substrate
- Added `BLFT.IngestGCFX`.
- Added `BLFT.IngestGCFY`.
- Locked down x-side and y-side BLFT GCF rewrite laws with focused tests.
- Added exact bounded BLFT helpers:
  - `ApplyComposedGCFXBLFTToTailsExact`
  - `ApplyComposedGCFYBLFTToTailsExact`
  - `ApplyComposedGCFXYBLFTToTailsExact`
- Locked down two-sided BLFT ingestion consistency and x/y order-independence on exact evaluation.

### DiagBLFT exact finite-GCF substrate
- Confirmed existing `DiagBLFT.IngestGCF` against exact rewrite-law tests.
- Added `ApplyComposedGCFDiagBLFTToTailExact`.
- Locked down consistency between:
  - `DiagFromBLFT(base).ApplyRat(x)`
  - `base.ApplyRat(x,x)`
- Locked down consistency between:
  - diagonal ingestion through `DiagBLFT`
  - two-sided equal-source ingestion through `BLFT` followed by diagonal specialization.
- Added `GCFDiagStream` for the current narrow milestone:
  - finite GCF prefix ingestion
  - exact tail evidence only
  - regular-CF digit emission from exact rational result

## Important current conclusions
- Lambert specialization can improve cadence, but that is not currently the highest-leverage project direction.
- Brouncker source-specific tail certification remains exploratory and is consuming disproportionate attention.
- Tail certification is not necessary for correctness; it is only extra power for early certification from partial infinite named-source ingestion.
- The current project focus should remain on more central Gosper-style capabilities:
  - generic infinite `(p,q)` streaming,
  - unary operator support,
  - BLFT ingestion support,
  - diagonal support,
  - and `sqrt`.

## Change of focus
- Set aside named-source tail-certification work for now.
- Treat the current tail-certification framework as exploratory/incomplete rather than an active top priority.
- Do not expand Brouncker/Lambert source-specific certification further at this time.
- Revisit named-source tail certification only after generic infinite GCF streaming is healthier without source-specific lookahead advantages.

## Highest-leverage current production area
- The exact finite-GCF algebraic substrate for:
  - ULFT
  - BLFT
  - DiagBLFT
- This exact substrate is now materially healthier and better factored than before.
- The next central production challenge is no longer exact finite-prefix algebra; it is the transition from exact finite-prefix support toward mathematically justified non-exact / infinite GCF support.

## Recommended immediate work items
1. Review the new GCF-native exact finite-prefix path and keep it cohesive.
2. Decide the next narrow step for non-exact / infinite GCF support without reviving source-specific tail-certification work.
3. Keep BLFT and Diag work aligned with the same bounded exact-first discipline used for ULFT.
4. Resume `sqrt` work only after the required diagonal/GCF substrate is healthy enough to support it cleanly.
5. Keep the obsolete regular-CF unary path and its tests until the new GCF-native path is sufficiently complete and trusted.

## Recommended next production steps
1. Add a shared exact-tail stream abstraction helper if duplication between `GCFULFTStream` and `GCFDiagStream` begins to grow.
2. Decide whether the next step should be:
   - a shared exact GCF ingestion helper for future stream implementations, or
   - the first mathematically justified non-exact / infinite GCF interface.
3. Only after that, begin designing the smallest honest interface for uncertified vs certified infinite GCF progress.
4. Keep all tests bounded; do not add any infinite-stream test that can wait forever.

## Secondary cleanup tasks
- Clarify evidence hierarchy comments near `GCFTailEvidence` and related interfaces.
- Rename any remaining Brouncker tests whose names still imply infinite behavior when they now use finite-prefix wrappers.
- Add a compact note or matrix documenting which evidence layers each named source currently uses.
- Revisit API elegance after unary, BLFT, diagonal, and `sqrt` work are in a healthier state.
- Consider later consolidation if exact-tail stream implementations begin to duplicate too much structure.

## Deferred for now
- Further Brouncker/Lambert tail-certification expansion.
- Additional named-source proof helpers.
- More cadence tuning for named-source specialized evidence.
- Broad API simplification driven by tail-certification abstractions.
- Large API documentation rewrite.
- More documentation updates unrelated to current unary/BLFT/diag/`sqrt` production progress.
- Infinite uncertified GCF streaming that lacks a mathematically justified contract.

## Ongoing design principles
- Prefer fewer medium-sized cohesive files over many tiny files.
- Keep methods around one screen when practical; extract helpers for clarity.
- Bound tests and avoid any test that can wait forever on infinite uncertified streaming.
- Do not mistake “does not hang” for “mathematically justified”.
- Cash out infrastructure work into central library capabilities rather than source-specific side machinery.
- Keep the generic Gosper engine conceptually separate from source-specific proof/certification experiments.
- Keep the obsolete regular-CF path intact until the new GCF-native path is meaningfully complete.
- Prefer exact finite-prefix algebra first, then bounded interfaces, then carefully justified non-exact / infinite behavior.

## Current named-source characterization
- Lambert:
  - useful as a named-source experiment
  - not a current top-priority optimization target
- Brouncker:
  - useful as a stress case
  - source-specific certification remains exploratory
  - not the current focus

## Current architectural perspective
- Core engine problem:
  - ingest `(p,q)` terms generically
  - maintain transforms correctly
  - emit only when certification is mathematically justified
- Current exact finite-prefix milestone:
  - ULFT exact bounded GCF ingestion is working
  - BLFT exact bounded GCF ingestion is working
  - DiagBLFT exact bounded GCF ingestion is working
  - exact-tail regular-CF output streams exist for ULFT and DiagBLFT
- Separate experimental problem:
  - source-specific tail certification for partial infinite named-source ingestion
- The project should currently prioritize the first problem over the second.

## Suggested next-chat bootstrap
Ask for current contents of:
- the relevant new GCF-native production files
- the relevant exact finite-GCF helper files
- the relevant ULFT / BLFT / DiagBLFT test files
- any current `sqrt` production/test files if they are directly involved

Then focus on:
- mathematical correctness first
- bounded tests only
- exact finite-prefix algebra and interfaces first
- honest handling of non-exact / infinite GCF support second
- `sqrt` after the required substrate is ready
- not on named-source tail certification first.