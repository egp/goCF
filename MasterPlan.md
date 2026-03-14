# MasterPlan.md

# goCF Master Plan

## Current strategic goals
- Keep mathematical correctness first.
- Improve streaming certifiability for generalized continued fractions.
- Prefer source-specific mathematically justified evidence over weak generic fallback.
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
- Added Brouncker prefix-2 specialized tail evidence.
- Added Brouncker prefix-3 specialized tail evidence.
- Added Brouncker prefix-4 specialized tail evidence.
- Updated Brouncker approximation and inspect expectations to match stronger evidence.
- Added true lower-bound-only Brouncker wrapper/baseline comparisons.
- Traced a bad Brouncker third-digit path to generic `lower-bound-ray` fallback.
- Disabled weak Brouncker lower-bound-ray fallback by raising `LowerBoundMinPrefix` to an effectively-disabled sentinel.
- Converted several Brouncker infinite-stream tests to bounded finite-prefix versions to avoid hangs after fallback disablement.
- Established that Brouncker specialization is currently “no worse than” lower-bound-only for early digits, but its main benefit so far is correctness/safety rather than visible first-two-digit cadence improvement.

### Test organization
- Split oversized `gcf_stream_tail_evidence_test.go` into smaller focused files.
- Continued moving named-source-specific tests into named-source-specific files where appropriate.
- Reworked hanging tests into bounded forms where needed.

## Important current conclusions
- Lambert is paying rent: source-specific evidence improves real streaming cadence.
- Brouncker is not yet healthy as an infinite stream after weak fallback disablement:
  - correctness improved,
  - but some early infinite digits are no longer readily certifiable.
- Therefore, the next leverage is in Brouncker source math/evidence, not more generic infrastructure.

## Highest-leverage next production task
- Restore useful infinite Brouncker streaming by adding stronger explicit source-specific evidence.
- Target:
  - make the infinite Brouncker source safely certify at least the first two digits again,
  - and then work toward the third digit,
  - without re-enabling weak generic lower-bound-ray fallback.

## Recommended immediate work items
1. Review current Brouncker source/tail files and current named-source stream tests.
2. Identify which current Brouncker early-digit tests still use finite-prefix wrappers because infinite certification is too weak.
3. Improve Brouncker source-specific evidence, likely via one or more of:
   - stronger `TailEvidence()`
   - `CandidateTailEvidence()`
   - `PostEmitTailEvidence()`
4. Reintroduce true infinite-source Brouncker tests only when the source can certify them safely and boundedly.
5. Keep generic lower-bound-ray fallback effectively disabled for Brouncker unless a mathematically justified restricted policy is found.

## Secondary cleanup tasks
- Clarify evidence hierarchy comments near `GCFTailEvidence` and related interfaces.
- Rename any remaining Brouncker tests whose names still imply infinite behavior when they now use finite-prefix wrappers.
- Add a compact note or matrix documenting which evidence layers each named source currently uses.
- Revisit API elegance after Brouncker infinite streaming is in a healthier state.

## Deferred for now
- Large API documentation rewrite.
- Broad API simplification pass.
- Additional generic evidence machinery.
- More Brouncker prefix helpers if they are only crude upper-bound extensions without clear payoff.
- More documentation updates unrelated to current production/test progress.

## Ongoing design principles
- Prefer fewer medium-sized cohesive files over many tiny files.
- Keep methods around one screen when practical; extract helpers for clarity.
- Bound tests and avoid any test that can wait forever on infinite uncertified streaming.
- Do not mistake “does not hang” for “mathematically justified”.
- Cash out infrastructure work into real named-source improvements whenever possible.

## Current named-source characterization
- Lambert:
  - cadence-oriented specialization
  - prefix-specific evidence visibly improves streaming
- Brouncker:
  - safety-oriented specialization
  - weak generic fallback was found to be too permissive
  - explicit source evidence needs to be strengthened to recover healthy infinite early-digit certification

## Suggested next-chat bootstrap
Ask for current contents of:
- `cf/brouncker_pi_tail.go`
- `cf/brouncker_pi_gcf_test.go`
- `cf/gcf_stream_named_sources_test.go`
- `cf/gcf_stream_tail_evidence.go`

Then focus on:
- strengthening Brouncker explicit evidence for infinite early digits,
- not on generic framework growth,
- not on docs first.
