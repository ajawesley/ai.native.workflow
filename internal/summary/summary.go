package summary

import "encoding/json"

type Result struct {
	TotalValid   int `json:"total_valid"`
	TotalInvalid int `json:"total_invalid"`
}

func Compute(valid []json.RawMessage, invalidCount int) Result {
	return Result{
		TotalValid:   len(valid),
		TotalInvalid: invalidCount,
	}
}
