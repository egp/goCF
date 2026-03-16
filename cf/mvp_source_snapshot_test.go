// mvp_source_snapshot_test.go v1
package cf

import "testing"

func TestMVPDefaultESourceFunc_UsesCanonicalSource(t *testing.T) {
	src, err := MVPNilSafeGCFSourceFromFunc(MVPDefaultESourceFunc())
	if err != nil {
		t.Fatalf("MVPNilSafeGCFSourceFromFunc failed: %v", err)
	}

	got := collectPQ(src, 7)
	want := [][2]int64{
		{2, 1},
		{1, 1},
		{2, 1},
		{1, 1},
		{1, 1},
		{4, 1},
		{1, 1},
	}
	if !equalPQ(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestMVPNilSafeGCFSourceFromFunc_RejectsNil(t *testing.T) {
	_, err := MVPNilSafeGCFSourceFromFunc(nil)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMVPApproxSnapshotFromSourceFunc_EPath(t *testing.T) {
	got, err := MVPApproxSnapshotFromSourceFunc(
		MVPDefaultESourceFunc(),
		6,
	)
	if err != nil {
		t.Fatalf("MVPApproxSnapshotFromSourceFunc failed: %v", err)
	}

	want, err := GCFSourceConvergent(NewECFGSource(), 6)
	if err != nil {
		t.Fatalf("GCFSourceConvergent failed: %v", err)
	}

	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}

func TestMVPApproxSnapshotFromSourceFunc_RejectsBadPrefixTerms(t *testing.T) {
	_, err := MVPApproxSnapshotFromSourceFunc(
		MVPDefaultFourOverPiSourceFunc(),
		0,
	)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestMVPDefaultEApproxSnapshot_MatchesEConvergent(t *testing.T) {
	got, err := MVPDefaultEApproxSnapshot(6)
	if err != nil {
		t.Fatalf("MVPDefaultEApproxSnapshot failed: %v", err)
	}

	want, err := GCFSourceConvergent(NewECFGSource(), 6)
	if err != nil {
		t.Fatalf("GCFSourceConvergent failed: %v", err)
	}

	if got.Convergent.Cmp(want) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, want)
	}
}
