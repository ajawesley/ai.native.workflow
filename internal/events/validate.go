package events

import (
	"encoding/json"
	"errors"
)

func Validate(raw []byte) error {
	var e Event
	if err := json.Unmarshal(raw, &e); err != nil {
		return err
	}

	if e.ID == "" {
		return errors.New("missing id")
	}

	switch e.Type {
	case "credit", "debit":
		// ok
	default:
		return errors.New("invalid type")
	}

	if e.Amount <= 0 {
		return errors.New("amount must be > 0")
	}

	if e.Timestamp.IsZero() {
		return errors.New("missing or invalid timestamp")
	}

	return nil
}
