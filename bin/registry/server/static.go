package server

import (
	"embed"
)

//go:embed static
var assets embed.FS

//go:embed errors
var errors embed.FS
