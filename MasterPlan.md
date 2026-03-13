# MasterPlan.md

## Project Goal

Build a mathematically correct, test-driven continued-fraction arithmetic library centered on Gosper-style streaming transforms, exact rational foundations, conservative enclosure reasoning, and eventual exact-real style workflows.

Stretch goal: make Gosper smile.

---

## Current Status

### Core arithmetic and representations
- Exact `Rational` foundation is in place.
- Regular CF finite streaming via `RationalCF` is working.
- Core interval/enclosure type `Range` is in place, with inside/outside semantics and refinement metrics.
- `Bounder` for regular CF prefixes is implemented and well tested.
- `GCFBounder` exists and supports:
  - exact convergents for finite GCF prefixes
  - placeholder point ranges
  - lower-bound ray enclosure
  - explicit tail-range enclosure
- ULFT / BLFT / DiagBLFT exact rational application is implemented.

### Streaming engines
- `ULFTStream`, `BLFTStream`, and `DiagBLFTStream` exist.
- Recent work substantially improved correctness for:
  - exact-point termination
  - remainder-pole clean exhaustion after emitted digits
  - distinguish real exact-input singularities from clean exhaustion
  - refine-before-fail behavior on coarse intervals
  - shared refine-budget accounting
- Stream `Next()` methods have started to be refactored into smaller helpers.

### Generalized continued fractions
- GCF finite prefix ingestion is implemented.
- GCF exact convergents and bounded inspection helpers are implemented.
- Prefix-aware specialized enclosure helpers exist for:
  - Lambert pi/4
  - Brouncker 4/pi
- Specialized inspection helpers are implemented.
- Generic and specialized GCF helper layers have begun to be consolidated.
- `GCFStream` finite-first behavior is implemented and currently passes tests for:
  - finite `SliceGCF` equivalence vs `EvaluateFiniteGCF`
  - adapted regular CF round-trip through `AdaptCFToGCF`
  - bounded exact termination semantics
  - generic stable `TailRange()`-driven earlier emission
- `GCFStream` now queries `TailRange()` dynamically from source state rather than caching it at construction time.

### Sources
- Regular CF sources:
  - `SliceCF`
  - periodic irrational-style sources including sqrt examples
- GCF sources:
  - `SliceGCF`
  - `PeriodicGCF`
  - adapted regular CF to GCF bridge
  - Lambert pi/4 source
  - Brouncker 4/pi source
  - algorithmic/unit-pattern GCF sources
- Note: many juicy future GCF sources remain, including pi and other constants.
  - https://en.wikipedia.org/wiki/Physical_constant

### Testing and quality
- Coverage is now above 81%.
- Property tests exist for:
  - ULFTStream vs exact rational image
  - BLFTStream vs exact rational image
  - DiagBLFTStream vs exact rational image (within currently supported subclass)
- Recent focused regression tests were added for:
  - exact input poles
  - remainder-pole clean exhaustion
  - refine-budget guards
  - GCFBounder range semantics
  - RationalCF
  - Rational.Div
  - BLFT denominator helpers
  - Range floor/string/refinement behavior
  - Bounder semantics
  - `GCFStream` generic stable `TailRange()` cadence / early emission semantics
- First fuzz target added for `ULFTStream`.
- Future note: expand fuzz coverage later, especially for BLFTStream, DiagBLFTStream, and other core arithmetic/streaming components.

### Documentation and API
- User guide skeleton exists.
- `api_spec.md` and `UserGuide.md` currently overlap; later either combine them or sharpen the distinction and remove the redundancy.
- Hold off on more documentation updates for now while production code and tests improve.
- Public API is not considered frozen yet; continue improving client ergonomics where it does not compromise mathematical clarity.

---

## Design / Style Rules

- Favor mathematical correctness over cleverness.
- Prefer test-first work when appropriate.
- Under normal circumstances, keep methods around 30–40 lines so they fit on one screen.
- Break longer methods into smaller helpers for readability, self-documentation, testability, and safer refactoring.
- Continue reviewing production code for DRY opportunities and consolidation where it improves maintainability without obscuring the math.
- Treat current source and tests as ground truth over stale documentation.
- Prefer dynamic, prefix-sensitive tail evidence over cached metadata when source contracts are stateful.

---

## Recent Progress Highlights

