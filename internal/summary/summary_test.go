package summary_test

import (
	"encoding/json"
	"testing"
	"time"

	"ai.native.workflow/internal/summary"
)

// makeRaw marshals v into a json.RawMessage or fails the test.
func makeRaw(t *testing.T, v any) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("makeRaw: %v", err)
	}
	return json.RawMessage(b)
}

// eventPayload mirrors the event shape for test isolation.
type eventPayload struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

func TestCompute_HappyPath(t *testing.T) {
	now := time.Now()

	events := []json.RawMessage{
		makeRaw(t, eventPayload{ID: "1", Type: "credit", Amount: 100.00, Timestamp: now}),
		makeRaw(t, eventPayload{ID: "2", Type: "debit", Amount: 40.00, Timestamp: now}),
	}

	got := summary.Compute(events)

	if got.TotalValid != 2 {
		t.Errorf("TotalValid: want 2, got %d", got.TotalValid)
	}
	if got.TotalInvalid != 0 {
		t.Errorf("TotalInvalid: want 0, got %d", got.TotalInvalid)
	}
}

func TestCompute_EmptyInput(t *testing.T) {
	got := summary.Compute([]json.RawMessage{})

	if got.TotalValid != 0 {
		t.Errorf("TotalValid: want 0, got %d", got.TotalValid)
	}
	if got.TotalInvalid != 0 {
		t.Errorf("TotalInvalid: want 0, got %d", got.TotalInvalid)
	}
}

func TestCompute_NilInput(t *testing.T) {
	got := summary.Compute(nil)

	if got.TotalValid != 0 {
		t.Errorf("TotalValid: want 0, got %d", got.TotalValid)
	}
	if got.TotalInvalid != 0 {
		t.Errorf("TotalInvalid: want 0, got %d", got.TotalInvalid)
	}
}

func TestCompute_MixedCreditDebit(t *testing.T) {
	now := time.Now()

	events := []json.RawMessage{
		makeRaw(t, eventPayload{ID: "c1", Type: "credit", Amount: 200.00, Timestamp: now}),
		makeRaw(t, eventPayload{ID: "c2", Type: "credit", Amount: 50.00, Timestamp: now}),
		makeRaw(t, eventPayload{ID: "d1", Type: "debit", Amount: 75.00, Timestamp: now}),
	}

	got := summary.Compute(events)

	if got.TotalValid != 3 {
		t.Errorf("TotalValid: want 3, got %d", got.TotalValid)
	}
	if got.TotalInvalid != 0 {
		t.Errorf("TotalInvalid: want 0, got %d", got.TotalInvalid)
	}
}
