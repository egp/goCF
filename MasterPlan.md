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

## Important current conclusions
- Lambert specialization can improve cadence, but that is not currently the highest-leverage project direction.
- Brouncker source-specific tail certification remains exploratory and is consuming disproportionate attention.
- Tail certification is not necessary for correctness; it is only extra power for early certification from partial infinite named-source ingestion.
- The current project focus should return to more central Gosper-style capabilities:
  - generic infinite `(p,q)` streaming,
  - unary operator support,
  - and `sqrt`.

## Change of focus
- Set aside named-source tail-certification work for now.
- Treat the current tail-certification framework as exploratory/incomplete rather than an active top priority.
- Do not expand Brouncker/Lambert source-specific certification further at this time.
- Revisit named-source tail certification only after generic infinite GCF streaming is healthier without source-specific lookahead advantages.

## Highest-leverage next production task
- Refocus on unary/LFT machinery (`ULFT`) as the next central production area.
- Use unary support as the path back toward more interesting and foundational capabilities.
- After unary/LFT work is in a healthier state, continue `sqrt`.

## Recommended immediate work items
1. Review the current ULFT/unary code and tests.
2. Identify the highest-leverage remaining ULFT correctness/design issue.
3. Continue unary operator work with correctness and bounded tests first.
4. Resume `sqrt` work once the unary path is in better shape.
5. Improve generic infinite `(p,q)` stream behavior without relying on named-source lookahead or source-specific tail-certification advantages.

## Secondary cleanup tasks
- Clarify evidence hierarchy comments near `GCFTailEvidence` and related interfaces.
- Rename any remaining Brouncker tests whose names still imply infinite behavior when they now use finite-prefix wrappers.
- Add a compact note or matrix documenting which evidence layers each named source currently uses.
- Revisit API elegance after unary and `sqrt` work are in a healthier state.

## Deferred for now
- Further Brouncker/Lambert tail-certification expansion.
- Additional named-source proof helpers.
- More cadence tuning for named-source specialized evidence.
- Broad API simplification driven by tail-certification abstractions.
- Large API documentation rewrite.
- More documentation updates unrelated to current unary/`sqrt` production progress.

## Ongoing design principles
- Prefer fewer medium-sized cohesive files over many tiny files.
- Keep methods around one screen when practical; extract helpers for clarity.
- Bound tests and avoid any test that can wait forever on infinite uncertified streaming.
- Do not mistake “does not hang” for “mathematically justified”.
- Cash out infrastructure work into central library capabilities rather than source-specific side machinery.
- Keep the generic Gosper engine conceptually separate from source-specific proof/certification experiments.

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
- Separate experimental problem:
  - source-specific tail certification for partial infinite named-source ingestion
- The project should currently prioritize the first problem over the second.

## Suggested next-chat bootstrap
Ask for current contents of:
- the relevant ULFT production file(s)
- the relevant ULFT test file(s)
- any current `sqrt` production/test files if they are directly involved

Then focus on:
- ULFT correctness first
- bounded tests
- generic infinite stream behavior where relevant
- `sqrt` second
- not on named-source tail certification first.