package ingest_test

import (
	"encoding/json"
	"strings"
	"testing"

	"ai.native.workflow/internal/ingest"
	"github.com/google/go-cmp/cmp"
)

func TestReadNDJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantValid   []json.RawMessage
		wantInvalid []string
		wantErr     bool
	}{
		{
			name:  "single valid line",
			input: `{"id":1,"type":"alpha"}`,
			wantValid: []json.RawMessage{
				json.RawMessage(`{"id":1,"type":"alpha"}`),
			},
			wantInvalid: nil,
			wantErr:     false,
		},
		{
			name: "multiple valid lines",
			input: `{"id":1,"type":"alpha"}
{"id":2,"type":"beta"}`,
			wantValid: []json.RawMessage{
				json.RawMessage(`{"id":1,"type":"alpha"}`),
				json.RawMessage(`{"id":2,"type":"beta"}`),
			},
			wantInvalid: nil,
			wantErr:     false,
		},
		{
			name: "skip blank lines",
			input: `
{"id":1,"type":"alpha"}

{"id":2,"type":"beta"}
`,
			wantValid: []json.RawMessage{
				json.RawMessage(`{"id":1,"type":"alpha"}`),
				json.RawMessage(`{"id":2,"type":"beta"}`),
			},
			wantInvalid: nil,
			wantErr:     false,
		},
		{
			name: "invalid JSON line",
			input: `{"id":1,"type":"alpha"}
not-json
{"id":2,"type":"beta"}`,
			wantValid: []json.RawMessage{
				json.RawMessage(`{"id":1,"type":"alpha"}`),
				json.RawMessage(`{"id":2,"type":"beta"}`),
			},
			wantInvalid: []string{
				"not-json",
			},
			wantErr: false,
		},
		{
			name: "all invalid lines",
			input: `nope
still bad
1234`,
			wantValid: nil,
			wantInvalid: []string{
				"nope",
				"still bad",
				"1234",
			},
			wantErr: false,
		},
		{
			name:        "empty input",
			input:       ``,
			wantValid:   nil,
			wantInvalid: nil,
			wantErr:     false,
		},
		{
			name: "whitespace only",
			input: `

            
            
`,
			wantValid:   nil,
			wantInvalid: nil,
			wantErr:     false,
		},
		{
			name: "malformed but non-fatal",
			input: `{"id":1}
{"id":2}
{bad json}`,
			wantValid: []json.RawMessage{
				json.RawMessage(`{"id":1}`),
				json.RawMessage(`{"id":2}`),
			},
			wantInvalid: []string{
				"{bad json}",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			valid, invalid, err := ingest.ReadNDJSON(strings.NewReader(tt.input))

			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error state: got err=%v wantErr=%v", err, tt.wantErr)
			}

			if diff := cmp.Diff(valid, tt.wantValid); diff != "" {
				t.Fatalf("valid events mismatch (-got +want):\n%s", diff)
			}

			if diff := cmp.Diff(invalid, tt.wantInvalid); diff != "" {
				t.Fatalf("invalid events mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
