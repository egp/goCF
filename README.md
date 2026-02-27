# goCF

Streaming arithmetic for continued fractions (Gosper / HAKMEM 101B), ported from a Scala prototype into Go with a correctness-first mindset.

The goal is to perform arithmetic directly on continued fractions *as streams* (ULFT/BLFT transforms), emitting digits only when they are **provably** correct using conservative range reasoning (including Gosper’s “toroidal” wrap-around semantics).

## Status
Early-stage but already includes:
- `Rational` with checked `int64` arithmetic (`ErrOverflow` on overflow)
- `Range` with:
  - inside vs outside intervals (toroidal model)
  - open/closed endpoint flags
  - `Contains` and `ContainsZero` semantics used by denom guards
- `ULFT` + `ULFTStream` (unary transforms + safe digit emission)
- `BLFT` + `BLFTStream` (binary transforms + denom guard + bounded rational finalization)
- `Bounder` (range shrinkage while ingesting CF terms)
- Property-based testing via `pgregory.net/rapid`

## Fast local loop
Run the full local pipeline:

```bash
./tools/check.sh