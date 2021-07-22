package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type LogEvent struct {
	id      int
	message string
}

var (
	resourceCreatedMessage = LogEvent{1001, "Resource created.\nContent:\n%s"}
	resourceUpdatedMessage = LogEvent{1002, "Resource updated.\nBefore:\n%s\nAfter:\n%s"}
	resourceReadMessage    = LogEvent{1003, "Resource read.\nContent:\n%s"}
	resourceDeletedMessage = LogEvent{1004, "Resource %s deleted."}
)

type ResourceCreated struct{ Resource string }
type ResourceUpdated struct {
	PastResource string
	NewResource  string
}
type ResourceRead struct{ Resource string }
type ResourceDeleted struct{ Resource string }

func init() {
	log.Logger = zerolog.New(os.Stdout).Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFieldFormat}).With().Timestamp().Logger()
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func (e ResourceCreated) Info() string {
	return fmt.Sprintf(resourceCreatedMessage.message, e.Resource)
}

func (e ResourceUpdated) Info() string {
	return fmt.Sprintf(resourceUpdatedMessage.message, e.PastResource, e.NewResource)
}

func (e ResourceRead) Info() string {
	return fmt.Sprintf(resourceReadMessage.message, e.Resource)
}

func (e ResourceDeleted) Info() string {
	return fmt.Sprintf(resourceReadMessage.message, e.Resource)
}
