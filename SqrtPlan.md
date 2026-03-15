# MasterPlan.md

# goCF Master Plan

## Project mission
Build a mathematically correct, testable, Gosper-inspired continued-fraction arithmetic library with strong support for generalized continued fractions (GCFs), unary/binary transform operators, and eventually a real streaming `sqrt` operator that emits certified digits instead of merely collapsing bounded approximations to rationals.

---

## Current strategic position

### What is now substantially complete
The project has crossed an important threshold:

- the old sprawling sqrt experimentation has largely been consolidated onto a newer canonical/internal path
- bounded and inspectable sqrt streams now exist for:
  - CF prefix input
  - GCF prefix input
  - GCF exact-tail input
- exact-tail transform streams exist for:
  - ULFT
  - DiagBLFT
  - BLFT
- a public unary reciprocal operator now exists for:
  - GCF exact-tail input
  - GCF bounded-prefix input
- a first proof-safe conservative sqrt enclosure engine exists
- a first certified-progressive CF sqrt stream exists and can refine input and continue emitting certified digits on the CF side

### What is not yet complete
The library does **not yet** have a finished general streaming sqrt operator.

The current gap is now sharply defined:

- CF-side certified-progressive sqrt has meaningful progress
- GCF-side certified-progressive sqrt exists in draft form but currently hangs / fails to progress cleanly
- repeated certified emission still rebuilds/replays rather than carrying a principled persistent transformed remainder state
- diagonal / transform-driven sqrt has not yet become the main engine
- naming/API cleanup and final retirement of compatibility layers remain unfinished

This is now a much better problem than before:
the remaining work is not “mess cleanup,” it is “finish the real operator.”

---

## Completed accomplishments

### 1. Exact-tail transform substrate
Completed exact-tail bounded transform streams and associated support for:

- `GCFULFTStream`
- `GCFDiagStream`
- `GCFBLFTStream`

Including:
- exact-tail source plumbing
- bounded ingestion
- shared helper extraction to reduce duplication
- separate x/y ingest bounds for BLFT
- stronger tests for exact rational image matching

### 2. sqrt canonicalization and migration
Created/established the newer canonical sqrt substrate:

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

### 3. Legacy simplification / retirement
Substantial retirement of duplicate or obsolete sqrt code/tests:

- old exact/Newton implementation bodies reduced to wrappers or retired
- old range implementations reduced to wrappers or retired
- dead helpers such as `exactSqrtBig` removed
- obsolete/duplicated old tests retired
- remaining legacy tests trimmed toward public-surface regression coverage
- several misleading test filenames renamed to reflect current purpose

### 4. Public bounded sqrt operator surfaces
Public/user-facing bounded sqrt stream constructors now exist:

- `SqrtStream(...)`
- `SqrtGCFStream(...)`
- `SqrtGCFExactTailStream(...)`

All with inspectable stream state via snapshots and stream status.

### 5. Stream introspection / semantics
Added inspectable stream interfaces and snapshots that expose:

- started / prefix terms
- bounded approximation used
- input approximation object where applicable
- stream status

Current meaningful statuses include:

- `unstarted`
- `exact_input`
- `bounded_collapse`
- `certified_progressive`
- `failed`

### 6. Proof-safe sqrt enclosure engine
A first proof-safe conservative sqrt enclosure engine now exists:

- `SqrtLowerBoundRational`
- `SqrtUpperBoundRational`
- `SqrtRangeConservative`

Current implementation:
- exact-square fast paths
- negative rejection
- proof-safe scaled integer bracketing for non-square rationals
- conservative enclosure over nonnegative inside ranges

This was the key blocker on the critical path and is now in place.

### 7. Unary reciprocal operator
Unary reciprocal is now on the operator path and should be considered a permanent work item in the design:

- public exact-tail GCF reciprocal stream
- public bounded GCF-prefix reciprocal stream
- reciprocal range operator for proof-safe range work

This operator is strategically important because sqrt remainder updates naturally involve reciprocal structure.

### 8. Certified-progressive sqrt (CF side)
A first real certified-progressive sqrt stream exists for CF-prefix input:

- conservative sqrt range
- certified digit extraction from ranges
- persistent certified CF range emitter
- refinement of input when current certification is exhausted
- continued emission without collapsing immediately to a single rational approximation

This is the first genuine step beyond “bounded collapse.”

---

## Current status summary

### Stable and working
These areas are in good shape:

