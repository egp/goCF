// gcf_stream_trace_test.go v1
package cf

import "testing"

func TestGCFStream_Trace_ShowsPoleFallbackThenIngest(t *testing.T) {
	src := &reusableOracleTailRangeGCFSource{}
	var events []string

	s := NewGCFStream(src, GCFStreamOptions{
		Trace: func(event string) {
			events = append(events, event)
		},
	})

	got := collectTerms(s, 3)
	if len(got) != 3 {
		t.Fatalf("expected 3 digits, got=%v err=%v", got, s.Err())
	}

	foundPoleFallback := false
	foundIngestAfterPole := false
	for i, ev := range events {
		if ev == "tail-evidence/pole-fallback" {
			foundPoleFallback = true
			for j := i + 1; j < len(events); j++ {
				if events[j] == "tail-evidence/ingest" {
					foundIngestAfterPole = true
					break
				}
			}
			break
		}
	}

	if !foundPoleFallback {
		t.Fatalf("expected trace to include tail-evidence/pole-fallback, got %v", events)
	}
	if !foundIngestAfterPole {
		t.Fatalf("expected ingest after pole fallback, got %v", events)
	}

	if err := s.Err(); err != nil {
		t.Fatalf("unexpected stream err: %v", err)
	}
}

// gcf_stream_trace_test.go v1
