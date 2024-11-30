package messages

import "github.com/MatthieuCoder/OrionV3/internal/state"

type Event struct {
	Kind    string      `json:"k"`
	Content state.Event `json:"e"`
}
