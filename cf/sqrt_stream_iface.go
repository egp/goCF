// sqrt_stream_iface.go v1
package cf

type SqrtApproxStream interface {
	ContinuedFraction
	Err() error
	Snapshot() SqrtApproxStreamSnapshot
}

// sqrt_stream_iface.go v1
