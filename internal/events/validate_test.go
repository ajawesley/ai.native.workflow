package events_test

import (
	"encoding/json"
	"testing"
	"time"

	"ai.native.workflow/internal/events"
)

func mustMarshal(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("mustMarshal: %v", err)
	}
	return b
}

func TestValidate(t *testing.T) {
	now := time.Now()

	validEvent := events.Event{
		ID:        "evt-001",
		Type:      "credit",
		Amount:    99.99,
		Timestamp: now,
	}

	tests := []struct {
		name    string
		event   events.Event
		wantErr bool
	}{
		{
			name:    "valid credit event",
			event:   validEvent,
			wantErr: false,
		},
		{
			name:    "valid debit event",
			event:   events.Event{ID: "evt-002", Type: "debit", Amount: 10.00, Timestamp: now},
			wantErr: false,
		},
		{
			name:    "empty ID",
			event:   events.Event{ID: "", Type: "credit", Amount: 50.00, Timestamp: now},
			wantErr: true,
		},
		{
			name:    "invalid type",
			event:   events.Event{ID: "evt-003", Type: "refund", Amount: 50.00, Timestamp: now},
			wantErr: true,
		},
		{
			name:    "unknown type empty string",
			event:   events.Event{ID: "evt-004", Type: "", Amount: 50.00, Timestamp: now},
			wantErr: true,
		},
		{
			name:    "zero amount",
			event:   events.Event{ID: "evt-005", Type: "debit", Amount: 0, Timestamp: now},
			wantErr: true,
		},
		{
			name:    "negative amount",
			event:   events.Event{ID: "evt-006", Type: "credit", Amount: -1.00, Timestamp: now},
			wantErr: true,
		},
		{
			name:    "zero-value timestamp",
			event:   events.Event{ID: "evt-007", Type: "debit", Amount: 5.00, Timestamp: time.Time{}},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			raw := mustMarshal(t, tc.event)
			err := events.Validate(raw)
			if tc.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("expected no error, got: %v", err)
			}
		})
	}
}
