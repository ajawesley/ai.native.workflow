package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"ai.native.workflow/internal/events"
	"ai.native.workflow/internal/summary"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: processor <path-to-ndjson-file>")
	}

	inputPath := os.Args[1]

	f, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("open input: %v", err)
	}
	defer f.Close()

	var (
		validEvents   []json.RawMessage
		invalidEvents []json.RawMessage
	)

	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()

		if len(line) == 0 {
			continue // skip blank lines
		}

		// Defensive copy — scanner reuses its buffer.
		raw := make([]byte, len(line))
		copy(raw, line)

		if err := events.Validate(raw); err != nil {
			invalidEvents = append(invalidEvents, raw)
			continue
		}

		validEvents = append(validEvents, raw)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("scan input: %v", err)
	}

	// --- write invalid events ---
	if len(invalidEvents) > 0 {
		invalidPath := filepath.Join(
			filepath.Dir(inputPath),
			"invalid_events.jsonl",
		)
		if err := writeNDJSON(invalidPath, invalidEvents); err != nil {
			log.Fatalf("write invalid events: %v", err)
		}
		fmt.Fprintf(os.Stderr, "wrote %d invalid event(s) to %s\n", len(invalidEvents), invalidPath)
	}

	// --- compute + print summary ---
	result := summary.Compute(validEvents)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		log.Fatalf("encode summary: %v", err)
	}
}

// writeNDJSON writes each entry as a newline-delimited JSON record.
func writeNDJSON(path string, records []json.RawMessage) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, r := range records {
		if _, err := w.Write(r); err != nil {
			return err
		}
		if err := w.WriteByte('\n'); err != nil {
			return err
		}
	}
	return w.Flush()
}
