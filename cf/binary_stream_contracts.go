// binary_stream_contracts.go v1
package cf

type binaryClassifiedStream interface {
	binaryClass() binaryStreamClass
}

func classifyBinaryStream(stream any) (binaryStreamClass, bool) {
	s, ok := stream.(binaryClassifiedStream)
	if !ok {
		return binaryStreamClass{}, false
	}
	return s.binaryClass(), true
}
