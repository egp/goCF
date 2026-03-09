// periodic_gcf_test.go v1
package cf

import "testing"

func TestPeriodicGCF_PrefixThenPeriod(t *testing.T) {
	g := NewPeriodicGCF(
		[][2]int64{{3, 2}, {5, 7}},
		[][2]int64{{11, 13}, {17, 19}},
	)

	want := [][2]int64{
		{3, 2},
		{5, 7},
		{11, 13},
		{17, 19},
		{11, 13},
		{17, 19},
	}

	for i, w := range want {
		p, q, ok := g.NextPQ()
		if !ok {
			t.Fatalf("expected term %d", i)
		}
		if p != w[0] || q != w[1] {
			t.Fatalf("term %d: got (%d,%d), want (%d,%d)", i, p, q, w[0], w[1])
		}
	}
}

func TestPeriodicGCF_EmptyPrefix(t *testing.T) {
	g := NewPeriodicGCF(
		nil,
		[][2]int64{{2, 1}},
	)

	for i := 0; i < 4; i++ {
		p, q, ok := g.NextPQ()
		if !ok {
			t.Fatalf("expected term %d", i)
		}
		if p != 2 || q != 1 {
			t.Fatalf("term %d: got (%d,%d), want (2,1)", i, p, q)
		}
	}
}

func TestPeriodicGCF_CopiesInput(t *testing.T) {
	prefix := [][2]int64{{3, 2}}
	period := [][2]int64{{5, 7}}

	g := NewPeriodicGCF(prefix, period)

	prefix[0] = [2]int64{99, 99}
	period[0] = [2]int64{88, 88}

	p, q, ok := g.NextPQ()
	if !ok {
		t.Fatalf("expected first term")
	}
	if p != 3 || q != 2 {
		t.Fatalf("got (%d,%d), want copied original (3,2)", p, q)
	}

	p, q, ok = g.NextPQ()
	if !ok {
		t.Fatalf("expected second term")
	}
	if p != 5 || q != 7 {
		t.Fatalf("got (%d,%d), want copied original (5,7)", p, q)
	}
}

// periodic_gcf_test.go v1
