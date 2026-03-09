// gcf_named_enclosure_test.go v1
package cf

import "testing"

func TestBrouncker4OverPiGCFSource_EnclosureContainsConvergent_ForSeveralPrefixes(t *testing.T) {
	for prefix := 1; prefix <= 6; prefix++ {
		b, err := IngestGCFPrefix(NewBrouncker4OverPiGCFSource(), prefix)
		if err != nil {
			t.Fatalf("prefix %d: IngestGCFPrefix failed: %v", prefix, err)
		}

		conv, err := b.Convergent()
		if err != nil {
			t.Fatalf("prefix %d: Convergent failed: %v", prefix, err)
		}

		r, ok, err := b.Range()
		if err != nil {
			t.Fatalf("prefix %d: Range failed: %v", prefix, err)
		}
		if !ok {
			t.Fatalf("prefix %d: expected ok=true", prefix)
		}
		if !r.IsInside() {
			t.Fatalf("prefix %d: expected inside range, got %v", prefix, r)
		}
		if !r.Contains(conv) {
			t.Fatalf("prefix %d: range %v does not contain convergent %v", prefix, r, conv)
		}
	}
}

func TestLambertPiOver4GCFSource_EnclosureContainsConvergent_ForSeveralPrefixes(t *testing.T) {
	for prefix := 1; prefix <= 6; prefix++ {
		b, err := IngestGCFPrefix(NewLambertPiOver4GCFSource(), prefix)
		if err != nil {
			t.Fatalf("prefix %d: IngestGCFPrefix failed: %v", prefix, err)
		}

		conv, err := b.Convergent()
		if err != nil {
			t.Fatalf("prefix %d: Convergent failed: %v", prefix, err)
		}

		r, ok, err := b.Range()
		if err != nil {
			t.Fatalf("prefix %d: Range failed: %v", prefix, err)
		}
		if !ok {
			t.Fatalf("prefix %d: expected ok=true", prefix)
		}
		if !r.IsInside() {
			t.Fatalf("prefix %d: expected inside range, got %v", prefix, r)
		}
		if !r.Contains(conv) {
			t.Fatalf("prefix %d: range %v does not contain convergent %v", prefix, r, conv)
		}
	}
}

func TestBrouncker4OverPiGCFSource_Prefix2And3ConcreteRanges(t *testing.T) {
	tests := []struct {
		prefix int
		wantLo Rational
		wantHi Rational
	}{
		// Prefix 2: x = 1 + 1/(2 + 1/tail), tail >= 1
		// Range = [4/3, 3/2]
		{2, mustRat(4, 3), mustRat(3, 2)},

		// Prefix 3: x = 1 + 1/(2 + 1/(2 + 9/tail)), tail >= 1
		// tail=1 => 19/14, tail=∞ => 7/5
		{3, mustRat(7, 5), mustRat(34, 23)},
	}

	for _, tc := range tests {
		b, err := IngestGCFPrefix(NewBrouncker4OverPiGCFSource(), tc.prefix)
		if err != nil {
			t.Fatalf("prefix %d: IngestGCFPrefix failed: %v", tc.prefix, err)
		}

		r, ok, err := b.Range()
		if err != nil {
			t.Fatalf("prefix %d: Range failed: %v", tc.prefix, err)
		}
		if !ok {
			t.Fatalf("prefix %d: expected ok=true", tc.prefix)
		}

		if r.Lo.Cmp(tc.wantLo) != 0 || r.Hi.Cmp(tc.wantHi) != 0 {
			t.Fatalf("prefix %d: got [%v,%v] want [%v,%v]", tc.prefix, r.Lo, r.Hi, tc.wantLo, tc.wantHi)
		}
	}
}

func TestLambertPiOver4GCFSource_Prefix2And3ConcreteRanges(t *testing.T) {
	tests := []struct {
		prefix int
		wantLo Rational
		wantHi Rational
	}{
		// Prefix 2: x = 0 + 1/(1 + 1/tail), tail >= 1
		// Range = [1/2, 1]
		{2, mustRat(1, 2), mustRat(1, 1)},

		// Prefix 3: x = 0 + 1/(1 + 1/(3 + 4/tail)), tail >= 1
		// tail=1 => 4/5, tail=∞ => 3/4, so range [3/4, 4/5]
		{3, mustRat(3, 4), mustRat(7, 8)},
	}

	for _, tc := range tests {
		b, err := IngestGCFPrefix(NewLambertPiOver4GCFSource(), tc.prefix)
		if err != nil {
			t.Fatalf("prefix %d: IngestGCFPrefix failed: %v", tc.prefix, err)
		}

		r, ok, err := b.Range()
		if err != nil {
			t.Fatalf("prefix %d: Range failed: %v", tc.prefix, err)
		}
		if !ok {
			t.Fatalf("prefix %d: expected ok=true", tc.prefix)
		}

		if r.Lo.Cmp(tc.wantLo) != 0 || r.Hi.Cmp(tc.wantHi) != 0 {
			t.Fatalf("prefix %d: got [%v,%v] want [%v,%v]", tc.prefix, r.Lo, r.Hi, tc.wantLo, tc.wantHi)
		}
	}
}

func TestNamedGCFEnclosures_AreBoundedAndPositive_ForSeveralPrefixes(t *testing.T) {
	sources := []struct {
		name string
		new  func() GCFSource
	}{
		{"Brouncker", func() GCFSource { return NewBrouncker4OverPiGCFSource() }},
		{"Lambert", func() GCFSource { return NewLambertPiOver4GCFSource() }},
	}

	for _, src := range sources {
		for prefix := 2; prefix <= 8; prefix++ {
			b, err := IngestGCFPrefix(src.new(), prefix)
			if err != nil {
				t.Fatalf("%s prefix %d: IngestGCFPrefix failed: %v", src.name, prefix, err)
			}

			r, ok, err := b.Range()
			if err != nil {
				t.Fatalf("%s prefix %d: Range failed: %v", src.name, prefix, err)
			}
			if !ok {
				t.Fatalf("%s prefix %d: expected ok=true", src.name, prefix)
			}
			if !r.IsInside() {
				t.Fatalf("%s prefix %d: expected inside range, got %v", src.name, prefix, r)
			}
			if r.Lo.Cmp(intRat(0)) <= 0 {
				t.Fatalf("%s prefix %d: expected positive lower bound, got %v", src.name, prefix, r.Lo)
			}
		}
	}
}

// gcf_named_enclosure_test.go v1