- exact-tail transform streams
- canonical bounded sqrt substrate
- public bounded sqrt stream APIs
- proof-safe sqrt enclosure engine
- reciprocal exact-tail / bounded-prefix operators
- CF-side certified-progressive sqrt prototype

### Unstable / current blocker
The current unstable area is:

- **GCF-side certified-progressive sqrt stream**

Latest attempt:
- stream drafted
- compile issues resolved by switching to `NextPQ()` / `IngestPQ()`
- then test `TestSqrtCertifiedGCFPrefixStream_Sqrt2_RefinesAndContinues` hung

Most likely cause:
- current unfinished-tail enclosure quality on the GCF side is not yet strong or monotone enough for this progressive-certification strategy
- replay/rebuild logic may be repeatedly rebuilding with no new certifiable progress
- this is a signal that the next work should be on **GCF unfinished-tail enclosure strength / progress guarantees**, not further patching at the surface

### Meaning of the current moment
This is the right time for a clean break:
- broad refactoring is no longer the right focus
- the architecture is good enough
- the real remaining work is to finish general sqrt as an operator

---

## The path to general sqrt

By “general sqrt” here, the target is:

- streaming sqrt over CF and GCF input
- digits emitted because they are certified
- not merely because a bounded rational approximation was first collapsed and then expanded
- principled remainder/update progression
- eventually diagonal / transform-driven implementation rather than repeated ad hoc replay

### Required layers for general sqrt

#### Layer A — proof-safe enclosure
Needed:
- conservative one-sided bounds for `sqrt(r)` for rational `r`
- conservative `sqrt(range)`

Status:
- **done (first usable version)**

#### Layer B — certified digit extraction from a conservative output range
Needed:
- repeatedly extract CF digits from a certified range
- remainder update via subtract + reciprocal

Status:
- **done (generic certified range certifier / emitter exists)**

#### Layer C — unary reciprocal as a real operator
Needed because CF remainder update naturally includes reciprocal structure.

Status:
- **done (exact-tail and bounded GCF-prefix versions exist, plus reciprocal range operator)**

#### Layer D — persistent progressive emission with refinement
Needed:
- keep emitting certified digits
- when current range cannot certify more, refine input
- resume without semantic break

Status:
- **partially done on CF side**
- **not yet robust on GCF side**

#### Layer E — stable GCF unfinished-tail enclosure and progress discipline
Needed:
- GCF-side progressive certification must use unfinished-tail metadata/ranges that are strong enough to produce monotone useful progress
- must avoid hangs / rebuild loops with no useful progress

Status:
- **not done**
- this is now the main blocker for general sqrt across GCF input

#### Layer F — persistent transformed remainder state
Needed:
- avoid replaying certified prefix from scratch on each refinement
- carry a principled remainder/operator state forward

Status:
- **not done**
- current implementation replays/rebuilds

#### Layer G — diagonal / transform-driven sqrt engine
Needed:
- make the remainder/update and emission engine live naturally in transform/diagonal machinery
- move from “range certifier with rebuild” toward a true operator engine

Status:
- **not done**
- this is the long-term elegant implementation path

#### Layer H — public API convergence / cleanup
Needed:
- finalize naming between old names / `*2` names / canonical names
- retire wrappers once final semantics are stable

Status:
- **not done**
- intentionally deferred until the operator settles

---

## Critical path to complete the entire Master Plan

This is the current priority order for the actual critical path.

### Priority 1
**Stabilize and strengthen GCF unfinished-tail enclosure / progress guarantees**

This is now the most important blocker because:
- CF-side progressive sqrt already exists
- GCF-side progressive sqrt hangs
- general sqrt over the intended substrate needs GCF-side progressive refinement to be trustworthy

Concretely, this probably means:
- inspect current `GCFBounder.Range()` behavior on unfinished sources
- improve unfinished-tail enclosure quality
- ensure refinement cannot loop forever without either certifying more or exhausting the allowed bound

### Priority 2
**Finish multi-digit certified-progressive sqrt on the GCF side**

After the unfinished-tail enclosure issue is fixed:
- make `SqrtCertifiedGCFPrefixStream` actually refine and continue
- ensure prefix stability and monotone progress
- get it to emit meaningful multi-digit prefixes on representative sources like adapted `sqrt(2)`

### Priority 3
**Replace replay-on-refinement with persistent transformed remainder state**

Current progressive streams rebuild and replay certification after refinement.

Need:
- a proper carried remainder state
- less recomputation
- better structure for later diagonal integration

### Priority 4
**Drive remainder/update and emission through diagonal / transform machinery**

