// gcf_specialized_vs_generic_test.go v1
package cf

import "testing"

func rangeSpan(r Range) Rational {
	span, err := r.Hi.Sub(r.Lo)
	if err != nil {
		panic(err)
	}
	return span
}

func TestLambertSpecializedVsGeneric_Prefix1(t *testing.T) {
	spec, err := LambertPiOver4ApproxFromPrefix(1)
	if err != nil {
		t.Fatalf("LambertPiOver4ApproxFromPrefix failed: %v", err)
	}
	gen, err := GCFApproxFromPrefix(NewLambertPiOver4GCFSource(), 1)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}

	if spec.Range == nil || gen.Range == nil {
		t.Fatalf("expected both ranges to be non-nil")
	}
	if !spec.Range.IsInside() || !gen.Range.IsInside() {
		t.Fatalf("expected both ranges to be inside; spec=%v gen=%v", spec.Range, gen.Range)
	}

	specSpan := rangeSpan(*spec.Range)
	genSpan := rangeSpan(*gen.Range)

	// Prefix-aware Lambert should be tighter here.
	if specSpan.Cmp(genSpan) >= 0 {
		t.Fatalf("expected specialized Lambert range tighter than generic: spec=%v gen=%v", specSpan, genSpan)
	}

	// Concrete ranges:
	// specialized: [3/4,1]
	// generic:     [0,1]
	wantSpecLo := mustRat(3, 4)
	wantSpecHi := mustRat(1, 1)
	if spec.Range.Lo.Cmp(wantSpecLo) != 0 || spec.Range.Hi.Cmp(wantSpecHi) != 0 {
		t.Fatalf("specialized got [%v,%v] want [%v,%v]", spec.Range.Lo, spec.Range.Hi, wantSpecLo, wantSpecHi)
	}

	wantGenLo := mustRat(0, 1)
	wantGenHi := mustRat(1, 1)
	if gen.Range.Lo.Cmp(wantGenLo) != 0 || gen.Range.Hi.Cmp(wantGenHi) != 0 {
		t.Fatalf("generic got [%v,%v] want [%v,%v]", gen.Range.Lo, gen.Range.Hi, wantGenLo, wantGenHi)
	}
}

type lambertLowerBoundOnlySource struct {
	src *LambertPiOver4GCFSource
}

func newLambertLowerBoundOnlySource() *lambertLowerBoundOnlySource {
	return &lambertLowerBoundOnlySource{src: NewLambertPiOver4GCFSource()}
}

func (s *lambertLowerBoundOnlySource) NextPQ() (int64, int64, bool) {
	return s.src.NextPQ()
}

func (s *lambertLowerBoundOnlySource) TailLowerBound() Rational {
	return mustRat(1, 1)
}

func TestBrounckerSpecializedVsGeneric_Prefix1(t *testing.T) {
	spec, err := Brouncker4OverPiApproxFromPrefix(1)
	if err != nil {
		t.Fatalf("Brouncker4OverPiApproxFromPrefix failed: %v", err)
	}
	gen, err := GCFApproxFromPrefix(NewBrouncker4OverPiGCFSource(), 1)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix failed: %v", err)
	}

	if spec.Range == nil || gen.Range == nil {
		t.Fatalf("expected both ranges to be non-nil")
	}
	if !spec.Range.IsInside() || !gen.Range.IsInside() {
		t.Fatalf("expected both ranges to be inside; spec=%v gen=%v", spec.Range, gen.Range)
	}

	specSpan := rangeSpan(*spec.Range)
	genSpan := rangeSpan(*gen.Range)

	// Prefix-aware Brouncker should be tighter here.
	if specSpan.Cmp(genSpan) >= 0 {
		t.Fatalf("expected specialized Brouncker range tighter than generic: spec=%v gen=%v", specSpan, genSpan)
	}

	// specialized: [7/5,3/2]
	// generic:     [1,2]
	wantSpecLo := mustRat(7, 5)
	wantSpecHi := mustRat(3, 2)
	if spec.Range.Lo.Cmp(wantSpecLo) != 0 || spec.Range.Hi.Cmp(wantSpecHi) != 0 {
		t.Fatalf("specialized got [%v,%v] want [%v,%v]", spec.Range.Lo, spec.Range.Hi, wantSpecLo, wantSpecHi)
	}

	wantGenLo := mustRat(1, 1)
	wantGenHi := mustRat(2, 1)
	if gen.Range.Lo.Cmp(wantGenLo) != 0 || gen.Range.Hi.Cmp(wantGenHi) != 0 {
		t.Fatalf("generic got [%v,%v] want [%v,%v]", gen.Range.Lo, gen.Range.Hi, wantGenLo, wantGenHi)
	}
}
func TestLambertSpecializedVsGeneric_Prefix2_SpecializedGain(t *testing.T) {
	spec, err := LambertPiOver4ApproxFromPrefix(2)
	if err != nil {
		t.Fatalf("LambertPiOver4ApproxFromPrefix failed: %v", err)
	}

	gen, err := specializedGCFApproxFromPrefix(
		2,
		func() GCFSource { return newLambertLowerBoundOnlySource() },
		func(prefixTerms int) (Range, bool, error) { return Range{}, false, nil },
		LambertPiOver4TailLowerBoundAfterPrefix,
	)
	if err != nil {
		t.Fatalf("specializedGCFApproxFromPrefix failed: %v", err)
	}

	if spec.Range == nil || gen.Range == nil {
		t.Fatalf("expected both ranges to be non-nil")
	}

	specSpan, err := spec.Range.Hi.Sub(spec.Range.Lo)
	if err != nil {
		t.Fatalf("specialized span failed: %v", err)
	}

	genSpan, err := gen.Range.Hi.Sub(gen.Range.Lo)
	if err != nil {
		t.Fatalf("generic span failed: %v", err)
	}

	if specSpan.Cmp(genSpan) >= 0 {
		t.Fatalf("expected Lambert specialized prefix 2 span to be tighter: spec=%v gen=%v", specSpan, genSpan)
	}
}

// gcf_specialized_vs_generic_test.go v1
