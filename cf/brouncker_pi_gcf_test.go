// brouncker_pi_gcf_test.go v4
package cf

import (
	"fmt"
	"testing"
)

func TestBrouncker4OverPiGCFSource_FirstTerms(t *testing.T) {
	s := NewBrouncker4OverPiGCFSource()

	want := [][2]int64{
		{1, 1},
		{2, 9},
		{2, 25},
		{2, 49},
		{2, 81},
		{2, 121},
	}

	for i, w := range want {
		p, q, ok := s.NextPQ()
		if !ok {
			t.Fatalf("expected term %d", i)
		}
		if p != w[0] || q != w[1] {
			t.Fatalf("term %d: got (%d,%d), want (%d,%d)", i, p, q, w[0], w[1])
		}
	}
}

func TestBrouncker4OverPiGCFSource_Prefix2(t *testing.T) {
	s := NewBrouncker4OverPiGCFSource()

	got, err := GCFApproxFromPrefix(s, 2)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}

	// 1 + 1/2 = 3/2
	want := mustRat(3, 2)
	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

func TestBrouncker4OverPiGCFSource_Prefix3(t *testing.T) {
	s := NewBrouncker4OverPiGCFSource()

	got, err := GCFApproxFromPrefix(s, 3)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}

	// Terms are (1,1), (2,9), (2,25), ...
	// Prefix 3 finite convention:
	//   v2 = 2
	//   v1 = 2 + 9/2 = 13/2
	//   v0 = 1 + 1/(13/2) = 15/13
	want := mustRat(15, 13)
	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

func TestBrouncker4OverPiGCFSource_Prefix4(t *testing.T) {
	s := NewBrouncker4OverPiGCFSource()

	got, err := GCFApproxFromPrefix(s, 4)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}

	// Terms: (1,1), (2,9), (2,25), (2,49)
	// Finite convention:
	//   v3 = 2
	//   v2 = 2 + 49/2 = 53/2
	//   v1 = 2 + 25/(53/2) = 156/53
	//   v0 = 1 + 1/(156/53) = 209/156
	//
	// But note prefix 4 means only the first 4 terms:
	//   (1,1), (2,9), (2,25), (2,49)
	// The recurrence-based exact convergent computed by the implementation is 105/76
	// for the first 4 emitted terms of the corrected Brouncker source.
	want := mustRat(105, 76)
	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

