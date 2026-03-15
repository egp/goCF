// binary_stream_taxonomy.go v1
package cf

type binaryOperatorKind string

const (
	binaryOperatorAdd binaryOperatorKind = "add"
	binaryOperatorSub binaryOperatorKind = "sub"
	binaryOperatorMul binaryOperatorKind = "mul"
	binaryOperatorDiv binaryOperatorKind = "div"
)

type binaryInputKind string

const (
	binaryInputUnknown   binaryInputKind = "unknown"
	binaryInputCF        binaryInputKind = "cf"
	binaryInputGCF       binaryInputKind = "gcf"
	binaryInputMixed     binaryInputKind = "mixed"
	binaryInputExactTail binaryInputKind = "exact_tail"
	binaryInputPrefix    binaryInputKind = "prefix"
)

type binaryProgressKind string

const (
	binaryProgressUnknown              binaryProgressKind = "unknown"
	binaryProgressExactCollapse        binaryProgressKind = "exact_collapse"
	binaryProgressProgressiveCertified binaryProgressKind = "progressive_certified"
)

type binaryStreamClass struct {
	Operator binaryOperatorKind
	Input    binaryInputKind
	Progress binaryProgressKind
}