This is the real long-term operator goal:
- transform-driven remainder state
- principled digit emission
- not merely external range re-wrapping

### Priority 5
**Unify progressive sqrt architecture across CF and GCF input**

By this point:
- CF progressive sqrt
- GCF progressive sqrt
- common interfaces / semantics / snapshots
- minimal duplication

### Priority 6
**Finalize unary reciprocal as a first-class operator family**

Already started, but should be completed as part of the operator story:
- bounded/exact-tail paths polished
- integrated conceptually into unary operator architecture
- positioned as a reusable primitive for sqrt and future unary operators

### Priority 7
**Finalize naming / retire wrappers / retire duplicate tests / document semantics**

Only after the real operator shape is stable:
- final naming decisions
- compatibility cleanup
- documentation of exact semantics
- removal of obsolete transitional layers

---

## Short critical-path graph

The shortest dependency chain to sqrt/operator completion now looks like:

1. proof-safe sqrt enclosure engine — DONE
2. certified CF range emitter — DONE
3. unary reciprocal operator — DONE
4. stabilize GCF unfinished-tail enclosure / progress guarantees — TODO
5. working multi-digit progressive GCF sqrt — TODO
6. persistent transformed remainder state — TODO
7. diagonal / transform-driven sqrt engine — TODO
8. API convergence / wrapper retirement / docs — TODO

---

## Recommended immediate next technical focus

### Do next
Investigate and strengthen **GCF unfinished-tail enclosure behavior**.

This should likely include:
- focused tests on `GCFBounder.Range()` for unfinished adapted CF/GCF sources
- prove that refinement either:
  - certifies additional digits, or
  - narrows the enclosure, or
  - terminates cleanly at the configured bound
- no hangs / no progress loops

### Do not do next
Do **not**:
- spend time on more wrapper additions
- do another naming cleanup wave
- broaden documentation before the operator shape is stable
- add decimal-digit output yet

---

## Detailed plan to complete general sqrt

### Phase 1 — make progressive GCF sqrt real
1. Reproduce and isolate the hang in `TestSqrtCertifiedGCFPrefixStream_Sqrt2_RefinesAndContinues`
2. Inspect `GCFApprox` / `GCFBounder.Range()` on unfinished adapted `sqrt(2)` input
3. Improve GCF unfinished-tail enclosure strength or install progress checks
4. Get the GCF progressive sqrt stream to emit a stable multi-digit certified prefix

### Phase 2 — remove replay/rebuild weakness
5. Replace “rebuild emitter and replay certified prefix” with a carried transformed remainder state
6. Preserve monotone prefix correctness without full replay
7. Ensure refinement updates the carried state coherently

### Phase 3 — move toward transform-driven sqrt
8. Recast the certified remainder progression using operator/transform state
9. Integrate diagonal or transform machinery where it simplifies repeated remainder/update
10. Reduce ad hoc range plumbing around the progression loop

### Phase 4 — converge the public operator
11. Confirm unified semantics for CF and GCF progressive sqrt
12. Decide final public names and which compatibility wrappers remain
13. Retire redundant wrappers/tests
14. Document exact semantics, guarantees, and limitations

---

## Remaining non-critical but important work

### Public unary reciprocal
Keep in the long-term design:
- unary reciprocal (ingest GCF, emit reciprocal)
- it belongs before full multi-digit sqrt completion because sqrt remainder progression naturally depends on reciprocal structure

### Decimal digit emission
Add later, low priority:
- emit decimal digits instead of or in addition to CF digits
- not on the critical path

### Broader unary operator family
Later possible work:
- use reciprocal and sqrt structure to shape a wider unary-operator architecture

---

## What “sqrt MVP” should now mean

A reasonable sqrt MVP target is:

- proof-safe enclosure engine
- public bounded and progressive sqrt APIs
- CF and GCF progressive sqrt both working for representative nontrivial cases
- multi-digit certified emission
- no hangs / no hidden collapse-only behavior where progressive status is claimed
- inspectable semantics
- unary reciprocal operator available as part of the operator toolkit

This is now much closer than before, but not yet complete.

---

## Practical next-chat instructions for myself
When resuming:
- do not restart broad refactoring
- focus first on the GCF progressive sqrt hang
- inspect unfinished-tail enclosure quality before changing the progressive stream surface
- keep changes medium-sized and operator-focused
- prefer complete-function replacements for large files
- run the full suite frequently (`go test ./cf`) because it is fast

# EOF MasterPlan.md