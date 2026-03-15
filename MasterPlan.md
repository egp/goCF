# MasterPlan.md

# goCF Master Plan

## Project end goal
Build a mathematically correct, testable, GCF-native continued-fraction arithmetic library that can:

- ingest finite and infinite generalized continued fractions
- compose unary and binary operators conservatively
- emit regular CF terms only when certified
- support progressively richer exact-real expressions

Representative long-term target expression:

    sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))

This target is intentionally ambitious. It implies the architecture must support:

- named constant sources (`e`, `pi`, algebraic constants)
- unary operators (`sqrt`, later `sin`, `tanh`, etc.)
- binary operators (`+`, `-`, `*`, `/`)
- transform-based composition
- conservative refinement and certified emission
- a composition model that scales beyond one-off operator implementations

## Core strategic direction
- Mathematical correctness first.
- Design for Test (DfT) first.
- Use test-first / TDD for new production work.
- Keep the canonical implementation centered on GCF-native ingestion and transform-based arithmetic.
- Prefer bounded, testable, incremental progress over broad speculative refactors.
- Keep production code cohesive, and retire obsolete CF-first legacy code once the canonical path is proven.
- Since there are no library clients yet, improve the API deliberately rather than preserving migration shapes longer than necessary.

## Current architectural conclusion
The project has now moved far enough that the canonical internal substrate should be treated as:

- `GCFSource`
- exact finite-GCF semantics
- bounded GCF-prefix semantics
- transform-based arithmetic (`ULFT`, `BLFT`, `DiagBLFT`)
- conservative range enclosure
- certified regular-CF emission from enclosures / transform state

Older regular-CF-first paths should now be treated as compatibility layers, helpers, or retirement candidates unless they remain clearly useful.

## Immediate top priorities
1. Freeze exact finite-GCF semantics across all canonical finite paths.
2. Inventory production files into:
   - canonical
   - compatibility / migration
   - obsolete / retirement candidates
3. Promote `GCFStream` as the reference certifying emission engine for native GCF ingestion.
4. Make reciprocal the reference unary operator for the generic progressive-operator model.
5. Recast sqrt as a client of the generic unary framework rather than the driver of architecture.

## Why the priority shifted
Recent work on sqrt was useful and produced important machinery, but the true project bottleneck is no longer “make sqrt fancier.”

The real bottleneck is the operator substrate needed to eventually support expressions like:

    sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))

That substrate requires:

- unambiguous finite semantics
- native GCF ingestion
- reusable certified-progress machinery
- reliable binary composition
- generic unary-operator structure

So the focus now shifts from local sqrt evolution to the general transform/certification architecture.

## Current status summary
- The codebase now has a credible GCF-native core.
- Exact finite-GCF ingestion and evaluation exist.
- Bounded prefix ingestion exists.
- Transform-based arithmetic exists for:
  - `ULFT`
  - `BLFT`
  - `DiagBLFT`
- Certified-progress machinery now exists in reusable pieces:
  - conservative range propagation
  - certified CF range emitter
  - remainder update helpers
- A progressive GCF sqrt path exists and the recent hang has been fixed.
- The full suite is green again.

## Important recent conclusions
- The recent progressive sqrt hang was a stream-control bug, not a proof that the overall direction was wrong.
- The subsequent sqrt(2) test failure revealed a test expectation mismatch (`sqrt(sqrt(2))` versus `sqrt(2)`), which reinforces the value of exact semantic tests.
- The finite-core GCF semantics now appear strong enough to freeze as a foundational layer.
- Before more operator work, the project should explicitly lock down exact finite meaning and identify legacy code that no longer deserves to survive.

## Phase 1: freeze exact finite-GCF semantics
This is the highest-leverage immediate phase.

### Goals
- Make exact finite-GCF behavior completely explicit and test-backed.
- Ensure all canonical finite paths agree on the same value.
- Separate exact finite meaning from unfinished-tail family meaning.

### Canonical exact finite meaning
For a finite generalized continued fraction under:

    x = p + q/x'

the terminal convention is:

    the last emitted term contributes just p_last

Equivalently, if the finite prefix composes to a transform:

    T(x) = (A*x + B) / (C*x + D)

then the finite value is:

    T(infinity) = A/C

provided `C != 0`.

### Canonical finite paths that must agree
- `EvaluateFiniteGCF`
- `IngestAllGCF(...).Convergent()`
- `IngestAllGCF(...).Range()` after finish
- finite `GCFStream` fallback / emission behavior

### Immediate work items
1. Add semantic contract tests proving all finite paths agree on the same fixtures.
2. Strengthen tests around:
   - empty sources
   - single-term sources
   - multi-term sources
   - mixed-q fixtures
   - bounded prefixes
   - early exhaustion
3. Make the distinction explicit between:
   - finite truncation value
   - unfinished-tail family enclosure
4. Verify that unfinished placeholder point ranges are never mistaken for true family enclosures.

### Exit criteria
- Exact finite-GCF behavior is boring, explicit, and trusted.
- Any future disagreement among finite paths becomes a regression caught by tests.

## Phase 2: inventory and purge legacy orientation
### Goals
- Identify which files are still CF-first legacy artifacts.
- Decide what should remain canonical and what should be retired.

### Work items
- Inventory production files into:
  - canonical
  - compatibility / migration
  - obsolete
- Avoid deleting useful code merely because it is older.
- Remove code only when:
  - the canonical replacement is proven
  - tests cover the intended behavior
  - deletion simplifies the model rather than hiding it

### Exit criteria
- New work has one obvious canonical target.
- Legacy wrappers stop expanding.

