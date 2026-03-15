// sqrt_stream_status.go v1
package cf

type SqrtStreamStatus string

const (
	SqrtStreamStatusUnstarted       SqrtStreamStatus = "unstarted"
	SqrtStreamStatusBoundedCollapse SqrtStreamStatus = "bounded_collapse"
	SqrtStreamStatusFailed          SqrtStreamStatus = "failed"
)

// sqrt_stream_status.go v1
