package internal

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

// handleError logs the error message and returns an error with additional context.
func HandleError(err error, msg string) error {
	if err != nil {
		log.Error().Err(err).Msg(msg)
		return fmt.Errorf("%s: %w", msg, err)
	}
	return nil
}
