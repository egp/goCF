// cf/gcf_prefix_state_test.go v1
package cf

import "testing"

func TestNewGcfPrefixState_StartsEmpty(t *testing.T) {
	s := newGcfPrefixState()

	if s.hasApprox() {
		t.Fatalf("hasApprox: got true want false")
	}
}

func TestGcfPrefixState_OneIngest_ProducesOneTermApprox(t *testing.T) {
	s := newGcfPrefixState()

	if err := s.ingestOne(2, 1); err != nil {
		t.Fatalf("ingestOne failed: %v", err)
	}

	if !s.hasApprox() {
		t.Fatalf("hasApprox: got false want true")
	}

	got := s.approx()
	if got.Convergent.Cmp(mustRat(2, 1)) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, mustRat(2, 1))
	}
	if got.PrefixTerms != 1 {
		t.Fatalf("PrefixTerms: got %d want 1", got.PrefixTerms)
	}
}

func TestGcfPrefixState_TwoIngests_ProducesTwoTermApprox(t *testing.T) {
	s := newGcfPrefixState()

	if err := s.ingestOne(2, 1); err != nil {
		t.Fatalf("first ingestOne failed: %v", err)
	}
	if err := s.ingestOne(1, 1); err != nil {
		t.Fatalf("second ingestOne failed: %v", err)
	}

	got := s.approx()
	if got.Convergent.Cmp(mustRat(3, 1)) != 0 {
		t.Fatalf("got %v want %v", got.Convergent, mustRat(3, 1))
	}
	if got.PrefixTerms != 2 {
		t.Fatalf("PrefixTerms: got %d want 2", got.PrefixTerms)
	}
}
