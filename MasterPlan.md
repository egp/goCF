# MasterPlan.md

# GoCF Master Plan

## Mission

Build a Go library for mathematically correct continued-fraction and generalized-continued-fraction arithmetic, with emphasis on:

- exactness and certified bounds over floating-point approximation
- streaming / demand-driven computation
- reusable unary and binary operator plumbing
- clean separation between canonical production paths and temporary MVP exceptions
- strong tests, small TDD steps, and explicit architectural seams

---

## Immediate Project Goal: MVP

Deliver a tested MVP for the certified bounded expression:

    sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))
    [NOTE-add expected RCF, approx Rational, approx decimal]}

Requirements for MVP:

- production code computes a mathematically justified bounded result
- regular CF output is acceptable
- code is test-backed and green
- canonical source choices and default budgets are explicit
- temporary exceptions are explicit, documented, and fenced by tests
- architectural direction continues to move toward GCF-first reusable operators

---

## Current Status

### Overall
- The MVP formula path exists and is working in production code.
- Tests are green at the current checkpoint.
- The codebase has progressed from loose special-case work toward a clearer canonical MVP path with bounded parity hooks.

### Canonical source choices
- Canonical 4/pi family for MVP: **Brouncker**
- Alternate parity family available at subexpression layer: **Lambert**
- Canonical e source: `ECFGSource`
- Current exact 69° angle path remains an explicit temporary exception:
  - `MVP69DegreeGCFSource()`
  - `MVP69DegreeTail()`

### Canonical default budgets
- `MVPDefaultFourOverPiPrefixTerms = 6`
- `MVPDefaultEPrefixTerms = 8`
- `MVPNumeratorBridgePrefixTerms = 64`

### Numerator
Current numerator target:

    sqrt(3/pi^2 + e)

Current production layering:

- `MVPNumeratorRadicandApprox(...)`
- `MVPNumeratorRadicandBridgeSource(...)`
- `MVPNumeratorApprox(...)`

Interpretation:
- the numerator radicand is computed as a bounded rational subexpression
- the current bridge source is the explicit temporary seam
- sqrt is applied through the GCF-ingesting unary path after the bridge

This is the biggest remaining architectural compromise on the MVP path.

### Denominator
Current denominator target:

    tanh(sqrt(5)) - sin(69°)

Current status:
- denominator path is considered acceptable for MVP
- `tanh(sqrt(5))` uses a GCF-ingesting metadata-driven special path
- `sin(69°)` uses an explicit exact-tail angle exception
- denominator excludes zero and has certified positive lower bound

### Target assembly
Current target:

    sqrt(3/pi^2 + e) / (tanh(sqrt(5)) - sin(69°))

Status:
- target-level bound assembly exists
- target default path uses canonical Brouncker family and current default budgets
- target is intentionally bounded/non-point at current MVP stage

### Infinite-contract / exception audit
Audit tests exist and currently document:

1. canonical mathematical sources are infinite/algorithmic where expected
2. some current MVP helpers still rely on explicit finite/exact-tail exceptions
3. numerator bridge is finite by design and explicitly budgeted
4. current target still works despite these explicit exceptions

---

## Critical Path to MVP Completion

### 1. Freeze canonical MVP path
Status: largely done

Needed:
- keep Brouncker canonical for 4/pi
- keep Lambert as parity-only below numerator/target API level unless clearly needed
- keep current default budgets explicit and stable

### 2. Decide whether the narrowed numerator seam is acceptable for MVP
Status: active critical-path item

Question:
- is `MVPNumeratorRadicandBridgeSource(...)` an acceptable temporary MVP seam?
- or is one more source-driven replacement step justified before sign-off?

This is now the most important remaining decision.

### 3. Freeze denominator as accepted MVP mechanism
Status: essentially done

Current accepted temporary mechanism:
- special GCF-ingesting tanh path for sqrt(5)
- explicit exact-tail 69° path for sin(69°)

No further denominator redesign should happen before MVP unless a bug is found.

