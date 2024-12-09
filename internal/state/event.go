package state

import "encoding/json"

// Meta interface holding the different events
type Event interface{}

type JsonEvent struct {
	Kind    string          `json:"k"`
	Content json.RawMessage `json:"e"`
}
