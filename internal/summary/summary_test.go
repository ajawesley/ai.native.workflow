package summary_test

import (
	"encoding/json"
	"testing"

	"ai.native.workflow/internal/summary"
	"github.com/google/go-cmp/cmp"
)

func TestCompute(t *testing.T) {
	tests := []struct {
		name         string
		valid        []json.RawMessage
		invalidCount int
		want         summary.Result
	}{
		{
			name:         "nil valid slice",
			valid:        nil,
			invalidCount: 0,
			want: summary.Result{
				TotalValid:   0,
				TotalInvalid: 0,
			},
		},
		{
			name:         "empty valid slice",
			valid:        []json.RawMessage{},
			invalidCount: 0,
			want: summary.Result{
				TotalValid:   0,
				TotalInvalid: 0,
			},
		},
		{
			name: "one valid event",
			valid: []json.RawMessage{
				json.RawMessage(`{"id":1}`),
			},
			invalidCount: 0,
			want: summary.Result{
				TotalValid:   1,
				TotalInvalid: 0,
			},
		},
		{
			name: "multiple valid events",
			valid: []json.RawMessage{
				json.RawMessage(`{"id":1}`),
				json.RawMessage(`{"id":2}`),
				json.RawMessage(`{"id":3}`),
			},
			invalidCount: 0,
			want: summary.Result{
				TotalValid:   3,
				TotalInvalid: 0,
			},
		},
		{
			name:         "invalid only",
			valid:        nil,
			invalidCount: 5,
			want: summary.Result{
				TotalValid:   0,
				TotalInvalid: 5,
			},
		},
		{
			name: "mixed valid + invalid",
			valid: []json.RawMessage{
				json.RawMessage(`{"id":1}`),
				json.RawMessage(`{"id":2}`),
			},
			invalidCount: 3,
			want: summary.Result{
				TotalValid:   2,
				TotalInvalid: 3,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := summary.Compute(tt.valid, tt.invalidCount)

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Fatalf("summary mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
