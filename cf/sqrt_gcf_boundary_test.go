// cf/sqrt_gcf_boundary_test.go v3
package cf

import "testing"

func requireEqualRational(t *testing.T, got, want Rational, where string) {
	t.Helper()
	if got.Cmp(want) != 0 {
		t.Fatalf("%s: got %v want %v", where, got, want)
	}
}

func requireNonNilRange(t *testing.T, r *Range, where string) {
	t.Helper()
	if r == nil {
		t.Fatalf("%s: got nil range", where)
	}
}

func requireEqualApprox(t *testing.T, got, want GCFApprox, where string) {
	t.Helper()

	requireEqualRational(t, got.Convergent, want.Convergent, where+" convergent")

	if got.PrefixTerms != want.PrefixTerms {
		t.Fatalf("%s prefixTerms: got %d want %d", where, got.PrefixTerms, want.PrefixTerms)
	}

	switch {
	case got.Range == nil && want.Range == nil:
		return
	case got.Range == nil || want.Range == nil:
		t.Fatalf("%s range presence mismatch: got %#v want %#v", where, got.Range, want.Range)
	default:
		requireEqualRational(t, got.Range.Lo, want.Range.Lo, where+" range.lo")
		requireEqualRational(t, got.Range.Hi, want.Range.Hi, where+" range.hi")
		if got.Range.IncLo != want.Range.IncLo {
			t.Fatalf("%s range.incLo: got %v want %v", where, got.Range.IncLo, want.Range.IncLo)
		}
		if got.Range.IncHi != want.Range.IncHi {
			t.Fatalf("%s range.incHi: got %v want %v", where, got.Range.IncHi, want.Range.IncHi)
		}
	}
}

func TestSqrtGCFBoundary_SourceRangeSeedMatchesSnapshotPipeline(t *testing.T) {
	const prefixTerms = 8

	p := DefaultSqrtPolicy2()

	got, err := SqrtApproxFromGCFSourceRangeSeed2(NewECFGSource(), prefixTerms, p)
	if err != nil {
		t.Fatalf("SqrtApproxFromGCFSourceRangeSeed2 error: %v", err)
	}

	a, err := GCFApproxFromPrefix(NewECFGSource(), prefixTerms)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix error: %v", err)
	}

	want, err := SqrtApproxFromGCFApproxRangeSeed2(a, p)
	if err != nil {
		t.Fatalf("SqrtApproxFromGCFApproxRangeSeed2 error: %v", err)
	}

	requireEqualRational(t, got, want, "bounded sqrt path equivalence")
}

func TestSqrtGCFBoundary_PrefixStreamSnapshotMatchesDirectApprox(t *testing.T) {
	const prefixTerms = 5

	p := DefaultSqrtPolicy2()
	stream := NewSqrtGCFPrefixStream2(NewBrouncker4OverPiGCFSource(), prefixTerms, p)

	pre := stream.Snapshot()
	if pre.Started {
		t.Fatalf("pre-init Started: got true want false")
	}
	if pre.Approx != nil {
		t.Fatalf("pre-init Approx: got non-nil want nil")
	}
	if pre.GCFInputApprox != nil {
		t.Fatalf("pre-init GCFInputApprox: got non-nil want nil")
	}
	if pre.PrefixTerms != prefixTerms {
		t.Fatalf("pre-init PrefixTerms: got %d want %d", pre.PrefixTerms, prefixTerms)
	}

	_, _ = stream.Next()

	post := stream.Snapshot()
	if !post.Started {
		t.Fatalf("post-init Started: got false want true")
	}
	if post.Approx == nil {
		t.Fatalf("post-init Approx: got nil want non-nil")
	}
	if post.GCFInputApprox == nil {
		t.Fatalf("post-init GCFInputApprox: got nil want non-nil")
	}

	a, err := GCFApproxFromPrefix(NewBrouncker4OverPiGCFSource(), prefixTerms)
	if err != nil {
		t.Fatalf("GCFApproxFromPrefix error: %v", err)
	}
	wantApprox, err := SqrtApproxFromGCFApproxRangeSeed2(a, p)
	if err != nil {
		t.Fatalf("SqrtApproxFromGCFApproxRangeSeed2 error: %v", err)
	}

	requireEqualApprox(t, *post.GCFInputApprox, a, "stream input snapshot")
	requireEqualRational(t, *post.Approx, wantApprox, "stream approx")
}

