// stream_exact_point_test.go v1
package cf

import "testing"

func TestExactPointTermination_BeforeAnyDigitIsError(t *testing.T) {
	done, err := exactPointTermination("ULFTStream:", false, "denominator is zero at exact point x=1")
	if done {
		t.Fatalf("expected done=false")
	}
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
}

func TestExactPointTermination_AfterEmitIsCleanExhaustion(t *testing.T) {
	done, err := exactPointTermination("ULFTStream:", true, "denominator is zero at exact point x=1")
	if !done {
		t.Fatalf("expected done=true")
	}
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
}

// stream_exact_point_test.go v1
