package messages

import "github.com/MatthieuCoder/OrionV3/internal/state"

type Event struct {
	Kind    string      `json:"k"`
	Content state.Event `json:"e"`
}

type Hello struct {
	Message  string               `json:"message"`
	Identity state.RouterIdentity `json:"identity"`
	Version  string               `json:"version"`
	Commit   string               `json:"commit"`
	Session  string               `json:"session"`
}
