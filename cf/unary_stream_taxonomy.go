// unary_stream_taxonomy.go v1
package cf

type unaryOperatorKind string

const (
	unaryOperatorReciprocal unaryOperatorKind = "reciprocal"
	unaryOperatorSqrt       unaryOperatorKind = "sqrt"
)

type unaryInputKind string

const (
	unaryInputUnknown   unaryInputKind = "unknown"
	unaryInputGCFPrefix unaryInputKind = "gcf_prefix"
	unaryInputGCFExact  unaryInputKind = "gcf_exact_tail"
	unaryInputCFPrefix  unaryInputKind = "cf_prefix"
)

type unaryProgressKind string

const (
	unaryProgressUnknown              unaryProgressKind = "unknown"
	unaryProgressExactCollapse        unaryProgressKind = "exact_collapse"
	unaryProgressProgressiveCertified unaryProgressKind = "progressive_certified"
)

type unaryStreamClass struct {
	Operator unaryOperatorKind
	Input    unaryInputKind
	Progress unaryProgressKind
}

func classifyReciprocalSnapshot(s ReciprocalApproxStreamSnapshot) unaryStreamClass {
	if s.MaxIngestTerms > 0 {
		return unaryStreamClass{
			Operator: unaryOperatorReciprocal,
			Input:    unaryInputGCFExact,
			Progress: unaryProgressExactCollapse,
		}
	}
	if s.PrefixTerms > 0 {
		return unaryStreamClass{
			Operator: unaryOperatorReciprocal,
			Input:    unaryInputGCFPrefix,
			Progress: unaryProgressExactCollapse,
		}
	}
	return unaryStreamClass{
		Operator: unaryOperatorReciprocal,
		Input:    unaryInputUnknown,
		Progress: unaryProgressExactCollapse,
	}
}

func classifySqrtSnapshot(s SqrtApproxStreamSnapshot) unaryStreamClass {
	in := unaryInputUnknown
	switch {
	case s.GCFInputApprox != nil || s.PrefixTerms > 0:
		in = unaryInputGCFPrefix
	case s.CFInputApprox != nil:
		in = unaryInputCFPrefix
	}

	progress := unaryProgressUnknown
	switch s.Status {
	case SqrtStreamStatusExactInput, SqrtStreamStatusBoundedCollapse:
		progress = unaryProgressExactCollapse
	case SqrtStreamStatusCertifiedProgressive:
		progress = unaryProgressProgressiveCertified
	}

	return unaryStreamClass{
		Operator: unaryOperatorSqrt,
		Input:    in,
		Progress: progress,
	}
}

// unary_stream_taxonomy.go v1
