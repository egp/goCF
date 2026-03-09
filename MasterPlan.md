# MasterPlan.md v2

# CF2 / goCF Master Plan Snapshot

## Current status

The project is in a transitional state after migrating core exact arithmetic toward BigRat / BigInt.

What is currently true:

- Exact rational representation has moved away from exposed numerator/denominator fields.
- ULFT coefficients are now BigInt-based.
- Bounder has been migrated to BigInt recurrence state.
- Most tests pass.
- One important BLFT test was isolated and diagnosed:
  - multiplying two independent streams `sqrt(2) * sqrt(2)` does **not** become digit-safe under the current rectangular BLFT range model, because the engine knows only `x in xr` and `y in yr`, not `x == y`.
- Therefore the remaining failing sqrt-square BLFT test is not just a coding bug; it exposes a model limitation.

---

## Immediate conclusion

The current BLFT engine handles independent X and Y correctly, but it lacks a way to represent correlated inputs such as:

- `T(x, x)`
- `x * x`
- `f(x, x)` in general

That suggests the next major feature should be a unary specialization for diagonal binary transforms.

---

## Priority order

## P0. Stabilize the current codebase

### P0.1 Keep main branch green except for the known diagonal-correlation test
- Preserve the current passing state.
- Skip or demote the known failing `sqrt(n)^2` BLFT test until diagonal support exists.
- Keep refine guards enabled so failures are deterministic, not hangs.

### P0.2 Finish API migration cleanup
- Eliminate all remaining references to old Rational field access patterns.
- Eliminate remaining mixed int64 / BigInt / BigRat assumptions.
- Ensure all tests use constructor/helper-based Rational creation only.
- Ensure all ULFT/BLFT tests compare values, not pointer identity.

### P0.3 Audit all production code for hidden legacy assumptions
- Any code assuming:
  - `Rational.P / Rational.Q`
  - `ULFT int64 coefficients`
  - struct equality for BigInt-bearing structs
- Replace with exact-value comparisons and helper constructors.

---

## P1. Add diagonal binary-transform specialization

### P1.1 Introduce a unary diagonal specialization for `T(x, x)`
Purpose:

- support `x * x`
- support other diagonal forms such as `x + x`, `x / x`, etc.
- allow the engine to exploit the knowledge that both inputs are the same stream/value

Possible names:

- `DLFT`
- `DiagonalBLFT`
- `UnaryBLFT`
- `BLFTDiag`

Recommended semantics:

- represent a BLFT constrained to the diagonal:
  - `Tdiag(x) = T(x, x)`
- reduce:
  - `Axy + Bx + Cy + D`
  - on diagonal becomes
  - `A x^2 + (B + C) x + D`
- denominator similarly becomes
  - `E x^2 + (F + G) x + H`

### P1.2 Add exact evaluation for diagonal transforms
Required operations:

- apply exact to Rational
- apply range to InsideRange
- emit digit safely
- refinement loop for one source only

### P1.3 Add tests for diagonal specialization
Initial black-box tests:

- `sqrt(2)^2` begins with 2
- `sqrt(3)^2` begins with 3
- `sqrt(5)^2` begins with 5
- `sqrt(6)^2` begins with 6
- `sqrt(7)^2` begins with 7

Later stronger tests:

- exact rational diagonal cases
- x/x = 1 when defined
- x+x = 2x
- bounded refine-limit behavior

---

## P2. Revisit the one known failing BLFT test

Once diagonal support exists:

### P2.1 Move sqrt-square tests from BLFT rectangle engine to diagonal engine
Rationale:

- current rectangle BLFT cannot prove `x == y`
- diagonal engine is the right abstraction for that case

### P2.2 Keep one regression test documenting the rectangle limitation
Document current expected behavior:

- independent `BLFTStream` on identical source constructors does not imply correlation
- bounded refinement should fail fast rather than hang

That test is valuable documentation of the model boundary.

---

## P3. Strengthen enclosures and diagnostics

### P3.1 Improve tracing / white-box diagnostics
Add optional tracing hooks for:

- current x range
- current y range
- current image range
- floor bounds
- refine choice
- refine counters
- fingerprint history

### P3.2 Improve error messages
Especially for:

- exceeded refine limits
- denominator may cross zero
- cycle detected
- exact integer termination conditions

### P3.3 Keep fingerprinting and ring-buffer diagnostics
Already started and useful.
Continue integrating them consistently.

---

## P4. Formalize the new specification in CF2

This is now underway.

### P4.1 CF2 spec decisions already made
- generalized CF ingestion
- regular CF emission only
- arbitrary precision required
- separate InsideRange and OutsideRange
- strong pre/postconditions
- DfT hooks from the start

### P4.2 Use the README spec as the new project contract
The spec should drive:

- module boundaries
- tests
- invariants
- public API
- migration decisions

---

## P5. Square root operator track

