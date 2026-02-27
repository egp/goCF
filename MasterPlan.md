# goCF Master Plan

## Snapshot (current state)
- Language: Go (int64 + checked arithmetic; BigInt later)
- Domain: Gosper/HакMEM 101B continued-fraction arithmetic
- Core objects in place:
  - `Rational` (checked int64 arithmetic)
  - `Range` (inside/outside + open/closed endpoints; toroidal semantics)
  - `ULFT` / `ULFTStream` (unary transforms + safe digit emission)
  - `BLFT` / `BLFTStream` (binary transforms + denom guard)
  - `Bounder` (range shrinkage while ingesting CF terms)
  - PBT via `pgregory.net/rapid`

## Working discipline
- Tight loop:
  - `./tools/check.sh`
- When core APIs change, resync contract:
  - `mkdir -p tmp && go doc -all ./cf > tmp/api_cf.txt`

## Priorities (highest first)

### P0 — Correctness first (Gosper smiles)
1. Range correctness (inside/outside + open/closed) across:
   - `Contains`, `ContainsZero`, enclosure construction (must be conservative)
2. Denominator/pole guards (ULFT + BLFT)
   - `DenomRange` must enclose all corners
   - `DenomMayHitZero` must be conservative
3. Streaming safety rules
   - `SafeDigit`/emission logic must only emit when provably safe
4. Checked int64 everywhere on critical paths
   - overflow detection is correctness, not optimization

### P1 — Debuggability (reduce time-to-truth)
1. Fingerprinting (human-readable)
   - ULFT coeffs + BLFT coeffs + Range (incl flags) + input ranges
2. Cycle detection (ring buffer of fingerprints)
3. Structured trace hooks (enable/disable)

### P2 — Property-based tests (torus edge-cases)
1. Rapid suite targeting:
   - pole boundaries (denom near/at 0)
   - wrap-around outside ranges
   - “corner extremum” assumptions for bilinear denom and BLFT mapping
2. Golden tests for known transforms on rationals and periodic CFs

### P3 — Open/closed endpoint propagation
1. Ensure enclosures produced by range-mapping are CLOSED unless proven open
2. Decide semantics for how endpoint flags propagate through ULFT/BLFT maps

### P4 — Constants and higher-level API
1. Principled generators:
   - `sqrt(2)`, `phi`, `e` as infinite CF generators
2. `pi`:
   - start as finite prefix (OEIS A001203)
   - later replace with a spigot/source
3. Expression-level API (future):
   - compose ULFT/BLFT streams into expression DAG

## Stretch goal (end state)
Compute the following expression via streaming CF arithmetic (degrees):
$$\frac{\sqrt{3/\pi^2 + e}}{\tanh(\sqrt{5} - \sin(69^\circ))}$$

Which might eval to [2; 62, 1, 3, 1, 1, 5, 1, 1, 2, 1, 2, 11, 3, 3, 1, 2, 1, 12, 1, 6, 5, 2, 3, 13, 4, 1, 1, 3, 4, 129, 2, 1, 3, 1, 3, 1, 5, 1, 16, 1, 1, 6, 4, 9, 3, 1, 16, 1, 4, 1, 1, 1, 1, 4, 1, 2, 2, 1, 1, 1, 8, 3, 32, 1, 2, 3, 6, 1, 1, 1, 1, 2, 3, 1, 1, 5, 1, 4, 5, 2, 2, 7, 12, 1, 3, 1, 11, 1, 4, 6, 2, 15, 2, 12, 1, 1, 23, 2, 5, 1, 4, 167, 8, 2, 3, ...]