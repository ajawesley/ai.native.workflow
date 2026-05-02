package ingest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"

	"ai.native.workflow/internal/events"
)

func ReadNDJSON(r io.Reader) ([]json.RawMessage, []json.RawMessage, error) {
	scanner := bufio.NewScanner(r)

	var (
		valid   []json.RawMessage
		invalid []json.RawMessage
	)

	for scanner.Scan() {
		line := scanner.Bytes()

		// Skip blank or whitespace-only lines
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}

		// Defensive copy — scanner buffer is reused
		raw := make([]byte, len(line))
		copy(raw, line)

		// Validate by unmarshaling into events.Event
		if err := events.Validate(raw); err != nil {
			invalid = append(invalid, json.RawMessage(raw))
			continue
		}

		valid = append(valid, json.RawMessage(raw))
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return valid, invalid, nil
}
