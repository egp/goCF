# MasterPlan.md v3

# Master Plan

## Mission

Build a mathematically trustworthy continued-fraction arithmetic library in Go,
with a design trajectory that moves steadily toward Gosper-style streaming
exact-real arithmetic.

Current rule for prioritization:

- prefer highest-leverage infrastructure over isolated features
- prefer exactness and composability over convenience shortcuts
- prefer shared ingestion / transform machinery over one-off algorithms
- prefer improvements that bring us closer to safe streaming arithmetic

---

## P0. Guiding principles

### P0.1 Correctness first
- Exact rational arithmetic remains the ground truth.
- All transform rewrites must be algebraically verified.
- Streaming digit emission must remain conservative.

### P0.2 Ingest before emit
- For new mathematical objects, ingestion comes before output streaming.
- GCF work follows this rule explicitly.

### P0.3 Shared machinery over special cases
- Prefer reusable transform and bounder abstractions.
- Avoid proliferating narrowly useful one-off APIs unless they unlock a major path.

### P0.4 Fewer cohesive files
- Keep production code in medium-sized coherent files.
- Avoid fragmenting into many tiny files.

---

## P1. Highest-leverage current priority: GCF enclosure semantics

### Why this is first
We now have a real GCF subsystem:

- `GCFSource`
- `SliceGCF`, `PeriodicGCF`, `FuncGCFSource`
- regular CF to GCF adapter
- exact finite GCF evaluation
- forward ULFT / diagonal GCF ingestion
- `GCFBounder`
- `GCFApprox`
- named algorithmic GCF sources:
  - `e` via regular-CF-as-GCF
  - Brouncker for `4/pi`
  - Lambert for `pi/4`

But the main missing piece is still:

- `GCFBounder.Range()` is only a point placeholder, not a conservative enclosure
  for unfinished infinite GCF prefixes

This is the biggest obstacle between “interesting GCF experiments” and
“generalized streaming arithmetic that might make Gosper smile”.

### Work items
- Design a principled enclosure model for unfinished GCF prefixes.
- Decide whether GCF ranges should:
  - mirror regular CF inside-interval semantics when possible, or
  - use a more general transform-derived enclosure mechanism.
- Prototype conservative prefix enclosure logic for a useful subclass of GCFs.
- Upgrade `GCFBounder.Range()` from point placeholder to real enclosure where sound.
- Add explicit documentation for any subclass restrictions.

### Success criteria
- Infinite GCF prefixes produce conservative, explainable enclosures.
- The design works for at least one real named nontrivial GCF source.
- The resulting API remains compatible with later streaming transform engines.

---

## P2. Second priority: unify regular CF and GCF ingestion mentally and structurally

### Why this is next
Regular CF and GCF ingestion are now close enough that the architecture should
make their relationship obvious.

This is high leverage because it reduces duplication and clarifies what the
library is really about: transform-driven exact arithmetic on structured streams.

### Work items
- Revisit `Bounder` and `GCFBounder` side by side.
- Identify what can be unified conceptually, and what must remain separate.
- Introduce a common documentation story for:
  - regular CF prefixes
  - GCF prefixes
  - convergents
  - enclosures
  - finite termination
- Add more equivalence tests:
  - regular CF vs adapted GCF
  - finite evaluation vs prefix convergent paths
- Consider a shared snapshot vocabulary:
  - `CFApprox`
  - `GCFApprox`
  - possibly later a common approximation interface

### Success criteria
- A user can understand how regular CF and GCF differ without reading internals.
- Core ingestion APIs feel like one family, not two unrelated subsystems.

---

## P3. Third priority: tighten the public API around recommended vs advanced usage

### Why this matters
The library now has a lot of public surface area, especially in sqrt and
bounded prefix inspection.

This is not yet a code crisis, but it is a comprehension and usability risk.

### Work items
- Separate “recommended starting points” from “advanced / experimental APIs”.
- Review exported names for overlap and unnecessary duplication.
- Keep compatibility wrappers where useful, but document them clearly as wrappers.
- Add a user-facing map from goals to APIs:
  - “I want exact rational arithmetic”
  - “I want to inspect CF prefixes”
  - “I want bounded sqrt approximations”
  - “I want to experiment with GCFs”
  - “I want transform streaming”

### Success criteria
- New users can find the right entry points quickly.
- Advanced users can still reach the deeper machinery without guesswork.

---

## P4. Fourth priority: user documentation and API guide

### Why this moved up
The package is now large enough that documentation is a leverage multiplier.
Good docs reduce misuse, reduce future refactoring churn, and make advanced work
easier to validate.

### Work items
- Draft and maintain `UserGuide.md`.
- Add a proper API user guide covering:
  - public interfaces
  - usage notes
  - warnings
  - examples
- Document exact vs bounded vs heuristic semantics explicitly.
- Document the GCF convention:
  - `x = p + q/x'`
- Add worked examples for:
  - regular CF finite and periodic sources
  - `Bounder`
  - transform streams
  - sqrt APIs
  - Brouncker and Lambert prefix inspection
  - regular CF adapted to GCF

### Success criteria
- Users can get productive without reading source code.
- Experimental features carry explicit warnings.
- The guide makes the architecture legible.

---

## P5. Fifth priority: stream engine hardening

