package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type LogEvent struct {
	id      int
	message string
}

func init() {
	log.Logger = zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFieldFormat}).With().Timestamp().Logger()
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}
