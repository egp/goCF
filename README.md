[README.md.]: #

# goCF

goCF is an experimental Go package for exact and semi-exact arithmetic with
continued fractions, inspired by Gosper-style streaming arithmetic.

## What the package already does

- represents exact rationals with arbitrary precision
- streams finite rational values as continued fractions
- provides exact ULFT and BLFT transform machinery
- supports continued-fraction sources for:
  - finite slices
  - exact rationals
  - periodic sources such as φ and several sqrt(n) examples
- tracks shrinking enclosures with Range and Bounder
- supports diagonal specializations for transforms of T(x,x)
- includes algebraic shortcuts for sqrt(n) sources under narrow diagonal cases
- provides bounded Newton-based sqrt approximation for exact rationals
- exposes sqrt approximations as:
  - exact rationals
  - continued-fraction sources
  - finite CF term slices
- supports default policy, explicit policy, and explicit seed control for
  bounded sqrt approximation APIs

## Current design philosophy

- correctness before cleverness
- exact arithmetic before performance tuning
- conservative enclosures when proof is required
- small testable steps
- design for white-box and black-box testing
- avoid hidden assumptions such as x==y unless the type system says so

## What the package does not yet do

- true streaming sqrt for arbitrary irrational CF inputs
- full algebraic reasoning for general quadratic irrationals
- generalized continued-fraction ingestion and emission
- a complete exact-real operator stack for transcendental expressions
- full policy unification and simplification across all approximation APIs

## What the package is moving toward

- better use of CF prefix ranges, not just convergents
- safer bounded sqrt from CF-source input
- eventually a true streaming sqrt operator
- richer diagonal and algebraic source reasoning
- generalized CF support as a first-class design goal
- broader exact-real style composition of unary and binary operators

## Long-term direction

The long-term goal is a Gosper-friendly arithmetic toolkit that can combine
continued-fraction sources, exact rational transformations, conservative proof
machinery, and practical approximation strategies without losing mathematical
clarity.

The project is being grown incrementally: each small API layer is intended to
be correct, testable, and useful on its own before the next layer is added.


[End of README.md.]: #