## Phase 3: standardize tail evidence and certified progress
### Goals
Turn unfinished-tail evidence into a clean reusable foundation instead of an operator-specific trick.

### Work items
- Standardize the meaning of:
  - lower-bound tail evidence
  - explicit tail-range evidence
  - reusable versus one-shot evidence
  - same-prefix refinement evidence
  - post-emit evidence
- Strengthen tests around:
  - fallback ordering
  - pole handling
  - refinement stop conditions
  - same-prefix certification behavior

### Exit criteria
- Tail evidence is trustworthy and reusable across multiple streams/operators.

## Phase 4: define the generic progressive-operator model
### Goals
Create one reusable conceptual model for unary progressive operators.

### Generic operator model
A progressive unary operator should expose or internally manage:

1. current source approximation state
2. conservative enclosure for current value
3. certified output-prefix state
4. transformed remainder state
5. refinement / ingestion policy
6. clean stop / fail / exhaustion semantics

### Intended clients
- reciprocal
- sqrt
- later `sin`
- later `tanh`
- possibly other analytic operators

### Exit criteria
- Reciprocal and sqrt can both be described as instances of the same machine.

## Phase 5: make reciprocal the reference unary operator
### Why reciprocal first
- simpler than sqrt
- tightly tied to CF remainder logic
- a better reference model for generic unary progress
- directly useful to many later operators

### Work items
- Harden exact-tail reciprocal path
- Harden bounded-prefix reciprocal path
- Ensure same-prefix certification/refinement is test-backed
- Prefer persistent remainder state over replay where justified

### Exit criteria
- Reciprocal becomes the clean reference unary operator implementation.

## Phase 6: recast sqrt as a client of the unary framework
### Goals
Keep sqrt important, but stop letting it dominate architecture.

### Work items
- Preserve current green behavior
- Refactor only where it clarifies the shared unary-operator model
- Keep sqrt-specific logic local:
  - seed range
  - conservative sqrt enclosure
  - exact-square fast path

### Exit criteria
- Sqrt looks like one operator with a specific enclosure strategy, not a separate subsystem.

## Phase 7: harden binary transform arithmetic
### Goals
Make binary composition robust enough to support serious expression building.

### Work items
- Audit and strengthen:
  - `BLFT`
  - `DiagBLFT`
  - GCF exact-tail and bounded-prefix binary paths
- Ensure exact finite and unfinished-tail semantics remain explicit
- Consolidate duplicate binary logic where safe

### Exit criteria
- `+`, `-`, `*`, `/` are reliable on the canonical GCF-native substrate.

## Phase 8: add an internal expression-composition layer
### Goals
Support complex expressions without hand-wiring each operator chain.

### Work items
- Introduce internal concepts such as:
  - source node
  - unary node
  - binary node
- Keep this internal until the external API is mature
- Favor composability over cleverness

### Exit criteria
- Expressions can be assembled mechanically from reusable parts.

## Phase 9: constants needed for real expressions
### Immediate constant priorities
1. `e`
2. `pi`
3. selected algebraic constants
4. later derived constants

### Work items
- Harden `e` source behavior
- Choose / harden canonical `pi` strategy
- Provide useful tail evidence where reasonable

### Exit criteria
- The library can construct the constants required by the target expression.

## Phase 10: decide trig semantics before trig implementation
Before implementing `sin(69°)`, decide:

- radians versus degrees API policy
- exact representation of angles
- whether integer literals imply radians or require explicit angle constructors
- how trig operators accept inputs (rational, CF, GCF, range-seeded, etc.)

### Exit criteria
- Trig semantics are explicit before implementation begins.

## Phase 11: add transcendental unary operators
Only after the unary framework is convincing should the project add:

- `sin`
- `tanh`

### Rules for this phase
- bounded, conservative, test-first
- no speculative giant leap
- use stubs and failing tests first
- keep correctness and DfT dominant

### Exit criteria
- These operators extend the framework rather than distorting it.

## Phase 12: API finalization and cleanup
### Goals
Use the no-client window to improve the external API deliberately.

### Work items
- Retire migration wrappers that no longer pay rent
- Remove temporary naming such as transitional `*2` forms where appropriate
- Simplify exported surface around canonical abstractions
- Merge/remove obsolete tests after coverage is preserved
- Update user-facing docs only after the API shape stabilizes

### Exit criteria
- Public API matches the actual architecture.
- The codebase is smaller, clearer, and more canonical.

## Current critical path
1. Freeze exact finite-GCF semantics with stronger contract tests.
2. Inventory canonical versus obsolete code.
3. Promote `GCFStream` as the reference certifying engine.
4. Make reciprocal the reference unary operator.
5. Refit sqrt into the generic unary design.
6. Harden binary composition.
7. Add expression composition.
8. Harden `e` and `pi`.
9. Decide trig semantics.
10. Add `sin`, then `tanh`.
11. Finalize API and retire legacy scaffolding.

## Immediate next task
Create a new finite-core semantic contract test layer proving agreement among:

- `EvaluateFiniteGCF`
- `IngestAllGCF(...).Convergent()`
- exact finished `Range()`
- finite `GCFStream` emission

This should be done test-first, with no production changes unless the tests reveal semantic drift.

## Ongoing design principles
- Mathematical correctness first.
- DfT and TDD first.
- Use stubs where appropriate.
- Prefer fewer medium-sized cohesive files over many tiny files.
- Replace small files wholesale when changing them.
- For larger files, replace complete functions.
- Keep tests bounded and fast.
- Do not speculate about stale implementation details; resync from code.
- Prefer canonical internal improvements over preserving old wrapper shapes.
- Remove legacy code only after the new path is proven.

# EOF MasterPlan.md