func TestSqrtGCFBoundary_PrefixStreamStatusTracksSnapshotExactness(t *testing.T) {
	t.Run("exact finite input", func(t *testing.T) {
		stream := NewSqrtGCFPrefixStream2(
			NewSliceGCF([2]int64{9, 1}),
			1,
			DefaultSqrtPolicy2(),
		)

		_, _ = stream.Next()
		snap := stream.Snapshot()

		if snap.Status != SqrtStreamStatusExactInput {
			t.Fatalf("exact finite status: got %v want %v", snap.Status, SqrtStreamStatusExactInput)
		}
		if snap.GCFInputApprox == nil {
			t.Fatalf("exact finite GCFInputApprox: got nil want non-nil")
		}
		requireNonNilRange(t, snap.GCFInputApprox.Range, "exact finite input range")
		if snap.GCFInputApprox.Range.Lo.Cmp(snap.GCFInputApprox.Range.Hi) != 0 {
			t.Fatalf(
				"exact finite input range not point-valued: lo=%v hi=%v",
				snap.GCFInputApprox.Range.Lo,
				snap.GCFInputApprox.Range.Hi,
			)
		}
	})

	t.Run("bounded infinite input", func(t *testing.T) {
		stream := NewSqrtGCFPrefixStream2(
			NewECFGSource(),
			6,
			DefaultSqrtPolicy2(),
		)

		_, _ = stream.Next()
		snap := stream.Snapshot()

		if snap.Status != SqrtStreamStatusBoundedCollapse {
			t.Fatalf("bounded infinite status: got %v want %v", snap.Status, SqrtStreamStatusBoundedCollapse)
		}
		if snap.GCFInputApprox == nil {
			t.Fatalf("bounded infinite GCFInputApprox: got nil want non-nil")
		}
		requireNonNilRange(t, snap.GCFInputApprox.Range, "bounded infinite input range")
		if snap.GCFInputApprox.Range.Lo.Cmp(snap.GCFInputApprox.Range.Hi) == 0 {
			t.Fatalf(
				"bounded infinite input range unexpectedly exact: lo=%v hi=%v",
				snap.GCFInputApprox.Range.Lo,
				snap.GCFInputApprox.Range.Hi,
			)
		}
	})
}

func TestMVPRadicandBoundary_SnapshotWrapperMatchesCanonicalAssembly(t *testing.T) {
	const (
		fourOverPiPrefixTerms = 6
		ePrefixTerms          = 8
	)

	got, err := MVPRadicandSnapshot(fourOverPiPrefixTerms, ePrefixTerms)
	if err != nil {
		t.Fatalf("MVPRadicandSnapshot error: %v", err)
	}

	want, err := MVPRadicandAssembleSnapshot(fourOverPiPrefixTerms, ePrefixTerms)
	if err != nil {
		t.Fatalf("MVPRadicandAssembleSnapshot error: %v", err)
	}

	requireEqualApprox(t, got, want, "MVP snapshot wrapper equivalence")
}

func TestMVPRadicandBoundary_RootValueMatchesSnapshotThenSqrt(t *testing.T) {
	const (
		fourOverPiPrefixTerms = 6
		ePrefixTerms          = 8
	)

	p := DefaultSqrtPolicy2()

	got, err := MVPRadicandRootValue(fourOverPiPrefixTerms, ePrefixTerms, p)
	if err != nil {
		t.Fatalf("MVPRadicandRootValue error: %v", err)
	}

	a, err := MVPRadicandSnapshot(
		fourOverPiPrefixTerms,
		ePrefixTerms,
	)
	if err != nil {
		t.Fatalf("MVPRadicandSnapshot error: %v", err)
	}

	want, err := SqrtApproxFromGCFApproxRangeSeed2(a, p)
	if err != nil {
		t.Fatalf("SqrtApproxFromGCFApproxRangeSeed2 error: %v", err)
	}

	requireEqualRational(t, got, want, "MVP radicand root wrapper equivalence")
}

// cf/sqrt_gcf_boundary_test.go v3
