# NDJSON Event Processor

A small command‑line tool that ingests newline‑delimited JSON (NDJSON), validates
each line, writes invalid records to a sidecar file, and prints a summary of
valid/invalid counts as JSON.

This implementation is intentionally minimal and focused on correctness,
testability, and clarity — aligned with a 2‑hour coding challenge scope.

---

## Features

- Reads NDJSON from a file
- Skips blank/whitespace lines
- Validates each line using domain‑level validation (`events.Validate`)
- Collects valid events as `json.RawMessage`
- Collects invalid lines as raw strings
- Writes invalid lines to `invalid_events.jsonl` next to the input file
- Prints a JSON summary to stdout

---

## Usage

```bash
go run ./cmd/processor <path-to-ndjson-file>
