package models

import(
	"encoding/json"
)

type Event struct {
	Command string      `json:"command"`
	Entity  interface{} `json:"entity"`
}

func (e *Event) Marshal() (json.RawMessage, error) {
	return json.Marshal(e)
}