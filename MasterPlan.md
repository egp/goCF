# MasterPlan.md

# goCF Master Plan

## Purpose

Build a mathematically correct, test-driven continued-fraction arithmetic library in Go, centered on Gosper-style streaming arithmetic, generalized continued-fraction ingestion, and a path to exact-real style composition.

Current concrete MVP target:

    sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))

Current implementation strategy for the numerator:

    sqrt((3/16) * (4/pi)^2 + e)

using Brouncker's 4/pi source.

The full target formula remains assembled in test code for now. Production code should expose reusable sources, unary/binary machinery, bounded approximations, and certified ranges, not a public one-off formula API.

---

## Core Principles

1. Mathematical correctness first.
2. Test-first / DfT always.
3. Small, composable production pieces.
4. Prefer certified bounds over optimistic guesses.
5. Keep API flexible; no external clients yet, so improve the API while it is still cheap.
6. Remove legacy regular-CF-only assumptions when they no longer serve the GCF-first design.
7. Keep focus on the MVP critical path before broader cleanup or polish.

---

## Current Status Summary

### High-level status

The project has successfully crossed from “can ingest finite regular CFs” into a real GCF-first architecture with:

- finite GCF ingestion
- bounded GCF evaluation
- exact-tail GCF evaluation
- GCF-native unary reciprocal streams
- multiple sqrt paths
- source-specific infinite-source tail evidence
- first certified range for the full MVP target formula

The current blocker to a materially useful MVP result is **tightening the numerator and denominator bounds**, not basic architecture.

### Current MVP status

For the target:

    sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))

we now have:

- a bounded numerator point approximation
- a certified denominator range
- denominator range excludes zero
- a certified full-target quotient range

What is still missing for MVP quality:

- tighter denominator range
- tighter numerator approximation
- a materially sharper full-target certified range
- a final small end-to-end acceptance layer around the target formula

---

## What Has Been Accomplished Recently

### 1. GCF ingestion and evaluation foundation

Completed:

- `GCFBounder`
- `EvaluateFiniteGCF`
- `IngestAllGCF`
- `IngestGCFPrefix`
- shared bounded ingestion helper
- exact compose/eval with explicit tails

This is now solid enough to support unary and binary work.

### 2. GCF stream tail-evidence framework

Implemented:

- lower-bound ray support
- explicit tail range support
- candidate same-state evidence
- refined same-state evidence
- post-emit evidence
- source-driven tail evidence model

This removed the prior hang and gave a much better base for infinite GCF source behavior.

### 3. Brouncker and Lambert source evidence

Implemented and stabilized:

- Brouncker 4/pi source
- Lambert pi/4 source
- source-specific tail lower bounds
- source-specific tail ranges
- Brouncker lookahead-derived candidate evidence

Decision made for MVP:

- canonical reciprocal-pi source is **Brouncker 4/pi**

Reason:

- MVP needs `3/pi^2`
- easiest route is `(3/16) * (4/pi)^2`

### 4. Unary reciprocal framework

Completed current milestone:

- reciprocal GCF prefix stream
- reciprocal GCF exact-tail stream
- shared exact-collapse core
- snapshot semantics
- finite semantic alignment tests
- unary taxonomy and classification layer

Status:

- reciprocal exact-collapse unary path is in good shape for MVP

### 5. Unary sqrt framework

Completed current milestones:

- bounded sqrt approximation APIs
- GCF exact-tail sqrt path
- GCF range-seeded sqrt approximation path
- certified progressive sqrt prefix stream
- sqrt snapshot/status model
- unary taxonomy alignment with reciprocal

Status:

- enough sqrt machinery exists for MVP numerator work
- progressive certification exists but is not yet the bottleneck

### 6. Binary taxonomy / contract direction

Completed:

- initial binary taxonomy
- initial binary contract layer
- BLFT / GCF binary stream classification work
- organizational split of larger BLFT stream file

Status:

- binary streaming architecture is far enough along for MVP support work
- binary tightening can continue after MVP if needed

### 7. Angle / trig groundwork

Completed:

- explicit `Angle` type
- explicit degree vs radian semantics
- current MVP denominator fixed on degree input
- MasterPlan note added to generalize later to both radians and degrees

### 8. Sin degree path

Completed:

- exact degree points: 0°, 30°, 90°, 180°
- bounded certified range for 69°
- tightened conservative bound for `sin(69°)`

Current bound:

- `sin(69°) ∈ [6/7, 267/280]`

### 9. Tanh special-case path

Completed:

- conservative certified range for `tanh(sqrt(5))`

Current bound:

- `tanh(sqrt(5)) ∈ [21/22, 1]`

### 10. Denominator certification

Completed:

- certified range for:

      tanh(sqrt(5)) - sin(69°)

Current bound:

- denominator in `[3/3080, 1/7]`

This excludes zero, which was the critical breakthrough needed for the full quotient.

### 11. Numerator helper

Completed:

- bounded approximation helper for:

      sqrt(3/pi^2 + e)

implemented as:

      sqrt((3/16) * (4/pi)^2 + e)

### 12. Full target formula in test code

Completed:

- full target shape frozen in test code
- first certified approximation range for the full target formula

Production code now has reusable helpers for numerator, denominator, and quotient range, but the literal target formula remains intentionally a test-level assembly.

---

## Current Known Good Facts

### Numerator

Production helper exists:

    MVPNumeratorApprox(...)

This returns a bounded rational point approximation for:

    sqrt(3/pi^2 + e)

### Denominator

Production helper exists:

    MVPDenominatorBounds(...)

Current certified denominator range excludes zero:

    [3/3080, 1/7]

### Full target

Production helper exists:

    MVPTargetBounds(...)

This returns a certified quotient range for the full target.

### Tests

The suite is currently green.

---

## Immediate Next Steps Toward MVP

1. Tighten denominator bounds.
2. Tighten numerator approximation.
3. Tighten full target quotient range.
4. Add final MVP acceptance tests.
5. Review whether any minimal legacy purge is needed before declaring MVP.

---

## Remaining Tasks in Priority Order

## Priority 1 — Tighten denominator bounds

This is the highest-leverage remaining work.

### Why

Even though zero is excluded, the denominator interval is still wide. Since the full target quotient divides by this interval, denominator width dominates quotient width.

### Tasks

1. Tighten `sin(69°)` bound further.
2. Tighten `tanh(sqrt(5))` bound further.
3. Update denominator tests to assert improved range.
4. Propagate improvements into target-range tests.

### Candidate approaches

For `sin(69°)`:

- better first-quadrant trigonometric inequalities
- tangent/secant-style enclosure
- narrow polynomial/Taylor enclosure with rigorous remainder
- interval reasoning on radians after exact degree conversion

For `tanh(sqrt(5))`:

- tighter lower bound for `sqrt(5)`
- monotonicity of tanh
- better lower/upper enclosure for tanh on the resulting interval
- use identities involving exponentials only if bounds remain rigorous and simple

### Exit condition

- denominator still certified
- denominator excludes zero by a visibly better margin
- tests lock the tighter interval

---

## Priority 2 — Tighten numerator approximation

### Why

Numerator is currently a bounded rational point approximation, but its sharpness is limited by current prefix choices and sqrt policy.

### Tasks

1. Increase Brouncker prefix depth for `4/pi` where it clearly improves bounds.
2. Increase `e` prefix depth where it clearly improves bounds.
3. Evaluate whether current `DefaultSqrtPolicy2` is sufficient.
4. Add tests proving monotone or at least improved numerator sharpness.

### Exit condition

- numerator point approximation is materially sharper
- tests lock chosen prefix depths and policy

---

## Priority 3 — Tighten the full target quotient range

### Why

Once numerator and denominator improve, the quotient range should narrow substantially.

### Tasks

1. Update `MVPTargetBoundsDefault()` expectations.
2. Add a test asserting the new width is strictly narrower than the current width.
3. Add sanity assertions:
   - positive
   - finite
   - certified inside range
   - denominator zero exclusion preserved

### Exit condition

- full-target certified range is materially useful, not merely nontrivial

---

## Priority 4 — Add final MVP acceptance tests

### Tasks

1. Add one top-level MVP acceptance test for the full target range.
2. Add one regression test proving the denominator excludes zero.
3. Add one regression test proving the target range remains inside and positive.
4. Add one test locking the canonical pi source choice for MVP.
5. Add one test documenting that the literal target formula remains test-only, not public API.

### Exit condition

- one obvious test file shows “this is the MVP and it works”