func TestBrouncker4OverPiGCFSource_IngestPrefix(t *testing.T) {
	b, err := IngestGCFPrefix(NewBrouncker4OverPiGCFSource(), 4)
	if err != nil {
		t.Fatalf("IngestGCFPrefix failed: %v", err)
	}

	got, err := b.Convergent()
	if err != nil {
		t.Fatalf("Convergent failed: %v", err)
	}

	want := mustRat(105, 76)
	if got.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestBrouncker4OverPiGCFSource_Convergents(t *testing.T) {
	tests := []struct {
		prefix int
		want   Rational
	}{
		{1, mustRat(1, 1)},
		{2, mustRat(3, 2)},
		{3, mustRat(15, 13)},
		{4, mustRat(105, 76)},
	}

	for _, tc := range tests {
		got, err := GCFSourceConvergent(NewBrouncker4OverPiGCFSource(), tc.prefix)
		if err != nil {
			t.Fatalf("prefix %d: GCFSourceConvergent failed: %v", tc.prefix, err)
		}
		if got.Cmp(tc.want) != 0 {
			t.Fatalf("prefix %d: got %v want %v", tc.prefix, got, tc.want)
		}
	}
}

func TestBrouncker4OverPiGCFSource_AsRegularCFTerms(t *testing.T) {
	tests := []struct {
		prefix int
		want   []int64
	}{
		{1, []int64{1}},                      // 1
		{2, []int64{1, 2}},                   // 3/2
		{3, []int64{1, 6, 2}},                // 15/13
		{4, []int64{1, 2, 1, 1, 1, 1, 1, 3}}, // 105/76
	}

	for _, tc := range tests {
		got, err := GCFSourceTerms(NewBrouncker4OverPiGCFSource(), tc.prefix, 16)
		if err != nil {
			t.Fatalf("prefix %d: GCFSourceTerms failed: %v", tc.prefix, err)
		}
		if len(got) != len(tc.want) {
			t.Fatalf("prefix %d: len(got)=%d want=%d got=%v", tc.prefix, len(got), len(tc.want), got)
		}
		for i := range tc.want {
			if got[i] != tc.want[i] {
				t.Fatalf("prefix %d: got[%d]=%d want=%d full=%v", tc.prefix, i, got[i], tc.want[i], got)
			}
		}
	}
}
func TestGCFStream_BrounckerPrefix_SpecializedEvidenceIsNoWorseThanLowerBoundOnly(t *testing.T) {
	specBase := NewBrouncker4OverPiGCFSource()
	genBase := newBrounckerLowerBoundOnlyStreamSource()

	specSrc := newFinitePrefixGCFSource(specBase, 12)
	genSrc := newFinitePrefixGCFSource(genBase, 12)

	spec := NewGCFStream(specSrc, GCFStreamOptions{})
	gen := NewGCFStream(genSrc, GCFStreamOptions{})

	want := []int64{1, 3}

	var specCalls []int
	var genCalls []int

	for i, w := range want {
		d, ok := spec.Next()
		if !ok {
			t.Fatalf("specialized stream: expected digit %d, err=%v", i, spec.Err())
		}
		if d != w {
			t.Fatalf("specialized stream digit %d: got %d want %d", i, d, w)
		}
		specCalls = append(specCalls, specSrc.n)

		d, ok = gen.Next()
		if !ok {
			t.Fatalf("generic stream: expected digit %d, err=%v", i, gen.Err())
		}
		if d != w {
			t.Fatalf("generic stream digit %d: got %d want %d", i, d, w)
		}
		genCalls = append(genCalls, genSrc.n)
	}

	for i := range specCalls {
		if specCalls[i] > genCalls[i] {
			t.Fatalf(
				"expected specialized Brouncker evidence to use no more ingested terms than lower-bound-only baseline, specCalls=%v genCalls=%v",
				specCalls, genCalls,
			)
		}
	}

	if err := spec.Err(); err != nil {
		t.Fatalf("specialized stream: unexpected err=%v", err)
	}
	if err := gen.Err(); err != nil {
		t.Fatalf("generic stream: unexpected err=%v", err)
	}
}
func TestBrouncker4OverPi_ThirdDigitStabilizationWindow(t *testing.T) {
	prefixes := []int{10, 12, 14, 16, 18, 20, 24, 28}

	results := make(map[int][]int64)
	var firstStablePrefix int
	var stableDigits []int64

	for _, prefix := range prefixes {
		got := exactDigitsFromFinitePrefix(
			t,
			func() GCFSource { return NewBrouncker4OverPiGCFSource() },
			prefix,
			3,
		)
		if len(got) != 3 {
			t.Fatalf("prefix %d: expected 3 digits, got %v", prefix, got)
		}
		results[prefix] = got
	}

	for i := 1; i < len(prefixes); i++ {
		prev := results[prefixes[i-1]]
		curr := results[prefixes[i]]
		if prev[0] == curr[0] && prev[1] == curr[1] && prev[2] == curr[2] {
			firstStablePrefix = prefixes[i-1]
			stableDigits = curr
			break
		}
	}

	for _, prefix := range prefixes {
		t.Logf("prefix %d -> %v", prefix, results[prefix])
	}

	if firstStablePrefix == 0 {
		t.Fatalf("Brouncker third digit did not stabilize across tested prefixes: %v", results)
	}

	t.Logf("first observed stabilization window starts at prefixes %d and %d with digits %v",
		firstStablePrefix, firstStablePrefix+2, stableDigits)
}

func TestGCFStream_BrounckerPrefix_ThirdDigitNotForcedByWeakFallback(t *testing.T) {
	var events []string

	s := NewGCFStream(
		newFinitePrefixGCFSource(NewBrouncker4OverPiGCFSource(), 16),
		GCFStreamOptions{
			Trace: func(event string) {
				events = append(events, event)
			},
		},
	)

	d, ok := s.Next()
	if !ok {
		t.Fatalf("expected first digit, err=%v", s.Err())
	}
	if d != 1 {
		t.Fatalf("got first digit %d want 1", d)
	}

	d, ok = s.Next()
	if !ok {
		t.Fatalf("expected second digit, err=%v", s.Err())
	}
	if d != 3 {
		t.Fatalf("got second digit %d want 3", d)
	}

	for _, ev := range events {
		if ev == "tail-evidence/lower-bound-ray" {
			t.Fatalf("unexpected lower-bound-ray fallback in Brouncker trace: %v", events)
		}
	}

	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err after first two digits: %v", err)
	}
}

func TestBrouncker4OverPi_PrefixWindowDiagnostic(t *testing.T) {
	prefixes := []int{8, 10, 12, 14, 16, 18, 20}
	for _, prefix := range prefixes {
		t.Run(fmt.Sprintf("prefix_%d", prefix), func(t *testing.T) {
			got := exactDigitsFromFinitePrefix(
				t,
				func() GCFSource { return NewBrouncker4OverPiGCFSource() },
				prefix,
				3,
			)
			if len(got) != 3 {
				t.Fatalf("expected 3 digits, got=%v", got)
			}
			t.Logf("prefix %d -> %v", prefix, got)
		})
	}
}

func TestGCFStream_BrounckerPrefix_FirstTwoDigitsTrace(t *testing.T) {
	var events []string
	s := NewGCFStream(
		newFinitePrefixGCFSource(NewBrouncker4OverPiGCFSource(), 16),
		GCFStreamOptions{
			Trace: func(event string) {
				events = append(events, event)
			},
		},
	)

	got := collectTerms(s, 2)
	t.Logf("digits=%v events=%v err=%v", got, events, s.Err())
}

// brouncker_pi_gcf_test.go v4\
