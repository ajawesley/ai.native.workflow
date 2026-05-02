package summary

import "encoding/json"

// Result is the shape of the final summary output.
// TODO: add real fields.
type Result struct {
	TotalValid   int `json:"total_valid"`
	TotalInvalid int `json:"total_invalid"`
}

// Compute derives a summary from the valid event set.
// TODO: implement aggregation logic.
func Compute(events []json.RawMessage) Result {
	return Result{TotalValid: len(events)}
}