### Current status

Implemented and green:

- Exact rational sqrt detection:
  - `RationalSqrtExact`
- Exact Newton machinery:
  - `NewtonSqrtStep`
  - `NewtonSqrtIterates`
  - bounded rational approximation helpers
  - stop-on-exact helpers
  - stop-on-residual helpers
- Default seed and policy support:
  - `DefaultSqrtSeed`
  - `SqrtPolicy`
  - default / explicit policy APIs
  - explicit seed compatibility wrappers
- CF-facing bounded sqrt APIs:
  - rational -> CF wrappers
  - term extraction helpers
- CF-prefix bridge:
  - `CFApprox`
  - convergent-based approximation from CF prefix
  - range-seeded approximation from CF prefix
  - midpoint-targeted approximation from CF prefix
- Range-side sqrt support:
  - `SqrtRangeExact`
  - `SqrtRangeHeuristic`

### Current architecture

The sqrt subsystem is now consolidated into four coherent production files:

- `sqrt_newton.go`
- `sqrt_api.go`
- `sqrt_cf.go`
- `sqrt_range.go`

This is a good medium-sized layout and should remain stable unless the design changes substantially.

### Coverage checkpoint

Current package coverage is about 70%.

Largest meaningful remaining low-coverage targets:

- `BLFTStream.Next`
- `ULFTStream.Next`

These are now more important than adding new sqrt production code immediately.

### Next baby steps

1. Add targeted tests for `ULFTStream` branch/error paths.
2. Add one tiny smoke test covering:
   - `DiagBLFT.String`
   - `ULFT.String`
   - `Range.String`
3. Re-run coverage and reassess whether `BLFTStream` or `ULFTStream` still dominates the remaining risk.
4. Then return to production sqrt work.

### Next production step after the coverage detour

Compare the CF-prefix sqrt strategies on the same prefix:

- convergent-targeted
- range-seeded
- midpoint-targeted

using exact residual comparisons against the same rational target.

That comparison should tell us which bridge toward a future true streaming `sqrt()` is the most promising.


## P6. Generalized CF ingestion deeper work

### P6.1 Formal ingestion algebra
Still to be fully implemented from the spec:

- ingest `[p,q]` via:
  - `x = p + q / x'`
- coefficient rewrite for:
  - ULFT ingest
  - BLFT ingest X
  - BLFT ingest Y
  - diagonal ingest

### P6.2 Bounder for GCF terms
Current bounder is still shaped around regular CF ingestion.
Eventually it should be generalized so the new project can ingest GCF natively and exactly.

---

## P7. Replace fragile implementation habits

Lessons learned from this project:

### P7.1 Do not expose internal numeric fields
Always use constructors and methods.

### P7.2 Do not rely on struct equality for BigInt/BigRat-bearing types
Use value comparison helpers.

### P7.3 Keep correlation explicit in the type system
Do not expect a rectangular two-variable enclosure to prove diagonal identities.

### P7.4 Spec first
The lack of a sufficiently strong spec was the biggest mistake.
The new CF2 project should be driven by the spec from the first commit.

---

## P8. Stretch goal

Long-term stretch goal remains:

Compute exact-real style expressions such as:

- `sqrt(3/pi^2 + e) / tanh(sqrt(5) - sin(69 degrees))`
- `sqrt(3/π² + e) / tanh(sqrt(5) - sin(69))`
- $$\frac{\sqrt{3/\pi^2 + e}}{\tanh(\sqrt{5} - \sin(69^\circ))}$$

Which might eval to approx 1.77031957889 or as a CF [2; 62, 1, 3, 1, 1, 5, 1, 1, 2, 1, 2, 11, 3, 3, 1, 2, 1, 12, 1, 6, 5, 2, 3, 13, 4, 1, 1, 3, 4, 129, 2, 1, 3, 1, 3, 1, 5, 1, 16, 1, 1, 6, 4, 9, 3, 1, 16, 1, 4, 1, 1, 1, 1, 4, 1, 2, 2, 1, 1, 1, 8, 3, 32, 1, 2, 3, 6, 1, 1, 1, 1, 2, 3, 1, 1, 5, 1, 4, 5, 2, 2, 7, 12, 1, 3, 1, 11, 1, 4, 6, 2, 15, 2, 12, 1, 1, 23, 2, 5, 1, 4, 167, 8, 2, 3, ...]

This depends on:

- stable unary and binary operators
- diagonal support
- sqrt operator
- transcendental sources
- disciplined exact enclosure logic

This remains a stretch goal, not an immediate implementation target.

---

## Recommended next step

1. Commit the current mostly-green state.
2. Introduce a diagonal transform abstraction for `T(x,x)`.
3. Move `sqrt(n)^2` tests to the diagonal engine.
4. After diagonal support is stable, begin square-root operator work using the CF2 spec as the design authority.

# End of MasterPlan.md V2