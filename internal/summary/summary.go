package summary

import "encoding/json"

type Result struct {
	TotalValid   int `json:"total_valid"`
	TotalInvalid int `json:"total_invalid"`
}

func Compute(events []json.RawMessage) Result {
	if events == nil {
		return Result{}
	}
	return Result{
		TotalValid: len(events),
	}
}
