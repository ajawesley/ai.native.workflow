package ingest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
)

func ReadNDJSON(r io.Reader) ([]json.RawMessage, []string, error) {
	scanner := bufio.NewScanner(r)

	var (
		valid   []json.RawMessage
		invalid []string
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

		// JSON syntax validation
		var tmp interface{}
		if err := json.Unmarshal(raw, &tmp); err != nil {
			invalid = append(invalid, string(raw))
			continue
		}

		// Require JSON object only
		if _, ok := tmp.(map[string]interface{}); !ok {
			invalid = append(invalid, string(raw))
			continue
		}

		// Valid JSON object → keep raw message
		valid = append(valid, json.RawMessage(raw))
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return valid, invalid, nil
}
