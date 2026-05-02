package events

import "encoding/json"

// Validate returns a non-nil error if raw is not a valid event.
// TODO: implement field-level validation.
func Validate(raw []byte) error {
	var m map[string]any
	return json.Unmarshal(raw, &m) // bare minimum: must be a JSON object
}
