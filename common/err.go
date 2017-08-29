package common

import "encoding/json"

type SynqError struct {
	Name    string           `json:"name"`
	Message string           `json:"message"`
	Url     string           `json:"url"`
	Details *json.RawMessage `json:"details"`
}
