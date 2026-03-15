// binary_stream_taxonomy_test.go v1
package cf

import "testing"

func TestBinaryStreamTaxonomy_Literals(t *testing.T) {
	tests := []struct {
		name string
		got  binaryStreamClass
		want binaryStreamClass
	}{
		{
			name: "add exact collapse prefix",
			got: binaryStreamClass{
				Operator: binaryOperatorAdd,
				Input:    binaryInputPrefix,
				Progress: binaryProgressExactCollapse,
			},
			want: binaryStreamClass{
				Operator: binaryOperatorAdd,
				Input:    binaryInputPrefix,
				Progress: binaryProgressExactCollapse,
			},
		},
		{
			name: "div exact tail",
			got: binaryStreamClass{
				Operator: binaryOperatorDiv,
				Input:    binaryInputExactTail,
				Progress: binaryProgressExactCollapse,
			},
			want: binaryStreamClass{
				Operator: binaryOperatorDiv,
				Input:    binaryInputExactTail,
				Progress: binaryProgressExactCollapse,
			},
		},
		{
			name: "mul progressive certified",
			got: binaryStreamClass{
				Operator: binaryOperatorMul,
				Input:    binaryInputGCF,
				Progress: binaryProgressProgressiveCertified,
			},
			want: binaryStreamClass{
				Operator: binaryOperatorMul,
				Input:    binaryInputGCF,
				Progress: binaryProgressProgressiveCertified,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got != tc.want {
				t.Fatalf("got %+v want %+v", tc.got, tc.want)
			}
		})
	}
}
