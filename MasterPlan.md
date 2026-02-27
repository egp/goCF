# goCF Master Plan

## Ground rules
- Correctness > coverage > performance.
- Tight loop: `./tools/check.sh`
- When core APIs change, resync the contract:
  - `mkdir -p tmp && go doc -all ./cf > tmp/api_cf.txt`

## Highest priority (in order)

### P0 — Correctness (Gosper smiles)
1. **Range semantics** (inside/outside + open/closed)
   - `Contains` is the single truth for membership
   - `ContainsZero()` delegates to `Contains(0)`
2. **Denominator / pole guards** (ULFT + BLFT)
   - conservative `DenomRange`, conservative `DenomMayHitZero`
   - guards must respect Range open/closed flags correctly
3. **Conservative enclosures**
   - any computed enclosure must be **closed** unless we have a proof otherwise
4. **Checked int64 arithmetic** on all critical paths
   - overflow detection is correctness, not optimization

### P1 — Debuggability
1. **Fingerprints** (human-readable first)
   - ULFT/BLFT coeffs + input ranges (incl flags) + output range
2. **Ring-buffer cycle detection**
   - configurable size; detect repeats; emit diagnostic with recent fingerprints
3. Trace hooks (enable/disable)

### P2 — Property-based testing (Rapid)
1. Tight suite targeting torus edge cases:
   - outside ranges, endpoint openness, pole boundaries, denom near zero
2. BLFT streaming properties:
   - accepted transforms must enclose sampled points
   - stream prefix matches exact ApplyRat prefix where defined

### P3 — Constants and higher-level API
1. Principled infinite generators:
   - `sqrt(2)`, `phi`, `e`
2. `pi`:
   - finite prefix now (OEIS A001203)
   - later replace with a spigot
3. Expression-level composition API

## Stretch goal (end state)
Compute the following expression via streaming CF arithmetic (degrees):

$$\frac{\sqrt{3/\pi^2 + e}}{\tanh(\sqrt{5} - \sin(69^\circ))}$$

Which might eval to [2; 62, 1, 3, 1, 1, 5, 1, 1, 2, 1, 2, 11, 3, 3, 1, 2, 1, 12, 1, 6, 5, 2, 3, 13, 4, 1, 1, 3, 4, 129, 2, 1, 3, 1, 3, 1, 5, 1, 16, 1, 1, 6, 4, 9, 3, 1, 16, 1, 4, 1, 1, 1, 1, 4, 1, 2, 2, 1, 1, 1, 8, 3, 32, 1, 2, 3, 6, 1, 1, 1, 1, 2, 3, 1, 1, 5, 1, 4, 5, 2, 2, 7, 12, 1, 3, 1, 11, 1, 4, 6, 2, 15, 2, 12, 1, 1, 23, 2, 5, 1, 4, 167, 8, 2, 3, ...]