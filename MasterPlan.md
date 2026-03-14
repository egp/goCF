# MasterPlan.md

# goCF Master Plan

## Current strategic goals
- Keep mathematical correctness first.
- Continue shifting the library toward GCF-native ingestion and transform-based arithmetic.
- Build `sqrt` on top of the newer canonical substrate instead of expanding legacy RCF-oriented wrapper sprawl.
- Prefer bounded, testable, incremental operator-shaped progress.
- Keep production code cohesive and reduce obsolete duplication when the new path is proven.

## Current status summary
- A new canonical internal `sqrt` path now exists and is actively used by legacy wrappers.
- The library now has new-path bounded sqrt streams for:
  - exact-tail GCF input
  - bounded GCF-prefix input
  - bounded CF-prefix input
- Much of the older sqrt implementation has been turned into compatibility wrappers over the newer canonical path.
- A first round of obsolete sqrt production code and duplicated tests has already been retired.
- The next phase should prioritize actual operator progress for `sqrt`, not more broad refactoring.

## Completed recently

### GCF / transform substrate
- Completed new finite-GCF exact substrate for:
  - `ULFT`
  - `BLFT`
  - `DiagBLFT`
- Added exact-tail streams for:
  - `GCFULFTStream`
  - `GCFDiagStream`
  - `GCFBLFTStream`
- Reduced duplicated stream-state and bounded-ingestion logic with shared helpers.
- Added separate x/y ingest bounds for the BLFT exact-tail stream path.

### sqrt canonicalization
- Added parallel canonical sqrt core:
  - `sqrt_core_exact.go`
  - `sqrt_seed_range.go`
  - `sqrt_api2.go`
- Added canonical/internal delegation layer:
  - `sqrt_canonical_api.go`
  - `sqrt_canonical_source_api.go`
- Redirected old top-level sqrt API to the new path.
- Redirected old range-seeded and midpoint wrappers to the new path.
- Redirected old basic CF adapters to the new path.
- Turned old seed/range and Newton/exact implementations into compatibility wrappers over canonical code.
- Removed dead helper code such as `exactSqrtBig`.

### New sqrt operator-shaped progress
- Added exact-tail bounded sqrt over GCF:
  - `SqrtGCFStream2`
- Added bounded sqrt over GCF prefixes:
  - `SqrtGCFPrefixStream2`
- Added bounded sqrt over CF prefixes:
  - `SqrtCFPrefixStream2`
- Added new-path APIs for:
  - GCF exact-tail sqrt
  - GCF range-seeded sqrt
  - CF range-seeded sqrt
  - midpoint/range-seeded wrapper layers

### Test cleanup and validation
- Added new canonical tests for:
  - sqrt exact/Newton core
  - sqrt seed/range helpers
  - sqrt API `*2` path
  - GCF exact-tail sqrt APIs
  - GCF range-seeded sqrt APIs
  - bounded sqrt streams
- Retired obsolete duplicated old sqrt tests, including:
  - old Newton/core tests
  - old range tests
  - residual-specific duplicate tests
  - selected old edge/helper tests
- Renamed remaining legacy-oriented sqrt test files to better reflect current coverage.
- Full suite remains fast and green.

## Important current conclusions
- The new canonical sqrt path is now strong enough to serve as the implementation center.
- Legacy exported sqrt names should remain temporarily as compatibility wrappers until the final naming strategy is chosen.
- The `*2` names have served their migration purpose, but should not necessarily be the final public naming scheme.
- The best near-term progress is no longer broad refactoring; it is implementing more complete sqrt operator behavior.
- Current bounded sqrt streams still “collapse to an exact bounded approximation first”; they are not yet true progressively certifying sqrt streams.

## Change of focus
- Set aside broad refactoring for now.
- Use the current canonical path as the stable base.
- Focus next on making `sqrt` behave more like a genuine operator/stream and less like “approximate then wrap.”
- Continue baby-step development:
  - stub
  - focused failing tests
  - smallest passing implementation
  - keep full suite green

## Highest-leverage next production task
- Advance `sqrt` from bounded exact-collapse streams toward more genuine progressive streaming/certification.
- The strongest immediate candidate is to begin exploring:
  - a diagonal/operator-aware sqrt stream path that certifies output more directly,
  - or a stricter transform-based sqrt stream using the newer `DiagBLFT` / GCF substrate instead of merely collapsing bounded approximations.

## Recommended immediate work items
1. Identify the smallest next operator-shaped sqrt milestone beyond bounded exact-collapse.
2. Add stubs and focused tests for that milestone.
3. Keep the new canonical/internal path as the only implementation target for new work.
4. Avoid re-expanding old wrapper surfaces unless compatibility requires it.
5. Keep all new work bounded and testable.

## Secondary cleanup tasks
- Later decide final naming strategy:
  - keep legacy names and hide/remove `*2`
  - or rename canonical pieces back to cleaner exported names
- Reassess whether `sqrt_cf.go` should be split further or retired as compatibility-only.
- Reassess whether some remaining legacy sqrt tests should be merged or renamed again.
- Document current semantics differences between:
  - exact-tail GCF evaluation
  - bounded GCF prefix approximation
  - CF prefix approximation

## Deferred for now
- Full retirement of all legacy sqrt wrapper names.
- Large API renaming wave.
- True infinite certified sqrt streaming over uncertified tails.
- More aggressive transform-based sqrt certification beyond current bounded approximation collapse.
- Decimal digit emission support.
- Large documentation rewrite unrelated to current operator progress.

## Ongoing design principles
- Prefer fewer medium-sized cohesive files over many tiny files.
- When a changed file is small, replace the whole file; otherwise replace complete functions.
- If a function gets too large, break it into smaller helpers.
- Keep tests bounded and fast.
- Prefer mathematically justified behavior over convenience shortcuts.
- Use compatibility wrappers only when they clearly protect a public surface during migration.
- Migrate in baby steps and exploit the full test suite aggressively.

## Future work note
- Add support for emitting decimal digits in addition to CF terms.
- This is low priority and should remain near the end of the roadmap, after core sqrt/operator progress is in better shape.

## Current architectural perspective
- Canonical sqrt implementation now flows roughly as:
  - exact/Newton core
  - seed/range helpers
  - canonical internal sqrt API
  - canonical internal source/range-seeded API
  - compatibility wrappers / bounded stream surfaces
- New work should target the canonical/internal layers first, then expose wrapper APIs as needed.

## Suggested next-chat bootstrap
Ask for current contents of:
- the current canonical sqrt stream / sqrt operator files
- any newly added sqrt stream tests
- the current `DiagBLFT` / diagonal-stream related files if the next step will involve progressive certification

Then focus on:
- actual sqrt operator progress
- bounded, test-first baby steps
- not another large refactoring wave first

# EOF MasterPlan.md