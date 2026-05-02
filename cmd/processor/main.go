package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"ai.native.workflow/internal/ingest"
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

	// --- NDJSON ingestion ---
	validEvents, invalidEvents, err := ingest.ReadNDJSON(f)
	if err != nil {
		log.Fatalf("ingest: %v", err)
	}

	// --- write invalid events ---
	if len(invalidEvents) > 0 {
		invalidPath := filepath.Join(
			filepath.Dir(inputPath),
			"invalid_events.jsonl",
		)
		if err := writeInvalidLines(invalidPath, invalidEvents); err != nil {
			log.Fatalf("write invalid events: %v", err)
		}
		fmt.Fprintf(os.Stderr, "wrote %d invalid event(s) to %s\n", len(invalidEvents), invalidPath)
	}

	// --- compute + print summary ---
	result := summary.Compute(validEvents, len(invalidEvents))

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(result); err != nil {
		log.Fatalf("encode summary: %v", err)
	}
}

// writeInvalidLines writes each invalid line as-is to a .jsonl file.
func writeInvalidLines(path string, lines []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, line := range lines {
		if _, err := w.WriteString(line); err != nil {
			return err
		}
		if err := w.WriteByte('\n'); err != nil {
			return err
		}
	}
	return w.Flush()
}