### Why this remains important
Coverage work has improved confidence, but the stream engines are still among
the most delicate parts of the system:

- `ULFTStream`
- `BLFTStream`
- `DiagBLFTStream`

These are the pieces closest to Gosper-style streaming arithmetic.

### Work items
- Continue targeted test coverage improvements for:
  - hard error paths
  - cycle detection
  - exact integer termination cases
  - refinement guard behavior
- Revisit stream option semantics for clarity.
- Reduce any remaining surprising distinctions between clean termination and
  error termination.
- Keep diagnostics strong:
  - fingerprints
  - annotated errors
  - bounded refinement guards

### Success criteria
- Stream engines are easier to trust under failure and edge conditions.
- Diagnostics remain readable and actionable.
- Coverage continues to rise where it matters most.

---

## P6. Sixth priority: sqrt path refinement toward true streaming

### Current status
Implemented and green:

- exact rational sqrt detection
- Newton iteration machinery
- bounded residual / exact-stop policies
- bounded sqrt CF adapters
- CF-prefix-derived bounded sqrt paths
- exact and heuristic sqrt-range helpers

### Why it is not first right now
The sqrt subsystem is useful, but GCF enclosure and ingestion unification are
currently more leverage-rich for the overall architecture.

### Work items
- Compare bounded sqrt strategies more systematically:
  - convergent-targeted
  - range-seeded
  - midpoint-targeted
- Decide which prefix-based strategy is most promising.
- Keep heuristic range logic clearly separated from proof-safe logic.
- Revisit whether some sqrt APIs should be demoted to advanced/experimental docs.
- Later, design a true streaming sqrt operator if the transform/range machinery
  is strong enough.

### Success criteria
- The preferred bounded sqrt bridge is evidence-based.
- Future true streaming sqrt work starts from a clear best candidate path.

---

## P7. Seventh priority: more named published GCF sources

### Why this matters
Named sources provide:
- regression anchors
- mathematical credibility
- better design pressure than toy sources

We already have:
- `e` as regular-CF-as-GCF
- Brouncker for `4/pi`
- Lambert for `pi/4`

### Work items
- Add more published nontrivial GCF sources.
- Choose sources for:
  - mathematical importance
  - structural diversity
  - debug value
- Investigate GCF algorithms for:
  - `pi`
  - other constants
- Keep bounded prefix inspection tests for each named source.

### Success criteria
- The GCF subsystem is exercised by multiple genuinely different families.
- Prefix convergents and inspections are easy to compare and trust.

---

## P8. Convenience constructors and transform ergonomics

### Why this matters
Small ergonomic improvements can make experimentation much easier, especially
for transforms.

### Work items
- Brainstorm convenience constructors/helpers for common ULFT matrices, such as:
  - identity `(1,0,0,1)`
  - reciprocal
  - translate by integer
  - maybe other commonly used matrices
- Consider whether similar helpers make sense for diagonal transforms.
- Keep helpers small and mathematically transparent.

### Success criteria
- Common transform patterns are easier to express.
- Helper APIs reduce noise without hiding the math.

---

## P9. Coverage and trustworthiness maintenance

### Why this stays on the plan
Coverage is not the goal, but it is a useful pressure gauge.

### Work items
- Keep targeting low-coverage high-risk code, especially stream engines.
- Add smoke coverage for basic utility display functions where cheap.
- Prefer tests that cover:
  - exactness
  - edge conditions
  - failure semantics
  - equivalence across subsystems

### Success criteria
- Coverage rises in the right places.
- Tests improve confidence, not just percentages.

---

## P10. Stretch goal

Long-term stretch goal remains:

Compute exact-real style expressions such as:

- `sqrt(3/pi^2 + e) / tanh(sqrt(5) - sin(69 degrees))`
- `sqrt(3/π² + e) / tanh(sqrt(5) - sin(69))`
- $$\frac{\sqrt{3/\pi^2 + e}}{\tanh(\sqrt{5} - \sin(69^\circ))}$$

Which might eval to approx `1.77031957889` or as a CF like:

- `[2; 62, 1, 3, 1, 1, 5, 1, 1, 2, 1, 2, 11, 3, 3, 1, 2, 1, 12, 1, 6, 5, 2, 3, 13, 4, 1, 1, 3, 4, 129, 2, 1, 3, 1, 3, 1, 5, 1, 16, 1, 1, 6, 4, 9, 3, 1, 16, 1, 4, 1, 1, 1, 1, 4, 1, 2, 2, 1, 1, 1, 8, 3, 32, 1, 2, 3, 6, 1, 1, 1, 1, 2, 3, 1, 1, 5, 1, 4, 5, 2, 2, 7, 12, 1, 3, 1, 11, 1, 4, 6, 2, 15, 2, 12, 1, 1, 23, 2, 5, 1, 4, 167, 8, 2, 3, ...]`

This is the “Gosper smiling” target: exact-real style structured arithmetic with
credible streaming foundations.

---

## Immediate next actions

1. Design real conservative enclosure semantics for unfinished GCF prefixes.
2. Fold the resulting design back into `GCFBounder`.
3. Continue the user guide from skeleton to usable first draft.
4. Keep hardening stream engines with targeted tests.
5. Add another structurally different published GCF once enclosure semantics are
   clearer.

# MasterPlan.md v3