---

## Priority 5 — Small API cleanup before MVP freeze

### Tasks

1. Review current naming around “Approx”, “Bounds”, and “Terms”.
2. Ensure numerator/denominator/target helper names are consistent.
3. Remove or mark transitional aliases that are now misleading.
4. Review snapshots/status fields for clarity.
5. Purge any low-value legacy regular-CF-first naming that conflicts with GCF-first direction.

### Exit condition

- no obvious naming confusion in MVP-facing helpers

---

## Priority 6 — Purge legacy code that no longer serves the GCF-first design

This is important, but not before the MVP path is complete.

### Tasks

1. Identify regular-CF-specific legacy code no longer used.
2. Remove dead helpers.
3. Collapse redundant adapters if they no longer pay rent.
4. Trim old transitional APIs after test coverage confirms no value.
5. Split or reorganize medium/large files only when it improves clarity.

### Exit condition

- less conceptual clutter
- no regression in test clarity

---

## Priority 7 — Make trig infrastructure less ad hoc

After MVP, generalize beyond the special-case denominator path.

### Tasks

1. Add explicit support for both degree and radian angle semantics in trig APIs.
2. Introduce shared trig bound helpers instead of special-case files.
3. Add more exact/common-angle support.
4. Add general bounded `sin` path for non-special angles.
5. Add general bounded `tanh` path beyond the `sqrt(5)` special case.

### Exit condition

- trig is reusable, not MVP-special only

---

## Priority 8 — Strengthen binary streaming toward general composition

This is likely post-MVP.

### Tasks

1. Map remaining BLFT/DiagBLFT code cleanly onto binary taxonomy.
2. Strengthen exact-tail and bounded-prefix binary layers.
3. Improve GCF-native binary composition contracts.
4. Reduce overlap between one-off exact-collapse helpers and longer-term streaming operators.

### Exit condition

- binary composition is clearly ready for broader expression evaluation

---

## Priority 9 — Add a small internal expression evaluator

This should remain internal unless a public API genuinely emerges.

### Tasks

1. Add internal expression node types.
2. Support constants, sources, unary nodes, binary nodes.
3. Evaluate expression trees into bounded approximations and/or certified ranges.
4. Keep target formula assembly in tests until the evaluator shape is obviously stable.

### Exit condition

- internal expression assembly becomes easier than bespoke helpers

---

## Priority 10 — Documentation and cleanup after MVP

### Tasks

1. Update API user guide notes for GCF-first architecture.
2. Clarify exact-tail vs bounded-prefix vs progressive-certified terminology.
3. Reconcile overlapping docs:
   - `api_spec.md`
   - `UserGuide.md`
   - `MasterPlan.md`
4. Document canonical MVP target and what remains beyond MVP.

### Exit condition

- docs reflect the actual codebase, not the earlier architecture

---

## Explicit Non-Goals Until MVP Is Done

Do not get distracted by:

- broad trig generalization
- generic transcendental frameworks
- public expression DSL/API
- aggressive documentation work
- nonessential code reorganization
- performance tuning
- speculative exact-real expansion beyond what the MVP needs

---

## Critical Path to MVP

The shortest path from here is:

1. tighten `sin(69°)`
2. tighten `tanh(sqrt(5))`
3. tighten denominator range
4. tighten numerator approximation
5. tighten quotient range
6. add final MVP acceptance tests
7. perform small API cleanup
8. declare MVP complete

---

## Recommended Next Bite

1. Improve `sin(69°)` conservative bound.
2. Improve `tanh(sqrt(5))` conservative bound.
3. Update denominator tests with tighter expected interval.
4. Update full-target tests with improved range.

That is still the highest-leverage work item.

---

## Notes for Future Chats

When resuming in a fresh chat:

1. Ask for the current versions of:
   - `mvp_numerator.go`
   - `mvp_denominator.go`
   - `mvp_target_formula.go`
   - `sin_degrees.go`
   - `tanh_special.go`
   - associated tests

2. Reconfirm current certified denominator interval and whether zero is excluded.

3. Continue only on the highest-leverage blocker:
   - denominator tightening first
   - numerator tightening second
   - quotient tightening third

4. Keep the literal full target formula in tests, not public API.

5. Stay test-first, and use stubs when needed.

# End of MasterPlan.md