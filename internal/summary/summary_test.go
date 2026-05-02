package summary_test

import (
	"encoding/json"
	"testing"
	"time"

	"ai.native.workflow/internal/summary"
)

func makeRaw(t *testing.T, v any) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("makeRaw: %v", err)
	}
	return json.RawMessage(b)
}

type eventPayload struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

func TestCompute(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		events    []json.RawMessage
		wantValid int
	}{
		{
			name:      "happy path",
			wantValid: 2,
			events: []json.RawMessage{
				makeRaw(t, eventPayload{ID: "1", Type: "credit", Amount: 100, Timestamp: now}),
				makeRaw(t, eventPayload{ID: "2", Type: "debit", Amount: 40, Timestamp: now}),
			},
		},
		{
			name:      "empty slice",
			wantValid: 0,
			events:    []json.RawMessage{},
		},
		{
			name:      "nil slice",
			wantValid: 0,
			events:    nil,
		},
		{
			name:      "mixed credit/debit",
			wantValid: 3,
			events: []json.RawMessage{
				makeRaw(t, eventPayload{ID: "c1", Type: "credit", Amount: 200, Timestamp: now}),
				makeRaw(t, eventPayload{ID: "c2", Type: "credit", Amount: 50, Timestamp: now}),
				makeRaw(t, eventPayload{ID: "d1", Type: "debit", Amount: 75, Timestamp: now}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := summary.Compute(tt.events)

			if got.TotalValid != tt.wantValid {
				t.Errorf("TotalValid: want %d, got %d", tt.wantValid, got.TotalValid)
			}

			if got.TotalInvalid != 0 {
				t.Errorf("TotalInvalid: want 0, got %d", got.TotalInvalid)
			}
		})
	}
}
