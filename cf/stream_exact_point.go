// stream_exact_point.go v1
package cf

import "fmt"

// exactPointTermination classifies an exact-point singularity in a stream engine.
//
// Rules:
//   - if the stream has already emitted at least one digit, an exact-point
//     singularity means the remainder is exhausted => clean termination
//   - otherwise the original transform is undefined at the exact point => error
//
// Returns:
//   - done=true, err=nil   => clean exhaustion
//   - done=false, err!=nil => real error
func exactPointTermination(
	engineName string,
	emittedAny bool,
	msg string,
) (done bool, err error) {
	if emittedAny {
		return true, nil
	}
	return false, fmt.Errorf("%s %s", engineName, msg)
}

// stream_exact_point.go v1
