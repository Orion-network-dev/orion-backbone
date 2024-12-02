package messages

import (
	"github.com/MatthieuCoder/OrionV3/internal/state"
)

type Hello struct {
	Message  string               `json:"message"`
	Identity state.RouterIdentity `json:"identity"`
	Version  string               `json:"version"`
	Commit   string               `json:"commit"`
	Session  string               `json:"session"`
}
