// sqrt_stream_status.go v2
package cf

type SqrtStreamStatus string

const (
	SqrtStreamStatusUnstarted       SqrtStreamStatus = "unstarted"
	SqrtStreamStatusExactInput      SqrtStreamStatus = "exact_input"
	SqrtStreamStatusBoundedCollapse SqrtStreamStatus = "bounded_collapse"
	SqrtStreamStatusFailed          SqrtStreamStatus = "failed"
)

// sqrt_stream_status.go v2
