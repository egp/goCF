// stream_refine_budget_test.go v1
package cf

import "testing"

func TestConsumeRefineBudget_AllowsUnlimited(t *testing.T) {
	thisDigit := 0
	total := 0

	err := consumeRefineBudget("ULFTStream:", &thisDigit, &total, -1, -1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if thisDigit != 1 || total != 1 {
		t.Fatalf("got (%d,%d) want (1,1)", thisDigit, total)
	}
}

func TestConsumeRefineBudget_PerDigitExceeded(t *testing.T) {
	thisDigit := 0
	total := 0

	err := consumeRefineBudget("ULFTStream:", &thisDigit, &total, 0, -1)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestConsumeRefineBudget_TotalExceeded(t *testing.T) {
	thisDigit := 0
	total := 0

	err := consumeRefineBudget("ULFTStream:", &thisDigit, &total, -1, 0)
	if err == nil {
		t.Fatalf("expected error")
	}
}

// stream_refine_budget_test.go v1