### 4. Verify target-level readiness
Status: active but close

Needed:
- target-level tests continue to show:
  - numerator positive
  - denominator excludes zero
  - target bound is inside
  - canonical path is stable
- determine whether the current numerator seam is stable enough for MVP sign-off

### 5. Keep exception list short and explicit
Status: in progress

Current known MVP exceptions:
- exact-tail angle path for 69°
- finite numerator radicand bridge
- parity hooks not lifted into top-level numerator/target APIs

MVP should not add more exceptions.

---

## Must-Have Before Declaring MVP Complete

1. Canonical source family and defaults remain frozen
2. Numerator seam decision is made
3. Denominator is frozen as acceptable for MVP
4. Target assembly is verified and stable
5. Current exceptions are documented and test-backed
6. No hidden dependence on floating-point or unjustified heuristics

---

## Nice-To-Have But Not Required For MVP

- broader Lambert parity beyond subexpression layer
- lifting alternate 4/pi family hook upward into numerator API
- lifting alternate 4/pi family hook upward into target API
- replacing more temporary bridges before MVP sign-off
- expanding additional unary operator generality

---

## Definitely Post-MVP

### Infinite-GCF-oriented cleanup and verification
Goal after MVP:
- move prod ingestion toward an infinite-GCF-first model
- ensure prod operator paths do not rely on source exhaustion except where mathematically unavoidable and explicitly modeled
- retire current temporary finite/exact-tail MVP exceptions one by one

Verification work to add post-MVP:
- audit grep and review pass for finite/exact-tail assumptions
- contract tests around non-exhaustion / source behavior
- explicit whitelist of acceptable remaining exceptions until removed

### Additional well-known constant sources
Add more algorithmic constant/source families where worthwhile, including candidates such as:
- additional square roots / quadratic radicals
- more π-family constructions
- additional e-related constructions
- later candidates such as selected constants from Brouncker/Lambert/Ramanujan/Chudnovsky-style work where they fit the architecture

### Source-family parity
Post-MVP:
- strengthen Brouncker vs Lambert parity at the subexpression level
- decide whether source-family hooks should be elevated to numerator and target APIs
- add regression/parity tests ensuring both families remain viable where structurally appropriate

### Replace temporary numerator bridge
Post-MVP:
- replace `MVPNumeratorRadicandBridgeSource(...)` with a more genuinely source-driven path if practical
- or prove the current seam is the right permanent architectural boundary

### Replace temporary angle exception
Post-MVP:
- revisit whether sin(69°) should remain an exact-tail special case
- look for a cleaner long-term operator-facing representation if justified

---

## Current Next Recommended Step

Add an **MVP-readiness test layer** to answer the central remaining question:

- Is the current narrowed numerator seam acceptable for MVP sign-off?
- Or is one more source-driven replacement step worth doing before declaring the canonical MVP path done?

Recommended focus files:
- `mvp_numerator.go`
- `mvp_numerator_test.go`
- `mvp_target_formula.go`
- `mvp_target_formula_test.go`
- `mvp_sources.go`
- `mvp_sources_test.go`
- `infinite_gcf_contract_test.go`

The goal of that step is decision support, not redesign.

---

## Working Rules / Preferences

- mathematical correctness first
- TDD preferred
- small bite size, usually 3–4 files
- do not assume current file contents; resync from actual files
- keep methods readable and not overly large (fewer than 30 lines)
- prefer fewer coherent files over many tiny files
- regular CF output remains acceptable
- temporary exceptions must be explicit and tested
- keep focus on the MVP race; defer broad redesign unless it directly helps the critical path

---

## Summary

The project is close to an MVP.

The main remaining issue is no longer raw functionality; it is deciding whether the current narrowed numerator bridge seam is sufficiently acceptable for MVP sign-off, with the rest of the canonical path held stable. Create test cases to aid in decision making.

Everything else should be judged by whether it helps answer that question quickly and safely. Maintain the focus on getting to MVP.

# EOF MasterPlan.md