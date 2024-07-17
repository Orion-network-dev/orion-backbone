package internal

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func handleError(err error, msg string) error {
	if err != nil {
		log.Error().Err(err).Msg(msg)
		return fmt.Errorf("%s: %w", msg, err)
	}
	return nil
}
