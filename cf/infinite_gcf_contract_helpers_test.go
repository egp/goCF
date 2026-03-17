// infinite_gcf_contract_helpers_test.go v1
package cf

import "testing"

func mustReadNPQWithoutExhaustion(t *testing.T, src GCFSource, n int) [][2]int64 {
	t.Helper()

	out := make([][2]int64, 0, n)
	for i := 0; i < n; i++ {
		p, q, ok := src.NextPQ()
		if !ok {
			t.Fatalf("expected non-exhausting source for first %d terms, exhausted at i=%d", n, i)
		}
		out = append(out, [2]int64{p, q})
	}
	return out
}

// infinite_gcf_contract_helpers_test.go v1
