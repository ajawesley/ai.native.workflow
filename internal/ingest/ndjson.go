package ingest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"ai.native.workflow/internal/events"
)

// ReadNDJSON reads newline-delimited JSON from r and returns validated events
// alongside the raw bytes of any line that failed unmarshal or validation.
// I/O failures are returned as a non-nil error; per-line errors are not.
func ReadNDJSON(r io.Reader) (valid []events.Event, invalid []json.RawMessage, err error) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Bytes()

		if len(bytes.TrimSpace(line)) == 0 {
			continue // blank / whitespace-only — skip silently
		}

		// Defensive copy: scanner reuses its buffer on every Scan call.
		raw := make([]byte, len(line))
		copy(raw, line)

		var e events.Event
		if err := json.Unmarshal(raw, &e); err != nil {
			invalid = append(invalid, raw)
			continue
		}

		if err := events.Validate(raw); err != nil {
			invalid = append(invalid, raw)
			continue
		}

		valid = append(valid, e)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("ingest: read error: %w", err)
	}

	return valid, invalid, nil
}
