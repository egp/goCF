# MasterPlan.md v1

## Project goal
Implement mathematically correct, streaming continued-fraction arithmetic in Go, aligned with Gosper / HAKMEM 101B, with test-driven development, conservative range reasoning (toroidal/“wrap-around” semantics), and strong diagnostic tooling.

Stretch goal expression (degrees):
    sqrt(3/π² + e) / tanh(sqrt(5) - sin(69))

---

## Current snapshot (what exists today)

### Core numeric types
- Rational (int64, checked arithmetic + ErrOverflow)
- Range (Lo/Hi + IncLo/IncHi; inside vs outside semantics; Contains/ContainsZero aligned)
- Convergents / RationalFromTerms

### Continued-fraction sources
- RationalCF
- SliceCF
- PeriodicCF (Sqrt2CF, PhiCF)

### Transforms
- ULFT: (Ax+B)/(Cx+D)
- BLFT: (Axy+Bx+Cy+D)/(Exy+Fx+Gy+H)

### Streaming engines
- Bounder: ingests CF digits and maintains shrinking Range enclosure
- ULFTStream: range-safe digit emission via SafeDigit + refinement
  - strict progress guards
  - cycle detection (now via ring buffer of human-readable fingerprints)
  - errors annotated with FingerprintULFT
- BLFTStream: range-safe emission + refinement of X/Y
  - denom guard via DenomRange/DenomMayHitZero + ContainsZero
  - optional early finalize-to-rational tail
  - errors annotated with FingerprintBLFT

### Diagnostics / safety
- checked.go: add/sub/mul overflow detection
- fingerprints for ULFT/BLFT (canonicalized: gcd+sign)
- ring buffer window with dump for cycle errors
- Rapid PBT suites targeting toroidal edge cases & denom guards
- tools/check.sh pipeline (gofmt/goimports if present, go vet, staticcheck, go test)

---

## Priority roadmap (highest first)

### P0 — Correctness foundations and invariants
1. **Range endpoint semantics are first-class everywhere**
   - Ensure all Range-producing code chooses a conservative openness/closedness policy.
   - Default safe policy: computed enclosures are closed unless proven open/closed correctness.
2. **Denominator guard correctness**
   - BLFT: DenomRange exact-corner enclosure; DenomMayHitZero uses ContainsZero on the enclosure.
   - ULFT: analogous guard wherever denom safety matters (SafeDigit/ApplyULFT).
3. **Fingerprinting + cycle trace everywhere we stall**
   - ULFT: ring-buffer cycle detection (DONE)
   - BLFT: add similar cycle detection trace (TODO)

### P1 — Streaming “golden” behavior
4. **ULFT goldenization**
   - Keep ULFTStream stable with strict progress guards + good error messages.
   - Add targeted golden tests for key transforms and key constants.
5. **BLFT streaming robustness**
   - Strengthen BLFT streaming PBT to hunt denom edge cases and refinement stalls.
   - Add BLFT cycle trace with ring buffer fingerprints.

### P2 — Property-based testing expansion (toroidal focus)
6. Expand Rapid suites
   - Cover more BLFT transforms with known pole structures.
   - Include outside-ranges and wrap-around points more aggressively.
   - Add “no false accept”: if DenomMayHitZero says safe, sampled interior points should not hit denom=0.

### P3 — Fingerprint-driven debugging improvements
7. Ring-buffer diagnostics improvements
   - Include last N fingerprints in all “progress guard tripped” errors.
   - Optional: include a short “phase” tag (emit/refine/finalize) in fingerprints.
8. Human readable vs hash
   - Keep human-readable as canonical.
   - Optionally add a compact hash as a suffix for grep-ability.

### P4 — Constants / principled generators
9. Implement principled infinite CF generators:
   - sqrt2, phi (DONE)
   - e (TODO)
   - pi: keep finite prefix now; later a spigot algorithm (TODO)
10. Add a “constants” module with named constructors and tests.

### P5 — Unary operations beyond ULFT (higher-level algorithms)
11. **Reciprocal**
   - Defer for now; note equivalence of “prepend/strip 0” vs ULFT(0,1,1,0)
12. **Square root via Newton iteration with CF feedback**
   - Implement sqrt(x) using successive approximation:
     - x is a CF stream
     - produce a CF stream approximation y_n
     - feed y_n back into iteration for y_{n+1}
   - Must be bounded (progress/termination guards) and heavily tested.

### P6 — Generalized Continued Fractions (GCF)
13. **Ingest generalized CFs**
   - Add a GCF source interface/type supporting partial numerators/denominators.
   - Decide representation (e.g., (a0; b1/a1, b2/a2, ...) or equivalent).
   - Provide conversion paths:
     - GCF -> regular CF (when possible) for reuse of existing engines
     - Or extend bounder/range logic to directly ingest GCF terms
14. **Consider emitting generalized CFs**
   - Evaluate impact on:
     - API (do we expose both emit modes?)
     - range semantics
     - Gosper arithmetic clarity (he uses generalized CF heavily)
   - Decision: either “internal-only generalized, emit regular” vs “support both”.

### P7 — Docs and reproducibility
15. Create/maintain project docs
   - README.md describing goals and design principles
   - MasterPlan.md (this file) as living roadmap
16. API snapshot workflow
   - When needed: `go doc -all ./cf > tmp/api_cf.txt` (manual snapshot)

---

## Definition of “Gosper smiles”
- ULFT and BLFT streaming agree with exact rational evaluation on large random test domains.
- Denominator/pole guards are conservative (no missed poles) but not overly rejecting.
- Cycles/stalls produce actionable diagnostics (fingerprints + recent history).
- Principled generators (sqrt2/phi/e, later pi spigot) work with unary/binary transforms.
- Path exists toward generalized CF ingestion (and a clear decision on emission).

sqrt(3/π² + e) / tanh(sqrt(5) - sin(69))

$$\frac{\sqrt{3/\pi^2 + e}}{\tanh(\sqrt{5} - \sin(69^\circ))}$$

Which might eval to approx 1.77031957889 or as a CF [2; 62, 1, 3, 1, 1, 5, 1, 1, 2, 1, 2, 11, 3, 3, 1, 2, 1, 12, 1, 6, 5, 2, 3, 13, 4, 1, 1, 3, 4, 129, 2, 1, 3, 1, 3, 1, 5, 1, 16, 1, 1, 6, 4, 9, 3, 1, 16, 1, 4, 1, 1, 1, 1, 4, 1, 2, 2, 1, 1, 1, 8, 3, 32, 1, 2, 3, 6, 1, 1, 1, 1, 2, 3, 1, 1, 5, 1, 4, 5, 2, 2, 7, 12, 1, 3, 1, 11, 1, 4, 6, 2, 15, 2, 12, 1, 1, 23, 2, 5, 1, 4, 167, 8, 2, 3, ...]