package events

import "time"

// Event represents a single financial event in the system.
// It is intentionally simple and uses only standard library types.
type Event struct {
	ID        string
	Type      string
	Amount    float64
	Timestamp time.Time
}