### Stream correctness
- Fixed ULFT exact-integer termination issues.
- Fixed ULFT exact-point remainder-pole clean exhaustion semantics.
- Fixed BLFT transient rectangle-pole behavior by refining before failing.
- Fixed BLFT exact-point remainder-pole clean exhaustion semantics.
- Fixed DiagBLFT refine-before-fail behavior for coarse unsupported intervals.
- Fixed DiagBLFT exact-point remainder-pole clean exhaustion semantics.
- Fixed `GCFStream` to use current source-provided `TailRange()` evidence dynamically instead of stale constructor-cached range metadata.
- Corrected a `GCFStream` stable-tail-range cadence test fixture to match the actual regular CF of the `(1,1)`-forever GCF source (golden-ratio behavior).

### Refactoring progress
- Factored shared exact-point exhaustion policy.
- Factored shared refine-budget accounting across stream engines.
- Refactored `ULFTStream.Next()` into smaller private helpers.
- Refactored `BLFTStream.Next()` into smaller private helpers.
- Refactored `DiagBLFTStream.Next()` into smaller private helpers.
- Factored shared GCFApprox construction from prepared bounders.
- Factored shared GCFInspect construction from GCFApprox.
- Factored shared positive-prefix validation.
- Factored shared bounded-prefix ingestion loop.

---

## Immediate Next Milestone

### Harden `GCFStream` finite-first semantics and prepare for unfinished-tail support

`GCFStream` now exists and passes its current finite-first and early-emission tests. The next step is to harden the semantics around tail evidence, cadence, and source contracts before expanding unfinished-tail support.

### Why
Generalized continued fractions are now not only:
- inspectable
- prefix-comparable
- enclosure-aware

but also capable of finite-first lazy ordinary CF production.

The next leverage comes from making the finite-first implementation more robust and from clarifying the boundary between:
- exact finite fallback
- dynamic tail metadata
- weaker lower-bound-ray fallback
- future unfinished-tail enclosure streaming

### Phase 1a
Strengthen tests around current `GCFStream` behavior:
- dynamic prefix-sensitive `TailRange()` semantics
- preference of explicit `TailRange()` over lower-bound ray
- cadence rules for metadata-driven emission
- finite exact fallback after source exhaustion
- clean separation between reusable strong evidence and weaker conservative evidence

### Phase 1b
Review the public/client-facing `GCFStream` API:
- confirm whether `Next() (int64, bool)` + `Err()` is the right long-term surface
- compare ergonomics and semantics with the other stream engines
- add small helper constructors or naming improvements if they materially improve clarity

---

## Near-Term Plan After Finite GCFStream

1. Harden finite `GCFStream` semantics with more targeted tests.
2. Add additional adapted regular-CF round-trip and cadence tests.
3. Add named finite-prefix fixtures for Lambert and Brouncker.
4. Extend `GCFStream` from finite-only behavior toward unfinished-tail enclosure support.
5. Reuse existing GCF tail metadata ideas where mathematically justified.
6. Let named sources benefit automatically once unfinished-tail streaming works.

---

## Medium-Term Plan

### GCF / exact-real direction
- Extend `GCFStream` to infinite or unfinished sources using conservative tail enclosures.
- Improve GCF enclosure semantics where lower-bound ray logic is too weak.
- Add more published GCF sources for pi and other constants.
- Investigate additional generalized continued fraction algorithms for pi and other constants as future `GCFSource` implementations.

### Unary / client-facing operator support
- Consider useful client-facing unary operators where they give real leverage beyond current ULFT support.
- Brainstorm convenience constructors/helpers for common ULFT matrices such as identity `(1,0,0,1)` and other useful named transforms.

### sqrt direction
- Continue improving sqrt approximation / source workflows once GCF streaming is stronger.
- Preserve focus on exact rational and bounded semantics rather than opaque approximation.

---

## Coverage / Hardening Plan

Keep increasing coverage strategically, prioritizing:
- branch-heavy core streaming logic
- transform safety / pole behavior
- exact termination semantics
- finite rational equivalence to exact arithmetic
- GCF tail-evidence semantics and fallback boundaries

Future targeted fuzzing:
- BLFTStream
- DiagBLFTStream
- Range / BLFT denominator helpers
- additional arithmetic and streaming invariants
- later: `GCFStream` once its unfinished-tail contract stabilizes

---

## Documentation Later

When production code stabilizes further:
- refresh API surface view if needed
- expand `UserGuide.md`
- reconcile `api_spec.md` vs `UserGuide.md`
- add an API user guide covering public interfaces, usage notes, warnings, and examples

---

## Current Focus

Do not broaden documentation work right now.
Do not chase easy cosmetic coverage.
Keep making meaningful progress toward:
- true generalized-CF streaming
- stronger exact arithmetic foundations
- transform correctness
- client-usable API improvements where justified
- Gosper-utopia