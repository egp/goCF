// sqrt_stream_status.go v2
package cf

type SqrtStreamStatus string

const (
	SqrtStreamStatusUnstarted            SqrtStreamStatus = "unstarted"
	SqrtStreamStatusExactInput           SqrtStreamStatus = "exact_input"
	SqrtStreamStatusBoundedCollapse      SqrtStreamStatus = "bounded_collapse"
	SqrtStreamStatusCertifiedProgressive SqrtStreamStatus = "certified_progressive"
	SqrtStreamStatusFailed               SqrtStreamStatus = "failed"
)

// sqrt_stream_status.go v2
