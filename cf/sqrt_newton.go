// sqrt_newton.go v3
package cf

// RationalSqrtExact reports whether x has an exact rational square root.
// If so, it returns that root (the nonnegative one).
func RationalSqrtExact(x Rational) (Rational, bool, error) {
	return SqrtCoreRationalExact(x)
}

// NewtonSqrtStep computes one exact Newton iteration for sqrt(x):
//
//	y_(n+1) = (y_n + x / y_n) / 2
//
// Preconditions:
//   - x >= 0
//   - y != 0
//   - y > 0 (we stay on the nonnegative branch)
func NewtonSqrtStep(x, y Rational) (Rational, error) {
	return SqrtCoreNewtonStep(x, y)
}

// NewtonSqrtIterates returns successive exact Newton iterates for sqrt(x),
// starting from seed. The returned slice has length steps.
//
// Behavior:
//   - steps < 0 => error
//   - if x has an exact rational square root, that exact root is returned
//     repeated steps times
//   - otherwise, the sequence is produced by repeated NewtonSqrtStep calls
func NewtonSqrtIterates(x, seed Rational, steps int) ([]Rational, error) {
	return SqrtCoreNewtonIterates(x, seed, steps)
}

// SqrtApproxRational returns a single exact rational approximation to sqrt(x)
// after the requested number of Newton iterations from the given seed.
//
// Behavior:
//   - steps < 0 => error
//   - steps == 0 => returns seed (after domain/seed validation, unless x is exact square)
//   - if x has an exact rational square root, returns that exact root immediately
func SqrtApproxRational(x, seed Rational, steps int) (Rational, error) {
	return SqrtCoreApproxRational(x, seed, steps)
}

// SqrtResidual returns the exact residual y^2 - x.
func SqrtResidual(x, y Rational) (Rational, error) {
	return SqrtCoreResidual(x, y)
}

// SqrtResidualAbs returns the exact absolute residual |y^2 - x|.
func SqrtResidualAbs(x, y Rational) (Rational, error) {
	return SqrtCoreResidualAbs(x, y)
}

// SqrtApproxRationalUntilExact performs Newton iterations for sqrt(x),
// stopping early if an iterate becomes exact (residual == 0), or after
// maxSteps iterations.
//
// Returns:
//   - the last iterate examined/produced
//   - exact=true if that iterate is an exact square root of x
//   - error on invalid input
//
// Behavior:
//   - maxSteps < 0 => error
//   - if x has an exact rational square root, returns it immediately with exact=true
//   - otherwise, if maxSteps == 0, returns seed with exact=false after validation

func SqrtApproxRationalUntilExact(x, seed Rational, maxSteps int) (Rational, bool, error) {
	return SqrtCoreApproxRationalUntilExact(x, seed, maxSteps)
}

// SqrtApproxRationalUntilExactDefault performs bounded Newton iteration for
// sqrt(x), using DefaultSqrtSeed(x), and stops early if exact convergence occurs.
func SqrtApproxRationalUntilExactDefault(x Rational, maxSteps int) (Rational, bool, error) {
	seed, err := DefaultSqrtSeed(x)
	if err != nil {
		return Rational{}, false, err
	}
	return SqrtApproxRationalUntilExact(x, seed, maxSteps)
}

// SqrtApproxRationalUntilResidual performs Newton iterations for sqrt(x),
// stopping early once the exact absolute residual satisfies:
//
//	|y^2 - x| <= tol
//
// Returns:
//   - the last iterate examined/produced
//   - satisfied=true if the tolerance condition was met
//   - error on invalid input
//
// Behavior:
//   - maxSteps < 0 => error
//   - tol < 0      => error
//   - if x has an exact rational square root, returns it immediately with satisfied=true
//   - if maxSteps == 0, returns seed with satisfied determined from the residual
func SqrtApproxRationalUntilResidual(x, seed Rational, maxSteps int, tol Rational) (Rational, bool, error) {
	return SqrtCoreApproxRationalUntilResidual(x, seed, maxSteps, tol)
}

// SqrtApproxRationalUntilResidualDefault performs bounded Newton iteration for
// sqrt(x), using DefaultSqrtSeed(x), and stops early once |y^2 - x| <= tol.
func SqrtApproxRationalUntilResidualDefault(x Rational, maxSteps int, tol Rational) (Rational, bool, error) {
	seed, err := DefaultSqrtSeed(x)
	if err != nil {
		return Rational{}, false, err
	}
	return SqrtApproxRationalUntilResidual(x, seed, maxSteps, tol)
}

// sqrt_newton.go v3
