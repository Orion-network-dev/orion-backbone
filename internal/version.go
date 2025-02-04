package internal

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
	BuiltBy = "go"
)

func PrintVersionHeader() {
	log.Info().MsgFunc(func() string {
		message := ("********************************************\n")
		message += fmt.Sprintf(" Orion Version %s           \n", Version)
		message += fmt.Sprintf(" Commit: %s Build time: %s          \n", Commit, Date)
		message += ("********************************************\n")
		return message
	})
}
